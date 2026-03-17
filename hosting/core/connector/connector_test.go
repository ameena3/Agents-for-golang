// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package connector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ameena3/Agents-for-golang/activity"
)

type mockTokenProvider struct {
	token string
}

func (m *mockTokenProvider) GetAccessToken(_ context.Context, _ string) (string, error) {
	return m.token, nil
}

func TestNewConnectorClient(t *testing.T) {
	c := NewConnectorClient("https://example.com", nil)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.serviceURL != "https://example.com" {
		t.Errorf("unexpected serviceURL: %q", c.serviceURL)
	}
}

func TestConnectorClient_SendToConversation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": "resp123"})
	}))
	defer srv.Close()

	c := NewConnectorClient(srv.URL, &mockTokenProvider{token: "tok"})
	act := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "hello"}
	resp, err := c.SendToConversation(context.Background(), "conv1", act)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "resp123" {
		t.Errorf("expected resp123, got %q", resp.ID)
	}
}

func TestConnectorClient_ReplyToActivity(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": "reply1"})
	}))
	defer srv.Close()

	c := NewConnectorClient(srv.URL, nil)
	act := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "reply"}
	resp, err := c.ReplyToActivity(context.Background(), "conv1", "act1", act)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "reply1" {
		t.Errorf("expected reply1, got %q", resp.ID)
	}
}

func TestConnectorClient_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewConnectorClient(srv.URL, nil)
	_, err := c.SendToConversation(context.Background(), "conv1", &activity.Activity{})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestNewUserTokenClient(t *testing.T) {
	c := NewUserTokenClient("https://token.example.com")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
