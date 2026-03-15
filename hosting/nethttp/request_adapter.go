// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package nethttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/microsoft/agents-sdk-go/hosting/core/authorization"
)

// RequestAdapter wraps an *http.Request for use with the SDK.
type RequestAdapter struct {
	r *http.Request
}

// NewRequestAdapter creates a new RequestAdapter.
func NewRequestAdapter(r *http.Request) *RequestAdapter {
	return &RequestAdapter{r: r}
}

// Method returns the HTTP method.
func (a *RequestAdapter) Method() string {
	return a.r.Method
}

// Header returns a header value.
func (a *RequestAdapter) Header(name string) string {
	return a.r.Header.Get(name)
}

// Body reads and parses the JSON body into dest.
func (a *RequestAdapter) Body(dest interface{}) error {
	defer a.r.Body.Close()
	data, err := io.ReadAll(a.r.Body)
	if err != nil {
		return fmt.Errorf("nethttp: failed to read request body: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("nethttp: failed to parse request body as JSON: %w", err)
	}
	return nil
}

// GetClaimsIdentity extracts and validates the Bearer token from Authorization header.
// Returns an anonymous identity if no token is present (for local testing).
func (a *RequestAdapter) GetClaimsIdentity(ctx context.Context) (*authorization.ClaimsIdentity, error) {
	authHeader := a.r.Header.Get("Authorization")
	if authHeader == "" {
		return authorization.NewClaimsIdentity(false, "Anonymous", nil), nil
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil, fmt.Errorf("nethttp: authorization header must use Bearer scheme")
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return nil, fmt.Errorf("nethttp: Bearer token is empty")
	}

	validator := &authorization.JWTValidator{}
	identity, err := validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("nethttp: token validation failed: %w", err)
	}
	identity.SecurityToken = token
	return identity, nil
}
