// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package state

import (
	"context"
	"fmt"

	"github.com/ameena3/Agents-for-golang/hosting/core/storage"
)

// TurnState manages the three standard state scopes for a single conversation turn:
//   - Conversation: persisted per conversation
//   - User: persisted per user across conversations
//   - Temp: ephemeral, only lives for the current turn
//
// The generic parameter AppStateT allows callers to embed application-specific
// state alongside the standard scopes.
type TurnState[AppStateT any] struct {
	// Conversation holds per-conversation persistent state.
	Conversation *ConversationState
	// User holds per-user persistent state.
	User *UserState
	// Temp holds ephemeral per-turn data (never persisted).
	Temp *TempState
	// App is application-specific state threaded through every handler.
	// It is not persisted by TurnState itself.
	App AppStateT

	stor storage.Storage
}

// NewTurnState creates a TurnState backed by the given storage.
// If s is nil, Conversation and User state are in-memory only.
func NewTurnState[AppStateT any](s storage.Storage) *TurnState[AppStateT] {
	return &TurnState[AppStateT]{
		stor:         s,
		Conversation: NewConversationState(s),
		User:         NewUserState(s),
		Temp:         NewTempState(),
	}
}

// Load reads both Conversation and User state from storage.
// channelID, conversationID, and userID are extracted from the incoming Activity.
func (ts *TurnState[AppStateT]) Load(ctx context.Context, channelID, conversationID, userID string) error {
	if err := ts.Conversation.Load(ctx, channelID, conversationID); err != nil {
		return fmt.Errorf("turn state load conversation: %w", err)
	}
	if err := ts.User.Load(ctx, channelID, userID); err != nil {
		return fmt.Errorf("turn state load user: %w", err)
	}
	return nil
}

// Save writes any modified state back to storage.
// It saves both Conversation and User state; Temp state is never persisted.
func (ts *TurnState[AppStateT]) Save(ctx context.Context, channelID, conversationID, userID string) error {
	if err := ts.Conversation.Save(ctx, channelID, conversationID); err != nil {
		return fmt.Errorf("turn state save conversation: %w", err)
	}
	if err := ts.User.Save(ctx, channelID, userID); err != nil {
		return fmt.Errorf("turn state save user: %w", err)
	}
	return nil
}

// GetConversationValue is a convenience helper that retrieves a value from
// the Conversation scope. Returns (nil, false) if the key is not present.
func (ts *TurnState[AppStateT]) GetConversationValue(name string) (interface{}, bool) {
	return ts.Conversation.Get(name)
}

// SetConversationValue is a convenience helper that stores a value in the
// Conversation scope.
func (ts *TurnState[AppStateT]) SetConversationValue(name string, value interface{}) {
	ts.Conversation.Set(name, value)
}

// GetUserValue is a convenience helper that retrieves a value from the User scope.
func (ts *TurnState[AppStateT]) GetUserValue(name string) (interface{}, bool) {
	return ts.User.Get(name)
}

// SetUserValue is a convenience helper that stores a value in the User scope.
func (ts *TurnState[AppStateT]) SetUserValue(name string, value interface{}) {
	ts.User.Set(name, value)
}
