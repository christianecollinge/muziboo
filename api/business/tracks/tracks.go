// Package tracks provides business logic for music tracks.
package tracks

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/muziboo/api/foundation/supabase"
)

// Track represents a music track.
type Track struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	AudioURL    string `json:"audio_url"`
	ArtworkURL  string `json:"artwork_url"`
	Duration    int    `json:"duration"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Genre       string `json:"genre,omitempty"`
	AudioURL    string `json:"audio_url"`
	ArtworkURL  string `json:"artwork_url,omitempty"`
	Duration    int    `json:"duration,omitempty"`
}

// Core manages track operations.
type Core struct {
	client *supabase.Client
}

// NewCore creates a new tracks Core.
func NewCore(client *supabase.Client) *Core {
	return &Core{client: client}
}

// List returns all tracks ordered by newest first, with profile info.
func (c *Core) List(limit, offset int) ([]TrackWithProfile, error) {
	query := fmt.Sprintf(
		"select=*,profiles(username,display_name,avatar_url)&order=created_at.desc&limit=%d&offset=%d",
		limit, offset,
	)

	data, err := c.client.Query("tracks", query)
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
func (c *Core) GetByID(id string) (TrackWithProfile, error) {
	query := fmt.Sprintf(
		"select=*,profiles(username,display_name,avatar_url)&id=eq.%s",
		url.QueryEscape(id),
	)

	data, err := c.client.Query("tracks", query)
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

// GetByUserID returns all tracks for a specific user.
func (c *Core) GetByUserID(userID string) ([]Track, error) {
	query := fmt.Sprintf(
		"select=*&user_id=eq.%s&order=created_at.desc",
		url.QueryEscape(userID),
	)

	data, err := c.client.Query("tracks", query)
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
func (c *Core) Create(nt NewTrack) (Track, error) {
	data, err := c.client.Insert("tracks", nt)
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
func (c *Core) Delete(trackID, userID string) error {
	filter := fmt.Sprintf("id=eq.%s&user_id=eq.%s",
		url.QueryEscape(trackID),
		url.QueryEscape(userID),
	)

	if err := c.client.Delete("tracks", filter); err != nil {
		return fmt.Errorf("deleting track: %w", err)
	}

	return nil
}
