// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"context"
	"fmt"

	"github.com/microsoft/agents-sdk-go/activity"
)

// ActivityHandler is a base struct that routes activities to typed handler methods.
// Embed this in your agent struct and override the On* methods you need.
//
// Example:
//
//	type MyAgent struct {
//	    core.ActivityHandler
//	}
//
//	func (a *MyAgent) OnMessageActivity(ctx context.Context, tc *TurnContext) error {
//	    _, err := tc.SendActivity(ctx, activity.NewMessageActivity("Hello!"))
//	    return err
//	}
type ActivityHandler struct{}

// OnTurn dispatches the activity to the appropriate typed handler method based
// on the activity type. This implements the Agent interface when embedded.
func (h *ActivityHandler) OnTurn(ctx context.Context, tc *TurnContext) error {
	if tc == nil {
		return fmt.Errorf("ActivityHandler.OnTurn: turn context cannot be nil")
	}
	act := tc.Activity()
	if act == nil {
		return fmt.Errorf("ActivityHandler.OnTurn: turn context must have a non-nil activity")
	}
	if act.Type == "" {
		return fmt.Errorf("ActivityHandler.OnTurn: activity must have a non-empty type")
	}

	switch act.Type {
	case activity.ActivityTypeMessage:
		return h.OnMessageActivity(ctx, tc)
	case activity.ActivityTypeMessageUpdate:
		return h.OnMessageUpdateActivity(ctx, tc)
	case activity.ActivityTypeMessageDelete:
		return h.OnMessageDeleteActivity(ctx, tc)
	case activity.ActivityTypeConversationUpdate:
		return h.OnConversationUpdateActivity(ctx, tc)
	case activity.ActivityTypeMessageReaction:
		return h.OnMessageReactionActivity(ctx, tc)
	case activity.ActivityTypeEvent:
		return h.OnEventActivity(ctx, tc)
	case activity.ActivityTypeInvoke:
		return h.OnInvokeActivity(ctx, tc)
	case activity.ActivityTypeEndOfConversation:
		return h.OnEndOfConversationActivity(ctx, tc)
	case activity.ActivityTypeTyping:
		return h.OnTypingActivity(ctx, tc)
	case activity.ActivityTypeInstallationUpdate:
		return h.OnInstallationUpdateActivity(ctx, tc)
	default:
		return h.OnUnrecognizedActivityType(ctx, tc)
	}
}

// OnMessageActivity handles message activities.
// Override in a derived struct to provide agent response logic.
func (h *ActivityHandler) OnMessageActivity(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnMessageUpdateActivity handles messageUpdate activities.
func (h *ActivityHandler) OnMessageUpdateActivity(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnMessageDeleteActivity handles messageDelete activities.
func (h *ActivityHandler) OnMessageDeleteActivity(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnConversationUpdateActivity handles conversationUpdate activities.
// By default it calls OnMembersAdded or OnMembersRemoved as appropriate.
func (h *ActivityHandler) OnConversationUpdateActivity(ctx context.Context, tc *TurnContext) error {
	act := tc.Activity()
	if len(act.MembersAdded) > 0 {
		if err := h.OnMembersAdded(ctx, act.MembersAdded, tc); err != nil {
			return err
		}
	}
	if len(act.MembersRemoved) > 0 {
		if err := h.OnMembersRemoved(ctx, act.MembersRemoved, tc); err != nil {
			return err
		}
	}
	return nil
}

// OnMembersAdded handles members added to a conversation.
// Override in a derived struct to send welcome messages.
func (h *ActivityHandler) OnMembersAdded(ctx context.Context, membersAdded []*activity.ChannelAccount, tc *TurnContext) error {
	return nil
}

// OnMembersRemoved handles members removed from a conversation.
func (h *ActivityHandler) OnMembersRemoved(ctx context.Context, membersRemoved []*activity.ChannelAccount, tc *TurnContext) error {
	return nil
}

// OnMessageReactionActivity handles messageReaction activities.
// By default it calls OnReactionsAdded or OnReactionsRemoved as appropriate.
func (h *ActivityHandler) OnMessageReactionActivity(ctx context.Context, tc *TurnContext) error {
	act := tc.Activity()
	if len(act.ReactionsAdded) > 0 {
		if err := h.OnReactionsAdded(ctx, act.ReactionsAdded, tc); err != nil {
			return err
		}
	}
	if len(act.ReactionsRemoved) > 0 {
		if err := h.OnReactionsRemoved(ctx, act.ReactionsRemoved, tc); err != nil {
			return err
		}
	}
	return nil
}

// OnReactionsAdded handles reactions added to a message.
func (h *ActivityHandler) OnReactionsAdded(ctx context.Context, reactions []*activity.MessageReaction, tc *TurnContext) error {
	return nil
}

// OnReactionsRemoved handles reactions removed from a message.
func (h *ActivityHandler) OnReactionsRemoved(ctx context.Context, reactions []*activity.MessageReaction, tc *TurnContext) error {
	return nil
}

// OnEventActivity handles event activities.
// By default it dispatches to OnTokenResponseEvent for tokens/response events,
// or OnEvent for all other event activities.
func (h *ActivityHandler) OnEventActivity(ctx context.Context, tc *TurnContext) error {
	if tc.Activity().Name == "tokens/response" {
		return h.OnTokenResponseEvent(ctx, tc)
	}
	return h.OnEvent(ctx, tc)
}

// OnTokenResponseEvent handles tokens/response event activities (OAuth flows).
func (h *ActivityHandler) OnTokenResponseEvent(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnEvent handles general event activities that are not tokens/response.
func (h *ActivityHandler) OnEvent(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnInvokeActivity handles invoke activities.
// Override to process adaptive card actions, sign-in verification, etc.
func (h *ActivityHandler) OnInvokeActivity(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnTypingActivity handles typing activities.
func (h *ActivityHandler) OnTypingActivity(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnEndOfConversationActivity handles endOfConversation activities.
func (h *ActivityHandler) OnEndOfConversationActivity(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnInstallationUpdateActivity handles installationUpdate activities.
func (h *ActivityHandler) OnInstallationUpdateActivity(ctx context.Context, tc *TurnContext) error {
	return nil
}

// OnUnrecognizedActivityType handles any activity type not recognized by other handlers.
func (h *ActivityHandler) OnUnrecognizedActivityType(ctx context.Context, tc *TurnContext) error {
	return nil
}
