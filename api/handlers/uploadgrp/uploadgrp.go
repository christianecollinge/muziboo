// Package uploadgrp handles file upload HTTP requests.
package uploadgrp

import (
	"io"
	"net/http"
	"strings"

	"github.com/muziboo/api/business/managers/library"
	"github.com/muziboo/api/business/managers/social"
	"github.com/muziboo/api/foundation/web"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for upload handlers.
type Handlers struct {
	LibraryManager *library.Manager
	SocialManager  *social.Manager
}

// UploadTrack handles the complete track upload including files and metadata.
func (h *Handlers) UploadTrack(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 1. Parse Multipart Form
	if err := r.ParseMultipartForm(64 << 20); err != nil { // 64 MB max
		web.RespondError(w, "error parsing form", http.StatusBadRequest)
		return
	}

	// 2. Extract Data
	title := r.FormValue("title")
	description := r.FormValue("description")
	genre := r.FormValue("genre")

	// Audio file (required)
	audioFile, audioHeader, err := r.FormFile("audio")
	if err != nil {
		web.RespondError(w, "audio file is required", http.StatusBadRequest)
		return
	}
	defer audioFile.Close()

	// Artwork file (optional)
	var artworkFile io.Reader
	var artworkName, artworkType string
	var artworkSize int64

	art, artHeader, err := r.FormFile("artwork")
	if err == nil {
		defer art.Close()
		artworkFile = art
		artworkName = artHeader.Filename
		artworkType = artHeader.Header.Get("Content-Type")
		artworkSize = artHeader.Size
	}

	isColab := r.FormValue("is_colab") == "true"
	parentTrackID := r.FormValue("parent_track_id")
	isPublic := r.FormValue("is_public") == "true"

	invitedUsersRaw := r.FormValue("invited_users")
	var invitedUsers []string
	if invitedUsersRaw != "" {
		usernames := strings.Split(invitedUsersRaw, ",")
		for _, un := range usernames {
			un = strings.TrimSpace(un)
			if un == "" {
				continue
			}
			// Resolve username to profile ID
			if p, err := h.SocialManager.GetProfileByUsername(un); err == nil {
				invitedUsers = append(invitedUsers, p.ID)
			}
		}
	}

	var stems []library.UploadStemRequest
	for key, fileHeaders := range r.MultipartForm.File {
		if strings.HasPrefix(key, "stem_") && strings.HasSuffix(key, "_audio") {
			if len(fileHeaders) == 0 {
				continue
			}
			fh := fileHeaders[0]
			idxStr := strings.TrimSuffix(strings.TrimPrefix(key, "stem_"), "_audio")
			stemName := r.FormValue("stem_" + idxStr + "_name")
			if stemName == "" {
				stemName = fh.Filename
			}
			f, err := fh.Open()
			if err != nil {
				continue
			}
			stems = append(stems, library.UploadStemRequest{
				Name:      stemName,
				Audio:     f,
				AudioName: fh.Filename,
				AudioType: fh.Header.Get("Content-Type"),
				AudioSize: fh.Size,
			})
		}
	}

	defer func() {
		for _, s := range stems {
			if closer, ok := s.Audio.(io.ReadCloser); ok {
				closer.Close()
			}
		}
	}()

	// 3. Call Manager
	track, err := h.LibraryManager.UploadTrack(library.UploadTrackRequest{
		UserID:        userID,
		Title:         title,
		Description:   description,
		Genre:         genre,
		Audio:         audioFile,
		AudioName:     audioHeader.Filename,
		AudioType:     audioHeader.Header.Get("Content-Type"),
		AudioSize:     audioHeader.Size,
		Artwork:       artworkFile,
		ArtworkName:   artworkName,
		ArtworkType:   artworkType,
		ArtworkSize:   artworkSize,
		IsColab:       isColab,
		ParentTrackID: parentTrackID,
		IsPublic:      isPublic,
		Stems:         stems,
		InvitedUsers:  invitedUsers,
	})
	if err != nil {
		web.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.Respond(w, track, http.StatusCreated)
}

// UploadAudio handles simple audio file uploads (legacy/partial).
func (h *Handlers) UploadAudio(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		web.RespondError(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	publicURL, err := h.LibraryManager.UploadMedia(
		userID,
		"audio",
		header.Filename,
		header.Size,
		file,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		web.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.Respond(w, map[string]string{"url": publicURL}, http.StatusOK)
}

// UploadArtwork handles artwork image uploads (legacy/partial).
func (h *Handlers) UploadArtwork(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		web.RespondError(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine type based on query param
	mediaType := "artwork"
	if strings.ToLower(r.URL.Query().Get("type")) == "avatar" {
		mediaType = "avatar"
	}

	publicURL, err := h.LibraryManager.UploadMedia(
		userID,
		mediaType,
		header.Filename,
		header.Size,
		file,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		web.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.Respond(w, map[string]string{"url": publicURL}, http.StatusOK)
}
