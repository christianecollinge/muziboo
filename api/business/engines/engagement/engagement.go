// Package engagement provides business rules for social interactions and profile policies.
package engagement

import (
	"fmt"
	"strings"
)

// Engine encapsulates business rules for social engagement.
type Engine struct{}

// NewEngine creates a new engagement Engine.
func NewEngine() *Engine {
	return &Engine{}
}

// DeriveUsername generates a valid username from an email or metadata.
func (e *Engine) DeriveUsername(email, metadataUsername string) string {
	username := metadataUsername
	if username == "" {
		// Derive from email
		atIdx := strings.Index(email, "@")
		if atIdx > 0 {
			username = email[:atIdx]
		} else {
			username = email
		}
	}

	// Sanitize: keep only alphanumeric, underscores, hyphens
	sanitized := make([]byte, 0, len(username))
	for _, c := range username {
		c := rune(c)
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-' {
			sanitized = append(sanitized, byte(strings.ToLower(string(c))[0]))
		}
	}

	result := string(sanitized)
	if result == "" {
		return "user"
	}
	return result
}

// ValidateComment checks if the comment content satisfies business rules.
func (e *Engine) ValidateComment(content string) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return fmt.Errorf("comment cannot be empty")
	}
	if len(trimmed) > 1000 {
		return fmt.Errorf("comment too long (max 1000 characters)")
	}
	return nil
}

// CanVote determines if a user is allowed to vote on a track.
// In the future, this could implement rate limiting or karma requirements.
func (e *Engine) CanVote(trackID, userID string) error {
	return nil
}
