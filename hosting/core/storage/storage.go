// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package storage provides the storage abstractions for agent state persistence.
// Implementations include MemoryStorage for testing and in-process use,
// with Azure Blob and CosmosDB backends available in separate packages.
package storage

import "context"

// StoreItem is any value that can be persisted in storage.
// Implementations should be JSON-serializable structs or primitive values.
type StoreItem interface{}

// Storage is the interface for reading and writing agent state.
// All methods operate on multiple keys for efficiency.
type Storage interface {
	// Read retrieves items by keys. Returns a map of found items;
	// missing keys are omitted from the result rather than causing an error.
	Read(ctx context.Context, keys []string) (map[string]StoreItem, error)
	// Write persists items. Overwrites any existing items at the given keys.
	Write(ctx context.Context, changes map[string]StoreItem) error
	// Delete removes items by keys. Missing keys are silently ignored.
	Delete(ctx context.Context, keys []string) error
}
