// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package state provides per-turn state management for AgentApplication.
// It defines three state scopes: ConversationState, UserState, and TempState.
package state

import "context"

// StatePropertyAccessor provides typed get/set/delete access to a single property
// within a state scope, providing typed get/set access.
type StatePropertyAccessor[T any] interface {
	// Get returns the current value. If the property has not been set it returns
	// the zero value for T and a nil error.
	Get(ctx context.Context) (T, error)
	// Set stores value under this accessor's property name.
	Set(ctx context.Context, value T) error
	// Delete removes the property from the state scope.
	Delete(ctx context.Context) error
}
