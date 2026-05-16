// Package social provides orchestration for social interactions and user profiles.
package social

import (
	"fmt"

	engAccess "github.com/muziboo/api/business/access/engagement"
	"github.com/muziboo/api/business/access/library"
	"github.com/muziboo/api/business/access/profiles"
	engEngine "github.com/muziboo/api/business/engines/engagement"
)

// Manager orchestrates social and profile use cases.
type Manager struct {
	profileAccess    *profiles.Access
	libraryAccess    *library.Access
	engagementAccess *engAccess.Access
	engagementEngine *engEngine.Engine
}

// NewManager creates a new social Manager.
func NewManager(p *profiles.Access, l *library.Access, eAcc *engAccess.Access, eEng *engEngine.Engine) *Manager {
	return &Manager{
		profileAccess:    p,
		libraryAccess:    l,
		engagementAccess: eAcc,
		engagementEngine: eEng,
	}
}

// GetProfileWithTracks retrieves a profile and all tracks associated with it.
func (m *Manager) GetProfileWithTracks(username string) (profiles.Profile, []library.Track, error) {
	profile, err := m.profileAccess.GetByUsername(username)
	if err != nil {
		return profiles.Profile{}, nil, fmt.Errorf("profile access: %w", err)
	}

	tracks, err := m.libraryAccess.GetByUserID(profile.ID)
	if err != nil {
		return profile, []library.Track{}, nil // Non-fatal if tracks fail
	}

	return profile, tracks, nil
}

// GetOrCreateProfile retrieves a profile by ID or creates one from metadata.
func (m *Manager) GetOrCreateProfile(userID, email, metadataUsername, metadataDisplayName string) (profiles.Profile, error) {
	profile, err := m.profileAccess.GetByID(userID)
	if err == nil {
		return profile, nil
	}

	// Rule: Deriving username is an engine task
	username := m.engagementEngine.DeriveUsername(email, metadataUsername)
	displayName := metadataDisplayName
	if displayName == "" {
		displayName = username
	}

	newProfile, err := m.profileAccess.Create(profiles.NewProfile{
		ID:          userID,
		Username:    username,
		DisplayName: displayName,
	})
	if err != nil {
		return profiles.Profile{}, fmt.Errorf("creating profile: %w", err)
	}

	return newProfile, nil
}

// GetProfileByID retrieves a profile by user ID.
func (m *Manager) GetProfileByID(userID string) (profiles.Profile, error) {
	return m.profileAccess.GetByID(userID)
}

// UpdateProfile updates the profile data.
func (m *Manager) UpdateProfile(userID string, up profiles.UpdateProfile) (profiles.Profile, error) {
	return m.profileAccess.Update(userID, up)
}

// ToggleUpvote adds or removes an upvote for a track.
func (m *Manager) ToggleUpvote(trackID, userID string) (int, bool, error) {
	// 1. Check policy
	if err := m.engagementEngine.CanVote(trackID, userID); err != nil {
		return 0, false, err
	}

	// 2. Check if already voted
	existing, err := m.engagementAccess.GetVotesByTrack(trackID)
	if err != nil {
		return 0, false, fmt.Errorf("checking votes: %w", err)
	}

	hasVoted := false
	for _, v := range existing {
		if v.UserID == userID {
			hasVoted = true
			break
		}
	}

	// 3. Act
	if hasVoted {
		if err := m.engagementAccess.RemoveVote(trackID, userID); err != nil {
			return 0, false, err
		}
	} else {
		if _, err := m.engagementAccess.AddVote(engAccess.NewVote{TrackID: trackID, UserID: userID}); err != nil {
			return 0, false, err
		}
	}

	// 4. Return new count and state
	newVotes, _ := m.engagementAccess.GetVotesByTrack(trackID)
	return len(newVotes), !hasVoted, nil
}

// PostComment adds a new comment to a track.
func (m *Manager) PostComment(trackID, userID, content string) (engAccess.Comment, error) {
	// 1. Validate
	if err := m.engagementEngine.ValidateComment(content); err != nil {
		return engAccess.Comment{}, err
	}

	// 2. Persist
	return m.engagementAccess.AddComment(engAccess.NewComment{
		TrackID: trackID,
		UserID:  userID,
		Content: content,
	})
}

// GetTrackEngagement returns votes and comments for a track.
func (m *Manager) GetTrackEngagement(trackID string) ([]engAccess.Vote, []engAccess.Comment, error) {
	votes, err := m.engagementAccess.GetVotesByTrack(trackID)
	if err != nil {
		return nil, nil, err
	}
	comments, err := m.engagementAccess.GetCommentsByTrack(trackID)
	if err != nil {
		return votes, nil, nil
	}
	return votes, comments, nil
}
