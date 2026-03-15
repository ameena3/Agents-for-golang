// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blob

// Config holds configuration for Azure Blob Storage.
type Config struct {
	// ConnectionString is the Azure Storage connection string.
	// Mutually exclusive with AccountURL + Credential.
	ConnectionString string
	// AccountURL is the storage account URL (e.g., https://account.blob.core.windows.net).
	AccountURL string
	// ContainerName is the blob container to use for state storage.
	ContainerName string
	// UseManagedIdentity uses Azure Managed Identity for authentication.
	UseManagedIdentity bool
}
