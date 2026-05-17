// Package pagegrp handles HTTP requests for HTML pages served by Go templates.
package pagegrp

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	engagementAccess "github.com/muziboo/api/business/access/engagement"
	libraryAccess "github.com/muziboo/api/business/access/library"
	msgAccess "github.com/muziboo/api/business/access/messages"
	profileAccess "github.com/muziboo/api/business/access/profiles"
	libraryManager "github.com/muziboo/api/business/managers/library"
	socialManager "github.com/muziboo/api/business/managers/social"
	"github.com/muziboo/api/middleware"
)

// Handlers holds dependencies for page handlers.
type Handlers struct {
	LibraryManager  *libraryManager.Manager
	SocialManager   *socialManager.Manager
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
	Profile    *profileAccess.Profile
	Tracks     []libraryAccess.Track
	TrackCards []TrackCardData
	HasMore    bool
	NextOffset int
	// Track detail specific
	Comments       []engagementAccess.Comment
	Track          TrackCardData
	Stems          []libraryAccess.Stem
	Continuations  []TrackCardData
	IsInvitedColab bool
	// Messages
	Messages        []msgAccess.Message
	OtherUser       *profileAccess.Profile
}

// UserInfo is a simplified user struct for the nav bar.
type UserInfo struct {
	ID          string
	Username    string
	DisplayName string
}

