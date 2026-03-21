// Package pagegrp handles HTTP requests for HTML pages served by Go templates.
package pagegrp

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/muziboo/api/business/profiles"
	"github.com/muziboo/api/business/tracks"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for page handlers.
type Handlers struct {
	Profiles        *profiles.Core
	Tracks          *tracks.Core
	Templates       map[string]*template.Template // page name -> parsed template
	SupabaseURL     string
	SupabaseAnonKey string
}

// TemplateData holds common data passed to every template.
type TemplateData struct {
	Title           string
	Description     string
	CurrentPage     string
	SupabaseURL     string
	SupabaseAnonKey string
	User            *UserInfo
	// Page-specific data
	Profile    *profiles.Profile
	Tracks     []tracks.Track
	TrackCards []TrackCardData
	HasMore    bool
	NextOffset int
}

// UserInfo is a simplified user struct for the nav bar.
type UserInfo struct {
	ID          string
	Username    string
	DisplayName string
}

// TrackCardData holds the data for a single track card template.
type TrackCardData struct {
	Title          string
	Genre          string
	AudioURL       string
	ArtworkURL     string
	ArtistName     string
	ArtistUsername string
	TimeAgo        string
}

// =========================================================================
// Page Handlers
// =========================================================================

// Login renders the login page.
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	h.render(w, "login", TemplateData{
		Title:       "Log In",
		Description: "Log in to your Muziboo account",
		CurrentPage: "/app/login",
	})
}

// Signup renders the signup page.
func (h *Handlers) Signup(w http.ResponseWriter, r *http.Request) {
	h.render(w, "signup", TemplateData{
		Title:       "Sign Up",
		Description: "Create a Muziboo account",
		CurrentPage: "/app/signup",
	})
}

// Dashboard renders the authenticated user's dashboard.
func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Redirect(w, r, "/app/login", http.StatusSeeOther)
		return
	}

	profile, err := h.Profiles.GetByID(userID)
	if err != nil {
		// Auto-create profile
		user, userOk := middleware.GetUser(r.Context())
		if !userOk || user == nil {
			http.Redirect(w, r, "/app/login", http.StatusSeeOther)
			return
		}
		username := user.UserMetadata.Username
		displayName := user.UserMetadata.DisplayName
		if username == "" {
			email := user.Email
			if idx := strings.Index(email, "@"); idx > 0 {
				username = email[:idx]
			} else {
				username = email
			}
		}
		if displayName == "" {
			displayName = username
		}
		profile, err = h.Profiles.Create(profiles.NewProfile{
			ID:          userID,
			Username:    username,
			DisplayName: displayName,
		})
		if err != nil {
			http.Error(w, "Failed to create profile: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	userTracks, err := h.Tracks.GetByUserID(userID)
	if err != nil {
		userTracks = []tracks.Track{}
	}

	cards := make([]TrackCardData, len(userTracks))
	for i, t := range userTracks {
		cards[i] = TrackCardData{
			Title:          t.Title,
			Genre:          t.Genre,
			AudioURL:       t.AudioURL,
			ArtworkURL:     t.ArtworkURL,
			ArtistName:     profile.DisplayName,
			ArtistUsername: profile.Username,
			TimeAgo:        timeAgo(t.CreatedAt),
		}
	}

	h.render(w, "dashboard", TemplateData{
		Title:       "Dashboard",
		Description: "Manage your tracks and profile",
		CurrentPage: "/app/dashboard",
		User:        &UserInfo{ID: profile.ID, Username: profile.Username, DisplayName: profile.DisplayName},
		Profile:     &profile,
		Tracks:      userTracks,
		TrackCards:  cards,
	})
}

// Explore renders the public track feed.
func (h *Handlers) Explore(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	allTracks, err := h.Tracks.List(limit+1, offset) // +1 to check if there are more
	if err != nil {
		allTracks = []tracks.TrackWithProfile{}
	}

	hasMore := len(allTracks) > limit
	if hasMore {
		allTracks = allTracks[:limit]
	}

	cards := make([]TrackCardData, len(allTracks))
	for i, t := range allTracks {
		cards[i] = TrackCardData{
			Title:          t.Title,
			Genre:          t.Genre,
			AudioURL:       t.AudioURL,
			ArtworkURL:     t.ArtworkURL,
			ArtistName:     t.Profile.DisplayName,
			ArtistUsername: t.Profile.Username,
			TimeAgo:        timeAgo(t.CreatedAt),
		}
	}

	// Check for authenticated user for nav
	var user *UserInfo
	if userID, ok := middleware.GetUserID(r.Context()); ok {
		if p, err := h.Profiles.GetByID(userID); err == nil {
			user = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
		}
	}

	h.render(w, "explore", TemplateData{
		Title:       "Explore",
		Description: "Discover real music from real people on Muziboo",
		CurrentPage: "/app/explore",
		User:        user,
		TrackCards:  cards,
		HasMore:     hasMore,
		NextOffset:  offset + limit,
	})
}

// Upload renders the upload page.
func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Redirect(w, r, "/app/login", http.StatusSeeOther)
		return
	}

	var user *UserInfo
	if p, err := h.Profiles.GetByID(userID); err == nil {
		user = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
	}

	h.render(w, "upload", TemplateData{
		Title:       "Upload",
		Description: "Upload your music to Muziboo",
		CurrentPage: "/app/upload",
		User:        user,
	})
}

