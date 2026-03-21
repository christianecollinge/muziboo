// Package profilegrp handles HTTP requests for user profiles.
package profilegrp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/muziboo/api/business/profiles"
	"github.com/muziboo/api/business/tracks"
	"github.com/muziboo/api/foundation/web"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for profile handlers.
type Handlers struct {
	Profiles *profiles.Core
	Tracks   *tracks.Core
}

// ProfileWithTracks is the response for a public profile page.
type ProfileWithTracks struct {
	Profile profiles.Profile `json:"profile"`
	Tracks  []tracks.Track   `json:"tracks"`
}

// GetByUsername returns a public profile with their tracks.
func (h *Handlers) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	profile, err := h.Profiles.GetByUsername(username)
	if err != nil {
		web.RespondError(w, "user not found", http.StatusNotFound)
		return
	}

	userTracks, err := h.Tracks.GetByUserID(profile.ID)
	if err != nil {
		web.RespondError(w, "failed to fetch tracks", http.StatusInternalServerError)
		return
	}

	web.Respond(w, ProfileWithTracks{
		Profile: profile,
		Tracks:  userTracks,
	}, http.StatusOK)
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

	profile, err := h.Profiles.Update(userID, up)
	if err != nil {
		web.RespondError(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	web.Respond(w, profile, http.StatusOK)
}

// Me returns the authenticated user's own profile.
// If no profile exists yet, it auto-creates one using the user's auth metadata.
func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := h.Profiles.GetByID(userID)
	if err != nil {
		// Profile not found — auto-create using metadata from the JWT
		user, userOk := middleware.GetUser(r.Context())
		if !userOk || user == nil {
			web.RespondError(w, "profile not found and could not retrieve user info", http.StatusNotFound)
			return
		}

		// Build a username from auth metadata or fallback to email prefix
		username := user.UserMetadata.Username
		displayName := user.UserMetadata.DisplayName
		if username == "" {
			// Derive username from email (e.g. "john.doe@example.com" → "johndoe")
			email := user.Email
			atIdx := -1
			for i, c := range email {
				if c == '@' {
					atIdx = i
					break
				}
			}
			if atIdx > 0 {
				username = email[:atIdx]
			} else {
				username = email
			}
			// Sanitize: keep only alphanumeric, underscores, hyphens
			sanitized := make([]byte, 0, len(username))
			for _, c := range username {
				if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '-' {
					sanitized = append(sanitized, byte(c))
				}
			}
			username = string(sanitized)
			if username == "" {
				username = userID[:8]
			}
		}
		if displayName == "" {
			displayName = username
		}

		newProfile, createErr := h.Profiles.Create(profiles.NewProfile{
			ID:          userID,
			Username:    username,
			DisplayName: displayName,
		})
		if createErr != nil {
			web.RespondError(w, "profile not found: "+createErr.Error(), http.StatusInternalServerError)
			return
		}
		profile = newProfile
	}

	web.Respond(w, profile, http.StatusOK)
}
