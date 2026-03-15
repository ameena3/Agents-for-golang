// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

// ContentTypeAdaptiveCard is the content type for an adaptive card attachment.
const ContentTypeAdaptiveCard = "application/vnd.microsoft.card.adaptive"

// AdaptiveCardInvokeAction represents the action in an adaptiveCard/action invoke.
type AdaptiveCardInvokeAction struct {
	// Type is the type of the adaptive card invoke action.
	Type string `json:"type,omitempty"`
	// ID is the ID of the adaptive card invoke action.
	ID string `json:"id,omitempty"`
	// Verb is the verb of the adaptive card invoke action.
	Verb string `json:"verb,omitempty"`
	// Data contains the data for this adaptive card invoke action.
	Data map[string]interface{} `json:"data,omitempty"`
}

// TokenExchangeInvokeRequest represents a token exchange request for adaptive card invocations.
type TokenExchangeInvokeRequest struct {
	// ID is the ID of the token exchange request.
	ID string `json:"id,omitempty"`
	// ConnectionName is the connection name.
	ConnectionName string `json:"connectionName,omitempty"`
	// Token is the exchange token.
	Token string `json:"token,omitempty"`
}

// AdaptiveCardInvokeValue defines the structure arriving in Activity.Value for
// Invoke activity with Name of 'adaptiveCard/action'.
type AdaptiveCardInvokeValue struct {
	// Action is the adaptive card invoke action.
	Action *AdaptiveCardInvokeAction `json:"action,omitempty"`
	// Authentication is the token exchange invoke request.
	Authentication *TokenExchangeInvokeRequest `json:"authentication,omitempty"`
	// State is the 'state' or magic code for an OAuth flow.
	State string `json:"state,omitempty"`
}

// AdaptiveCardInvokeResponse defines the structure returned as the result of an
// Invoke activity with Name of 'adaptiveCard/action'.
type AdaptiveCardInvokeResponse struct {
	// StatusCode is the card action response status code.
	StatusCode int `json:"statusCode,omitempty"`
	// Type is the type of this card action response.
	Type string `json:"type,omitempty"`
	// Value is the JSON response object.
	Value map[string]interface{} `json:"value,omitempty"`
}
