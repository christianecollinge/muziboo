// Package library provides resource access to track data.
package library

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/muziboo/api/foundation/supabase"
)

// Track represents a music track.
type Track struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Genre         string `json:"genre"`
	AudioURL      string `json:"audio_url"`
	ArtworkURL    string `json:"artwork_url"`
	Duration      int    `json:"duration"`
	CreatedAt     string `json:"created_at"`
	IsColab       bool   `json:"is_colab"`
	ParentTrackID string `json:"parent_track_id,omitempty"`
	IsPublic      bool   `json:"is_public"`
}

// TrackWithProfile is a track with its owner's profile info.
type TrackWithProfile struct {
	Track
	Profile TrackProfile `json:"profiles"`
}

// TrackProfile is the embedded profile in a track query.
type TrackProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// NewTrack is the data needed to create a track.
type NewTrack struct {
	UserID        string `json:"user_id"`
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	Genre         string `json:"genre,omitempty"`
	AudioURL      string `json:"audio_url"`
	ArtworkURL    string `json:"artwork_url,omitempty"`
	Duration      int    `json:"duration,omitempty"`
	IsColab       bool   `json:"is_colab"`
	ParentTrackID string `json:"parent_track_id,omitempty"`
	IsPublic      bool   `json:"is_public"`
}

// Stem represents an individual audio layer of a Colab project.
type Stem struct {
	ID        string `json:"id"`
	TrackID   string `json:"track_id"`
	Name      string `json:"name"`
	AudioURL  string `json:"audio_url"`
	CreatedAt string `json:"created_at"`
}

// NewStem is the data needed to create a stem.
type NewStem struct {
	TrackID  string `json:"track_id"`
	Name     string `json:"name"`
	AudioURL string `json:"audio_url"`
}

// TrackInvite represents a user invited to a Colab.
type TrackInvite struct {
	ID            string `json:"id"`
	TrackID       string `json:"track_id"`
	InvitedUserID string `json:"invited_user_id"`
	CreatedAt     string `json:"created_at"`
}

// NewTrackInvite is the data needed to invite a user.
type NewTrackInvite struct {
	TrackID       string `json:"track_id"`
	InvitedUserID string `json:"invited_user_id"`
}

// Access manages track operations.
type Access struct {
	client *supabase.Client
}

// NewAccess creates a new library Access.
func NewAccess(client *supabase.Client) *Access {
	return &Access{client: client}
}

// objectPath extracts the bucket-relative object path from a stored
// public-form storage URL. Returns "" if the URL doesn't point at the bucket.
func objectPath(rawURL, bucket string) string {
	marker := "/storage/v1/object/public/" + bucket + "/"
	idx := strings.Index(rawURL, marker)
	if idx < 0 {
		return ""
	}
	return rawURL[idx+len(marker):]
}

// signURLs replaces stored public-form storage URLs with time-limited signed
// URLs so files in private buckets stay playable. URLs that don't point at
// the bucket are left untouched; on signing failure the originals remain
// (playback will fail, but pages still render).
func (a *Access) signURLs(bucket string, urls []*string) {
	ptrsByPath := make(map[string][]*string)
	var paths []string
	for _, u := range urls {
		if u == nil || *u == "" {
			continue
		}
		p := objectPath(*u, bucket)
		if p == "" {
			continue
		}
		if _, seen := ptrsByPath[p]; !seen {
			paths = append(paths, p)
		}
		ptrsByPath[p] = append(ptrsByPath[p], u)
	}

	if len(paths) == 0 {
		return
	}

	signed, err := a.client.SignURLs(bucket, paths, supabase.SignedURLTTL)
	if err != nil {
		log.Printf("signing %s urls: %v", bucket, err)
		return
	}

	for p, ptrs := range ptrsByPath {
		if s, ok := signed[p]; ok {
			for _, ptr := range ptrs {
				*ptr = s
			}
		}
	}
}

// signTrackAudio signs the audio URLs of a slice of tracks in one call.
func (a *Access) signTrackAudio(tracks []Track) {
	urls := make([]*string, 0, len(tracks))
	for i := range tracks {
		urls = append(urls, &tracks[i].AudioURL)
	}
	a.signURLs("audio", urls)
}

// List returns all live tracks ordered by newest first, with profile info.
func (a *Access) List(limit, offset int) ([]TrackWithProfile, error) {
	query := fmt.Sprintf(
		"select=*,profiles(username,display_name,avatar_url)&is_public=eq.true&order=created_at.desc&limit=%d&offset=%d",
		limit, offset,
	)

	data, err := a.client.Query("tracks", query)
	if err != nil {
		return nil, fmt.Errorf("querying tracks: %w", err)
	}

	var tracks []TrackWithProfile
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("decoding tracks: %w", err)
	}

	urls := make([]*string, 0, len(tracks))
	for i := range tracks {
		urls = append(urls, &tracks[i].AudioURL)
	}
	a.signURLs("audio", urls)

	return tracks, nil
}

// GetByID returns a single track by ID.
func (a *Access) GetByID(id string) (TrackWithProfile, error) {
	query := fmt.Sprintf(
		"select=*,profiles(username,display_name,avatar_url)&id=eq.%s",
		url.QueryEscape(id),
	)

	data, err := a.client.Query("tracks", query)
	if err != nil {
		return TrackWithProfile{}, fmt.Errorf("querying track: %w", err)
	}

	var tracks []TrackWithProfile
	if err := json.Unmarshal(data, &tracks); err != nil {
		return TrackWithProfile{}, fmt.Errorf("decoding track: %w", err)
	}

	if len(tracks) == 0 {
		return TrackWithProfile{}, fmt.Errorf("track not found")
	}

	a.signURLs("audio", []*string{&tracks[0].AudioURL})

	return tracks[0], nil
}

