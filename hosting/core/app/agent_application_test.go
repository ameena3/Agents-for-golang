// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/microsoft/agents-sdk-go/activity"
	"github.com/microsoft/agents-sdk-go/hosting/core"
	"github.com/microsoft/agents-sdk-go/hosting/core/app"
	"github.com/microsoft/agents-sdk-go/hosting/core/app/routes"
	"github.com/microsoft/agents-sdk-go/hosting/core/storage"
)

// --- helpers ---

// newTC builds a minimal TurnContext for testing.
func newTC(actType, text, name string) *core.TurnContext {
	act := &activity.Activity{
		Type:         actType,
		Text:         text,
		Name:         name,
		ChannelID:    "test-channel",
		Conversation: &activity.ConversationAccount{ID: "conv-1"},
		From:         &activity.ChannelAccount{ID: "user-1"},
	}
	return core.NewTurnContext(nil, act, nil)
}

// --- route matching ---

func TestOnMessage_MatchesAllMessages(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})
	tc := newTC(activity.ActivityTypeMessage, "hello", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if !called {
		t.Fatal("message handler was not called")
	}
}

func TestOnMessage_PatternMatch(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnMessage(`^hello`, func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})

	// should match
	tc := newTC(activity.ActivityTypeMessage, "hello world", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if !called {
		t.Fatal("handler should have been called for matching text")
	}

	// should not match
	called = false
	tc2 := newTC(activity.ActivityTypeMessage, "goodbye", "")
	if err := application.OnTurn(context.Background(), tc2); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if called {
		t.Fatal("handler should NOT have been called for non-matching text")
	}
}

func TestOnMessage_NonMessageActivity_NotCalled(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})
	tc := newTC(activity.ActivityTypeEvent, "", "someEvent")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if called {
		t.Fatal("message handler should not be called for event activity")
	}
}

func TestOnActivity_MatchesType(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnActivity(activity.ActivityTypeEvent, func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})
	tc := newTC(activity.ActivityTypeEvent, "", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if !called {
		t.Fatal("activity handler was not called")
	}
}

func TestOnInvoke_MatchesByName(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnInvoke("card/action", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})
	tc := newTC(activity.ActivityTypeInvoke, "", "card/action")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if !called {
		t.Fatal("invoke handler was not called")
	}
}

func TestOnInvoke_DifferentNameNotCalled(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnInvoke("card/action", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})
	tc := newTC(activity.ActivityTypeInvoke, "", "other/action")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if called {
		t.Fatal("invoke handler should not be called for different name")
	}
}

func TestOnConversationUpdate_Called(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnConversationUpdate(func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})
	tc := newTC(activity.ActivityTypeConversationUpdate, "", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if !called {
		t.Fatal("conversation update handler was not called")
	}
}

func TestOnMembersAdded_Called(t *testing.T) {
	called := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnMembersAdded(func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		called = true
		return nil
	})

	act := &activity.Activity{
		Type:         activity.ActivityTypeConversationUpdate,
		ChannelID:    "test-channel",
		Conversation: &activity.ConversationAccount{ID: "conv-1"},
		From:         &activity.ChannelAccount{ID: "user-1"},
		MembersAdded: []*activity.ChannelAccount{{ID: "new-user"}},
	}
	tc := core.NewTurnContext(nil, act, nil)

	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if !called {
		t.Fatal("members added handler was not called")
	}
}

// --- route priority ---

func TestRoutePriority_InvokeBeforeMessage(t *testing.T) {
	order := []string{}
	application := app.New[struct{}](app.AppOptions[struct{}]{})

	// Register message handler first, invoke handler second.
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		order = append(order, "message")
		return nil
	})
	application.OnInvoke("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		order = append(order, "invoke")
		return nil
	})

	// invoke activity should match the invoke route only
	tc := newTC(activity.ActivityTypeInvoke, "", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn returned error: %v", err)
	}
	if len(order) != 1 || order[0] != "invoke" {
		t.Fatalf("expected [invoke], got %v", order)
	}
}

// --- error handling ---

func TestOnError_CalledOnHandlerError(t *testing.T) {
	handlerErr := errors.New("boom")
	var capturedErr error
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		return handlerErr
	})
	application.OnError(func(_ context.Context, _ *core.TurnContext, err error) error {
		capturedErr = err
		return nil // consumed
	})

	tc := newTC(activity.ActivityTypeMessage, "test", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn should not return error when error handler consumes it: %v", err)
	}
	if capturedErr == nil {
		t.Fatal("error handler was not called")
	}
	if !errors.Is(capturedErr, handlerErr) {
		t.Fatalf("error handler received unexpected error: %v", capturedErr)
	}
}

