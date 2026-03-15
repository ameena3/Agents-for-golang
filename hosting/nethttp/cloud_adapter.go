// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package nethttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/microsoft/agents-sdk-go/activity"
	"github.com/microsoft/agents-sdk-go/hosting/core"
	"github.com/microsoft/agents-sdk-go/hosting/core/authorization"
)

// CloudAdapter processes incoming HTTP requests and routes them to an agent.
// It handles JWT authentication, activity parsing, and response formatting.
type CloudAdapter struct {
	core.ChannelServiceAdapter
	allowUnauthenticated bool
}

// NewCloudAdapter creates a new CloudAdapter.
// For local testing without auth, set allowUnauthenticated = true.
func NewCloudAdapter(allowUnauthenticated bool) *CloudAdapter {
	return &CloudAdapter{
		ChannelServiceAdapter: *core.NewChannelServiceAdapter(),
		allowUnauthenticated:  allowUnauthenticated,
	}
}

// Process handles a single HTTP request: parses the activity, validates auth,
// creates a TurnContext, and runs the agent pipeline.
func (a *CloudAdapter) Process(ctx context.Context, r *http.Request, w http.ResponseWriter, agent core.Agent) error {
	ra := NewRequestAdapter(r)

	// Parse activity from request body.
	var act activity.Activity
	if err := ra.Body(&act); err != nil {
		http.Error(w, fmt.Sprintf("bad request: %v", err), http.StatusBadRequest)
		return nil
	}

	// Extract claims identity from the Authorization header.
	var identity *authorization.ClaimsIdentity
	var err error
	identity, err = ra.GetClaimsIdentity(ctx)
	if err != nil {
		if !a.allowUnauthenticated {
			http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
			return nil
		}
		// Fall back to anonymous identity when unauthenticated requests are permitted.
		identity = authorization.NewClaimsIdentity(false, "Anonymous", nil)
	}

	// Process the activity through the pipeline.
	var tc *core.TurnContext
	pipelineErr := a.ChannelServiceAdapter.ProcessActivity(ctx, identity, &act, func(pCtx context.Context, turnCtx *core.TurnContext) error {
		tc = turnCtx
		return agent.OnTurn(pCtx, turnCtx)
	})

	if pipelineErr != nil {
		http.Error(w, fmt.Sprintf("internal server error: %v", pipelineErr), http.StatusInternalServerError)
		return pipelineErr
	}

	// Determine the appropriate HTTP status code.
	statusCode := http.StatusAccepted
	if tc != nil {
		statusCode = a.ChannelServiceAdapter.GetHTTPStatusCode(tc)
	}

	// If status is 200 and there are buffered replies, write them as JSON.
	if statusCode == http.StatusOK && tc != nil && len(tc.BufferedReplies) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if encErr := json.NewEncoder(w).Encode(tc.BufferedReplies); encErr != nil {
			return fmt.Errorf("nethttp: failed to encode buffered replies: %w", encErr)
		}
		return nil
	}

	w.WriteHeader(statusCode)
	return nil
}

// SendActivities sends activities back to the channel service.
// This implementation returns an error indicating no connector is configured;
// concrete adapters that need to send proactive messages should override this.
func (a *CloudAdapter) SendActivities(ctx context.Context, tc *core.TurnContext, activities []*activity.Activity) ([]*activity.ResourceResponse, error) {
	return a.ChannelServiceAdapter.SendActivities(ctx, tc, activities)
}

// UpdateActivity updates an activity via the channel service.
func (a *CloudAdapter) UpdateActivity(ctx context.Context, tc *core.TurnContext, act *activity.Activity) (*activity.ResourceResponse, error) {
	return a.ChannelServiceAdapter.UpdateActivity(ctx, tc, act)
}

// DeleteActivity deletes an activity via the channel service.
func (a *CloudAdapter) DeleteActivity(ctx context.Context, tc *core.TurnContext, activityID string) error {
	return a.ChannelServiceAdapter.DeleteActivity(ctx, tc, activityID)
}

// ContinueConversation proactively continues a conversation.
func (a *CloudAdapter) ContinueConversation(ctx context.Context, ref *activity.ConversationReference, handler func(context.Context, *core.TurnContext) error) error {
	return a.ChannelServiceAdapter.ContinueConversation(ctx, ref, handler)
}
