// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

// ContentTypeOAuthCard is the content type for an OAuth card attachment.
const ContentTypeOAuthCard = "application/vnd.microsoft.card.oauth"

// TokenExchangeResource contains information about a token exchange resource.
type TokenExchangeResource struct {
	// ID is the unique identifier for this token exchange resource.
	ID string `json:"id,omitempty"`
	// URI is the URI for this token exchange resource.
	URI string `json:"uri,omitempty"`
	// ProviderDisplayName is the display name of the provider.
	ProviderDisplayName string `json:"providerDisplayName,omitempty"`
}

// TokenPostResource contains information about a token post resource.
type TokenPostResource struct {
	// SasURL is the SAS URL for posting the token.
	SasURL string `json:"sasUrl,omitempty"`
}

// OAuthCard represents a request to perform a sign in via OAuth.
type OAuthCard struct {
	// Text is the text for the sign-in request.
	Text string `json:"text,omitempty"`
	// ConnectionName is the name of the registered connection.
	ConnectionName string `json:"connectionName,omitempty"`
	// Buttons is the list of actions to use to perform sign-in.
	Buttons []*CardAction `json:"buttons,omitempty"`
	// TokenExchangeResource contains information for token exchange.
	TokenExchangeResource *TokenExchangeResource `json:"tokenExchangeResource,omitempty"`
	// TokenPostResource contains information for posting the token.
	TokenPostResource *TokenPostResource `json:"tokenPostResource,omitempty"`
}
