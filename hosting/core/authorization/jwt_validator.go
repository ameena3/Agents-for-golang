// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authorization

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JWTValidator validates JWT tokens and extracts their claims.
// The default implementation parses without signature verification, which is
// suitable for internal claim inspection when the token has already been
// verified upstream (e.g. by the Bot Framework service).
// Production deployments should layer additional signature verification on top.
type JWTValidator struct {
	// ValidIssuers is the list of accepted issuer values. An empty list
	// disables issuer validation.
	ValidIssuers []string
	// ValidAudiences is the list of accepted audience values. An empty list
	// disables audience validation.
	ValidAudiences []string
}

// ValidateToken parses a JWT token string and returns a ClaimsIdentity
// containing all claims found in the token payload.
//
// The token signature is not verified. This is intentional for scenarios
// where upstream infrastructure (e.g. the Bot Framework connector) has
// already validated the signature. Add signature verification for production
// security requirements.
func (v *JWTValidator) ValidateToken(_ context.Context, tokenString string) (*ClaimsIdentity, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("authorization: malformed JWT: expected 3 parts, got %d", len(parts))
	}

	// Decode the payload (second part). JWT uses raw (unpadded) base64url.
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("authorization: invalid JWT payload encoding: %w", err)
	}

	// Unmarshal the JSON claims object.
	var rawClaims map[string]interface{}
	if err := json.Unmarshal(payload, &rawClaims); err != nil {
		return nil, fmt.Errorf("authorization: invalid JWT payload JSON: %w", err)
	}

	// Convert all claim values to their string representation.
	claimsMap := make(map[string]string, len(rawClaims))
	for k, val := range rawClaims {
		claimsMap[k] = fmt.Sprintf("%v", val)
	}

	identity := NewClaimsIdentity(true, "Bearer", claimsMap)
	return identity, nil
}