// UserProfile renders a public user profile page.
func (h *Handlers) UserProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	profile, err := h.Profiles.GetByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userTracks, err := h.Tracks.GetByUserID(profile.ID)
	if err != nil {
		userTracks = []tracks.Track{}
	}

	cards := make([]TrackCardData, len(userTracks))
	for i, t := range userTracks {
		cards[i] = TrackCardData{
			Title:          t.Title,
			Genre:          t.Genre,
			AudioURL:       t.AudioURL,
			ArtworkURL:     t.ArtworkURL,
			ArtistName:     profile.DisplayName,
			ArtistUsername: profile.Username,
			TimeAgo:        timeAgo(t.CreatedAt),
		}
	}

	// Check for authenticated user for nav
	var user *UserInfo
	if userID, ok := middleware.GetUserID(r.Context()); ok {
		if p, err := h.Profiles.GetByID(userID); err == nil {
			user = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
		}
	}

	h.render(w, "profile", TemplateData{
		Title:       profile.DisplayName,
		Description: profile.DisplayName + " on Muziboo",
		CurrentPage: "/app/user/" + username,
		User:        user,
		Profile:     &profile,
		Tracks:      userTracks,
		TrackCards:  cards,
	})
}

// =========================================================================
// Helpers
// =========================================================================

func (h *Handlers) render(w http.ResponseWriter, page string, data TemplateData) {
	data.SupabaseURL = h.SupabaseURL
	data.SupabaseAnonKey = h.SupabaseAnonKey

	tmpl, ok := h.Templates[page]
	if !ok {
		log.Printf("Template not found: %s", page)
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error (%s): %v", page, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func timeAgo(dateStr string) string {
	t, err := time.Parse(time.RFC3339Nano, dateStr)
	if err != nil {
		// try other common formats
		t, err = time.Parse("2006-01-02T15:04:05", dateStr)
		if err != nil {
			return dateStr
		}
	}

	seconds := time.Since(t).Seconds()
	if seconds < 60 {
		return "just now"
	}
	mins := int(seconds / 60)
	if mins < 60 {
		return strconv.Itoa(mins) + "m ago"
	}
	hours := int(math.Floor(float64(mins) / 60))
	if hours < 24 {
		return strconv.Itoa(hours) + "h ago"
	}
	days := int(math.Floor(float64(hours) / 24))
	if days < 30 {
		return strconv.Itoa(days) + "d ago"
	}
	months := int(math.Floor(float64(days) / 30))
	return strconv.Itoa(months) + "mo ago"
}
