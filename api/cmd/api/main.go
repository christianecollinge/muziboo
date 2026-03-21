// Muziboo API — Entry point
package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/muziboo/api/business/profiles"
	"github.com/muziboo/api/business/tracks"
	"github.com/muziboo/api/foundation/supabase"
	"github.com/muziboo/api/handlers/pagegrp"
	"github.com/muziboo/api/handlers/profilegrp"
	"github.com/muziboo/api/handlers/trackgrp"
	"github.com/muziboo/api/handlers/uploadgrp"
	"github.com/muziboo/api/middleware"
)

func main() {
	// Load .env file if it exists
	loadEnv(".env")

	// Read configuration from environment
	port := envOr("PORT", "8080")
	supabaseURL := mustEnv("SUPABASE_URL")
	supabaseAnon := mustEnv("SUPABASE_ANON_KEY")
	supabaseService := mustEnv("SUPABASE_SERVICE_ROLE_KEY")
	allowedOrigins := envOr("ALLOWED_ORIGINS", "http://localhost:4321")

	// Initialize Supabase client
	supa := supabase.New(supabaseURL, supabaseAnon, supabaseService)

	// Initialize business cores
	profileCore := profiles.NewCore(supa)
	trackCore := tracks.NewCore(supa)

	// Load HTML templates
	tmpl := loadTemplates()

	// Initialize handlers
	profileHandlers := &profilegrp.Handlers{
		Profiles: profileCore,
		Tracks:   trackCore,
	}

	trackHandlers := &trackgrp.Handlers{
		Tracks: trackCore,
	}

	uploadHandlers := &uploadgrp.Handlers{
		Supabase: supa,
	}

	pageHandlers := &pagegrp.Handlers{
		Profiles:        profileCore,
		Tracks:          trackCore,
		Templates:       tmpl,
		SupabaseURL:     supabaseURL,
		SupabaseAnonKey: supabaseAnon,
	}

	// Build router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(allowedOrigins))

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// =====================================================================
	// JSON API routes (unchanged)
	// =====================================================================

	// Public API routes
	r.Get("/api/tracks", trackHandlers.List)
	r.Get("/api/tracks/{id}", trackHandlers.GetByID)
	r.Get("/api/profiles/{username}", profileHandlers.GetByUsername)

	// Authenticated API routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(supa))

		r.Post("/api/tracks", trackHandlers.Create)
		r.Delete("/api/tracks/{id}", trackHandlers.Delete)

		r.Get("/api/me", profileHandlers.Me)
		r.Put("/api/profiles", profileHandlers.Update)

		r.Post("/api/upload/audio", uploadHandlers.UploadAudio)
		r.Post("/api/upload/artwork", uploadHandlers.UploadArtwork)
	})

	// =====================================================================
	// HTML App routes (Go templates + HTMX)
	// =====================================================================

	// Public pages (no auth required)
	r.Get("/app/login", pageHandlers.Login)
	r.Get("/app/signup", pageHandlers.Signup)
	r.Get("/app/explore", pageHandlers.Explore)
	r.Get("/app/user/{username}", pageHandlers.UserProfile)

	// Authenticated pages
	r.Group(func(r chi.Router) {
		r.Use(middleware.OptionalAuth(supa))

		r.Get("/app/dashboard", pageHandlers.Dashboard)
		r.Get("/app/upload", pageHandlers.Upload)
	})

	// Redirect root to explore
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/explore", http.StatusSeeOther)
	})

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting Muziboo API on %s", addr)
	log.Printf("Supabase: %s", supabaseURL)
	log.Printf("Allowed Origins: %s", allowedOrigins)
	log.Printf("App pages: http://localhost:%s/app/", port)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// =========================================================================
// Template loading
// =========================================================================

func loadTemplates() map[string]*template.Template {
	// Get the directory of this Go source file
	_, filename, _, _ := runtime.Caller(0)
	apiRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	templateDir := filepath.Join(apiRoot, "foundation", "web", "templates")

	base := filepath.Join(templateDir, "base.html")
	partials := filepath.Join(templateDir, "track_card.html")

	pages := []string{"login", "signup", "dashboard", "explore", "upload", "profile"}
	templates := make(map[string]*template.Template, len(pages))

	for _, page := range pages {
		pageFile := filepath.Join(templateDir, page+".html")
		t, err := template.ParseFiles(base, pageFile, partials)
		if err != nil {
			log.Fatalf("FATAL: failed to parse template %s: %v", page, err)
		}
		templates[page] = t
	}

	log.Printf("Loaded %d page templates from %s", len(templates), templateDir)
	return templates
}

// =========================================================================
// Helpers
// =========================================================================

// loadEnv reads a .env file and sets environment variables.
func loadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return // .env is optional
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

// mustEnv returns the environment variable or logs fatal.
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("FATAL: environment variable %s is required", key)
	}
	return val
}

// envOr returns the environment variable or a default value.
func envOr(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
