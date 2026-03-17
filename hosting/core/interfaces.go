// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"context"

	"github.com/ameena3/Agents-for-golang/activity"
)

// Agent is the primary interface that all agents must implement.
type Agent interface {
	// OnTurn is called for every incoming activity.
	OnTurn(ctx context.Context, turnCtx *TurnContext) error
}

// Middleware is an interface for processing activities in a pipeline.
type Middleware interface {
	// OnTurn processes the activity and calls next to continue the pipeline.
	OnTurn(ctx context.Context, turnCtx *TurnContext, next func(context.Context) error) error
}

// ChannelAdapter is the interface for adapters that communicate with channels.
type ChannelAdapter interface {
	// SendActivities sends a slice of activities to the channel.
	SendActivities(ctx context.Context, turnCtx *TurnContext, activities []*activity.Activity) ([]*activity.ResourceResponse, error)
	// UpdateActivity updates an existing activity.
	UpdateActivity(ctx context.Context, turnCtx *TurnContext, act *activity.Activity) (*activity.ResourceResponse, error)
	// DeleteActivity deletes an activity by its conversation reference.
	DeleteActivity(ctx context.Context, turnCtx *TurnContext, activityID string) error
	// ContinueConversation continues a conversation proactively.
	ContinueConversation(ctx context.Context, ref *activity.ConversationReference, handler func(context.Context, *TurnContext) error) error
	// Use registers one or more middleware handlers.
	Use(middleware ...Middleware)
}

// AccessTokenProvider can supply OAuth access tokens.
type AccessTokenProvider interface {
	// GetAccessToken returns an access token for the given resource URL and scopes.
	GetAccessToken(ctx context.Context, resourceURL string, scopes []string, forceRefresh bool) (string, error)
	// AcquireTokenOnBehalfOf acquires a token on behalf of a user.
	AcquireTokenOnBehalfOf(ctx context.Context, scopes []string, userAssertion string) (string, error)
}

// TurnContext holds the state for a single turn of conversation.
// The full implementation lives in turn_context.go.
// TurnContext is declared in turn_context.go.
