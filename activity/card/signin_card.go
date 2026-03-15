// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

// ContentTypeSigninCard is the content type for a sign-in card attachment.
const ContentTypeSigninCard = "application/vnd.microsoft.card.signin"

// SigninCard represents a request to sign in.
type SigninCard struct {
	// Text is the text for the sign-in request.
	Text string `json:"text,omitempty"`
	// Buttons is the list of actions to use to perform sign-in.
	Buttons []*CardAction `json:"buttons,omitempty"`
}
