// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package client

// SubscribeRequest is used to set up event subscriptions.
type SubscribeRequest struct {
	// EventTypes is the list of event types to subscribe to.
	EventTypes []string `json:"eventTypes,omitempty"`
}

// SubscribeResponse is returned from a subscribe request.
type SubscribeResponse struct {
	// SubscriptionID is the unique identifier for the subscription.
	SubscriptionID string `json:"subscriptionId,omitempty"`
}

// SubscribeEvent is an event received from a Copilot Studio subscription.
type SubscribeEvent struct {
	// Type is the event type.
	Type string `json:"type,omitempty"`
	// Payload is the event-specific data.
	Payload interface{} `json:"payload,omitempty"`
}
