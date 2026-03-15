// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package state

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/microsoft/agents-sdk-go/hosting/core/storage"
)

// UserState holds per-user persistent state, keyed across all conversations by
// "{channelID}/users/{userID}". Keyed by "{channelID}/users/{userID}" in storage.
type UserState struct {
	stor     storage.Storage
	data     map[string]interface{}
	modified bool
}

// NewUserState creates a UserState backed by the given storage.
// If s is nil the state is in-memory only and Save is a no-op.
func NewUserState(s storage.Storage) *UserState {
	return &UserState{
		stor: s,
		data: make(map[string]interface{}),
	}
}

// Get returns the value stored under name and whether it was found.
func (us *UserState) Get(name string) (interface{}, bool) {
	v, ok := us.data[name]
	return v, ok
}

// Set stores a value under name and marks the state as modified.
func (us *UserState) Set(name string, value interface{}) {
	us.data[name] = value
	us.modified = true
}

// Delete removes the value stored under name and marks the state as modified.
func (us *UserState) Delete(name string) {
	if _, ok := us.data[name]; ok {
		delete(us.data, name)
		us.modified = true
	}
}

// IsModified reports whether any values have changed since the last load or save.
func (us *UserState) IsModified() bool {
	return us.modified
}

// userStorageKey builds the storage key for the given channel / user identifiers.
func userStorageKey(channelID, userID string) string {
	return fmt.Sprintf("%s/users/%s", channelID, userID)
}

// Load reads the user state from storage using the provided identifiers.
// If the key is not found in storage the state is initialized to empty.
func (us *UserState) Load(ctx context.Context, channelID, userID string) error {
	if us.stor == nil {
		return nil
	}
	key := userStorageKey(channelID, userID)
	items, err := us.stor.Read(ctx, []string{key})
	if err != nil {
		return fmt.Errorf("user state load: %w", err)
	}
	if raw, ok := items[key]; ok {
		b, err := json.Marshal(raw)
		if err != nil {
			return fmt.Errorf("user state marshal: %w", err)
		}
		newData := make(map[string]interface{})
		if err := json.Unmarshal(b, &newData); err != nil {
			return fmt.Errorf("user state unmarshal: %w", err)
		}
		us.data = newData
	} else {
		us.data = make(map[string]interface{})
	}
	us.modified = false
	return nil
}

// Save writes the user state to storage.
// If the state has not been modified since the last load, this is a no-op.
func (us *UserState) Save(ctx context.Context, channelID, userID string) error {
	if us.stor == nil || !us.modified {
		return nil
	}
	key := userStorageKey(channelID, userID)
	if err := us.stor.Write(ctx, map[string]storage.StoreItem{key: us.data}); err != nil {
		return fmt.Errorf("user state save: %w", err)
	}
	us.modified = false
	return nil
}

// Clear removes all values from the in-memory state and marks it as modified.
func (us *UserState) Clear() {
	us.data = make(map[string]interface{})
	us.modified = true
}
