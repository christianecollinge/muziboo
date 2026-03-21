// Package profiles provides business logic for user profiles.
package profiles

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/muziboo/api/foundation/supabase"
)

// Profile represents a user profile.
type Profile struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// NewProfile contains fields for creating a new profile.
type NewProfile struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

// UpdateProfile contains fields that can be updated.
type UpdateProfile struct {
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Username    *string `json:"username,omitempty"`
}

// Core manages profile operations.
type Core struct {
	client *supabase.Client
}

// NewCore creates a new profiles Core.
func NewCore(client *supabase.Client) *Core {
	return &Core{client: client}
}

// Create inserts a new profile row.
func (c *Core) Create(np NewProfile) (Profile, error) {
	data, err := c.client.Insert("profiles", np)
	if err != nil {
		return Profile{}, fmt.Errorf("creating profile: %w", err)
	}

	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return Profile{}, fmt.Errorf("decoding profile: %w", err)
	}

	if len(profiles) == 0 {
		return Profile{}, fmt.Errorf("profile not returned after creation")
	}

	return profiles[0], nil
}

// GetByUsername retrieves a profile by username.
func (c *Core) GetByUsername(username string) (Profile, error) {
	query := fmt.Sprintf("select=*&username=eq.%s", url.QueryEscape(username))

	data, err := c.client.Query("profiles", query)
	if err != nil {
		return Profile{}, fmt.Errorf("querying profile: %w", err)
	}

	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return Profile{}, fmt.Errorf("decoding profile: %w", err)
	}

	if len(profiles) == 0 {
		return Profile{}, fmt.Errorf("profile not found")
	}

	return profiles[0], nil
}

// GetByID retrieves a profile by user ID.
func (c *Core) GetByID(id string) (Profile, error) {
	query := fmt.Sprintf("select=*&id=eq.%s", url.QueryEscape(id))

	data, err := c.client.Query("profiles", query)
	if err != nil {
		return Profile{}, fmt.Errorf("querying profile: %w", err)
	}

	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return Profile{}, fmt.Errorf("decoding profile: %w", err)
	}

	if len(profiles) == 0 {
		return Profile{}, fmt.Errorf("profile not found")
	}

	return profiles[0], nil
}

// Update updates a user's profile.
func (c *Core) Update(userID string, up UpdateProfile) (Profile, error) {
	filter := fmt.Sprintf("id=eq.%s", url.QueryEscape(userID))

	data, err := c.client.Update("profiles", filter, up)
	if err != nil {
		return Profile{}, fmt.Errorf("updating profile: %w", err)
	}

	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return Profile{}, fmt.Errorf("decoding profile: %w", err)
	}

	if len(profiles) == 0 {
		return Profile{}, fmt.Errorf("profile not found after update")
	}

	return profiles[0], nil
}
