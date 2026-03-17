// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"context"
	"errors"
	"testing"

	"github.com/ameena3/Agents-for-golang/activity"
	"github.com/ameena3/Agents-for-golang/hosting/core/authorization"
)

// --- stub adapter for testing ---

type stubAdapter struct {
	sendCalled   int
	updateCalled int
	deleteCalled int
	sendErr      error
	updateErr    error
	deleteErr    error
	sendRet      []*activity.ResourceResponse
}

func (s *stubAdapter) SendActivities(_ context.Context, _ *TurnContext, acts []*activity.Activity) ([]*activity.ResourceResponse, error) {
	s.sendCalled++
	if s.sendErr != nil {
		return nil, s.sendErr
	}
	if s.sendRet != nil {
		return s.sendRet, nil
	}
	resp := make([]*activity.ResourceResponse, len(acts))
	for i := range acts {
		resp[i] = &activity.ResourceResponse{ID: "resp-" + acts[i].Type}
	}
	return resp, nil
}

func (s *stubAdapter) UpdateActivity(_ context.Context, _ *TurnContext, act *activity.Activity) (*activity.ResourceResponse, error) {
	s.updateCalled++
	if s.updateErr != nil {
		return nil, s.updateErr
	}
	return &activity.ResourceResponse{ID: act.ID}, nil
}

func (s *stubAdapter) DeleteActivity(_ context.Context, _ *TurnContext, _ string) error {
	s.deleteCalled++
	return s.deleteErr
}

func (s *stubAdapter) ContinueConversation(_ context.Context, _ *activity.ConversationReference, _ func(context.Context, *TurnContext) error) error {
	return nil
}

func (s *stubAdapter) Use(_ ...Middleware) {}

// --- helpers ---

func newTestTurnContext(adapter ChannelAdapter) *TurnContext {
	act := &activity.Activity{
		Type:      activity.ActivityTypeMessage,
		ID:        "test-id",
		ChannelID: "test-channel",
		Conversation: &activity.ConversationAccount{
			ID: "conv-1",
		},
		From: &activity.ChannelAccount{
			ID:   "user-1",
			Name: "User",
		},
		Recipient: &activity.ChannelAccount{
			ID:   "agent-1",
			Name: "Agent",
		},
		ServiceURL: "https://example.com",
	}
	identity := authorization.NewClaimsIdentity(true, "Bearer", map[string]string{
		"aud": "app-id",
	})
	return NewTurnContext(adapter, act, identity)
}

// --- tests ---

func TestNewTurnContext(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	if tc == nil {
		t.Fatal("expected non-nil TurnContext")
	}
	if tc.Activity() == nil {
		t.Error("Activity() should not be nil")
	}
	if tc.Adapter() == nil {
		t.Error("Adapter() should not be nil")
	}
	if tc.Identity() == nil {
		t.Error("Identity() should not be nil")
	}
	if tc.TurnState() == nil {
		t.Error("TurnState() should not be nil")
	}
	if tc.Responded() {
		t.Error("Responded should be false initially")
	}
}

