# API Reference: storage/cosmos

Import path: `github.com/microsoft/agents-sdk-go/storage/cosmos`

The `storage/cosmos` package implements `storage.Storage` using Azure Cosmos DB.
Each state key is stored as a separate item in a configured container, using the
sanitized key as both the document `id` and the partition key.

## CosmosDBStorage

```go
type CosmosDBStorage struct { /* unexported */ }
```

### Constructor

```go
func NewCosmosDBStorage(ctx context.Context, cfg Config) (*CosmosDBStorage, error)
```

Returns an initialized `*CosmosDBStorage`. Returns an error if the configuration is
invalid or if the account is unreachable.

### Methods

`CosmosDBStorage` implements `storage.Storage`:

```go
// Read retrieves items by key. Missing keys are silently omitted.
func (s *CosmosDBStorage) Read(ctx context.Context, keys []string) (map[string]storage.StoreItem, error)

// Write upserts items. Each key is stored as a Cosmos document with id = sanitized key.
func (s *CosmosDBStorage) Write(ctx context.Context, changes map[string]storage.StoreItem) error

// Delete removes items by key. Missing keys are silently ignored.
func (s *CosmosDBStorage) Delete(ctx context.Context, keys []string) error
```

## Config

```go
type Config struct {
    // Endpoint is the Cosmos DB account URL. Required.
    // Example: https://myaccount.documents.azure.com:443/
    Endpoint string

    // Key is the Cosmos DB account key.
    // Required unless UseManagedIdentity is true (managed identity not yet implemented).
    Key string

    // DatabaseID is the Cosmos DB database name. Required.
    DatabaseID string

    // ContainerID is the Cosmos DB container name. Required.
    ContainerID string

    // PartitionKey is the partition key path. Defaults to "/id".
    PartitionKey string

    // UseManagedIdentity uses Azure Managed Identity. Not yet implemented.
    UseManagedIdentity bool
}
```

## Usage Example

### Account Key

```go
import (
    "context"
    "github.com/microsoft/agents-sdk-go/storage/cosmos"
)

store, err := cosmos.NewCosmosDBStorage(context.Background(), cosmos.Config{
    Endpoint:    "https://myaccount.documents.azure.com:443/",
    Key:         os.Getenv("COSMOS_KEY"),
    DatabaseID:  "AgentDB",
    ContainerID: "AgentState",
})
if err != nil {
    log.Fatal(err)
}
```

### Using with AgentApplication

```go
agent := app.New[AppState](app.AppOptions[AppState]{
    Storage: store, // *cosmos.CosmosDBStorage satisfies storage.Storage
})
```

## Document Structure

Each Cosmos item has the following shape:

```json
{
  "id": "<sanitized-key>",
  "data": { ... }
}
```

The `data` field holds the serialized state value. On read, the `data` field is
extracted and returned as the `StoreItem`. If a document does not have a `data`
field (e.g., documents written outside the SDK), the entire document is returned.

## Key Sanitization

The `SanitizeKey` function (exported) replaces characters not safe for Cosmos
document IDs with URL-percent-encoded equivalents:

```go
func SanitizeKey(key string) string
```

This is applied automatically by `Read`, `Write`, and `Delete`.

## Notes

- Partition key defaults to `/id`. This means each document is its own partition,
  which is appropriate for agent state (many small documents, each accessed
  individually). For high-throughput scenarios, consider a custom partition strategy.
- Managed identity support is listed in `Config` but is not yet implemented; the
  `Key` field is required for now.
- All state values must be JSON-serializable.
- Errors from `ReadItem` (including 404) are currently treated as not-found and
  silently skipped; other errors are also skipped. For stricter error handling,
  wrap `CosmosDBStorage` in a custom type.
