// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authorization

import "context"

// ConnectionConfig holds the OAuth configuration for a named connection.
// It is used to configure MSAL-based authentication for agent-to-service
// and service-to-service authentication flows.
type ConnectionConfig struct {
	// ClientID is the Azure AD application (client) ID.
	ClientID string
	// ClientSecret is the client secret for confidential-client flows.
	// Leave empty when using certificate-based authentication.
	ClientSecret string
	// TenantID is the Azure AD tenant ID. Use "common" or "organizations"
	// for multi-tenant applications.
	TenantID string
	// CertificatePath is the path to a PEM or PFX certificate file used
	// for certificate-based authentication. Optional.
	CertificatePath string
	// CertificatePassword is the passphrase for the certificate file.
	// Only required for encrypted PFX files.
	CertificatePassword string
	// Authority is an optional override for the authority URL
	// (e.g. https://login.microsoftonline.com/{tenantID}). When empty the
	// SDK constructs the authority from TenantID.
	Authority string
}

// ConnectionManager manages named OAuth connections and retrieves access
// tokens for downstream services on behalf of the agent.
type ConnectionManager interface {
	// GetToken retrieves an access token for the named connection.
	// The returned string is a bearer token ready for use in an
	// Authorization header.
	GetToken(ctx context.Context, connectionName string) (string, error)
}
