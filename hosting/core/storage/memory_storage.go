// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// MemoryStorage implements Storage using an in-memory map.
// It is safe for concurrent use but all data is lost when the process exits.
// It is intended for testing and development scenarios.
type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string]json.RawMessage
}

// NewMemoryStorage creates a new empty MemoryStorage instance.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]json.RawMessage),
	}
}

// Read retrieves items for the given keys. Returns only the keys that were found.
// Returns an error if keys is nil or empty, or if any key is an empty string.
func (m *MemoryStorage) Read(ctx context.Context, keys []string) (map[string]StoreItem, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("storage.Read: keys are required when reading")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]StoreItem, len(keys))
	for _, key := range keys {
		if key == "" {
			return nil, fmt.Errorf("storage.Read: key cannot be empty")
		}
		raw, ok := m.data[key]
		if !ok {
			continue
		}
		// Deserialize into a generic interface{} value.
		var val interface{}
		if err := json.Unmarshal(raw, &val); err != nil {
			return nil, fmt.Errorf("storage.Read: failed to deserialize key %q: %w", key, err)
		}
		result[key] = val
	}
	return result, nil
}

// Write persists the given changes. Existing values at the same keys are overwritten.
// Returns an error if changes is nil or empty, or if any key is an empty string.
func (m *MemoryStorage) Write(ctx context.Context, changes map[string]StoreItem) error {
	if len(changes) == 0 {
		return fmt.Errorf("storage.Write: changes are required when writing")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for key, val := range changes {
		if key == "" {
			return fmt.Errorf("storage.Write: key cannot be empty")
		}
		raw, err := json.Marshal(val)
		if err != nil {
			return fmt.Errorf("storage.Write: failed to serialize value for key %q: %w", key, err)
		}
		m.data[key] = raw
	}
	return nil
}

// Delete removes the items with the given keys. Missing keys are silently ignored.
// Returns an error if keys is nil or empty, or if any key is an empty string.
func (m *MemoryStorage) Delete(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return fmt.Errorf("storage.Delete: keys are required when deleting")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range keys {
		if key == "" {
			return fmt.Errorf("storage.Delete: key cannot be empty")
		}
		delete(m.data, key)
	}
	return nil
}
