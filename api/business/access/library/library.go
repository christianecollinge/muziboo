// Package library provides resource access to track data.
package library

import (
	"encoding/json"
	"fmt"
	"net/url"

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
	return tracks, nil
}
