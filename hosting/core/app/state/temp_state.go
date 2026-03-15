// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package state

// TempState holds ephemeral per-turn data.
// It is never persisted to storage and is discarded when the turn ends.
// TempState holds ephemeral data that exists only for the current turn.
type TempState struct {
	data map[string]interface{}
}

// NewTempState creates an empty TempState.
func NewTempState() *TempState {
	return &TempState{data: make(map[string]interface{})}
}

// Get returns the value stored under name and whether it was found.
func (ts *TempState) Get(name string) (interface{}, bool) {
	v, ok := ts.data[name]
	return v, ok
}

// Set stores a value under name.
func (ts *TempState) Set(name string, value interface{}) {
	ts.data[name] = value
}

// Delete removes the value stored under name. If name is not present this is a no-op.
func (ts *TempState) Delete(name string) {
	delete(ts.data, name)
}

// Keys returns a snapshot of all currently set property names.
func (ts *TempState) Keys() []string {
	keys := make([]string, 0, len(ts.data))
	for k := range ts.data {
		keys = append(keys, k)
	}
	return keys
}
