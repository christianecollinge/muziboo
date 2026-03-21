// Package uploadgrp handles file upload HTTP requests.
package uploadgrp

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/muziboo/api/foundation/supabase"
	"github.com/muziboo/api/foundation/web"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for upload handlers.
type Handlers struct {
	Supabase *supabase.Client
}

// allowedAudioTypes are the permitted audio MIME types.
var allowedAudioTypes = map[string]bool{
	"audio/mpeg":      true, // mp3
	"audio/wav":       true,
	"audio/x-wav":     true,
	"audio/flac":      true,
	"audio/x-flac":    true,
	"audio/mp4":       true, // m4a
	"audio/x-m4a":     true,
	"audio/ogg":       true,
	"audio/aac":       true,
	"audio/webm":      true,
}

// allowedImageTypes are the permitted image MIME types.
var allowedImageTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/webp":    true,
	"image/gif":     true,
}

const maxAudioSize = 50 << 20  // 50 MB
const maxImageSize = 10 << 20  // 10 MB

// UploadAudio handles audio file uploads.
func (h *Handlers) UploadAudio(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAudioSize)

	if err := r.ParseMultipartForm(maxAudioSize); err != nil {
		web.RespondError(w, "file too large (max 50MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		web.RespondError(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !allowedAudioTypes[contentType] {
		web.RespondError(w, "invalid audio format. Allowed: mp3, wav, flac, m4a, ogg, aac", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".mp3"
	}
	filename := fmt.Sprintf("%s/%s%s", userID, uuid.New().String(), ext)

	publicURL, err := h.Supabase.UploadFile("audio", filename, file, contentType)
	if err != nil {
		web.RespondError(w, "failed to upload audio", http.StatusInternalServerError)
		return
	}

	web.Respond(w, map[string]string{"url": publicURL}, http.StatusOK)
}

// UploadArtwork handles artwork image uploads.
func (h *Handlers) UploadArtwork(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		web.RespondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize)

	if err := r.ParseMultipartForm(maxImageSize); err != nil {
		web.RespondError(w, "file too large (max 10MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		web.RespondError(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		web.RespondError(w, "invalid image format. Allowed: jpg, png, webp, gif", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	// Determine bucket based on query param
	bucket := "artwork"
	if strings.ToLower(r.URL.Query().Get("type")) == "avatar" {
		bucket = "avatars"
	}

	filename := fmt.Sprintf("%s/%s%s", userID, uuid.New().String(), ext)

	publicURL, err := h.Supabase.UploadFile(bucket, filename, file, contentType)
	if err != nil {
		web.RespondError(w, "failed to upload image", http.StatusInternalServerError)
		return
	}

	web.Respond(w, map[string]string{"url": publicURL}, http.StatusOK)
}
