// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cosmos_test

import (
	"context"
	"testing"

	"github.com/ameena3/Agents-for-golang/storage/cosmos"
)

// --- SanitizeKey tests ---

func TestSanitizeKey_NoSpecialChars(t *testing.T) {
	input := "simple-key_123"
	result := cosmos.SanitizeKey(input)
	if result != input {
		t.Errorf("expected %q unchanged, got %q", input, result)
	}
}

func TestSanitizeKey_ForwardSlash(t *testing.T) {
	result := cosmos.SanitizeKey("conv/user")
	expected := "conv-FSLASH-user"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSanitizeKey_BackSlash(t *testing.T) {
	result := cosmos.SanitizeKey(`conv\user`)
	expected := "conv-BSLASH-user"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSanitizeKey_QuestionMark(t *testing.T) {
	result := cosmos.SanitizeKey("key?value")
	expected := "key-QMARK-value"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSanitizeKey_Hash(t *testing.T) {
	result := cosmos.SanitizeKey("key#fragment")
	expected := "key-HASH-fragment"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSanitizeKey_MultipleSpecialChars(t *testing.T) {
	result := cosmos.SanitizeKey("conv/user?id#1")
	expected := "conv-FSLASH-user-QMARK-id-HASH-1"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSanitizeKey_EmptyString(t *testing.T) {
	result := cosmos.SanitizeKey("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestSanitizeKey_Idempotent(t *testing.T) {
	// Sanitizing an already-sanitized key should not change it further.
	input := "conv-FSLASH-user-QMARK-id"
	result := cosmos.SanitizeKey(input)
	if result != input {
		t.Errorf("expected %q unchanged, got %q", input, result)
	}
}

// --- NewCosmosDBStorage validation tests ---

func TestNewCosmosDBStorage_EmptyEndpoint(t *testing.T) {
	_, err := cosmos.NewCosmosDBStorage(context.Background(), cosmos.Config{
		DatabaseID:  "db",
		ContainerID: "container",
		Key:         "somekey",
	})
	if err == nil {
		t.Fatal("expected error when Endpoint is empty")
	}
}

func TestNewCosmosDBStorage_EmptyDatabaseID(t *testing.T) {
	_, err := cosmos.NewCosmosDBStorage(context.Background(), cosmos.Config{
		Endpoint:    "https://myaccount.documents.azure.com:443/",
		ContainerID: "container",
		Key:         "somekey",
	})
	if err == nil {
		t.Fatal("expected error when DatabaseID is empty")
	}
}

func TestNewCosmosDBStorage_EmptyContainerID(t *testing.T) {
	_, err := cosmos.NewCosmosDBStorage(context.Background(), cosmos.Config{
		Endpoint:   "https://myaccount.documents.azure.com:443/",
		DatabaseID: "db",
		Key:        "somekey",
	})
	if err == nil {
		t.Fatal("expected error when ContainerID is empty")
	}
}

func TestNewCosmosDBStorage_EmptyKey(t *testing.T) {
	_, err := cosmos.NewCosmosDBStorage(context.Background(), cosmos.Config{
		Endpoint:    "https://myaccount.documents.azure.com:443/",
		DatabaseID:  "db",
		ContainerID: "container",
	})
	if err == nil {
		t.Fatal("expected error when Key is empty")
	}
}

func TestNewCosmosDBStorage_AllEmpty(t *testing.T) {
	_, err := cosmos.NewCosmosDBStorage(context.Background(), cosmos.Config{})
	if err == nil {
		t.Fatal("expected error for empty config")
	}
}

// TestConfig_Fields verifies the Config struct fields are accessible.
func TestConfig_Fields(t *testing.T) {
	cfg := cosmos.Config{
		Endpoint:           "https://myaccount.documents.azure.com:443/",
		Key:                "base64key==",
		DatabaseID:         "mydb",
		ContainerID:        "mycontainer",
		PartitionKey:       "/id",
		UseManagedIdentity: false,
	}
	if cfg.Endpoint == "" {
		t.Error("Endpoint should not be empty")
	}
	if cfg.Key == "" {
		t.Error("Key should not be empty")
	}
	if cfg.DatabaseID == "" {
		t.Error("DatabaseID should not be empty")
	}
	if cfg.ContainerID == "" {
		t.Error("ContainerID should not be empty")
	}
	if cfg.PartitionKey == "" {
		t.Error("PartitionKey should not be empty")
	}
}
