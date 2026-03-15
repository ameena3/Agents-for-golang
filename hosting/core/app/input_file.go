// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// InputFile represents a file attachment uploaded by the user.
type InputFile struct {
	// Name is the original filename.
	Name string
	// ContentType is the MIME type of the file.
	ContentType string
	// ContentURL is the URL where the file content can be downloaded.
	ContentURL string
	// Content holds the file bytes if already downloaded.
	Content []byte
}

// InputFileDownloader downloads file attachments from activity attachments.
type InputFileDownloader struct {
	httpClient *http.Client
	// TokenProvider optionally supplies a Bearer token for authenticated downloads.
	TokenProvider func(ctx context.Context, resource string) (string, error)
}

// NewInputFileDownloader creates a new InputFileDownloader.
func NewInputFileDownloader() *InputFileDownloader {
	return &InputFileDownloader{
		httpClient: &http.Client{},
	}
}

// Download fetches the content of the given InputFile from its ContentURL.
// Returns the file with Content populated.
func (d *InputFileDownloader) Download(ctx context.Context, file *InputFile) (*InputFile, error) {
	if file.ContentURL == "" {
		return nil, fmt.Errorf("input_file: ContentURL is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.ContentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("input_file: creating request: %w", err)
	}

	if d.TokenProvider != nil {
		tok, err := d.TokenProvider(ctx, file.ContentURL)
		if err == nil && tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("input_file: downloading %s: %w", file.ContentURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("input_file: server returned %d for %s", resp.StatusCode, file.ContentURL)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("input_file: reading response body: %w", err)
	}

	result := *file
	result.Content = data
	if result.ContentType == "" && resp.Header.Get("Content-Type") != "" {
		result.ContentType = resp.Header.Get("Content-Type")
	}
	return &result, nil
}
