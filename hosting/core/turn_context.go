// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"context"
	"sync"

	"github.com/ameena3/Agents-for-golang/activity"
	"github.com/ameena3/Agents-for-golang/hosting/core/authorization"
)

// InvokeResponseKey is the TurnState key used to store an invoke response activity.
const InvokeResponseKey = "TurnContext.InvokeResponse"

// SendActivitiesHandler is a middleware hook called before activities are sent.
type SendActivitiesHandler func(ctx context.Context, activities []*activity.Activity, next func(context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error)) ([]*activity.ResourceResponse, error)

// UpdateActivityHandler is a middleware hook called before an activity is updated.
type UpdateActivityHandler func(ctx context.Context, act *activity.Activity, next func(context.Context, *activity.Activity) (*activity.ResourceResponse, error)) (*activity.ResourceResponse, error)

// DeleteActivityHandler is a middleware hook called before an activity is deleted.
type DeleteActivityHandler func(ctx context.Context, activityID string, next func(context.Context, string) error) error

// TurnContext holds all information for a single conversation turn.
// It is created by the adapter for each incoming activity and passed to the agent.
type TurnContext struct {
	mu              sync.Mutex
	incomingActivity *activity.Activity
	adapter         ChannelAdapter
	identity        *authorization.ClaimsIdentity
	turnState       map[string]interface{}
	sendHandlers    []SendActivitiesHandler
	updateHandlers  []UpdateActivityHandler
	deleteHandlers  []DeleteActivityHandler
	responded       bool
	BufferedReplies []*activity.Activity
}

// NewTurnContext creates a new TurnContext for the given adapter, activity, and identity.
func NewTurnContext(adapter ChannelAdapter, act *activity.Activity, identity *authorization.ClaimsIdentity) *TurnContext {
	return &TurnContext{
		incomingActivity: act,
		adapter:         adapter,
		identity:        identity,
		turnState:       make(map[string]interface{}),
		sendHandlers:    nil,
		updateHandlers:  nil,
		deleteHandlers:  nil,
		responded:       false,
		BufferedReplies: nil,
	}
}

// Activity returns the incoming activity for this turn.
func (tc *TurnContext) Activity() *activity.Activity {
	return tc.incomingActivity
}

// Adapter returns the channel adapter.
func (tc *TurnContext) Adapter() ChannelAdapter {
	return tc.adapter
}

// Identity returns the claims identity for the current turn.
func (tc *TurnContext) Identity() *authorization.ClaimsIdentity {
	return tc.identity
}

// TurnState returns the turn-scoped service cache (key/value bag).
func (tc *TurnContext) TurnState() map[string]interface{} {
	return tc.turnState
}

// Responded returns true if a non-trace reply has been sent during this turn.
func (tc *TurnContext) Responded() bool {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return tc.responded
}

// SendActivity sends a single activity to the channel.
func (tc *TurnContext) SendActivity(ctx context.Context, act *activity.Activity) (*activity.ResourceResponse, error) {
	responses, err := tc.SendActivities(ctx, []*activity.Activity{act})
	if err != nil {
		return nil, err
	}
	if len(responses) == 0 {
		return &activity.ResourceResponse{}, nil
	}
	return responses[0], nil
}

