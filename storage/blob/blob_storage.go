// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blob

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/microsoft/agents-sdk-go/hosting/core/storage"
)

// BlobStorage implements storage.Storage using Azure Blob Storage.
// Each key is stored as a separate blob named after the key.
type BlobStorage struct {
	client        *azblob.Client
	containerName string
}

// NewBlobStorage creates and initializes a BlobStorage.
// It creates the container if it does not exist.
func NewBlobStorage(ctx context.Context, cfg Config) (*BlobStorage, error) {
	if cfg.ContainerName == "" {
		return nil, fmt.Errorf("ContainerName is required")
	}

	var client *azblob.Client
	var err error

	switch {
	case cfg.ConnectionString != "":
		client, err = azblob.NewClientFromConnectionString(cfg.ConnectionString, nil)
		if err != nil {
			return nil, fmt.Errorf("creating blob client from connection string: %w", err)
		}
	case cfg.AccountURL != "" && cfg.UseManagedIdentity:
		cred, credErr := azidentity.NewManagedIdentityCredential(nil)
		if credErr != nil {
			return nil, fmt.Errorf("creating managed identity credential: %w", credErr)
		}
		client, err = azblob.NewClient(cfg.AccountURL, cred, nil)
		if err != nil {
			return nil, fmt.Errorf("creating blob client with managed identity: %w", err)
		}
	case cfg.AccountURL != "":
		return nil, fmt.Errorf("AccountURL requires UseManagedIdentity=true or use ConnectionString")
	default:
		return nil, fmt.Errorf("ConnectionString or AccountURL with UseManagedIdentity is required")
	}

	bs := &BlobStorage{client: client, containerName: cfg.ContainerName}
	if err := bs.initialize(ctx); err != nil {
		return nil, err
	}
	return bs, nil
}

func (s *BlobStorage) initialize(ctx context.Context) error {
	_, err := s.client.CreateContainer(ctx, s.containerName, nil)
	if err != nil {
		if respErr := errors.AsType[*azcore.ResponseError](err); respErr != nil && respErr.ErrorCode == string(bloberror.ContainerAlreadyExists) {
			return nil
		}
		return fmt.Errorf("creating container %s: %w", s.containerName, err)
	}
	return nil
}

// Read retrieves items by key from blob storage.
// Missing keys are silently omitted from the result.
func (s *BlobStorage) Read(ctx context.Context, keys []string) (map[string]storage.StoreItem, error) {
	result := make(map[string]storage.StoreItem)
	for _, key := range keys {
		blobName := sanitizeKey(key)
		resp, err := s.client.DownloadStream(ctx, s.containerName, blobName, nil)
		if err != nil {
			if respErr := errors.AsType[*azcore.ResponseError](err); respErr != nil && respErr.StatusCode == 404 {
				continue // key not found is OK
			}
			return nil, fmt.Errorf("reading key %s: %w", key, err)
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading body for key %s: %w", key, err)
		}
		var item interface{}
		if err := json.Unmarshal(data, &item); err != nil {
			return nil, fmt.Errorf("unmarshaling key %s: %w", key, err)
		}
		result[key] = item
	}
	return result, nil
}

// Write stores items to blob storage.
// Overwrites any existing blobs at the given keys.
func (s *BlobStorage) Write(ctx context.Context, changes map[string]storage.StoreItem) error {
	for key, value := range changes {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshaling key %s: %w", key, err)
		}
		blobName := sanitizeKey(key)
		_, err = s.client.UploadStream(ctx, s.containerName, blobName, bytes.NewReader(data), nil)
		if err != nil {
			return fmt.Errorf("writing key %s: %w", key, err)
		}
	}
	return nil
}

// Delete removes items from blob storage.
// Missing keys are silently ignored.
func (s *BlobStorage) Delete(ctx context.Context, keys []string) error {
	for _, key := range keys {
		blobName := sanitizeKey(key)
		_, err := s.client.DeleteBlob(ctx, s.containerName, blobName, nil)
		if err != nil {
			if respErr := errors.AsType[*azcore.ResponseError](err); respErr != nil && respErr.StatusCode == 404 {
				continue // already gone
			}
			return fmt.Errorf("deleting key %s: %w", key, err)
		}
	}
	return nil
}

// sanitizeKey replaces characters that are not valid in blob names with percent-encoded equivalents.
func sanitizeKey(key string) string {
	result := make([]byte, 0, len(key))
	for _, c := range []byte(key) {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9',
			c == '-', c == '_', c == '.', c == '/':
			result = append(result, c)
		default:
			result = append(result, []byte(fmt.Sprintf("%%%02X", c))...)
		}
	}
	return string(result)
}
