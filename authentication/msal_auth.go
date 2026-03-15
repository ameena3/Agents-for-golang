// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authentication

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// MsalAuth implements AccessTokenProvider using MSAL for Go.
// It caches tokens and refreshes them automatically before expiry.
type MsalAuth struct {
	config     Config
	client     *confidential.Client
	miCred     *azidentity.ManagedIdentityCredential
	mu         sync.Mutex
	tokenCache map[string]cachedToken
}

type cachedToken struct {
	token     string
	expiresAt time.Time
}

// NewMsalAuth creates a new MsalAuth from the given Config.
// Returns an error if the configuration is invalid.
func NewMsalAuth(cfg Config) (*MsalAuth, error) {
	if cfg.TenantID == "" || cfg.ClientID == "" {
		return nil, fmt.Errorf("TenantID and ClientID are required")
	}

	authority := cfg.Authority
	if authority == "" {
		authority = fmt.Sprintf("https://login.microsoftonline.com/%s", cfg.TenantID)
	}

	if cfg.UseManagedIdentity {
		var miOpts *azidentity.ManagedIdentityCredentialOptions
		if cfg.UserAssignedClientID != "" {
			miOpts = &azidentity.ManagedIdentityCredentialOptions{
				ID: azidentity.ClientID(cfg.UserAssignedClientID),
			}
		}
		miCred, err := azidentity.NewManagedIdentityCredential(miOpts)
		if err != nil {
			return nil, fmt.Errorf("creating managed identity credential: %w", err)
		}
		return &MsalAuth{
			config:     cfg,
			miCred:     miCred,
			tokenCache: make(map[string]cachedToken),
		}, nil
	}

	var cred confidential.Credential
	var err error

	if cfg.ClientSecret != "" {
		cred, err = confidential.NewCredFromSecret(cfg.ClientSecret)
		if err != nil {
			return nil, fmt.Errorf("creating secret credential: %w", err)
		}
	} else if cfg.CertificatePath != "" {
		certData, readErr := os.ReadFile(cfg.CertificatePath)
		if readErr != nil {
			return nil, fmt.Errorf("reading certificate file: %w", readErr)
		}
		certs, privateKey, parseErr := confidential.CertFromPEM(certData, cfg.CertificatePassword)
		if parseErr != nil {
			return nil, fmt.Errorf("parsing certificate: %w", parseErr)
		}
		cred, err = confidential.NewCredFromCert(certs, privateKey)
		if err != nil {
			return nil, fmt.Errorf("creating certificate credential: %w", err)
		}
	} else {
		return nil, fmt.Errorf("one of ClientSecret, CertificatePath, or UseManagedIdentity must be set")
	}

	client, err := confidential.New(authority, cfg.ClientID, cred)
	if err != nil {
		return nil, fmt.Errorf("creating MSAL client: %w", err)
	}

	return &MsalAuth{
		config:     cfg,
		client:     &client,
		tokenCache: make(map[string]cachedToken),
	}, nil
}

// GetAccessToken returns an access token for the given resource URL and scopes.
// Implements AccessTokenProvider.
func (m *MsalAuth) GetAccessToken(ctx context.Context, resourceURL string, scopes []string, forceRefresh bool) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cacheKey := resourceURL
	if len(scopes) > 0 {
		cacheKey = scopes[0]
	}

	// Check cache unless forceRefresh is requested
	if !forceRefresh {
		if cached, ok := m.tokenCache[cacheKey]; ok {
			if time.Now().Before(cached.expiresAt.Add(-5 * time.Minute)) {
				return cached.token, nil
			}
		}
	}

	resolvedScopes := scopes
	if len(resolvedScopes) == 0 {
		resolvedScopes = m.config.Scopes
	}
	if len(resolvedScopes) == 0 {
		resolvedScopes = []string{resourceURL + "/.default"}
	}

	if m.miCred != nil {
		// Managed identity path via azidentity
		tokenReq, err := m.miCred.GetToken(ctx, azidentity.GetTokenOptions(struct {
			Scopes []string
		}{Scopes: resolvedScopes}))
		if err != nil {
			return "", fmt.Errorf("acquiring managed identity token for %s: %w", resourceURL, err)
		}
		m.tokenCache[cacheKey] = cachedToken{
			token:     tokenReq.Token,
			expiresAt: tokenReq.ExpiresOn,
		}
		return tokenReq.Token, nil
	}

	result, err := m.client.AcquireTokenByCredential(ctx, resolvedScopes)
	if err != nil {
		return "", fmt.Errorf("acquiring token for %s: %w", resourceURL, err)
	}

	m.tokenCache[cacheKey] = cachedToken{
		token:     result.AccessToken,
		expiresAt: result.ExpiresOn,
	}

	return result.AccessToken, nil
}

// AcquireTokenOnBehalfOf acquires a token on behalf of a user using the OBO flow.
// Implements AccessTokenProvider.
func (m *MsalAuth) AcquireTokenOnBehalfOf(ctx context.Context, scopes []string, userAssertion string) (string, error) {
	if m.miCred != nil {
		return "", fmt.Errorf("on-behalf-of flow is not supported with managed identity authentication")
	}
	if m.client == nil {
		return "", fmt.Errorf("MSAL client is not initialized")
	}

	result, err := m.client.AcquireTokenOnBehalfOf(ctx, userAssertion, scopes)
	if err != nil {
		return "", fmt.Errorf("acquiring on-behalf-of token: %w", err)
	}

	return result.AccessToken, nil
}

// NewMsalAuthFromEnv creates a MsalAuth from standard environment variables:
// AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET (or AZURE_CERTIFICATE_PATH).
func NewMsalAuthFromEnv() (*MsalAuth, error) {
	cfg := Config{
		TenantID:            os.Getenv("AZURE_TENANT_ID"),
		ClientID:            os.Getenv("AZURE_CLIENT_ID"),
		ClientSecret:        os.Getenv("AZURE_CLIENT_SECRET"),
		CertificatePath:     os.Getenv("AZURE_CERTIFICATE_PATH"),
		CertificatePassword: os.Getenv("AZURE_CERTIFICATE_PASSWORD"),
		Authority:           os.Getenv("AZURE_AUTHORITY"),
	}

	if os.Getenv("AZURE_USE_MANAGED_IDENTITY") == "true" {
		cfg.UseManagedIdentity = true
		cfg.UserAssignedClientID = os.Getenv("AZURE_USER_ASSIGNED_CLIENT_ID")
	}

	return NewMsalAuth(cfg)
}
