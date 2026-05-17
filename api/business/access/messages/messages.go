// Package messages provides resource access to direct message data.
package messages

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/muziboo/api/foundation/supabase"
)

// Message represents a direct message between two users.
type Message struct {
	ID          string `json:"id"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
	ReadAt      string `json:"read_at"`
	// Embedded sender profile (joined from profiles table)
	Sender MessageProfile `json:"sender"`
}

// MessageProfile is the embedded profile in a message query.
type MessageProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// NewMessage contains the data needed to send a message.
type NewMessage struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
}

// Access manages message database operations.
type Access struct {
	client *supabase.Client
}

// NewAccess creates a new messages Access.
func NewAccess(client *supabase.Client) *Access {
	return &Access{client: client}
}

// Send inserts a new message.
func (a *Access) Send(nm NewMessage) (Message, error) {
	data, err := a.client.Insert("messages", nm)
	if err != nil {
		return Message{}, fmt.Errorf("sending message: %w", err)
	}

	var msgs []Message
	if err := json.Unmarshal(data, &msgs); err != nil {
		return Message{}, fmt.Errorf("decoding message: %w", err)
	}

	if len(msgs) == 0 {
		return Message{}, fmt.Errorf("no message returned after insert")
	}

	return msgs[0], nil
}

// GetConversation retrieves all messages between two users, oldest first.
func (a *Access) GetConversation(userA, userB string) ([]Message, error) {
	// Get messages where (sender=A AND recipient=B) OR (sender=B AND recipient=A)
	// Supabase PostgREST supports OR via the `or` param
	filter := fmt.Sprintf(
		"select=*,sender:profiles!messages_sender_id_fkey(username,display_name,avatar_url)&or=(and(sender_id.eq.%s,recipient_id.eq.%s),and(sender_id.eq.%s,recipient_id.eq.%s))&order=created_at.asc",
		url.QueryEscape(userA), url.QueryEscape(userB),
		url.QueryEscape(userB), url.QueryEscape(userA),
	)

	data, err := a.client.Query("messages", filter)
	if err != nil {
		return nil, fmt.Errorf("querying conversation: %w", err)
	}

	var msgs []Message
	if err := json.Unmarshal(data, &msgs); err != nil {
		return nil, fmt.Errorf("decoding messages: %w", err)
	}

	return msgs, nil
}
