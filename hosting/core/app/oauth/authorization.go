// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package oauth

import (
	"context"
	"fmt"
	"time"

	"github.com/microsoft/agents-sdk-go/activity"
	"github.com/microsoft/agents-sdk-go/hosting/core"
)

const defaultSignInExpiry = 10 * time.Minute

// OAuthFlow manages the state machine for OAuth sign-in flows.
// It sends OAuthCards, validates magic codes, and retrieves tokens.
type OAuthFlow struct {
	connectionName string
	title          string
	text           string
}

// NewOAuthFlow creates a new OAuthFlow for the given connection.
func NewOAuthFlow(connectionName, title, text string) *OAuthFlow {
	return &OAuthFlow{
		connectionName: connectionName,
		title:          title,
		text:           text,
	}
}

// Begin starts the OAuth flow by sending an OAuthCard to the user.
func (o *OAuthFlow) Begin(ctx context.Context, tc *core.TurnContext) (*SignInState, error) {
	act := tc.Activity()

	var userID string
	if act.From != nil {
		userID = act.From.ID
	}

	var conversationID string
	if act.Conversation != nil {
		conversationID = act.Conversation.ID
	}

	state := &SignInState{
		ConnectionName: o.connectionName,
		UserID:         userID,
		ConversationID: conversationID,
		Expiry:         time.Now().Add(defaultSignInExpiry),
		FlowType:       "user",
	}

	oauthCard := &activity.Activity{
		Type: activity.ActivityTypeMessage,
		Attachments: []*activity.Attachment{
			{
				ContentType: "application/vnd.microsoft.card.oauth",
				Content: map[string]interface{}{
					"connectionName": o.connectionName,
					"title":          o.title,
					"text":           o.text,
				},
			},
		},
	}

	if _, err := tc.SendActivity(ctx, oauthCard); err != nil {
		return nil, fmt.Errorf("oauth: sending sign-in card: %w", err)
	}
	return state, nil
}

// Continue checks whether the current activity completes a pending sign-in.
// Returns (token, true, nil) if auth is complete, ("", false, nil) if still pending.
func (o *OAuthFlow) Continue(ctx context.Context, tc *core.TurnContext, state *SignInState) (string, bool, error) {
	if state == nil || state.IsExpired() {
		return "", false, fmt.Errorf("oauth: sign-in state is nil or expired")
	}

	act := tc.Activity()

	// Handle token response event
	if act.Type == activity.ActivityTypeEvent && act.Name == "tokens/response" {
		if tokenMap, ok := act.Value.(map[string]interface{}); ok {
			if token, ok := tokenMap["token"].(string); ok && token != "" {
				return token, true, nil
			}
		}
		return "", false, fmt.Errorf("oauth: token response event missing token value")
	}

	// Handle magic code message
	if act.Type == activity.ActivityTypeMessage && act.Text != "" {
		// Magic codes are typically 6 digits
		if len(act.Text) >= 6 && len(act.Text) <= 12 {
			return act.Text, true, nil
		}
	}

	return "", false, nil
}
