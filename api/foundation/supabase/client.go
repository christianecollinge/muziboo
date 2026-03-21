// Package supabase provides a client for interacting with Supabase services.
package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Client wraps HTTP calls to the Supabase REST API and Storage API.
type Client struct {
	URL            string
	AnonKey        string
	ServiceRoleKey string
	HTTPClient     *http.Client
}

// New creates a new Supabase client.
func New(url, anonKey, serviceRoleKey string) *Client {
	return &Client{
		URL:            url,
		AnonKey:        anonKey,
		ServiceRoleKey: serviceRoleKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// =========================================================================
// Database (PostgREST) operations
// =========================================================================

// Query performs a GET request against the PostgREST API.
// table: the table name
// query: query params like "select=*&username=eq.john"
func (c *Client) Query(table, query string) ([]byte, error) {
	url := fmt.Sprintf("%s/rest/v1/%s?%s", c.URL, table, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("apikey", c.ServiceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Insert performs a POST request to insert a row.
func (c *Client) Insert(table string, data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshaling data: %w", err)
	}

	url := fmt.Sprintf("%s/rest/v1/%s", c.URL, table)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("apikey", c.ServiceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Update performs a PATCH request to update rows matching a filter.
func (c *Client) Update(table, filter string, data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshaling data: %w", err)
	}

	url := fmt.Sprintf("%s/rest/v1/%s?%s", c.URL, table, filter)

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("apikey", c.ServiceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Delete performs a DELETE request.
func (c *Client) Delete(table, filter string) error {
	url := fmt.Sprintf("%s/rest/v1/%s?%s", c.URL, table, filter)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("apikey", c.ServiceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// =========================================================================
// Storage operations
// =========================================================================

// UploadFile uploads a file to a Supabase Storage bucket.
// Returns the public URL of the uploaded file.
func (c *Client) UploadFile(bucket, path string, file multipart.File, contentType string) (string, error) {
	body, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.URL, bucket, path)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("apikey", c.ServiceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("storage error (%d): %s", resp.StatusCode, string(respBody))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.URL, bucket, path)
	return publicURL, nil
}

// DeleteFile deletes a file from a Supabase Storage bucket.
func (c *Client) DeleteFile(bucket, path string) error {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.URL, bucket, path)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("apikey", c.ServiceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// =========================================================================
// Auth operations
// =========================================================================

// UserMetadata holds user metadata from Supabase Auth.
type UserMetadata struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

// User represents a Supabase authenticated user.
type User struct {
	ID           string       `json:"id"`
	Email        string       `json:"email"`
	UserMetadata UserMetadata `json:"user_metadata"`
}

// GetUser verifies a JWT token by fetching the user from Supabase Auth.
func (c *Client) GetUser(token string) (*User, error) {
	url := fmt.Sprintf("%s/auth/v1/user", c.URL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("apikey", c.AnonKey)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("auth error (%d)", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &user, nil
}