func TestSendActivity_CallsAdapter(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	resp, err := tc.SendActivity(context.Background(), &activity.Activity{
		Type: activity.ActivityTypeMessage,
		Text: "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if stub.sendCalled != 1 {
		t.Errorf("expected adapter.SendActivities to be called once, got %d", stub.sendCalled)
	}
	if !tc.Responded() {
		t.Error("expected Responded to be true after sending a message")
	}
}

func TestSendActivities_TraceDoesNotSetResponded(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	_, err := tc.SendActivities(context.Background(), []*activity.Activity{
		{Type: activity.ActivityTypeTrace, Name: "debug"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tc.Responded() {
		t.Error("expected Responded to be false after sending only a trace activity")
	}
}

func TestSendActivities_AdapterError(t *testing.T) {
	stub := &stubAdapter{sendErr: errors.New("send failed")}
	tc := newTestTurnContext(stub)

	_, err := tc.SendActivities(context.Background(), []*activity.Activity{
		{Type: activity.ActivityTypeMessage, Text: "hi"},
	})
	if err == nil {
		t.Fatal("expected error from adapter but got nil")
	}
}

func TestOnSendActivities_HookCalledBeforeAdapter(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	var hookCalled bool
	tc.OnSendActivities(func(ctx context.Context, acts []*activity.Activity, next func(context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error)) ([]*activity.ResourceResponse, error) {
		hookCalled = true
		return next(ctx, acts)
	})

	_, err := tc.SendActivity(context.Background(), &activity.Activity{
		Type: activity.ActivityTypeMessage,
		Text: "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hookCalled {
		t.Error("expected OnSendActivities hook to be called")
	}
	if stub.sendCalled != 1 {
		t.Errorf("expected adapter.SendActivities to be called once, got %d", stub.sendCalled)
	}
}

func TestOnSendActivities_HookCanInterceptAndSkipAdapter(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	fakeResponse := []*activity.ResourceResponse{{ID: "intercepted"}}
	tc.OnSendActivities(func(ctx context.Context, acts []*activity.Activity, next func(context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error)) ([]*activity.ResourceResponse, error) {
		// Do not call next — intercept completely.
		return fakeResponse, nil
	})

	resp, err := tc.SendActivity(context.Background(), &activity.Activity{
		Type: activity.ActivityTypeMessage,
		Text: "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.ID != "intercepted" {
		t.Errorf("expected intercepted response, got %v", resp)
	}
	if stub.sendCalled != 0 {
		t.Error("adapter.SendActivities should not have been called when hook intercepts")
	}
}

func TestOnSendActivities_MultipleHooksCalledInOrder(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	var order []int
	tc.OnSendActivities(func(ctx context.Context, acts []*activity.Activity, next func(context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error)) ([]*activity.ResourceResponse, error) {
		order = append(order, 1)
		return next(ctx, acts)
	})
	tc.OnSendActivities(func(ctx context.Context, acts []*activity.Activity, next func(context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error)) ([]*activity.ResourceResponse, error) {
		order = append(order, 2)
		return next(ctx, acts)
	})

	_, err := tc.SendActivity(context.Background(), &activity.Activity{Type: activity.ActivityTypeMessage})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Errorf("expected hooks called in order [1,2], got %v", order)
	}
}

func TestGetConversationReference(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	ref := tc.GetConversationReference()
	if ref == nil {
		t.Fatal("expected non-nil ConversationReference")
	}
	if ref.ChannelID != tc.Activity().ChannelID {
		t.Errorf("expected ChannelID %q, got %q", tc.Activity().ChannelID, ref.ChannelID)
	}
	if ref.ServiceURL != tc.Activity().ServiceURL {
		t.Errorf("expected ServiceURL %q, got %q", tc.Activity().ServiceURL, ref.ServiceURL)
	}
	if ref.Conversation == nil || ref.Conversation.ID != "conv-1" {
		t.Error("expected Conversation to be set on ConversationReference")
	}
}

func TestTurnState_StoreAndRetrieve(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	tc.TurnState()["myKey"] = "myValue"
	val, ok := tc.TurnState()["myKey"]
	if !ok {
		t.Fatal("expected key to be present in TurnState")
	}
	if val != "myValue" {
		t.Errorf("expected 'myValue', got %v", val)
	}
}

func TestDeleteActivity_CallsAdapter(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	err := tc.DeleteActivity(context.Background(), "activity-to-delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.deleteCalled != 1 {
		t.Errorf("expected adapter.DeleteActivity to be called once, got %d", stub.deleteCalled)
	}
}

func TestUpdateActivity_CallsAdapter(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	resp, err := tc.UpdateActivity(context.Background(), &activity.Activity{
		ID:   "activity-to-update",
		Text: "updated text",
		Type: activity.ActivityTypeMessage,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response from UpdateActivity")
	}
	if stub.updateCalled != 1 {
		t.Errorf("expected adapter.UpdateActivity to be called once, got %d", stub.updateCalled)
	}
}

func TestSendActivities_BufferedReplyMode(t *testing.T) {
	stub := &stubAdapter{}
	tc := newTestTurnContext(stub)

	// Switch to expectReplies delivery mode.
	tc.Activity().DeliveryMode = activity.DeliveryModeExpectReplies

	_, err := tc.SendActivities(context.Background(), []*activity.Activity{
		{Type: activity.ActivityTypeMessage, Text: "buffered"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Adapter should NOT be called in buffered mode.
	if stub.sendCalled != 0 {
		t.Errorf("expected adapter.SendActivities NOT to be called in buffered mode, got %d calls", stub.sendCalled)
	}
	if len(tc.BufferedReplies) != 1 {
		t.Errorf("expected 1 buffered reply, got %d", len(tc.BufferedReplies))
	}
}
