// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"context"
	"time"

	"github.com/ameena3/Agents-for-golang/activity"
)

// TypingIndicatorMiddleware sends a typing indicator at the start of each turn.
type TypingIndicatorMiddleware struct {
	delay time.Duration
}

// NewTypingIndicatorMiddleware creates a new TypingIndicatorMiddleware.
// delay specifies how long to wait before sending the indicator (0 = immediately).
func NewTypingIndicatorMiddleware(delay time.Duration) *TypingIndicatorMiddleware {
	return &TypingIndicatorMiddleware{delay: delay}
}

// OnTurn implements Middleware. It sends a typing indicator activity before
// calling the next handler in the pipeline.
func (t *TypingIndicatorMiddleware) OnTurn(ctx context.Context, tc *TurnContext, next func(context.Context) error) error {
	// Only send typing indicators for message activities.
	if tc.Activity() != nil && tc.Activity().Type == activity.ActivityTypeMessage {
		if t.delay > 0 {
			select {
			case <-time.After(t.delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		typingAct := &activity.Activity{Type: activity.ActivityTypeTyping}
		_, _ = tc.SendActivity(ctx, typingAct)
	}
	return next(ctx)
}
