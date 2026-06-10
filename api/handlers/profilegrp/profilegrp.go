// Package profilegrp handles HTTP requests for user profiles and social engagement.
package profilegrp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/muziboo/api/business/access/profiles"
	"github.com/muziboo/api/business/managers/social"
	"github.com/muziboo/api/foundation/web"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for profile handlers.
type Handlers struct {
	SocialManager *social.Manager
}

// Update updates the authenticated user's profile.
func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var up profiles.UpdateProfile
	if err := web.Decode(r, &up); err != nil {
		web.RespondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	profile, err := h.SocialManager.UpdateProfile(userID, up)
	if err != nil {
		web.RespondError(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	web.Respond(w, profile, http.StatusOK)
}

// GetByUsername returns a public profile with their tracks.
// Drafts are included only when the requester is the profile owner.
func (h *Handlers) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	viewerID, _ := middleware.GetUserID(r.Context())

	profile, tracks, err := h.SocialManager.GetProfileWithTracks(username, viewerID)
	if err != nil {
		web.RespondError(w, "user not found", http.StatusNotFound)
		return
	}

	web.Respond(w, struct {
		Profile profiles.Profile `json:"profile"`
		Tracks  any              `json:"tracks"`
	}{
		Profile: profile,
		Tracks:  tracks,
	}, http.StatusOK)
}

// Me returns the authenticated user's own profile.
func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, userOk := middleware.GetUser(r.Context())
	if !userOk || user == nil {
		web.RespondError(w, "could not retrieve user info", http.StatusNotFound)
		return
	}

	profile, err := h.SocialManager.GetOrCreateProfile(
		userID,
		user.Email,
		user.UserMetadata.Username,
		user.UserMetadata.DisplayName,
	)
	if err != nil {
		web.RespondError(w, "failed to resolve profile", http.StatusInternalServerError)
		return
	}

	web.Respond(w, profile, http.StatusOK)
}

// Vote toggles an upvote for a track.
func (h *Handlers) Vote(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	trackID := chi.URLParam(r, "id")
	count, active, err := h.SocialManager.ToggleUpvote(trackID, userID)
	if err != nil {
		web.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.Respond(w, map[string]any{
		"count":  count,
		"active": active,
	}, http.StatusOK)
}

// Comment adds a new comment to a track.
func (h *Handlers) Comment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	trackID := chi.URLParam(r, "id")

	var body struct {
		Content string `json:"content"`
	}
	if err := web.Decode(r, &body); err != nil {
		web.RespondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.SocialManager.PostComment(trackID, userID, body.Content)
	if err != nil {
		web.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.Respond(w, comment, http.StatusCreated)
}
