// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authorization

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// ---- ClaimsIdentity tests ----

func TestNewClaimsIdentity(t *testing.T) {
	claims := map[string]string{"aud": "my-app", "iss": "https://sts.windows.net/tenant/"}
	ci := NewClaimsIdentity(true, "Bearer", claims)
	if !ci.IsAuthenticated {
		t.Error("expected IsAuthenticated=true")
	}
	if ci.AuthenticationType != "Bearer" {
		t.Errorf("AuthenticationType: got %q, want Bearer", ci.AuthenticationType)
	}
	if ci.Claims["aud"] != "my-app" {
		t.Errorf("Claims[aud]: got %q, want my-app", ci.Claims["aud"])
	}
}

func TestNewClaimsIdentityNilClaims(t *testing.T) {
	ci := NewClaimsIdentity(false, "Anonymous", nil)
	if ci.Claims == nil {
		t.Error("expected non-nil claims map when nil was passed")
	}
}

func TestClaimsIdentityHasClaim(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{"sub": "user123"})
	if !ci.HasClaim("sub") {
		t.Error("expected HasClaim(sub) to be true")
	}
	if ci.HasClaim("nonexistent") {
		t.Error("expected HasClaim(nonexistent) to be false")
	}
}

func TestClaimsIdentityGetClaimValue(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{"tid": "tenant-abc"})
	if ci.GetClaimValue("tid") != "tenant-abc" {
		t.Errorf("GetClaimValue(tid): got %q, want tenant-abc", ci.GetClaimValue("tid"))
	}
	if ci.GetClaimValue("missing") != "" {
		t.Errorf("GetClaimValue(missing): expected empty string")
	}
}

func TestClaimsIdentityGetAppIDV1(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		AppIDClaim: "app-v1-id",
	})
	// No aud claim — should fall back to appid
	if ci.GetAppID() != "app-v1-id" {
		t.Errorf("GetAppID v1: got %q, want app-v1-id", ci.GetAppID())
	}
}

func TestClaimsIdentityGetAppIDV2AudTakesPrecedence(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		AudienceClaim: "aud-app-id",
		AppIDClaim:    "appid-fallback",
	})
	if ci.GetAppID() != "aud-app-id" {
		t.Errorf("GetAppID v2: got %q, want aud-app-id", ci.GetAppID())
	}
}

func TestClaimsIdentityGetOutgoingAppIDV1(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		AppIDClaim:   "outgoing-v1",
		VersionClaim: "1.0",
	})
	if ci.GetOutgoingAppID() != "outgoing-v1" {
		t.Errorf("GetOutgoingAppID v1: got %q, want outgoing-v1", ci.GetOutgoingAppID())
	}
}

func TestClaimsIdentityGetOutgoingAppIDV2(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		AppIDV2Claim: "outgoing-v2",
		VersionClaim: "2.0",
	})
	if ci.GetOutgoingAppID() != "outgoing-v2" {
		t.Errorf("GetOutgoingAppID v2: got %q, want outgoing-v2", ci.GetOutgoingAppID())
	}
}

func TestClaimsIdentityIsAgentClaim(t *testing.T) {
	// Agent claim: version set, aud != azp
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		VersionClaim:  "1.0",
		AudienceClaim: "audience-app",
		AppIDClaim:    "caller-app",
	})
	if !ci.IsAgentClaim() {
		t.Error("expected IsAgentClaim=true when appID != audience")
	}
}

func TestClaimsIdentityIsNotAgentClaimWhenSame(t *testing.T) {
	// Not an agent claim: appid == aud
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		VersionClaim:  "1.0",
		AudienceClaim: "same-app",
		AppIDClaim:    "same-app",
	})
	if ci.IsAgentClaim() {
		t.Error("expected IsAgentClaim=false when appID == audience")
	}
}

func TestClaimsIdentityGetTokenAudienceAgentClaim(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		VersionClaim:  "1.0",
		AudienceClaim: "audience-app",
		AppIDClaim:    "caller-app",
	})
	want := "app://caller-app"
	if ci.GetTokenAudience() != want {
		t.Errorf("GetTokenAudience: got %q, want %q", ci.GetTokenAudience(), want)
	}
}

func TestClaimsIdentityGetTokenAudienceNonAgent(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{})
	if ci.GetTokenAudience() != AgentsSDKScope {
		t.Errorf("GetTokenAudience: got %q, want %q", ci.GetTokenAudience(), AgentsSDKScope)
	}
}

func TestClaimsIdentityGetTokenScopesAgentClaim(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{
		VersionClaim:  "1.0",
		AudienceClaim: "audience-app",
		AppIDClaim:    "caller-app",
	})
	scopes := ci.GetTokenScopes()
	if len(scopes) != 1 || scopes[0] != "caller-app/.default" {
		t.Errorf("GetTokenScopes agent: got %v", scopes)
	}
}

