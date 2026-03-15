// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authentication

// Config holds the configuration for MSAL authentication.
type Config struct {
	// TenantID is the Azure AD tenant ID.
	TenantID string
	// ClientID is the application (client) ID.
	ClientID string
	// ClientSecret is the application secret (for client credential flow).
	// Mutually exclusive with CertificatePath.
	ClientSecret string
	// CertificatePath is the path to a PEM/PFX certificate file.
	// Mutually exclusive with ClientSecret.
	CertificatePath string
	// CertificatePassword is the certificate password (for PFX files).
	CertificatePassword string
	// Authority is the AAD authority URL. Defaults to https://login.microsoftonline.com/{TenantID}
	Authority string
	// UseManagedIdentity uses Azure Managed Identity instead of client credentials.
	UseManagedIdentity bool
	// UserAssignedClientID is the client ID for user-assigned managed identity.
	UserAssignedClientID string
	// Scopes are the default OAuth scopes to request.
	Scopes []string
}
