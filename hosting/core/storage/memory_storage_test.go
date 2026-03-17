// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package storage_test

import (
	"context"
	"sync"
	"testing"

	"github.com/ameena3/Agents-for-golang/hosting/core/storage"
)

func TestMemoryStorage_WriteAndRead(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	changes := map[string]storage.StoreItem{
		"key1": map[string]interface{}{"value": "hello"},
		"key2": 42,
	}
	if err := s.Write(ctx, changes); err != nil {
		t.Fatalf("Write: unexpected error: %v", err)
	}

	result, err := s.Read(ctx, []string{"key1", "key2"})
	if err != nil {
		t.Fatalf("Read: unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Read: expected 2 items, got %d", len(result))
	}
	if _, ok := result["key1"]; !ok {
		t.Error("Read: expected key1 to be present")
	}
	if _, ok := result["key2"]; !ok {
		t.Error("Read: expected key2 to be present")
	}
}

func TestMemoryStorage_ReadMissingKey(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	result, err := s.Read(ctx, []string{"nonexistent"})
	if err != nil {
		t.Fatalf("Read: unexpected error for missing key: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Read: expected empty map for missing key, got %d items", len(result))
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	if err := s.Write(ctx, map[string]storage.StoreItem{"key1": "val1", "key2": "val2"}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if err := s.Delete(ctx, []string{"key1"}); err != nil {
		t.Fatalf("Delete: unexpected error: %v", err)
	}

	result, err := s.Read(ctx, []string{"key1", "key2"})
	if err != nil {
		t.Fatalf("Read after delete: %v", err)
	}
	if _, ok := result["key1"]; ok {
		t.Error("Delete: key1 should have been deleted")
	}
	if _, ok := result["key2"]; !ok {
		t.Error("Delete: key2 should still be present")
	}
}

func TestMemoryStorage_DeleteMissingKey(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	// Deleting a key that doesn't exist should not return an error.
	if err := s.Delete(ctx, []string{"ghost"}); err != nil {
		t.Fatalf("Delete of missing key: unexpected error: %v", err)
	}
}

func TestMemoryStorage_WriteOverwrite(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	if err := s.Write(ctx, map[string]storage.StoreItem{"k": "original"}); err != nil {
		t.Fatalf("first Write: %v", err)
	}
	if err := s.Write(ctx, map[string]storage.StoreItem{"k": "updated"}); err != nil {
		t.Fatalf("second Write: %v", err)
	}

	result, err := s.Read(ctx, []string{"k"})
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	// The stored value is round-tripped through JSON so the raw value is a string.
	got, ok := result["k"]
	if !ok {
		t.Fatal("Read: expected key k to be present")
	}
	// JSON unmarshals strings as string, numbers as float64, etc.
	if got != "updated" {
		t.Errorf("Read: expected \"updated\", got %v (%T)", got, got)
	}
}

func TestMemoryStorage_ReadEmptyKeys(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	_, err := s.Read(ctx, nil)
	if err == nil {
		t.Error("Read with nil keys: expected error, got nil")
	}

	_, err = s.Read(ctx, []string{})
	if err == nil {
		t.Error("Read with empty keys: expected error, got nil")
	}
}

func TestMemoryStorage_WriteEmptyChanges(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	err := s.Write(ctx, nil)
	if err == nil {
		t.Error("Write with nil changes: expected error, got nil")
	}

	err = s.Write(ctx, map[string]storage.StoreItem{})
	if err == nil {
		t.Error("Write with empty changes: expected error, got nil")
	}
}

func TestMemoryStorage_DeleteEmptyKeys(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	err := s.Delete(ctx, nil)
	if err == nil {
		t.Error("Delete with nil keys: expected error, got nil")
	}

	err = s.Delete(ctx, []string{})
	if err == nil {
		t.Error("Delete with empty keys: expected error, got nil")
	}
}

func TestMemoryStorage_ReadEmptyStringKey(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	_, err := s.Read(ctx, []string{""})
	if err == nil {
		t.Error("Read with empty-string key: expected error, got nil")
	}
}

func TestMemoryStorage_WriteEmptyStringKey(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	err := s.Write(ctx, map[string]storage.StoreItem{"": "v"})
	if err == nil {
		t.Error("Write with empty-string key: expected error, got nil")
	}
}

func TestMemoryStorage_DeleteEmptyStringKey(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	err := s.Delete(ctx, []string{""})
	if err == nil {
		t.Error("Delete with empty-string key: expected error, got nil")
	}
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	s := storage.NewMemoryStorage()
	ctx := context.Background()

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Half writers, half readers running simultaneously.
	for i := 0; i < goroutines; i++ {
		go func(n int) {
			defer wg.Done()
			if n%2 == 0 {
				// Writer
				_ = s.Write(ctx, map[string]storage.StoreItem{
					"shared": n,
				})
			} else {
				// Reader – ignore not-found, we just want no races.
				_, _ = s.Read(ctx, []string{"shared"})
			}
		}(i)
	}

	wg.Wait()
}

func TestMemoryStorage_ImplementsStorageInterface(t *testing.T) {
	// Compile-time assertion: *MemoryStorage must implement Storage.
	var _ storage.Storage = storage.NewMemoryStorage()
}
