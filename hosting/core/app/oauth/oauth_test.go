// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package oauth

import (
	"testing"
	"time"
)

func TestSignInState_IsExpired(t *testing.T) {
	past := &SignInState{Expiry: time.Now().Add(-1 * time.Minute)}
	if !past.IsExpired() {
		t.Error("expected expired")
	}

	future := &SignInState{Expiry: time.Now().Add(10 * time.Minute)}
	if future.IsExpired() {
		t.Error("expected not expired")
	}
}

func TestNewOAuthFlow(t *testing.T) {
	f := NewOAuthFlow("myconn", "Sign In", "Please sign in")
	if f.connectionName != "myconn" {
		t.Errorf("unexpected connectionName: %q", f.connectionName)
	}
	if f.title != "Sign In" {
		t.Errorf("unexpected title: %q", f.title)
	}
}

func TestSignInState_Fields(t *testing.T) {
	s := &SignInState{
		ConnectionName: "conn1",
		UserID:         "user1",
		ConversationID: "conv1",
		FlowType:       "user",
		Expiry:         time.Now().Add(5 * time.Minute),
	}
	if s.ConnectionName != "conn1" {
		t.Errorf("unexpected ConnectionName: %q", s.ConnectionName)
	}
	if s.FlowType != "user" {
		t.Errorf("unexpected FlowType: %q", s.FlowType)
	}
	if s.IsExpired() {
		t.Error("should not be expired")
	}
}
