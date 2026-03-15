// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package authorization provides authentication and authorization types for
// the Microsoft 365 Agents SDK, including JWT validation, claims handling,
// and OAuth connection management.
package authorization

import "fmt"

// ClaimsIdentity holds the authentication claims parsed from a JWT token.
type ClaimsIdentity struct {
	// IsAuthenticated indicates whether this identity represents an
	// authenticated principal. False for anonymous identities.
	IsAuthenticated bool
	// Claims holds all token claims as string key-value pairs.
	Claims map[string]string
	// AuthenticationType is the authentication scheme used (e.g. "Bearer", "Anonymous").
	AuthenticationType string
	// SecurityToken is the raw token string, if available.
	SecurityToken string
}

// NewClaimsIdentity creates a new ClaimsIdentity with the given values.
// If claims is nil an empty map is allocated.
func NewClaimsIdentity(isAuthenticated bool, authenticationType string, claims map[string]string) *ClaimsIdentity {
	if claims == nil {
		claims = make(map[string]string)
	}
	return &ClaimsIdentity{
		IsAuthenticated:    isAuthenticated,
		AuthenticationType: authenticationType,
		Claims:             claims,
	}
}

// GetClaimValue returns the value of the named claim, or an empty string if
// the claim is not present.
func (c *ClaimsIdentity) GetClaimValue(claimType string) string {
	return c.Claims[claimType]
}

// HasClaim returns true if a claim with the given type exists.
func (c *ClaimsIdentity) HasClaim(claimType string) bool {
	_, ok := c.Claims[claimType]
	return ok
}

// GetAppID returns the application ID from the claims. It checks the audience
// claim first (v2 tokens) and falls back to the appid claim (v1 tokens).
// Returns an empty string if neither claim is present.
func (c *ClaimsIdentity) GetAppID() string {
	if aud := c.Claims[AudienceClaim]; aud != "" {
		return aud
	}
	return c.Claims[AppIDClaim]
}

// GetOutgoingAppID returns the application ID to use for outgoing requests.
// For v1 tokens it reads the appid claim; for v2 tokens the azp claim.
func (c *ClaimsIdentity) GetOutgoingAppID() string {
	version := c.Claims[VersionClaim]
	if version == "" || version == "1.0" {
		return c.Claims[AppIDClaim]
	}
	if version == "2.0" {
		return c.Claims[AppIDV2Claim]
	}
	return ""
}

// IsAgentClaim returns true when the claims represent an agent-to-agent call
// rather than a call from ABS/SMBA infrastructure.
func (c *ClaimsIdentity) IsAgentClaim() bool {
	version := c.Claims[VersionClaim]
	if version == "" {
		return false
	}
	audience := c.Claims[AudienceClaim]
	if audience == "" {
		return false
	}
	appID := c.GetOutgoingAppID()
	if appID == "" {
		return false
	}
	return appID != audience
}

// GetTokenAudience returns the audience to use when requesting a token for
// the outgoing call described by these claims.
func (c *ClaimsIdentity) GetTokenAudience() string {
	if c.IsAgentClaim() {
		return fmt.Sprintf("app://%s", c.GetOutgoingAppID())
	}
	return AgentsSDKScope
}

// GetTokenScopes returns the OAuth scopes to request for the outgoing call.
func (c *ClaimsIdentity) GetTokenScopes() []string {
	if c.IsAgentClaim() {
		return []string{c.GetOutgoingAppID() + "/.default"}
	}
	return []string{AgentsSDKScope + "/.default"}
}
