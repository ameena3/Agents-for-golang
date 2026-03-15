// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package state

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/microsoft/agents-sdk-go/hosting/core/storage"
)

// ConversationState holds per-conversation persistent state.
// The storage key is "{channelID}/conversations/{conversationID}", matching
// the Python SDK's ConversationState.get_storage_key implementation.
type ConversationState struct {
	stor     storage.Storage
	data     map[string]interface{}
	modified bool
}

// NewConversationState creates a ConversationState backed by the given storage.
// If s is nil the state is in-memory only and Save is a no-op.
func NewConversationState(s storage.Storage) *ConversationState {
	return &ConversationState{
		stor: s,
		data: make(map[string]interface{}),
	}
}

// Get returns the value stored under name and whether it was found.
func (cs *ConversationState) Get(name string) (interface{}, bool) {
	v, ok := cs.data[name]
	return v, ok
}

// Set stores a value under name and marks the state as modified.
func (cs *ConversationState) Set(name string, value interface{}) {
	cs.data[name] = value
	cs.modified = true
}

// Delete removes the value stored under name and marks the state as modified.
func (cs *ConversationState) Delete(name string) {
	if _, ok := cs.data[name]; ok {
		delete(cs.data, name)
		cs.modified = true
	}
}

// IsModified reports whether any values have changed since the last load or save.
func (cs *ConversationState) IsModified() bool {
	return cs.modified
}

// storageKey builds the storage key for the given channel / conversation identifiers.
func conversationStorageKey(channelID, conversationID string) string {
	return fmt.Sprintf("%s/conversations/%s", channelID, conversationID)
}

// Load reads the conversation state from storage using the provided identifiers.
// If the key is not found in storage the state is initialized to empty.
func (cs *ConversationState) Load(ctx context.Context, channelID, conversationID string) error {
	if cs.stor == nil {
		return nil
	}
	key := conversationStorageKey(channelID, conversationID)
	items, err := cs.stor.Read(ctx, []string{key})
	if err != nil {
		return fmt.Errorf("conversation state load: %w", err)
	}
	if raw, ok := items[key]; ok {
		// StoreItem is interface{}; marshal/unmarshal to get map[string]interface{}
		b, err := json.Marshal(raw)
		if err != nil {
			return fmt.Errorf("conversation state marshal: %w", err)
		}
		newData := make(map[string]interface{})
		if err := json.Unmarshal(b, &newData); err != nil {
			return fmt.Errorf("conversation state unmarshal: %w", err)
		}
		cs.data = newData
	} else {
		cs.data = make(map[string]interface{})
	}
	cs.modified = false
	return nil
}

// Save writes the conversation state to storage.
// If the state has not been modified since the last load, this is a no-op.
func (cs *ConversationState) Save(ctx context.Context, channelID, conversationID string) error {
	if cs.stor == nil || !cs.modified {
		return nil
	}
	key := conversationStorageKey(channelID, conversationID)
	if err := cs.stor.Write(ctx, map[string]storage.StoreItem{key: cs.data}); err != nil {
		return fmt.Errorf("conversation state save: %w", err)
	}
	cs.modified = false
	return nil
}

// Clear removes all values from the in-memory state and marks it as modified.
func (cs *ConversationState) Clear() {
	cs.data = make(map[string]interface{})
	cs.modified = true
}
