// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package nethttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microsoft/agents-sdk-go/activity"
	"github.com/microsoft/agents-sdk-go/hosting/core"
)

// echoAgent is a simple agent that echoes back the received activity type via TurnContext.
type echoAgent struct {
	receivedActivity *activity.Activity
}

func (a *echoAgent) OnTurn(_ context.Context, tc *core.TurnContext) error {
	a.receivedActivity = tc.Activity()
	return nil
}

// TestMessageHandler_MethodNotAllowed verifies that non-POST requests return 405.
func TestMessageHandler_MethodNotAllowed(t *testing.T) {
	adapter := NewCloudAdapter(true)
	agent := &echoAgent{}
	handler := MessageHandler(adapter, agent)

	req := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestMessageHandler_ValidActivity verifies that a POST with valid JSON returns 202.
func TestMessageHandler_ValidActivity(t *testing.T) {
	adapter := NewCloudAdapter(true)
	agent := &echoAgent{}
	handler := MessageHandler(adapter, agent)

	act := activity.Activity{
		Type: activity.ActivityTypeMessage,
		Text: "Hello",
		Conversation: &activity.ConversationAccount{
			ID: "conv1",
		},
	}
	body, _ := json.Marshal(act)

	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	// Message activities return 202 Accepted.
	if w.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, w.Code)
	}
}

// TestCloudAdapter_ParsesActivity verifies that the activity is parsed and forwarded correctly.
func TestCloudAdapter_ParsesActivity(t *testing.T) {
	adapter := NewCloudAdapter(true)
	agent := &echoAgent{}

	act := activity.Activity{
		Type: activity.ActivityTypeMessage,
		Text: "ping",
		ID:   "act-001",
		Conversation: &activity.ConversationAccount{
			ID: "conv-abc",
		},
	}
	body, _ := json.Marshal(act)

	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := MessageHandler(adapter, agent)
	handler(w, req)

	if agent.receivedActivity == nil {
		t.Fatal("expected agent to receive an activity, got nil")
	}
	if agent.receivedActivity.Text != "ping" {
		t.Errorf("expected activity text %q, got %q", "ping", agent.receivedActivity.Text)
	}
	if agent.receivedActivity.Type != activity.ActivityTypeMessage {
		t.Errorf("expected activity type %q, got %q", activity.ActivityTypeMessage, agent.receivedActivity.Type)
	}
}

// TestMessageHandler_InvalidJSON verifies that malformed JSON returns 400.
func TestMessageHandler_InvalidJSON(t *testing.T) {
	adapter := NewCloudAdapter(true)
	agent := &echoAgent{}
	handler := MessageHandler(adapter, agent)

	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
