// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package http

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/microsoft/agents-sdk-go/hosting/core"
	"github.com/microsoft/agents-sdk-go/hosting/core/authorization"
)

// ChannelServiceRoutes registers standard Bot Framework routes on the
// given ServeMux. The default messages path is /api/messages.
func ChannelServiceRoutes(mux *http.ServeMux, base *HttpAdapterBase, agent core.Agent, messagesPath string) {
	if messagesPath == "" {
		messagesPath = "/api/messages"
	}
	mux.HandleFunc(messagesPath, func(w http.ResponseWriter, r *http.Request) {
		req := NewStdRequestAdapter(r)
		base.Process(r.Context(), req, agent, w)
	})
}

// StdRequestAdapter adapts a *http.Request to HttpRequestProtocol.
type StdRequestAdapter struct {
	r *http.Request
}

// NewStdRequestAdapter wraps a *http.Request as an HttpRequestProtocol.
func NewStdRequestAdapter(r *http.Request) *StdRequestAdapter {
	return &StdRequestAdapter{r: r}
}

// Method returns the HTTP method.
func (a *StdRequestAdapter) Method() string { return a.r.Method }

// Header returns the value of the named request header.
func (a *StdRequestAdapter) Header(name string) string { return a.r.Header.Get(name) }

// Body returns the request body as an io.Reader.
func (a *StdRequestAdapter) Body() io.Reader { return a.r.Body }

// GetClaimsIdentity extracts a ClaimsIdentity from the Authorization header.
// Returns nil if the header is absent or validation fails.
func (a *StdRequestAdapter) GetClaimsIdentity(ctx context.Context) *authorization.ClaimsIdentity {
	authHeader := a.r.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return nil
	}

	validator := &authorization.JWTValidator{}
	identity, err := validator.ValidateToken(ctx, token)
	if err != nil {
		return nil
	}
	identity.SecurityToken = token
	return identity
}
