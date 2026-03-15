// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package http

import (
	"context"
	"io"

	"github.com/microsoft/agents-sdk-go/hosting/core/authorization"
)

// HttpRequestProtocol abstracts an incoming HTTP request so that
// framework adapters (net/http, fasthttp, etc.) can share the same
// channel adapter base logic.
type HttpRequestProtocol interface {
	// Method returns the HTTP method (e.g. "POST").
	Method() string
	// Header returns the value of the named request header.
	Header(name string) string
	// Body returns the request body reader.
	Body() io.Reader
	// GetClaimsIdentity extracts and returns a ClaimsIdentity from
	// the Authorization header, or nil if no Bearer token is present.
	GetClaimsIdentity(ctx context.Context) *authorization.ClaimsIdentity
}
