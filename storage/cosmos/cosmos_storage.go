// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cosmos

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/microsoft/agents-sdk-go/hosting/core/storage"
)

// CosmosDBStorage implements storage.Storage using Azure Cosmos DB.
// Each state key is stored as a separate item in the configured container.
type CosmosDBStorage struct {
	client *azcosmos.ContainerClient
	config Config
}

// NewCosmosDBStorage creates a new CosmosDBStorage and verifies connectivity.
func NewCosmosDBStorage(ctx context.Context, cfg Config) (*CosmosDBStorage, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("cosmos: Endpoint is required")
	}
	if cfg.DatabaseID == "" || cfg.ContainerID == "" {
		return nil, fmt.Errorf("cosmos: DatabaseID and ContainerID are required")
	}
	if cfg.PartitionKey == "" {
		cfg.PartitionKey = "/id"
	}

	if cfg.Key == "" {
		return nil, fmt.Errorf("cosmos: Key is required (managed identity not yet implemented)")
	}

	cred, err := azcosmos.NewKeyCredential(cfg.Key)
	if err != nil {
		return nil, fmt.Errorf("cosmos: creating key credential: %w", err)
	}

	cosmosClient, err := azcosmos.NewClientWithKey(cfg.Endpoint, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("cosmos: creating client: %w", err)
	}

	container, err := cosmosClient.NewContainer(cfg.DatabaseID, cfg.ContainerID)
	if err != nil {
		return nil, fmt.Errorf("cosmos: getting container client: %w", err)
	}

	return &CosmosDBStorage{client: container, config: cfg}, nil
}

// Read retrieves items by key. Missing keys are silently skipped.
func (s *CosmosDBStorage) Read(ctx context.Context, keys []string) (map[string]storage.StoreItem, error) {
	result := make(map[string]storage.StoreItem)
	for _, key := range keys {
		safeKey := SanitizeKey(key)
		pk := azcosmos.NewPartitionKeyString(safeKey)
		resp, err := s.client.ReadItem(ctx, pk, safeKey, nil)
		if err != nil {
			// Treat any error as not-found for now (404 and others)
			continue
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(resp.Value, &doc); err != nil {
			return nil, fmt.Errorf("cosmos: unmarshaling key %s: %w", key, err)
		}
		// Extract the "data" field written by Write
		if data, ok := doc["data"]; ok {
			result[key] = data
		} else {
			result[key] = doc
		}
	}
	return result, nil
}

// Write persists items. Each item is upserted as a Cosmos document with "id" = sanitized key.
func (s *CosmosDBStorage) Write(ctx context.Context, changes map[string]storage.StoreItem) error {
	for key, value := range changes {
		safeKey := SanitizeKey(key)
		doc := map[string]interface{}{
			"id":   safeKey,
			"data": value,
		}
		data, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("cosmos: marshaling key %s: %w", key, err)
		}
		pk := azcosmos.NewPartitionKeyString(safeKey)
		if _, err := s.client.UpsertItem(ctx, pk, data, nil); err != nil {
			return fmt.Errorf("cosmos: writing key %s: %w", key, err)
		}
	}
	return nil
}

// Delete removes items. Missing items are silently ignored.
func (s *CosmosDBStorage) Delete(ctx context.Context, keys []string) error {
	for _, key := range keys {
		safeKey := SanitizeKey(key)
		pk := azcosmos.NewPartitionKeyString(safeKey)
		if _, err := s.client.DeleteItem(ctx, pk, safeKey, nil); err != nil {
			continue // ignore errors (e.g. 404)
		}
	}
	return nil
}
