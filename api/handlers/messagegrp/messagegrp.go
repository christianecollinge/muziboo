// Package messagegrp handles API requests for direct messages.
package messagegrp

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/muziboo/api/business/managers/social"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for message handlers.
type Handlers struct {
	SocialManager *social.Manager
}

type sendMessageRequest struct {
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
}

// Send handles POST /api/messages — send a direct message.
func (h *Handlers) Send(w http.ResponseWriter, r *http.Request) {
	senderID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RecipientID == "" || req.Content == "" {
		http.Error(w, "recipient_id and content are required", http.StatusBadRequest)
		return
	}

	msg, err := h.SocialManager.SendMessage(senderID, req.RecipientID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// GetConversation handles GET /api/messages/{userID} — fetch conversation with a user.
func (h *Handlers) GetConversation(w http.ResponseWriter, r *http.Request) {
	myID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	otherUserID := chi.URLParam(r, "userID")
	if otherUserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	msgs, err := h.SocialManager.GetConversation(myID, otherUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}
