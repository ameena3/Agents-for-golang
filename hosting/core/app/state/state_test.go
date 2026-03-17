// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package state

import (
	"context"
	"testing"

	"github.com/ameena3/Agents-for-golang/hosting/core/storage"
)

// ---- ConversationState tests ----

func TestConversationStateGetSet(t *testing.T) {
	cs := NewConversationState(nil)
	cs.Set("key1", "value1")
	v, ok := cs.Get("key1")
	if !ok {
		t.Error("expected key1 to be found")
	}
	if v != "value1" {
		t.Errorf("expected %q, got %v", "value1", v)
	}
}

func TestConversationStateGetMissing(t *testing.T) {
	cs := NewConversationState(nil)
	_, ok := cs.Get("nonexistent")
	if ok {
		t.Error("expected nonexistent key to not be found")
	}
}

func TestConversationStateSetModified(t *testing.T) {
	cs := NewConversationState(nil)
	if cs.IsModified() {
		t.Error("expected IsModified to be false initially")
	}
	cs.Set("k", "v")
	if !cs.IsModified() {
		t.Error("expected IsModified to be true after Set")
	}
}

func TestConversationStateDelete(t *testing.T) {
	cs := NewConversationState(nil)
	cs.Set("k", "v")
	cs.Delete("k")
	_, ok := cs.Get("k")
	if ok {
		t.Error("expected key to be deleted")
	}
	if !cs.IsModified() {
		t.Error("expected IsModified to be true after Delete")
	}
}

func TestConversationStateClear(t *testing.T) {
	cs := NewConversationState(nil)
	cs.Set("a", 1)
	cs.Set("b", 2)
	cs.Clear()
	_, ok := cs.Get("a")
	if ok {
		t.Error("expected state to be cleared")
	}
	if !cs.IsModified() {
		t.Error("expected IsModified to be true after Clear")
	}
}

func TestConversationStateStorageKey(t *testing.T) {
	key := conversationStorageKey("msteams", "conv-123")
	want := "msteams/conversations/conv-123"
	if key != want {
		t.Errorf("conversationStorageKey: got %q, want %q", key, want)
	}
}

func TestConversationStateLoadSave(t *testing.T) {
	mem := storage.NewMemoryStorage()
	cs := NewConversationState(mem)

	cs.Set("counter", float64(42))
	if err := cs.Save(context.Background(), "ch", "conv"); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	if cs.IsModified() {
		t.Error("expected IsModified to be false after Save")
	}

	cs2 := NewConversationState(mem)
	if err := cs2.Load(context.Background(), "ch", "conv"); err != nil {
		t.Fatalf("Load error: %v", err)
	}
	v, ok := cs2.Get("counter")
	if !ok {
		t.Fatal("expected counter to be present after Load")
	}
	if v != float64(42) {
		t.Errorf("expected counter=42, got %v", v)
	}
}

func TestConversationStateSaveNoOpWhenNotModified(t *testing.T) {
	mem := storage.NewMemoryStorage()
	cs := NewConversationState(mem)
	// Never call Set, so IsModified should be false
	if err := cs.Save(context.Background(), "ch", "conv"); err != nil {
		t.Fatalf("Save should be no-op: %v", err)
	}
}

func TestConversationStateLoadMissingKey(t *testing.T) {
	mem := storage.NewMemoryStorage()
	cs := NewConversationState(mem)
	// Load from empty storage; should not error, should produce empty state
	if err := cs.Load(context.Background(), "ch", "new-conv"); err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cs.IsModified() {
		t.Error("expected IsModified to be false after loading missing key")
	}
}

// ---- UserState tests ----

func TestUserStateGetSet(t *testing.T) {
	us := NewUserState(nil)
	us.Set("name", "Alice")
	v, ok := us.Get("name")
	if !ok {
		t.Error("expected name to be found")
	}
	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
}

func TestUserStateStorageKey(t *testing.T) {
	key := userStorageKey("msteams", "user-456")
	want := "msteams/users/user-456"
	if key != want {
		t.Errorf("userStorageKey: got %q, want %q", key, want)
	}
}

