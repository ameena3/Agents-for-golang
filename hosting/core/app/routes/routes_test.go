// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package routes

import (
	"context"
	"testing"

	"github.com/ameena3/Agents-for-golang/activity"
)

type testState struct{}

func nopHandler(_ context.Context, _ interface{}, _ testState) error { return nil }

func TestNewRouteListEmpty(t *testing.T) {
	rl := NewRouteList[testState]()
	if rl == nil {
		t.Fatal("expected non-nil RouteList")
	}
	if rl.Len() != 0 {
		t.Errorf("expected len 0, got %d", rl.Len())
	}
}

func TestRouteListAdd(t *testing.T) {
	rl := NewRouteList[testState]()
	r := NewMessageRoute[testState]("", nopHandler)
	rl.Add(r)
	if rl.Len() != 1 {
		t.Errorf("expected len 1, got %d", rl.Len())
	}
}

func TestNewMessageRouteMatchesAllMessages(t *testing.T) {
	r := NewMessageRoute[testState]("", nopHandler)

	act := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "hello world"}
	if !r.Selector(context.Background(), act) {
		t.Error("empty-pattern message route should match any message activity")
	}

	nonMsg := &activity.Activity{Type: activity.ActivityTypeInvoke}
	if r.Selector(context.Background(), nonMsg) {
		t.Error("message route should not match non-message activity")
	}
}

func TestNewMessageRouteWithPattern(t *testing.T) {
	r := NewMessageRoute[testState](`^hello`, nopHandler)

	matching := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "hello world"}
	if !r.Selector(context.Background(), matching) {
		t.Error("expected route to match 'hello world'")
	}

	nonMatching := &activity.Activity{Type: activity.ActivityTypeMessage, Text: "goodbye"}
	if r.Selector(context.Background(), nonMatching) {
		t.Error("expected route not to match 'goodbye'")
	}
}

func TestNewActivityTypeRoute(t *testing.T) {
	r := NewActivityTypeRoute[testState]("event", nopHandler)

	match := &activity.Activity{Type: "event"}
	if !r.Selector(context.Background(), match) {
		t.Error("expected route to match 'event' activity")
	}

	noMatch := &activity.Activity{Type: activity.ActivityTypeMessage}
	if r.Selector(context.Background(), noMatch) {
		t.Error("expected route not to match message activity")
	}
}

func TestNewActivityTypeRouteCaseInsensitive(t *testing.T) {
	r := NewActivityTypeRoute[testState]("Event", nopHandler)
	act := &activity.Activity{Type: "event"}
	if !r.Selector(context.Background(), act) {
		t.Error("activity type matching should be case-insensitive")
	}
}

func TestNewInvokeRouteMatchesAllInvokes(t *testing.T) {
	r := NewInvokeRoute[testState]("", nopHandler)

	match := &activity.Activity{Type: activity.ActivityTypeInvoke}
	if !r.Selector(context.Background(), match) {
		t.Error("empty-name invoke route should match any invoke activity")
	}

	noMatch := &activity.Activity{Type: activity.ActivityTypeMessage}
	if r.Selector(context.Background(), noMatch) {
		t.Error("invoke route should not match message activity")
	}
}

func TestNewInvokeRouteWithName(t *testing.T) {
	r := NewInvokeRoute[testState]("adaptiveCard/action", nopHandler)

	match := &activity.Activity{Type: activity.ActivityTypeInvoke, Name: "adaptiveCard/action"}
	if !r.Selector(context.Background(), match) {
		t.Error("expected route to match named invoke")
	}

	wrongName := &activity.Activity{Type: activity.ActivityTypeInvoke, Name: "task/fetch"}
	if r.Selector(context.Background(), wrongName) {
		t.Error("expected route not to match different invoke name")
	}
}

func TestNewConversationUpdateRoute(t *testing.T) {
	r := NewConversationUpdateRoute[testState](nopHandler)

	match := &activity.Activity{Type: activity.ActivityTypeConversationUpdate}
	if !r.Selector(context.Background(), match) {
		t.Error("expected route to match conversationUpdate")
	}

	noMatch := &activity.Activity{Type: activity.ActivityTypeMessage}
	if r.Selector(context.Background(), noMatch) {
		t.Error("expected route not to match message")
	}
}

func TestNewMembersAddedRoute(t *testing.T) {
	r := NewMembersAddedRoute[testState](nopHandler)

	match := &activity.Activity{
		Type:         activity.ActivityTypeConversationUpdate,
		MembersAdded: []*activity.ChannelAccount{{ID: "user1"}},
	}
	if !r.Selector(context.Background(), match) {
		t.Error("expected route to match when MembersAdded is non-empty")
	}

	noMembers := &activity.Activity{
		Type:         activity.ActivityTypeConversationUpdate,
		MembersAdded: nil,
	}
	if r.Selector(context.Background(), noMembers) {
		t.Error("expected route not to match when MembersAdded is empty")
	}
}

func TestRouteListFindRouteInvokePriorityOverMessage(t *testing.T) {
	rl := NewRouteList[testState]()

	var matched string
	msgRoute := NewMessageRoute[testState]("", func(_ context.Context, _ interface{}, _ testState) error {
		matched = "message"
		return nil
	})
	invokeRoute := NewInvokeRoute[testState]("", func(_ context.Context, _ interface{}, _ testState) error {
		matched = "invoke"
		return nil
	})

	// Add message first, then invoke — invoke should win for invoke activities
	rl.Add(msgRoute)
	rl.Add(invokeRoute)

	act := &activity.Activity{Type: activity.ActivityTypeInvoke}
	found := rl.FindRoute(context.Background(), act)
	if found == nil {
		t.Fatal("expected a route to be found")
	}
	// Execute it to record which handler was selected
	found.Handler(context.Background(), nil, testState{})
	if matched != "invoke" {
		t.Errorf("expected invoke handler to be selected, got %q", matched)
	}
}

func TestRouteListFindRouteNoMatch(t *testing.T) {
	rl := NewRouteList[testState]()
	rl.Add(NewMessageRoute[testState]("", nopHandler))

	act := &activity.Activity{Type: activity.ActivityTypeInvoke}
	found := rl.FindRoute(context.Background(), act)
	if found != nil {
		t.Error("expected nil when no route matches")
	}
}

func TestRouteListFindRouteFirstMatchWins(t *testing.T) {
	rl := NewRouteList[testState]()

	var order []string
	r1 := NewMessageRoute[testState]("", func(_ context.Context, _ interface{}, _ testState) error {
		order = append(order, "first")
		return nil
	})
	r2 := NewMessageRoute[testState]("", func(_ context.Context, _ interface{}, _ testState) error {
		order = append(order, "second")
		return nil
	})
	rl.Add(r1)
	rl.Add(r2)

	act := &activity.Activity{Type: activity.ActivityTypeMessage}
	found := rl.FindRoute(context.Background(), act)
	if found == nil {
		t.Fatal("expected a route to be found")
	}
	found.Handler(context.Background(), nil, testState{})
	if len(order) != 1 || order[0] != "first" {
		t.Errorf("expected first route to be selected, got %v", order)
	}
}

func TestRouteRankConstants(t *testing.T) {
	if RouteRankFirst != 0 {
		t.Errorf("RouteRankFirst: got %d, want 0", RouteRankFirst)
	}
	if RouteRankDefault != 32767 {
		t.Errorf("RouteRankDefault: got %d, want 32767", RouteRankDefault)
	}
	if RouteRankLast != 65535 {
		t.Errorf("RouteRankLast: got %d, want 65535", RouteRankLast)
	}
}
