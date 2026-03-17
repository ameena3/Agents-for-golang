// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package client

import (
	"context"

	"github.com/ameena3/Agents-for-golang/activity"
)

// TokenProvider supplies OAuth access tokens.
type TokenProvider interface {
	GetAccessToken(ctx context.Context, resource string) (string, error)
}

// CopilotClientProtocol defines the interface for Copilot Studio clients.
type CopilotClientProtocol interface {
	// StartConversation initializes a new conversation with the agent.
	StartConversation(ctx context.Context) (*StartConversationResponse, error)
	// ExecuteTurn sends an activity and retrieves the bot's responses.
	ExecuteTurn(ctx context.Context, act *activity.Activity) (*ExecuteTurnResponse, error)
	// SendText is a convenience method that sends a plain-text message.
	SendText(ctx context.Context, text string) (*ExecuteTurnResponse, error)
}
