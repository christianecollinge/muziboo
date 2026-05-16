// Package blob provides resource access to binary storage.
package blob

import (
	"fmt"
	"io"

	"github.com/muziboo/api/foundation/supabase"
)

// Access manages blob storage operations.
type Access struct {
	client *supabase.Client
}

// NewAccess creates a new blob Access.
func NewAccess(client *supabase.Client) *Access {
	return &Access{client: client}
}

// Upload uploads a file to a bucket and returns the public URL.
func (a *Access) Upload(bucket, filename string, reader io.Reader, contentType string) (string, error) {
	publicURL, err := a.client.UploadFile(bucket, filename, reader, contentType)
	if err != nil {
		return "", fmt.Errorf("uploading file to %s: %w", bucket, err)
	}

	return publicURL, nil
}
