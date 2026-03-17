// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"context"
	"net/http"

	"github.com/ameena3/Agents-for-golang/activity"
	"github.com/ameena3/Agents-for-golang/hosting/core/authorization"
)

// ChannelServiceAdapter extends ChannelAdapterBase with channel service communication.
// Concrete adapters (e.g., the aiohttp or FastAPI cloud adapters) embed this struct
// and supply the actual connector client via TurnState.
type ChannelServiceAdapter struct {
	ChannelAdapterBase
}

// NewChannelServiceAdapter creates a new ChannelServiceAdapter.
func NewChannelServiceAdapter() *ChannelServiceAdapter {
	return &ChannelServiceAdapter{
		ChannelAdapterBase: *NewChannelAdapterBase(),
	}
}

// ProcessActivity creates a TurnContext for the given identity and activity, then
// runs the middleware pipeline with the supplied handler.
func (a *ChannelServiceAdapter) ProcessActivity(
	ctx context.Context,
	identity *authorization.ClaimsIdentity,
	act *activity.Activity,
	handler func(context.Context, *TurnContext) error,
) error {
	if act == nil {
		return &adapterError{"ProcessActivity: activity cannot be nil"}
	}
	tc := NewTurnContext(a, act, identity)
	return a.RunPipeline(ctx, tc, handler)
}

// SendActivities sends activities to the channel service.
// Concrete adapters should override this method to integrate with a connector client.
// This base implementation returns an error indicating no connector is configured.
func (a *ChannelServiceAdapter) SendActivities(
	ctx context.Context,
	tc *TurnContext,
	activities []*activity.Activity,
) ([]*activity.ResourceResponse, error) {
	return nil, &adapterError{"SendActivities: no connector client configured — embed ChannelServiceAdapter and override SendActivities"}
}

// UpdateActivity updates an existing activity in the channel.
// Concrete adapters should override this to use a connector client.
func (a *ChannelServiceAdapter) UpdateActivity(
	ctx context.Context,
	tc *TurnContext,
	act *activity.Activity,
) (*activity.ResourceResponse, error) {
	return nil, &adapterError{"UpdateActivity: no connector client configured — embed ChannelServiceAdapter and override UpdateActivity"}
}

// DeleteActivity deletes an activity from the channel.
// Concrete adapters should override this to use a connector client.
func (a *ChannelServiceAdapter) DeleteActivity(
	ctx context.Context,
	tc *TurnContext,
	activityID string,
) error {
	return &adapterError{"DeleteActivity: no connector client configured — embed ChannelServiceAdapter and override DeleteActivity"}
}

// ContinueConversation proactively continues a conversation using a ConversationReference.
// It creates a new TurnContext from the reference and runs the middleware pipeline.
func (a *ChannelServiceAdapter) ContinueConversation(
	ctx context.Context,
	ref *activity.ConversationReference,
	handler func(context.Context, *TurnContext) error,
) error {
	if ref == nil {
		return &adapterError{"ContinueConversation: conversation reference cannot be nil"}
	}

	// Create a continuation activity from the reference.
	continuationActivity := &activity.Activity{
		Type:         activity.ActivityTypeEvent,
		Name:         "ContinueConversation",
		ChannelID:    ref.ChannelID,
		ServiceURL:   ref.ServiceURL,
		Conversation: ref.Conversation,
		Locale:       ref.Locale,
	}
	if ref.Agent != nil {
		from := *ref.Agent
		continuationActivity.From = &from
	}
	if ref.User != nil {
		recipient := *ref.User
		continuationActivity.Recipient = &recipient
	}
	if ref.ActivityID != "" {
		continuationActivity.ReplyToID = ref.ActivityID
	}

	tc := NewTurnContext(a, continuationActivity, nil)
	return a.RunPipeline(ctx, tc, handler)
}

// GetHTTPStatusCode returns the HTTP status code to send back for the turn.
// For invoke activities with a response stored in TurnState it returns 200;
// for all others it returns 202 (Accepted).
func (a *ChannelServiceAdapter) GetHTTPStatusCode(tc *TurnContext) int {
	if tc.Activity().Type == activity.ActivityTypeInvoke {
		if _, ok := tc.TurnState()[InvokeResponseKey]; ok {
			return http.StatusOK
		}
		return http.StatusNotImplemented
	}
	return http.StatusAccepted
}

// adapterError is a simple error type for adapter-level errors.
type adapterError struct {
	msg string
}

func (e *adapterError) Error() string {
	return e.msg
}
