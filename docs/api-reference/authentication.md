# API Reference: authentication

Import path: `github.com/ameena3/Agents-for-golang/authentication`

The `authentication` package provides MSAL-based OAuth authentication for agents.
`MsalAuth` implements `core.AccessTokenProvider` and handles token caching,
auto-refresh, and multiple credential types.

## MsalAuth

```go
type MsalAuth struct { /* unexported */ }
```

`MsalAuth` caches tokens in memory and automatically refreshes them 5 minutes before
expiry.

### Constructors

```go
// Create from an explicit Config struct.
func NewMsalAuth(cfg Config) (*MsalAuth, error)

// Create from standard Azure environment variables.
// Reads AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_CERTIFICATE_PATH, AZURE_CERTIFICATE_PASSWORD, AZURE_AUTHORITY,
// AZURE_USE_MANAGED_IDENTITY, AZURE_USER_ASSIGNED_CLIENT_ID.
func NewMsalAuthFromEnv() (*MsalAuth, error)
```

### Methods

```go
// GetAccessToken returns a token for the given resource and scopes.
// If forceRefresh is true, skips the cache and acquires a fresh token.
// Implements core.AccessTokenProvider.
func (m *MsalAuth) GetAccessToken(
    ctx context.Context,
    resourceURL string,
    scopes []string,
    forceRefresh bool,
) (string, error)

// AcquireTokenOnBehalfOf acquires a token on behalf of a user (OBO flow).
// Not supported with managed identity.
// Implements core.AccessTokenProvider.
func (m *MsalAuth) AcquireTokenOnBehalfOf(
    ctx context.Context,
    scopes []string,
    userAssertion string,
) (string, error)
```

## Config

```go
type Config struct {
    // TenantID is the Azure AD tenant ID. Required.
    TenantID string

    // ClientID is the application (client) ID. Required.
    ClientID string

    // ClientSecret is the application secret (client credential flow).
    // Mutually exclusive with CertificatePath and UseManagedIdentity.
    ClientSecret string

    // CertificatePath is the path to a PEM or PFX certificate file.
    // Mutually exclusive with ClientSecret and UseManagedIdentity.
    CertificatePath string

    // CertificatePassword is the PFX certificate password (if required).
    CertificatePassword string

    // Authority is the AAD authority URL.
    // Defaults to https://login.microsoftonline.com/{TenantID}.
    Authority string

    // UseManagedIdentity enables Azure Managed Identity authentication.
    // Mutually exclusive with ClientSecret and CertificatePath.
    UseManagedIdentity bool

    // UserAssignedClientID is the client ID for a user-assigned managed identity.
    // Only used when UseManagedIdentity is true.
    UserAssignedClientID string

    // Scopes are the default OAuth scopes to request.
    // If empty, defaults to {resourceURL}/.default.
    Scopes []string
}
```

## Authentication Flows

### Client Secret (most common for server agents)

```go
auth, err := authentication.NewMsalAuth(authentication.Config{
    TenantID:     os.Getenv("AZURE_TENANT_ID"),
    ClientID:     os.Getenv("AZURE_CLIENT_ID"),
    ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
})
```

### Certificate

```go
auth, err := authentication.NewMsalAuth(authentication.Config{
    TenantID:            os.Getenv("AZURE_TENANT_ID"),
    ClientID:            os.Getenv("AZURE_CLIENT_ID"),
    CertificatePath:     "/etc/secrets/bot-cert.pem",
    CertificatePassword: "",  // leave empty for PEM without passphrase
})
```

### Managed Identity (recommended for Azure-hosted agents)

System-assigned managed identity:

```go
auth, err := authentication.NewMsalAuth(authentication.Config{
    TenantID:           os.Getenv("AZURE_TENANT_ID"),
    ClientID:           os.Getenv("AZURE_CLIENT_ID"),
    UseManagedIdentity: true,
})
```

User-assigned managed identity:

```go
auth, err := authentication.NewMsalAuth(authentication.Config{
    TenantID:             os.Getenv("AZURE_TENANT_ID"),
    ClientID:             os.Getenv("AZURE_CLIENT_ID"),
    UseManagedIdentity:   true,
    UserAssignedClientID: os.Getenv("AZURE_USER_ASSIGNED_CLIENT_ID"),
})
```

### Environment Variables

`NewMsalAuthFromEnv` reads the following variables:

| Variable | Config Field |
|---|---|
| `AZURE_TENANT_ID` | `TenantID` |
| `AZURE_CLIENT_ID` | `ClientID` |
| `AZURE_CLIENT_SECRET` | `ClientSecret` |
| `AZURE_CERTIFICATE_PATH` | `CertificatePath` |
| `AZURE_CERTIFICATE_PASSWORD` | `CertificatePassword` |
| `AZURE_AUTHORITY` | `Authority` |
| `AZURE_USE_MANAGED_IDENTITY` | `UseManagedIdentity` (set to `"true"`) |
| `AZURE_USER_ASSIGNED_CLIENT_ID` | `UserAssignedClientID` |

## ConnectionManager

```go
type ConnectionManager struct { /* unexported */ }
```

`ConnectionManager` manages multiple named OAuth connections. It is used in
multi-skill and agent-to-agent scenarios where different connections require
different credentials.

```go
func NewConnectionManager() *ConnectionManager
func (m *ConnectionManager) Add(name string, provider core.AccessTokenProvider)
func (m *ConnectionManager) Get(name string) (core.AccessTokenProvider, bool)
```

## Notes

- `MsalAuth` caches tokens and refreshes them 5 minutes before the `ExpiresOn`
  timestamp returned by MSAL. Token cache is per-`MsalAuth` instance; create one
  instance per app lifecycle and share it.
- The OBO flow (`AcquireTokenOnBehalfOf`) is not supported with managed identity.
- For local development, set `AllowUnauthenticated: true` in `nethttp.ServerConfig`
  to skip JWT validation entirely. Never use this in production.
- `MsalAuth` implements `core.AccessTokenProvider`, which is the interface expected
  by `RestChannelServiceClientFactory` and connector clients for service-to-service
  authentication.