func TestOnError_PropagatesWhenNoHandler(t *testing.T) {
	handlerErr := errors.New("unhandled")
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		return handlerErr
	})

	tc := newTC(activity.ActivityTypeMessage, "test", "")
	err := application.OnTurn(context.Background(), tc)
	if err == nil {
		t.Fatal("expected error to propagate when no error handler registered")
	}
	if !errors.Is(err, handlerErr) {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- before/after hooks ---

func TestBeforeTurn_CalledBeforeHandler(t *testing.T) {
	order := []string{}
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.BeforeTurn(func(_ context.Context, _ *core.TurnContext) error {
		order = append(order, "before")
		return nil
	})
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		order = append(order, "handler")
		return nil
	})

	tc := newTC(activity.ActivityTypeMessage, "hi", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn error: %v", err)
	}
	if len(order) != 2 || order[0] != "before" || order[1] != "handler" {
		t.Fatalf("unexpected call order: %v", order)
	}
}

func TestAfterTurn_CalledAfterHandler(t *testing.T) {
	order := []string{}
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		order = append(order, "handler")
		return nil
	})
	application.AfterTurn(func(_ context.Context, _ *core.TurnContext) error {
		order = append(order, "after")
		return nil
	})

	tc := newTC(activity.ActivityTypeMessage, "hi", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn error: %v", err)
	}
	if len(order) != 2 || order[0] != "handler" || order[1] != "after" {
		t.Fatalf("unexpected call order: %v", order)
	}
}

func TestBeforeTurn_AbortsTurnOnError(t *testing.T) {
	beforeErr := errors.New("before error")
	handlerCalled := false
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	application.BeforeTurn(func(_ context.Context, _ *core.TurnContext) error {
		return beforeErr
	})
	application.OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error {
		handlerCalled = true
		return nil
	})

	tc := newTC(activity.ActivityTypeMessage, "hi", "")
	err := application.OnTurn(context.Background(), tc)
	if err == nil {
		t.Fatal("expected error from before hook to propagate")
	}
	if handlerCalled {
		t.Fatal("handler should not be called when before hook errors")
	}
}

// --- state management ---

func TestStateLoadSave_WithMemoryStorage(t *testing.T) {
	mem := storage.NewMemoryStorage()
	application := app.New[struct{}](app.AppOptions[struct{}]{Storage: mem})

	var savedKey string
	application.OnMessage("", func(_ context.Context, tc *core.TurnContext, _ struct{}) error {
		savedKey = "was-here"
		return nil
	})

	tc := newTC(activity.ActivityTypeMessage, "ping", "")
	if err := application.OnTurn(context.Background(), tc); err != nil {
		t.Fatalf("OnTurn error: %v", err)
	}
	// If we got here without error, Load/Save completed without panicking.
	if savedKey != "was-here" {
		t.Fatal("handler was not called")
	}
}

// --- method chaining ---

func TestMethodChaining(t *testing.T) {
	// Verify all builder methods return the same *AgentApplication.
	application := app.New[struct{}](app.AppOptions[struct{}]{})
	result := application.
		OnMessage("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error { return nil }).
		OnActivity("event", func(_ context.Context, _ *core.TurnContext, _ struct{}) error { return nil }).
		OnInvoke("", func(_ context.Context, _ *core.TurnContext, _ struct{}) error { return nil }).
		OnConversationUpdate(func(_ context.Context, _ *core.TurnContext, _ struct{}) error { return nil }).
		OnMembersAdded(func(_ context.Context, _ *core.TurnContext, _ struct{}) error { return nil }).
		BeforeTurn(func(_ context.Context, _ *core.TurnContext) error { return nil }).
		AfterTurn(func(_ context.Context, _ *core.TurnContext) error { return nil }).
		OnError(func(_ context.Context, _ *core.TurnContext, _ error) error { return nil })
	if result != application {
		t.Fatal("method chaining broke: returned different instance")
	}
}

// --- RouteList tests ---

func TestRouteList_FIFO_SameRank(t *testing.T) {
	// Two message routes, same rank. FIFO: first registered should match first.
	first := false
	second := false

	rl := routes.NewRouteList[struct{}]()
	rl.Add(routes.NewMessageRoute[struct{}]("", func(_ context.Context, _ interface{}, _ struct{}) error {
		first = true
		return nil
	}))
	rl.Add(routes.NewMessageRoute[struct{}]("", func(_ context.Context, _ interface{}, _ struct{}) error {
		second = true
		return nil
	}))

	act := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "hi"}
	r := rl.FindRoute(context.Background(), act)
	if r == nil {
		t.Fatal("FindRoute returned nil")
	}
	_ = r.Handler(context.Background(), nil, struct{}{})
	if !first {
		t.Fatal("first registered route should have been selected")
	}
	if second {
		t.Fatal("second route should not be called (FindRoute returns first match)")
	}
}

func TestRouteList_NoMatch_ReturnsNil(t *testing.T) {
	rl := routes.NewRouteList[struct{}]()
	rl.Add(routes.NewInvokeRoute[struct{}]("card/action", func(_ context.Context, _ interface{}, _ struct{}) error {
		return nil
	}))

	act := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "hello"}
	r := rl.FindRoute(context.Background(), act)
	if r != nil {
		t.Fatal("expected nil route for non-matching activity")
	}
}

// Ensure AgentApplication satisfies the core.Agent interface at compile time.
var _ core.Agent = (*app.AgentApplication[struct{}])(nil)
