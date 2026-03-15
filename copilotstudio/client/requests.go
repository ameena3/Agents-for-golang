// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package client

import "github.com/microsoft/agents-sdk-go/activity"

// StartConversationResponse is returned when a conversation is started.
type StartConversationResponse struct {
	// ConversationID is the ID of the new conversation.
	ConversationID string `json:"conversationId,omitempty"`
	// Token is the Direct Line token for this conversation.
	Token string `json:"token,omitempty"`
	// ExpiresIn is the token expiry time in seconds.
	ExpiresIn int `json:"expires_in,omitempty"`
	// StreamURL is the WebSocket stream URL (if streaming is enabled).
	StreamURL string `json:"streamUrl,omitempty"`
}

// ExecuteTurnRequest sends a message to an existing conversation.
type ExecuteTurnRequest struct {
	// Activity is the activity to send.
	Activity *activity.Activity `json:"activity,omitempty"`
}

// ExecuteTurnResponse holds the bot's response activities.
type ExecuteTurnResponse struct {
	// Activities are the response activities from the bot.
	Activities []*activity.Activity `json:"activities,omitempty"`
	// Watermark is used for polling continuations.
	Watermark string `json:"watermark,omitempty"`
}
