// Package engagement provides resource access to engagement data (votes, comments).
package engagement

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/muziboo/api/foundation/supabase"
)

// Vote represents a user's upvote on a track.
type Vote struct {
	ID        string `json:"id"`
	TrackID   string `json:"track_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// Comment represents a user's comment on a track.
type Comment struct {
	ID        string         `json:"id"`
	TrackID   string         `json:"track_id"`
	UserID    string         `json:"user_id"`
	Content   string         `json:"content"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Profile   CommentProfile `json:"profiles"`
}

// CommentProfile is the embedded profile in a comment query.
type CommentProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// NewVote is the data needed to create a vote.
type NewVote struct {
	TrackID string `json:"track_id"`
	UserID  string `json:"user_id"`
}

// NewComment is the data needed to create a comment.
type NewComment struct {
	TrackID string `json:"track_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

// Access manages engagement operations.
type Access struct {
	client *supabase.Client
}

// NewAccess creates a new engagement Access.
func NewAccess(client *supabase.Client) *Access {
	return &Access{client: client}
}

// AddVote inserts a new vote.
func (a *Access) AddVote(nv NewVote) (Vote, error) {
	data, err := a.client.Insert("votes", nv)
	if err != nil {
		return Vote{}, fmt.Errorf("inserting vote: %w", err)
	}

	var votes []Vote
	if err := json.Unmarshal(data, &votes); err != nil {
		return Vote{}, fmt.Errorf("decoding vote: %w", err)
	}

	if len(votes) == 0 {
		return Vote{}, fmt.Errorf("no vote returned after insert")
	}

	return votes[0], nil
}

// RemoveVote deletes a vote.
func (a *Access) RemoveVote(trackID, userID string) error {
	filter := fmt.Sprintf("track_id=eq.%s&user_id=eq.%s",
		url.QueryEscape(trackID),
		url.QueryEscape(userID),
	)

	if err := a.client.Delete("votes", filter); err != nil {
		return fmt.Errorf("deleting vote: %w", err)
	}

	return nil
}

// GetVotesByTrack retrieves all votes for a track.
func (a *Access) GetVotesByTrack(trackID string) ([]Vote, error) {
	query := fmt.Sprintf("select=*&track_id=eq.%s", url.QueryEscape(trackID))

	data, err := a.client.Query("votes", query)
	if err != nil {
		return nil, fmt.Errorf("querying votes: %w", err)
	}

	var votes []Vote
	if err := json.Unmarshal(data, &votes); err != nil {
		return nil, fmt.Errorf("decoding votes: %w", err)
	}

	return votes, nil
}

// AddComment inserts a new comment.
func (a *Access) AddComment(nc NewComment) (Comment, error) {
	data, err := a.client.Insert("comments", nc)
	if err != nil {
		return Comment{}, fmt.Errorf("inserting comment: %w", err)
	}

	var comments []Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return Comment{}, fmt.Errorf("decoding comment: %w", err)
	}

	if len(comments) == 0 {
		return Comment{}, fmt.Errorf("no comment returned after insert")
	}

	return comments[0], nil
}

// GetCommentsByTrack retrieves all comments for a track with profile info.
func (a *Access) GetCommentsByTrack(trackID string) ([]Comment, error) {
	query := fmt.Sprintf(
		"select=*,profiles(username,display_name,avatar_url)&track_id=eq.%s&order=created_at.asc",
		url.QueryEscape(trackID),
	)

	data, err := a.client.Query("comments", query)
	if err != nil {
		return nil, fmt.Errorf("querying comments: %w", err)
	}

	var comments []Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, fmt.Errorf("decoding comments: %w", err)
	}

	return comments, nil
}