func TestClaimsIdentityGetTokenScopesNonAgent(t *testing.T) {
	ci := NewClaimsIdentity(true, "Bearer", map[string]string{})
	scopes := ci.GetTokenScopes()
	want := AgentsSDKScope + "/.default"
	if len(scopes) != 1 || scopes[0] != want {
		t.Errorf("GetTokenScopes non-agent: got %v, want [%q]", scopes, want)
	}
}

// ---- Auth constants tests ----

func TestAuthConstantValues(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"AudienceClaim", AudienceClaim, "aud"},
		{"IssuerClaim", IssuerClaim, "iss"},
		{"AppIDClaim", AppIDClaim, "appid"},
		{"AppIDV2Claim", AppIDV2Claim, "azp"},
		{"VersionClaim", VersionClaim, "ver"},
		{"ServiceURLClaim", ServiceURLClaim, "serviceurl"},
		{"TenantIDClaim", TenantIDClaim, "tid"},
		{"KeyIDHeader", KeyIDHeader, "kid"},
		{"AgentsSDKScope", AgentsSDKScope, "https://api.botframework.com"},
		{"AnonymousSkillAppID", AnonymousSkillAppID, "anonymous"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.value != tc.want {
				t.Errorf("got %q, want %q", tc.value, tc.want)
			}
		})
	}
}

func TestTurnStateKeyConstants(t *testing.T) {
	if AgentIdentityKey == "" {
		t.Error("AgentIdentityKey should not be empty")
	}
	if OAuthScopeKey == "" {
		t.Error("OAuthScopeKey should not be empty")
	}
	if InvokeResponseKey == "" {
		t.Error("InvokeResponseKey should not be empty")
	}
}

// ---- JWTValidator tests ----

// makeTestJWT builds a minimal unsigned JWT with the given payload claims.
func makeTestJWT(claims map[string]interface{}) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload, _ := json.Marshal(claims)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	return fmt.Sprintf("%s.%s.", header, encodedPayload)
}

func TestJWTValidatorValidToken(t *testing.T) {
	v := &JWTValidator{}
	token := makeTestJWT(map[string]interface{}{
		"aud": "my-audience",
		"iss": "https://sts.windows.net/tenant/",
		"sub": "user-123",
	})

	identity, err := v.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if identity == nil {
		t.Fatal("expected non-nil identity")
	}
	if !identity.IsAuthenticated {
		t.Error("expected IsAuthenticated=true")
	}
	if identity.AuthenticationType != "Bearer" {
		t.Errorf("AuthenticationType: got %q, want Bearer", identity.AuthenticationType)
	}
	if identity.GetClaimValue("aud") != "my-audience" {
		t.Errorf("aud claim: got %q", identity.GetClaimValue("aud"))
	}
}

func TestJWTValidatorMalformedToken(t *testing.T) {
	v := &JWTValidator{}
	_, err := v.ValidateToken(context.Background(), "not.a.valid.jwt.token")
	if err == nil {
		t.Error("expected error for malformed JWT with 5 parts")
	}
}

func TestJWTValidatorTwoPartToken(t *testing.T) {
	v := &JWTValidator{}
	_, err := v.ValidateToken(context.Background(), "header.payload")
	if err == nil {
		t.Error("expected error for JWT with only 2 parts")
	}
}

func TestJWTValidatorInvalidBase64Payload(t *testing.T) {
	v := &JWTValidator{}
	// Use invalid base64 in the payload portion
	_, err := v.ValidateToken(context.Background(), "header.!!!invalid!!!.sig")
	if err == nil {
		t.Error("expected error for invalid base64 payload")
	}
}

func TestJWTValidatorInvalidJSONPayload(t *testing.T) {
	v := &JWTValidator{}
	// Valid base64 but not JSON
	payload := base64.RawURLEncoding.EncodeToString([]byte("not-json"))
	token := strings.Join([]string{"header", payload, "sig"}, ".")
	_, err := v.ValidateToken(context.Background(), token)
	if err == nil {
		t.Error("expected error for non-JSON payload")
	}
}

func TestJWTValidatorExtractsMultipleClaims(t *testing.T) {
	v := &JWTValidator{}
	token := makeTestJWT(map[string]interface{}{
		"ver":    "2.0",
		"aud":    "app-aud",
		"azp":    "caller-azp",
		"tid":    "my-tenant",
		"appid":  "app-id-v1",
	})

	identity, err := v.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if !identity.HasClaim("ver") {
		t.Error("expected ver claim")
	}
	if !identity.HasClaim("azp") {
		t.Error("expected azp claim")
	}
	if !identity.HasClaim("tid") {
		t.Error("expected tid claim")
	}
}
