// Package trackgrp handles HTTP requests for tracks.
package trackgrp

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/muziboo/api/business/tracks"
	"github.com/muziboo/api/foundation/web"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for track handlers.
type Handlers struct {
	Tracks *tracks.Core
}

// List returns all tracks (public).
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	list, err := h.Tracks.List(limit, offset)
	if err != nil {
		web.RespondError(w, "failed to fetch tracks", http.StatusInternalServerError)
		return
	}

	web.Respond(w, list, http.StatusOK)
}

// GetByID returns a single track (public).
func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	track, err := h.Tracks.GetByID(id)
	if err != nil {
		web.RespondError(w, "track not found", http.StatusNotFound)
		return
	}

	web.Respond(w, track, http.StatusOK)
}

// Create adds a new track (authenticated).
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var nt tracks.NewTrack
	if err := web.Decode(r, &nt); err != nil {
		web.RespondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	nt.UserID = userID

	if nt.Title == "" || nt.AudioURL == "" {
		web.RespondError(w, "title and audio_url are required", http.StatusBadRequest)
		return
	}

	track, err := h.Tracks.Create(nt)
	if err != nil {
		log.Printf("Track creation failed: %v", err)
		web.RespondError(w, "failed to create track", http.StatusInternalServerError)
		return
	}

	web.Respond(w, track, http.StatusCreated)
}

// Delete removes a track (authenticated, owner only).
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	if err := h.Tracks.Delete(id, userID); err != nil {
		web.RespondError(w, "failed to delete track", http.StatusInternalServerError)
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}