// GetByUserID returns tracks for a specific user. When onlyPublic is true,
// drafts are excluded (use for anyone who is not the owner).
func (a *Access) GetByUserID(userID string, onlyPublic bool) ([]Track, error) {
	visibility := ""
	if onlyPublic {
		visibility = "&is_public=eq.true"
	}
	query := fmt.Sprintf(
		"select=*&user_id=eq.%s%s&order=created_at.desc",
		url.QueryEscape(userID), visibility,
	)

	data, err := a.client.Query("tracks", query)
	if err != nil {
		return nil, fmt.Errorf("querying tracks: %w", err)
	}

	var tracks []Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("decoding tracks: %w", err)
	}

	a.signTrackAudio(tracks)

	return tracks, nil
}

// Create inserts a new track.
func (a *Access) Create(nt NewTrack) (Track, error) {
	data, err := a.client.Insert("tracks", nt)
	if err != nil {
		return Track{}, fmt.Errorf("inserting track: %w", err)
	}

	var tracks []Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return Track{}, fmt.Errorf("decoding track: %w", err)
	}

	if len(tracks) == 0 {
		return Track{}, fmt.Errorf("no track returned after insert")
	}

	a.signTrackAudio(tracks[:1])

	return tracks[0], nil
}

// Delete removes a track. Returns error if the track doesn't belong to the user.
func (a *Access) Delete(trackID, userID string) error {
	filter := fmt.Sprintf("id=eq.%s&user_id=eq.%s",
		url.QueryEscape(trackID),
		url.QueryEscape(userID),
	)

	if err := a.client.Delete("tracks", filter); err != nil {
		return fmt.Errorf("deleting track: %w", err)
	}

	return nil
}

// SetVisibility updates the is_public flag of a track owned by userID.
// Returns an error if the track doesn't exist or isn't owned by the user.
func (a *Access) SetVisibility(trackID, userID string, isPublic bool) (Track, error) {
	filter := fmt.Sprintf("id=eq.%s&user_id=eq.%s",
		url.QueryEscape(trackID),
		url.QueryEscape(userID),
	)

	data, err := a.client.Update("tracks", filter, map[string]bool{"is_public": isPublic})
	if err != nil {
		return Track{}, fmt.Errorf("updating track visibility: %w", err)
	}

	var tracks []Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return Track{}, fmt.Errorf("decoding track: %w", err)
	}

	if len(tracks) == 0 {
		return Track{}, fmt.Errorf("track not found")
	}

	return tracks[0], nil
}

// GetStems returns stems for a specific track.
func (a *Access) GetStems(trackID string) ([]Stem, error) {
	query := fmt.Sprintf("select=*&track_id=eq.%s", url.QueryEscape(trackID))
	data, err := a.client.Query("stems", query)
	if err != nil {
		return nil, fmt.Errorf("querying stems: %w", err)
	}

	var stems []Stem
	if err := json.Unmarshal(data, &stems); err != nil {
		return nil, fmt.Errorf("decoding stems: %w", err)
	}

	urls := make([]*string, 0, len(stems))
	for i := range stems {
		urls = append(urls, &stems[i].AudioURL)
	}
	a.signURLs("stems", urls)

	return stems, nil
}

// AddStems inserts new stems.
func (a *Access) AddStems(ns []NewStem) error {
	if len(ns) == 0 {
		return nil
	}
	_, err := a.client.Insert("stems", ns)
	if err != nil {
		return fmt.Errorf("inserting stems: %w", err)
	}
	return nil
}

// GetInvites returns invited users for a track.
func (a *Access) GetInvites(trackID string) ([]TrackInvite, error) {
	query := fmt.Sprintf("select=*&track_id=eq.%s", url.QueryEscape(trackID))
	data, err := a.client.Query("track_invites", query)
	if err != nil {
		return nil, fmt.Errorf("querying track_invites: %w", err)
	}

	var invites []TrackInvite
	if err := json.Unmarshal(data, &invites); err != nil {
		return nil, fmt.Errorf("decoding track_invites: %w", err)
	}
	return invites, nil
}

// AddInvites inserts track invites.
func (a *Access) AddInvites(ni []NewTrackInvite) error {
	if len(ni) == 0 {
		return nil
	}
	_, err := a.client.Insert("track_invites", ni)
	if err != nil {
		return fmt.Errorf("inserting track invites: %w", err)
	}
	return nil
}

// GetContinuations returns live tracks that branched off a specific track.
func (a *Access) GetContinuations(parentTrackID string) ([]TrackWithProfile, error) {
	query := fmt.Sprintf(
		"select=*,profiles(username,display_name,avatar_url)&parent_track_id=eq.%s&is_public=eq.true&order=created_at.desc",
		url.QueryEscape(parentTrackID),
	)

	data, err := a.client.Query("tracks", query)
	if err != nil {
		return nil, fmt.Errorf("querying continuations: %w", err)
	}

	var tracks []TrackWithProfile
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("decoding continuations: %w", err)
	}

	urls := make([]*string, 0, len(tracks))
	for i := range tracks {
		urls = append(urls, &tracks[i].AudioURL)
	}
	a.signURLs("audio", urls)

	return tracks, nil
}
