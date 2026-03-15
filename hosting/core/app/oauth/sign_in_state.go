// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package oauth

import "time"

// SignInState tracks ongoing OAuth sign-in flows.
type SignInState struct {
	// ConnectionName is the OAuth connection to use.
	ConnectionName string `json:"connectionName"`
	// UserID is the user initiating the sign-in.
	UserID string `json:"userId"`
	// ConversationID tracks the conversation awaiting completion.
	ConversationID string `json:"conversationId"`
	// Expiry is when this sign-in attempt expires.
	Expiry time.Time `json:"expiry"`
	// FlowType is "user" or "agentic".
	FlowType string `json:"flowType"`
}

// IsExpired returns true if the sign-in state has expired.
func (s *SignInState) IsExpired() bool {
	return time.Now().After(s.Expiry)
}