func TestUserStateLoadSave(t *testing.T) {
	mem := storage.NewMemoryStorage()
	us := NewUserState(mem)
	us.Set("score", float64(100))
	if err := us.Save(context.Background(), "ch", "u1"); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	us2 := NewUserState(mem)
	if err := us2.Load(context.Background(), "ch", "u1"); err != nil {
		t.Fatalf("Load error: %v", err)
	}
	v, ok := us2.Get("score")
	if !ok {
		t.Fatal("expected score to be present")
	}
	if v != float64(100) {
		t.Errorf("expected score=100, got %v", v)
	}
}

func TestUserStateDelete(t *testing.T) {
	us := NewUserState(nil)
	us.Set("x", 1)
	us.Delete("x")
	_, ok := us.Get("x")
	if ok {
		t.Error("expected x to be deleted")
	}
}

func TestUserStateClear(t *testing.T) {
	us := NewUserState(nil)
	us.Set("a", 1)
	us.Clear()
	_, ok := us.Get("a")
	if ok {
		t.Error("expected state to be cleared")
	}
	if !us.IsModified() {
		t.Error("expected IsModified=true after Clear")
	}
}

// ---- TempState tests ----

func TestTempStateGetSet(t *testing.T) {
	ts := NewTempState()
	ts.Set("tmp", "data")
	v, ok := ts.Get("tmp")
	if !ok {
		t.Error("expected tmp to be found")
	}
	if v != "data" {
		t.Errorf("expected data, got %v", v)
	}
}

func TestTempStateDelete(t *testing.T) {
	ts := NewTempState()
	ts.Set("k", "v")
	ts.Delete("k")
	_, ok := ts.Get("k")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestTempStateKeys(t *testing.T) {
	ts := NewTempState()
	ts.Set("a", 1)
	ts.Set("b", 2)
	keys := ts.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestTempStateDeleteNoOp(t *testing.T) {
	ts := NewTempState()
	// Deleting a non-existent key should not panic
	ts.Delete("nonexistent")
}

// ---- TurnState tests ----

func TestNewTurnState(t *testing.T) {
	ts := NewTurnState[struct{}](nil)
	if ts == nil {
		t.Fatal("expected non-nil TurnState")
	}
	if ts.Conversation == nil {
		t.Error("expected Conversation to be non-nil")
	}
	if ts.User == nil {
		t.Error("expected User to be non-nil")
	}
	if ts.Temp == nil {
		t.Error("expected Temp to be non-nil")
	}
}

func TestTurnStateConversationHelpers(t *testing.T) {
	ts := NewTurnState[struct{}](nil)
	ts.SetConversationValue("foo", "bar")
	v, ok := ts.GetConversationValue("foo")
	if !ok {
		t.Error("expected foo to be found")
	}
	if v != "bar" {
		t.Errorf("expected bar, got %v", v)
	}
}

func TestTurnStateUserHelpers(t *testing.T) {
	ts := NewTurnState[struct{}](nil)
	ts.SetUserValue("age", 30)
	v, ok := ts.GetUserValue("age")
	if !ok {
		t.Error("expected age to be found")
	}
	if v != 30 {
		t.Errorf("expected 30, got %v", v)
	}
}

func TestTurnStateLoadSave(t *testing.T) {
	mem := storage.NewMemoryStorage()
	ts := NewTurnState[struct{}](mem)

	ts.SetConversationValue("msg", "hello")
	ts.SetUserValue("pref", "dark")

	if err := ts.Save(context.Background(), "ch", "conv1", "user1"); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	ts2 := NewTurnState[struct{}](mem)
	if err := ts2.Load(context.Background(), "ch", "conv1", "user1"); err != nil {
		t.Fatalf("Load error: %v", err)
	}

	v, ok := ts2.GetConversationValue("msg")
	if !ok || v != "hello" {
		t.Errorf("conversation value: got %v, want hello", v)
	}
	v2, ok2 := ts2.GetUserValue("pref")
	if !ok2 || v2 != "dark" {
		t.Errorf("user value: got %v, want dark", v2)
	}
}
