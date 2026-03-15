// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package cosmos provides Azure Cosmos DB backend for Microsoft 365 Agents SDK state.
// It implements the Storage interface from hosting/core/storage.
//
// Usage:
//
//	cfg := cosmos.Config{
//	    Endpoint:    os.Getenv("COSMOS_ENDPOINT"),
//	    Key:         os.Getenv("COSMOS_KEY"),
//	    DatabaseID:  "agent-db",
//	    ContainerID: "agent-state",
//	}
//	store, err := cosmos.NewCosmosDBStorage(ctx, cfg)
package cosmos
