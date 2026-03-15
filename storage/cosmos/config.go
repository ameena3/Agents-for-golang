// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cosmos

// Config holds configuration for Azure Cosmos DB storage.
type Config struct {
	// Endpoint is the Cosmos DB account endpoint URL.
	Endpoint string
	// Key is the Cosmos DB account key. If empty, uses managed identity.
	Key string
	// DatabaseID is the Cosmos DB database name.
	DatabaseID string
	// ContainerID is the Cosmos DB container name.
	ContainerID string
	// PartitionKey is the partition key path. Defaults to "/id".
	PartitionKey string
	// UseManagedIdentity uses Azure Managed Identity.
	UseManagedIdentity bool
}
