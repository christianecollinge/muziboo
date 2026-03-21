// Package web provides HTTP helper utilities.
package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Respond sends a JSON response with the given status code.
func Respond(w http.ResponseWriter, data any, statusCode int) {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// RespondError sends a JSON error response.
func RespondError(w http.ResponseWriter, msg string, statusCode int) {
	Respond(w, map[string]string{"error": msg}, statusCode)
}

// Decode reads the request body into the given value.
func Decode(r *http.Request, v any) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 10<<20)) // 10MB limit
	if err != nil {
		return fmt.Errorf("reading body: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decoding json: %w", err)
	}

	return nil
}

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = errors.New("not found")

// ErrUnauthorized is returned when authentication fails.
var ErrUnauthorized = errors.New("unauthorized")

// ErrForbidden is returned when the user doesn't have permission.
var ErrForbidden = errors.New("forbidden")