// SendActivities sends multiple activities to the channel, running them through
// all registered OnSendActivities hooks before calling the adapter.
func (tc *TurnContext) SendActivities(ctx context.Context, activities []*activity.Activity) ([]*activity.ResourceResponse, error) {
	// Prepare activities: set conversation reference and defaults.
	ref := tc.incomingActivity.GetConversationReference()
	sentNonTrace := false

	output := make([]*activity.Activity, 0, len(activities))
	for _, act := range activities {
		// Clone to avoid mutating the caller's activity.
		clone := *act
		clone.ApplyConversationReference(ref, false)
		clone.ID = "" // Clear the ID; the channel assigns it.
		if clone.Type == "" {
			clone.Type = activity.ActivityTypeMessage
		}
		if clone.Type != activity.ActivityTypeTrace {
			sentNonTrace = true
		}
		if clone.InputHint == "" {
			clone.InputHint = activity.InputHintAcceptingInput
		}
		output = append(output, &clone)
	}

	// Build the terminal function that actually sends.
	var finalSend func(context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error)
	finalSend = func(sendCtx context.Context, acts []*activity.Activity) ([]*activity.ResourceResponse, error) {
		// Handle buffered (expectReplies) mode.
		if tc.incomingActivity.DeliveryMode == activity.DeliveryModeExpectReplies {
			responses := make([]*activity.ResourceResponse, len(acts))
			for i, a := range acts {
				tc.mu.Lock()
				tc.BufferedReplies = append(tc.BufferedReplies, a)
				if a.Type == activity.ActivityTypeInvokeResponse {
					tc.turnState[InvokeResponseKey] = a
				}
				tc.mu.Unlock()
				responses[i] = &activity.ResourceResponse{}
			}
			return responses, nil
		}
		return tc.adapter.SendActivities(sendCtx, tc, acts)
	}

	// Chain through send handlers in registration order.
	handlers := tc.sendHandlers
	var chain func(int, context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error)
	chain = func(i int, chainCtx context.Context, acts []*activity.Activity) ([]*activity.ResourceResponse, error) {
		if i >= len(handlers) {
			return finalSend(chainCtx, acts)
		}
		return handlers[i](chainCtx, acts, func(nextCtx context.Context, nextActs []*activity.Activity) ([]*activity.ResourceResponse, error) {
			return chain(i+1, nextCtx, nextActs)
		})
	}

	responses, err := chain(0, ctx, output)
	if err != nil {
		return nil, err
	}

	if sentNonTrace {
		tc.mu.Lock()
		tc.responded = true
		tc.mu.Unlock()
	}

	return responses, nil
}

// UpdateActivity updates an existing activity in the conversation.
func (tc *TurnContext) UpdateActivity(ctx context.Context, act *activity.Activity) (*activity.ResourceResponse, error) {
	ref := tc.incomingActivity.GetConversationReference()
	// Apply conversation reference (not incoming) to the activity being updated.
	clone := *act
	clone.ApplyConversationReference(ref, false)

	// Build terminal function.
	finalUpdate := func(updateCtx context.Context, a *activity.Activity) (*activity.ResourceResponse, error) {
		return tc.adapter.UpdateActivity(updateCtx, tc, a)
	}

	handlers := tc.updateHandlers
	var chain func(int, context.Context, *activity.Activity) (*activity.ResourceResponse, error)
	chain = func(i int, chainCtx context.Context, a *activity.Activity) (*activity.ResourceResponse, error) {
		if i >= len(handlers) {
			return finalUpdate(chainCtx, a)
		}
		return handlers[i](chainCtx, a, func(nextCtx context.Context, nextAct *activity.Activity) (*activity.ResourceResponse, error) {
			return chain(i+1, nextCtx, nextAct)
		})
	}

	return chain(0, ctx, &clone)
}

// DeleteActivity deletes an activity by its activity ID.
func (tc *TurnContext) DeleteActivity(ctx context.Context, activityID string) error {
	finalDelete := func(deleteCtx context.Context, id string) error {
		return tc.adapter.DeleteActivity(deleteCtx, tc, id)
	}

	handlers := tc.deleteHandlers
	var chain func(int, context.Context, string) error
	chain = func(i int, chainCtx context.Context, id string) error {
		if i >= len(handlers) {
			return finalDelete(chainCtx, id)
		}
		return handlers[i](chainCtx, id, func(nextCtx context.Context, nextID string) error {
			return chain(i+1, nextCtx, nextID)
		})
	}

	return chain(0, ctx, activityID)
}

// OnSendActivities registers a handler to be called before activities are sent.
// Handlers are called in the order they are registered.
func (tc *TurnContext) OnSendActivities(handler SendActivitiesHandler) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.sendHandlers = append(tc.sendHandlers, handler)
}

// OnUpdateActivity registers a handler to be called before an activity is updated.
func (tc *TurnContext) OnUpdateActivity(handler UpdateActivityHandler) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.updateHandlers = append(tc.updateHandlers, handler)
}

// OnDeleteActivity registers a handler to be called before an activity is deleted.
func (tc *TurnContext) OnDeleteActivity(handler DeleteActivityHandler) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.deleteHandlers = append(tc.deleteHandlers, handler)
}

// GetConversationReference returns a ConversationReference that can be used for
// proactive messaging later.
func (tc *TurnContext) GetConversationReference() *activity.ConversationReference {
	return tc.incomingActivity.GetConversationReference()
}
