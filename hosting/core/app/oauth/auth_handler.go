// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package oauth

import (
	"context"
	"fmt"

	"github.com/microsoft/agents-sdk-go/hosting/core"
)

// AuthHandler wraps route handlers that require authentication.
// It automatically starts the OAuth flow if the user is not signed in
// and calls the wrapped handler once a token is obtained.
type AuthHandler[StateT any] struct {
	flow    *OAuthFlow
	handler func(ctx context.Context, tc *core.TurnContext, token string) error
}

// NewAuthHandler creates an AuthHandler that wraps the given handler
// with an OAuth sign-in flow using the specified connection name.
func NewAuthHandler[StateT any](connectionName, title, text string, handler func(ctx context.Context, tc *core.TurnContext, token string) error) *AuthHandler[StateT] {
	return &AuthHandler[StateT]{
		flow:    NewOAuthFlow(connectionName, title, text),
		handler: handler,
	}
}

// Handle processes the turn. If no token exists, starts sign-in.
// If a sign-in is pending, tries to complete it. If complete, calls the handler.
func (a *AuthHandler[StateT]) Handle(ctx context.Context, tc *core.TurnContext, state *SignInState) error {
	if state == nil {
		// Start sign-in
		newState, err := a.flow.Begin(ctx, tc)
		if err != nil {
			return fmt.Errorf("auth_handler: beginning oauth flow: %w", err)
		}
		_ = newState // caller should persist this state
		return nil
	}

	token, complete, err := a.flow.Continue(ctx, tc, state)
	if err != nil {
		return fmt.Errorf("auth_handler: continuing oauth flow: %w", err)
	}

	if !complete {
		return nil // still waiting
	}

	return a.handler(ctx, tc, token)
}