// TrackCardData holds the data for a single track card template.
type TrackCardData struct {
	ID             string
	Title          string
	Description    string
	Genre          string
	AudioURL       string
	ArtworkURL     string
	ArtistName     string
	ArtistUsername string
	TimeAgo        string
	VoteCount      int
	CommentCount   int
	HasVoted       bool
	IsColab        bool
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

	user, userOk := middleware.GetUser(r.Context())
	if !userOk || user == nil {
		http.Redirect(w, r, "/app/login", http.StatusSeeOther)
		return
	}

	profile, err := h.SocialManager.GetOrCreateProfile(
		userID,
		user.Email,
		user.UserMetadata.Username,
		user.UserMetadata.DisplayName,
	)
	if err != nil {
		http.Error(w, "Failed to resolve profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, tracks, err := h.SocialManager.GetProfileWithTracks(profile.Username)
	if err != nil {
		tracks = []libraryAccess.Track{}
	}

	cards := make([]TrackCardData, len(tracks))
	for i, t := range tracks {
		votes, comments, _ := h.SocialManager.GetTrackEngagement(t.ID)
		hasVoted := false
		for _, v := range votes {
			if v.UserID == userID {
				hasVoted = true
				break
			}
		}

		cards[i] = TrackCardData{
			ID:             t.ID,
			Title:          t.Title,
			Genre:          t.Genre,
			AudioURL:       t.AudioURL,
			ArtworkURL:     t.ArtworkURL,
			ArtistName:     profile.DisplayName,
			ArtistUsername: profile.Username,
			TimeAgo:        timeAgo(t.CreatedAt),
			VoteCount:      len(votes),
			CommentCount:   len(comments),
			HasVoted:       hasVoted,
			IsColab:        t.IsColab,
		}
	}

	h.render(w, "dashboard", TemplateData{
		Title:       "Dashboard",
		Description: "Manage your tracks and profile",
		CurrentPage: "/app/dashboard",
		User:        &UserInfo{ID: profile.ID, Username: profile.Username, DisplayName: profile.DisplayName},
		Profile:     &profile,
		Tracks:      tracks,
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

	allTracks, err := h.LibraryManager.ListTracks(limit+1, offset)
	if err != nil {
		allTracks = []libraryAccess.TrackWithProfile{}
	}

	hasMore := len(allTracks) > limit
	if hasMore {
		allTracks = allTracks[:limit]
	}

	var currentUserID string
	if userID, ok := middleware.GetUserID(r.Context()); ok {
		currentUserID = userID
	}

	cards := make([]TrackCardData, len(allTracks))
	for i, t := range allTracks {
		votes, comments, _ := h.SocialManager.GetTrackEngagement(t.ID)
		hasVoted := false
		if currentUserID != "" {
			for _, v := range votes {
				if v.UserID == currentUserID {
					hasVoted = true
					break
				}
			}
		}

		cards[i] = TrackCardData{
			ID:             t.ID,
			Title:          t.Title,
			Genre:          t.Genre,
			AudioURL:       t.AudioURL,
			ArtworkURL:     t.ArtworkURL,
			ArtistName:     t.Profile.DisplayName,
			ArtistUsername: t.Profile.Username,
			TimeAgo:        timeAgo(t.CreatedAt),
			VoteCount:      len(votes),
			CommentCount:   len(comments),
			HasVoted:       hasVoted,
			IsColab:        t.IsColab,
		}
	}

	var userInfo *UserInfo
	if currentUserID != "" {
		if p, err := h.SocialManager.GetProfileByID(currentUserID); err == nil {
			userInfo = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
		}
	}

	h.render(w, "explore", TemplateData{
		Title:       "Explore",
		Description: "Discover real music from real people on Muziboo",
		CurrentPage: "/app/explore",
		User:        userInfo,
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

	var userInfo *UserInfo
	if p, err := h.SocialManager.GetProfileByID(userID); err == nil {
		userInfo = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
	}

	h.render(w, "upload", TemplateData{
		Title:       "Upload",
		Description: "Upload your music to Muziboo",
		CurrentPage: "/app/upload",
		User:        userInfo,
	})
}

// TrackDetail renders a single track with its comments.
func (h *Handlers) TrackDetail(w http.ResponseWriter, r *http.Request) {
	trackID := chi.URLParam(r, "id")

	track, err := h.LibraryManager.GetTrack(trackID)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	votes, comments, _ := h.SocialManager.GetTrackEngagement(trackID)

	var currentUserID string
	if userID, ok := middleware.GetUserID(r.Context()); ok {
		currentUserID = userID
	}

	hasVoted := false
	if currentUserID != "" {
		for _, v := range votes {
			if v.UserID == currentUserID {
				hasVoted = true
				break
			}
		}
	}

	card := TrackCardData{
		ID:             track.ID,
		Title:          track.Title,
		Description:    track.Description,
		Genre:          track.Genre,
		AudioURL:       track.AudioURL,
		ArtworkURL:     track.ArtworkURL,
		ArtistName:     track.Profile.DisplayName,
		ArtistUsername: track.Profile.Username,
		TimeAgo:        timeAgo(track.CreatedAt),
		VoteCount:      len(votes),
		CommentCount:   len(comments),
		HasVoted:       hasVoted,
		IsColab:        track.IsColab,
	}

	var userInfo *UserInfo
	isInvitedColab := false
	if currentUserID != "" {
		if p, err := h.SocialManager.GetProfileByID(currentUserID); err == nil {
			userInfo = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
		}

		if track.IsColab {
			if track.UserID == currentUserID {
				isInvitedColab = true
			} else {
				invites, _ := h.LibraryManager.GetInvites(track.ID)
				for _, inv := range invites {
					if inv.InvitedUserID == currentUserID {
						isInvitedColab = true
						break
					}
				}
			}
		}
	}

	var stems []libraryAccess.Stem
	if isInvitedColab {
		stems, _ = h.LibraryManager.GetStems(track.ID)
	}

	var continuations []TrackCardData
	if conts, err := h.LibraryManager.GetContinuations(track.ID); err == nil {
		for _, c := range conts {
			cVotes, cComments, _ := h.SocialManager.GetTrackEngagement(c.ID)
			cHasVoted := false
			if currentUserID != "" {
				for _, v := range cVotes {
					if v.UserID == currentUserID {
						cHasVoted = true
						break
					}
				}
			}
			continuations = append(continuations, TrackCardData{
				ID:             c.ID,
				Title:          c.Title,
				Genre:          c.Genre,
				AudioURL:       c.AudioURL,
				ArtworkURL:     c.ArtworkURL,
				ArtistName:     c.Profile.DisplayName,
				ArtistUsername: c.Profile.Username,
				TimeAgo:        timeAgo(c.CreatedAt),
				VoteCount:      len(cVotes),
				CommentCount:   len(cComments),
				HasVoted:       cHasVoted,
				IsColab:        c.IsColab,
			})
		}
	}

	h.render(w, "track_detail", TemplateData{
		Title:          track.Title,
		Description:    track.Description,
		CurrentPage:    "/app/tracks/" + trackID,
		User:           userInfo,
		Track:          card,
		Comments:       comments,
		Stems:          stems,
		Continuations:  continuations,
		IsInvitedColab: isInvitedColab,
	})
}

// UserProfile renders a public user profile page.
func (h *Handlers) UserProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	profile, tracks, err := h.SocialManager.GetProfileWithTracks(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var currentUserID string
	if userID, ok := middleware.GetUserID(r.Context()); ok {
		currentUserID = userID
	}

	cards := make([]TrackCardData, len(tracks))
	for i, t := range tracks {
		votes, comments, _ := h.SocialManager.GetTrackEngagement(t.ID)
		hasVoted := false
		if currentUserID != "" {
			for _, v := range votes {
				if v.UserID == currentUserID {
					hasVoted = true
					break
				}
			}
		}

		cards[i] = TrackCardData{
			ID:             t.ID,
			Title:          t.Title,
			Genre:          t.Genre,
			AudioURL:       t.AudioURL,
			ArtworkURL:     t.ArtworkURL,
			ArtistName:     profile.DisplayName,
			ArtistUsername: profile.Username,
			TimeAgo:        timeAgo(t.CreatedAt),
			VoteCount:      len(votes),
			CommentCount:   len(comments),
			HasVoted:       hasVoted,
			IsColab:        t.IsColab,
		}
	}

	var userInfo *UserInfo
	if currentUserID != "" {
		if p, err := h.SocialManager.GetProfileByID(currentUserID); err == nil {
			userInfo = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
		}
	}

	h.render(w, "profile", TemplateData{
		Title:       profile.DisplayName,
		Description: profile.DisplayName + " on Muziboo",
		CurrentPage: "/app/user/" + username,
		User:        userInfo,
		Profile:     &profile,
		Tracks:      tracks,
		TrackCards:  cards,
	})
}

// Messages renders the direct message conversation page with another user.
func (h *Handlers) Messages(w http.ResponseWriter, r *http.Request) {
	myID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Redirect(w, r, "/app/login", http.StatusSeeOther)
		return
	}

	username := chi.URLParam(r, "username")
	otherProfile, err := h.SocialManager.GetProfileByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	msgs, err := h.SocialManager.GetConversation(myID, otherProfile.ID)
	if err != nil {
		msgs = nil // non-fatal, show empty thread
	}

	var myProfile *UserInfo
	if p, err := h.SocialManager.GetProfileByID(myID); err == nil {
		myProfile = &UserInfo{ID: p.ID, Username: p.Username, DisplayName: p.DisplayName}
	}

	h.render(w, "messages", TemplateData{
		Title:       "Message " + otherProfile.DisplayName,
		Description: "Direct messages with " + otherProfile.DisplayName,
		CurrentPage: "/app/messages/" + username,
		User:        myProfile,
		OtherUser:   &otherProfile,
		Messages:    msgs,
	})
}

// =========================================================================
// Helpers
// =========================================================================

func (h *Handlers) render(w http.ResponseWriter, page string, data TemplateData) {
	tmpl, ok := h.Templates[page]
	if !ok {
		log.Printf("Template not found: %s", page)
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	// Inject Supabase credentials if not already present
	if data.SupabaseURL == "" {
		data.SupabaseURL = h.SupabaseURL
	}
	if data.SupabaseAnonKey == "" {
		data.SupabaseAnonKey = h.SupabaseAnonKey
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
