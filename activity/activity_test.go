// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity_test

import (
	"encoding/json"
	"testing"

	"github.com/microsoft/agents-sdk-go/activity"
)

func TestNewMessageActivity(t *testing.T) {
	a := activity.NewMessageActivity("hello")
	if a.Type != activity.ActivityTypeMessage {
		t.Errorf("expected type %q, got %q", activity.ActivityTypeMessage, a.Type)
	}
	if a.Text != "hello" {
		t.Errorf("expected text %q, got %q", "hello", a.Text)
	}
}

func TestNewEventActivity(t *testing.T) {
	a := activity.NewEventActivity("myEvent", "myValue")
	if a.Type != activity.ActivityTypeEvent {
		t.Errorf("expected type %q, got %q", activity.ActivityTypeEvent, a.Type)
	}
	if a.Name != "myEvent" {
		t.Errorf("expected name %q, got %q", "myEvent", a.Name)
	}
	if a.Value != "myValue" {
		t.Errorf("expected value %q, got %v", "myValue", a.Value)
	}
}

func TestActivityIsType(t *testing.T) {
	a := &activity.Activity{Type: "message"}
	if !a.IsType(activity.ActivityTypeMessage) {
		t.Error("expected IsType(message) to return true")
	}
	if a.IsType(activity.ActivityTypeEvent) {
		t.Error("expected IsType(event) to return false")
	}
}

func TestActivityIsTypeSubtype(t *testing.T) {
	a := &activity.Activity{Type: "invoke/action"}
	if !a.IsType("invoke") {
		t.Error("expected IsType(invoke) to return true for invoke/action")
	}
}

func TestAsMessageActivity(t *testing.T) {
	msg := activity.NewMessageActivity("hello")
	if msg.AsMessageActivity() == nil {
		t.Error("expected AsMessageActivity() to return non-nil for message activity")
	}
	if msg.AsEventActivity() != nil {
		t.Error("expected AsEventActivity() to return nil for message activity")
	}
}

