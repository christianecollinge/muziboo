// Package middleware provides HTTP middleware for the API.
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/muziboo/api/foundation/supabase"
	"github.com/muziboo/api/foundation/web"
)

// contextKey is used for context values.
type contextKey string

const userIDKey contextKey = "userID"
const userKey contextKey = "user"

// GetUserID extracts the authenticated user ID from the request context.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

// GetUser extracts the full authenticated user from the request context.
func GetUser(ctx context.Context) (*supabase.User, bool) {
	user, ok := ctx.Value(userKey).(*supabase.User)
	return user, ok
}

// Auth validates the Supabase token and injects the user ID into the context.
func Auth(supa *supabase.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				web.RespondError(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				web.RespondError(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Verify the token by fetching the user profile from Supabase
			user, err := supa.GetUser(tokenString)
			if err != nil {
				log.Printf("Auth error: failed to verify token with Supabase: %v", err)
				web.RespondError(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			if user == nil || user.ID == "" {
				log.Printf("Auth error: missing user id from Supabase response")
				web.RespondError(w, "missing user id in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, user.ID)
			ctx = context.WithValue(ctx, userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth tries to authenticate the user but does not block if there is no token.
// Use this for pages that need to show user info in the nav but also work for unauthenticated users.
func OptionalAuth(supa *supabase.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				next.ServeHTTP(w, r)
				return
			}

			user, err := supa.GetUser(parts[1])
			if err != nil || user == nil || user.ID == "" {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, user.ID)
			ctx = context.WithValue(ctx, userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
