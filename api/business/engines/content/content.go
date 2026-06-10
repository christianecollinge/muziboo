// Package content provides business rules for music content processing.
package content

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Engine encapsulates business rules for content.
type Engine struct{}

// NewEngine creates a new content Engine.
func NewEngine() *Engine {
	return &Engine{}
}

// AllowedAudioTypes are the permitted audio MIME types.
var allowedAudioTypes = map[string]bool{
	"audio/mpeg":   true,
	"audio/wav":    true,
	"audio/x-wav":  true,
	"audio/flac":   true,
	"audio/x-flac": true,
	"audio/mp4":    true,
	"audio/x-m4a":  true,
	"audio/ogg":    true,
	"audio/aac":    true,
	"audio/webm":   true,
}

// AllowedImageTypes are the permitted image MIME types.
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

const maxAudioSize = 50 << 20 // 50 MB
const maxImageSize = 10 << 20 // 10 MB

// ValidateAudio checks if the audio file satisfies the business rules.
func (e *Engine) ValidateAudio(contentType string, size int64) error {
	if !allowedAudioTypes[contentType] {
		return fmt.Errorf("invalid audio format")
	}
	if size > maxAudioSize {
		return fmt.Errorf("audio file too large")
	}
	return nil
}

// ValidateImage checks if the image file satisfies the business rules.
func (e *Engine) ValidateImage(contentType string, size int64) error {
	if !allowedImageTypes[contentType] {
		return fmt.Errorf("invalid image format")
	}
	if size > maxImageSize {
		return fmt.Errorf("image file too large")
	}
	return nil
}

// GenerateFilename creates a unique filename following the business policy.
func (e *Engine) GenerateFilename(userID, originalFilename, defaultExt string) string {
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = defaultExt
	}
	return fmt.Sprintf("%s/%s%s", userID, uuid.New().String(), ext)
}

// ResolveBucket determines the storage bucket based on type.
func (e *Engine) ResolveBucket(mediaType string) string {
	switch strings.ToLower(mediaType) {
	case "avatar":
		return "avatars"
	case "artwork":
		return "artwork"
	case "audio":
		return "audio"
	case "stem":
		return "stems"
	default:
		return "misc"
	}
}
