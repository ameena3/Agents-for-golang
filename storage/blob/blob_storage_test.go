// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blob_test

import (
	"context"
	"testing"

	"github.com/microsoft/agents-sdk-go/storage/blob"
)

func TestNewBlobStorage_EmptyConfig(t *testing.T) {
	_, err := blob.NewBlobStorage(context.Background(), blob.Config{})
	if err == nil {
		t.Fatal("expected error for empty config")
	}
}

func TestNewBlobStorage_EmptyContainerName(t *testing.T) {
	_, err := blob.NewBlobStorage(context.Background(), blob.Config{
		ConnectionString: "DefaultEndpointsProtocol=https;AccountName=test;AccountKey=dGVzdA==;EndpointSuffix=core.windows.net",
		ContainerName:    "",
	})
	if err == nil {
		t.Fatal("expected error when ContainerName is empty")
	}
}

func TestNewBlobStorage_AccountURLWithoutManagedIdentity(t *testing.T) {
	_, err := blob.NewBlobStorage(context.Background(), blob.Config{
		AccountURL:         "https://myaccount.blob.core.windows.net",
		ContainerName:      "my-container",
		UseManagedIdentity: false,
	})
	if err == nil {
		t.Fatal("expected error when AccountURL is set but UseManagedIdentity is false")
	}
}

func TestNewBlobStorage_NoConnectionStringNoAccountURL(t *testing.T) {
	_, err := blob.NewBlobStorage(context.Background(), blob.Config{
		ContainerName: "my-container",
	})
	if err == nil {
		t.Fatal("expected error when neither ConnectionString nor AccountURL is provided")
	}
}

func TestNewBlobStorage_InvalidConnectionString(t *testing.T) {
	// An invalid connection string should cause an error from the Azure SDK.
	_, err := blob.NewBlobStorage(context.Background(), blob.Config{
		ConnectionString: "invalid-connection-string",
		ContainerName:    "my-container",
	})
	if err == nil {
		t.Fatal("expected error for invalid connection string")
	}
}

// TestConfig_Fields verifies the Config struct fields are accessible.
func TestConfig_Fields(t *testing.T) {
	cfg := blob.Config{
		ConnectionString:   "some-string",
		AccountURL:         "https://account.blob.core.windows.net",
		ContainerName:      "container",
		UseManagedIdentity: true,
	}
	if cfg.ConnectionString == "" {
		t.Error("ConnectionString should not be empty")
	}
	if cfg.AccountURL == "" {
		t.Error("AccountURL should not be empty")
	}
	if cfg.ContainerName == "" {
		t.Error("ContainerName should not be empty")
	}
	if !cfg.UseManagedIdentity {
		t.Error("UseManagedIdentity should be true")
	}
}
