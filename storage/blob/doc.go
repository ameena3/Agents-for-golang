// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package blob provides Azure Blob Storage backend for Microsoft 365 Agents SDK state.
// It implements the Storage interface from hosting/core/storage.
//
// Usage:
//
//	cfg := blob.Config{
//	    ConnectionString: os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
//	    ContainerName:    "agent-state",
//	}
//	store, err := blob.NewBlobStorage(ctx, cfg)
package blob
