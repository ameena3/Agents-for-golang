# API Reference: storage/blob

Import path: `github.com/ameena3/Agents-for-golang/storage/blob`

The `storage/blob` package implements `storage.Storage` using Azure Blob Storage.
Each state key is stored as a separate blob within a configured container. The
container is created automatically if it does not exist.

## BlobStorage

```go
type BlobStorage struct { /* unexported */ }
```

### Constructor

```go
func NewBlobStorage(ctx context.Context, cfg Config) (*BlobStorage, error)
```

Returns an initialized `*BlobStorage`. The container is created if it does not exist.
Returns an error if the configuration is invalid or if the storage account is
unreachable.

### Methods

`BlobStorage` implements `storage.Storage`:

```go
// Read retrieves items by key. Missing keys are silently omitted from the result.
func (s *BlobStorage) Read(ctx context.Context, keys []string) (map[string]storage.StoreItem, error)

// Write persists items. Each key is stored as a separate blob. Overwrites existing blobs.
func (s *BlobStorage) Write(ctx context.Context, changes map[string]storage.StoreItem) error

// Delete removes blobs by key. Missing keys are silently ignored.
func (s *BlobStorage) Delete(ctx context.Context, keys []string) error
```

## Config

```go
type Config struct {
    // ConnectionString is the Azure Storage connection string.
    // Mutually exclusive with AccountURL + UseManagedIdentity.
    ConnectionString string

    // AccountURL is the storage account endpoint URL.
    // Example: https://myaccount.blob.core.windows.net
    // Requires UseManagedIdentity=true when used.
    AccountURL string

    // ContainerName is the blob container for state storage. Required.
    ContainerName string

    // UseManagedIdentity uses Azure Managed Identity for authentication.
    // Must be true when AccountURL is set.
    UseManagedIdentity bool
}
```

Exactly one of `ConnectionString` or (`AccountURL` + `UseManagedIdentity: true`) must
be provided.

## Usage Examples

### Connection String (development / testing)

```go
import (
    "context"
    "github.com/ameena3/Agents-for-golang/storage/blob"
)

store, err := blob.NewBlobStorage(context.Background(), blob.Config{
    ConnectionString: "DefaultEndpointsProtocol=https;AccountName=...",
    ContainerName:    "agent-state",
})
if err != nil {
    log.Fatal(err)
}
```

The local Azurite emulator connection string for local development:

```go
blob.Config{
    ConnectionString: "UseDevelopmentStorage=true",
    ContainerName:    "agent-state",
}
```

### Managed Identity (production)

```go
store, err := blob.NewBlobStorage(context.Background(), blob.Config{
    AccountURL:         "https://myaccount.blob.core.windows.net",
    ContainerName:      "agent-state",
    UseManagedIdentity: true,
})
```

### Using with AgentApplication

```go
agent := app.New[AppState](app.AppOptions[AppState]{
    Storage: store, // *blob.BlobStorage satisfies storage.Storage
})
```

## Key Sanitization

Storage keys from the SDK (e.g., `conv/msteams/19:xxx`) may contain characters not
allowed in blob names. `BlobStorage` percent-encodes characters that are not in the
safe set `[a-zA-Z0-9-_./ ]` before using them as blob names. This is transparent to
callers; you always use the original key strings.

## Notes

- Values are stored as JSON. All state values must be JSON-serializable.
- Concurrency: the underlying `azblob` SDK handles HTTP-level retries. For optimistic
  concurrency, consider implementing ETag-based checks in a custom storage wrapper.
- The container is created on the first call to `NewBlobStorage` if it does not exist.
  Ensure the configured credential has `Storage Blob Data Contributor` role or
  equivalent.
