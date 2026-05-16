// Package library provides orchestration for track management.
package library

import (
	"fmt"
	"io"

	"github.com/muziboo/api/business/access/blob"
	"github.com/muziboo/api/business/access/library"
	"github.com/muziboo/api/business/engines/content"
)

// Manager orchestrates library use cases.
type Manager struct {
	libraryAccess *library.Access
	blobAccess    *blob.Access
	contentEngine *content.Engine
}

// NewManager creates a new library Manager.
func NewManager(l *library.Access, b *blob.Access, e *content.Engine) *Manager {
	return &Manager{
		libraryAccess: l,
		blobAccess:    b,
		contentEngine: e,
	}
}

// UploadMedia handles uploading any media file (audio or image).
func (m *Manager) UploadMedia(userID, mediaType, filename string, size int64, reader io.Reader, contentType string) (string, error) {
	// 1. Validate based on engine rules
	var err error
	if mediaType == "audio" {
		err = m.contentEngine.ValidateAudio(contentType, size)
	} else {
		err = m.contentEngine.ValidateImage(contentType, size)
	}
	if err != nil {
		return "", fmt.Errorf("content validation: %w", err)
	}

	// 2. Resolve bucket and filename
	bucket := m.contentEngine.ResolveBucket(mediaType)
	defaultExt := ".mp3"
	if mediaType != "audio" {
		defaultExt = ".jpg"
	}
	destFilename := m.contentEngine.GenerateFilename(userID, filename, defaultExt)

	// 3. Save to storage via accessor
	publicURL, err := m.blobAccess.Upload(bucket, destFilename, reader, contentType)
	if err != nil {
		return "", fmt.Errorf("storing blob: %w", err)
	}

	return publicURL, nil
}

// UploadStemRequest holds data for a single stem upload.
type UploadStemRequest struct {
	Name      string
	Audio     io.Reader
	AudioName string
	AudioType string
	AudioSize int64
}

// UploadTrackRequest holds data for a complete track upload.
type UploadTrackRequest struct {
	UserID        string
	Title         string
	Description   string
	Genre         string
	Audio         io.Reader
	AudioName     string
	AudioType     string
	AudioSize     int64
	Artwork       io.Reader
	ArtworkName   string
	ArtworkType   string
	ArtworkSize   int64
	IsColab       bool
	ParentTrackID string
	Stems         []UploadStemRequest
	InvitedUsers  []string // profile IDs of invited users
}

// UploadTrack orchestrates the full sequence of uploading files and saving metadata.
func (m *Manager) UploadTrack(req UploadTrackRequest) (library.Track, error) {
	// 1. Upload Audio
	audioURL, err := m.UploadMedia(req.UserID, "audio", req.AudioName, req.AudioSize, req.Audio, req.AudioType)
	if err != nil {
		return library.Track{}, fmt.Errorf("audio upload: %w", err)
	}

	// 2. Upload Artwork (optional)
	artworkURL := ""
	if req.Artwork != nil {
		url, err := m.UploadMedia(req.UserID, "artwork", req.ArtworkName, req.ArtworkSize, req.Artwork, req.ArtworkType)
		if err != nil {
			return library.Track{}, fmt.Errorf("artwork upload: %w", err)
		}
		artworkURL = url
	}

	// 3. Upload Stems (if any)
	var uploadedStems []library.NewStem
	for _, s := range req.Stems {
		stemURL, err := m.UploadMedia(req.UserID, "stem", s.AudioName, s.AudioSize, s.Audio, s.AudioType)
		if err != nil {
			return library.Track{}, fmt.Errorf("stem upload failed (%s): %w", s.Name, err)
		}
		uploadedStems = append(uploadedStems, library.NewStem{
			Name:     s.Name,
			AudioURL: stemURL,
		})
	}

	// 4. Save Metadata
	track, err := m.libraryAccess.Create(library.NewTrack{
		UserID:        req.UserID,
		Title:         req.Title,
		Description:   req.Description,
		Genre:         req.Genre,
		AudioURL:      audioURL,
		ArtworkURL:    artworkURL,
		IsColab:       req.IsColab,
		ParentTrackID: req.ParentTrackID,
	})
	if err != nil {
		return library.Track{}, fmt.Errorf("saving metadata: %w", err)
	}

	// 5. Save Stems Metadata
	if len(uploadedStems) > 0 {
		for i := range uploadedStems {
			uploadedStems[i].TrackID = track.ID
		}
		if err := m.libraryAccess.AddStems(uploadedStems); err != nil {
			return track, fmt.Errorf("saving stems: %w", err)
		}
	}

	// 6. Save Invites
	if req.IsColab && len(req.InvitedUsers) > 0 {
		var invites []library.NewTrackInvite
		for _, uid := range req.InvitedUsers {
			invites = append(invites, library.NewTrackInvite{
				TrackID:       track.ID,
				InvitedUserID: uid,
			})
		}
		if err := m.libraryAccess.AddInvites(invites); err != nil {
			return track, fmt.Errorf("saving invites: %w", err)
		}
	}

	return track, nil
}

// SaveTrack records the track metadata.
func (m *Manager) SaveTrack(nt library.NewTrack) (library.Track, error) {
	return m.libraryAccess.Create(nt)
}

// GetTrack retrieves a track by its ID.
func (m *Manager) GetTrack(id string) (library.TrackWithProfile, error) {
	return m.libraryAccess.GetByID(id)
}

// ListTracks returns a list of tracks.
func (m *Manager) ListTracks(limit, offset int) ([]library.TrackWithProfile, error) {
	return m.libraryAccess.List(limit, offset)
}

// DeleteTrack removes a track from the library.
func (m *Manager) DeleteTrack(trackID, userID string) error {
	return m.libraryAccess.Delete(trackID, userID)
}

// GetStems retrieves stems for a track.
func (m *Manager) GetStems(trackID string) ([]library.Stem, error) {
	return m.libraryAccess.GetStems(trackID)
}

// GetContinuations retrieves continuations for a track.
func (m *Manager) GetContinuations(trackID string) ([]library.TrackWithProfile, error) {
	return m.libraryAccess.GetContinuations(trackID)
}

// GetInvites retrieves invites for a track.
func (m *Manager) GetInvites(trackID string) ([]library.TrackInvite, error) {
	return m.libraryAccess.GetInvites(trackID)
}