func TestActivityJSONMarshal(t *testing.T) {
	a := &activity.Activity{
		Type: activity.ActivityTypeMessage,
		ID:   "test-id",
		Text: "Hello World",
		From: &activity.ChannelAccount{
			ID:   "user1",
			Name: "User One",
		},
		Recipient: &activity.ChannelAccount{
			ID:   "bot1",
			Name: "Bot",
		},
		Conversation: &activity.ConversationAccount{
			ID: "conv1",
		},
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if result["type"] != activity.ActivityTypeMessage {
		t.Errorf("expected type %q, got %v", activity.ActivityTypeMessage, result["type"])
	}
	if result["id"] != "test-id" {
		t.Errorf("expected id %q, got %v", "test-id", result["id"])
	}
	if result["text"] != "Hello World" {
		t.Errorf("expected text %q, got %v", "Hello World", result["text"])
	}
}

func TestActivityJSONUnmarshal(t *testing.T) {
	raw := `{
		"type": "message",
		"id": "abc123",
		"text": "test message",
		"channelId": "msteams",
		"serviceUrl": "https://smba.trafficmanager.net/",
		"from": {"id": "user1", "name": "Alice", "role": "user"},
		"recipient": {"id": "bot1", "name": "MyBot", "role": "bot"},
		"conversation": {"id": "conv1", "isGroup": false}
	}`

	var a activity.Activity
	if err := json.Unmarshal([]byte(raw), &a); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if a.Type != "message" {
		t.Errorf("expected type %q, got %q", "message", a.Type)
	}
	if a.ID != "abc123" {
		t.Errorf("expected id %q, got %q", "abc123", a.ID)
	}
	if a.Text != "test message" {
		t.Errorf("expected text %q, got %q", "test message", a.Text)
	}
	if a.ChannelID != "msteams" {
		t.Errorf("expected channelId %q, got %q", "msteams", a.ChannelID)
	}
	if a.From == nil || a.From.ID != "user1" {
		t.Error("expected From to be set with id user1")
	}
	if a.Recipient == nil || a.Recipient.ID != "bot1" {
		t.Error("expected Recipient to be set with id bot1")
	}
	if a.Conversation == nil || a.Conversation.ID != "conv1" {
		t.Error("expected Conversation to be set with id conv1")
	}
}

func TestGetConversationReference(t *testing.T) {
	a := &activity.Activity{
		Type:       activity.ActivityTypeMessage,
		ID:         "act1",
		ChannelID:  "msteams",
		ServiceURL: "https://example.com",
		From:       &activity.ChannelAccount{ID: "user1", Name: "User"},
		Recipient:  &activity.ChannelAccount{ID: "bot1", Name: "Bot"},
		Conversation: &activity.ConversationAccount{
			ID: "conv1",
		},
	}

	ref := a.GetConversationReference()
	if ref.ActivityID != "act1" {
		t.Errorf("expected ActivityID %q, got %q", "act1", ref.ActivityID)
	}
	if ref.ChannelID != "msteams" {
		t.Errorf("expected ChannelID %q, got %q", "msteams", ref.ChannelID)
	}
	if ref.ServiceURL != "https://example.com" {
		t.Errorf("expected ServiceURL %q, got %q", "https://example.com", ref.ServiceURL)
	}
	if ref.User == nil || ref.User.ID != "user1" {
		t.Error("expected User to be set with id user1")
	}
	if ref.Agent == nil || ref.Agent.ID != "bot1" {
		t.Error("expected Agent to be set with id bot1")
	}
}

func TestCreateReply(t *testing.T) {
	a := &activity.Activity{
		Type:       activity.ActivityTypeMessage,
		ID:         "original-id",
		ChannelID:  "msteams",
		ServiceURL: "https://example.com",
		Locale:     "en-US",
		From:       &activity.ChannelAccount{ID: "user1", Name: "User"},
		Recipient:  &activity.ChannelAccount{ID: "bot1", Name: "Bot"},
		Conversation: &activity.ConversationAccount{
			ID: "conv1",
		},
	}

	reply := a.CreateReply("reply text")
	if reply.Type != activity.ActivityTypeMessage {
		t.Errorf("expected reply type message, got %q", reply.Type)
	}
	if reply.Text != "reply text" {
		t.Errorf("expected reply text %q, got %q", "reply text", reply.Text)
	}
	if reply.ReplyToID != "original-id" {
		t.Errorf("expected ReplyToID %q, got %q", "original-id", reply.ReplyToID)
	}
	if reply.From == nil || reply.From.ID != "bot1" {
		t.Error("expected reply From to be the original recipient (bot)")
	}
	if reply.Recipient == nil || reply.Recipient.ID != "user1" {
		t.Error("expected reply Recipient to be the original sender (user)")
	}
}

func TestHasContent(t *testing.T) {
	empty := &activity.Activity{Type: activity.ActivityTypeMessage}
	if empty.HasContent() {
		t.Error("expected HasContent() to return false for empty activity")
	}

	withText := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "hello"}
	if !withText.HasContent() {
		t.Error("expected HasContent() to return true for activity with text")
	}

	withAttachment := &activity.Activity{
		Type:        activity.ActivityTypeMessage,
		Attachments: []*activity.Attachment{{ContentType: "text/plain"}},
	}
	if !withAttachment.HasContent() {
		t.Error("expected HasContent() to return true for activity with attachments")
	}
}

func TestIsFromStreamingConnection(t *testing.T) {
	http := &activity.Activity{ServiceURL: "https://example.com"}
	if http.IsFromStreamingConnection() {
		t.Error("expected IsFromStreamingConnection() to return false for https URL")
	}

	streaming := &activity.Activity{ServiceURL: "wss://example.com"}
	if !streaming.IsFromStreamingConnection() {
		t.Error("expected IsFromStreamingConnection() to return true for wss URL")
	}

	empty := &activity.Activity{}
	if empty.IsFromStreamingConnection() {
		t.Error("expected IsFromStreamingConnection() to return false for empty ServiceURL")
	}
}

func TestActivityOmitEmpty(t *testing.T) {
	a := &activity.Activity{Type: "message"}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// Only "type" should be set; other fields should be omitted
	if len(result) != 1 {
		t.Errorf("expected only 1 field in JSON, got %d: %v", len(result), result)
	}
}

func TestIsAgenticRequest(t *testing.T) {
	nonAgentic := &activity.Activity{
		Recipient: &activity.ChannelAccount{Role: activity.RoleTypeAgent},
	}
	if nonAgentic.IsAgenticRequest() {
		t.Error("expected IsAgenticRequest() false for bot role")
	}

	agentic := &activity.Activity{
		Recipient: &activity.ChannelAccount{Role: activity.RoleTypeAgenticIdentity},
	}
	if !agentic.IsAgenticRequest() {
		t.Error("expected IsAgenticRequest() true for agenticAppInstance role")
	}
}
