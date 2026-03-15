// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"context"
)

// ChannelAdapterBase provides a reusable base for ChannelAdapter implementations.
// Embed this struct in concrete adapter implementations to inherit middleware
// registration and pipeline execution.
type ChannelAdapterBase struct {
	middlewareSet *MiddlewareSet
	// OnTurnError is an optional error handler called when the middleware pipeline
	// returns an error. If nil, the error is returned to the caller unchanged.
	OnTurnError func(ctx context.Context, tc *TurnContext, err error) error
}

// NewChannelAdapterBase creates a new ChannelAdapterBase with an empty middleware set.
func NewChannelAdapterBase() *ChannelAdapterBase {
	return &ChannelAdapterBase{
		middlewareSet: NewMiddlewareSet(),
	}
}

// Use registers one or more middleware handlers on this adapter.
func (a *ChannelAdapterBase) Use(middleware ...Middleware) {
	a.middlewareSet.Use(middleware...)
}

// RunPipeline executes the registered middleware pipeline and then calls the
// given agent handler function. If the pipeline or handler returns an error and
// OnTurnError is set, it is invoked; otherwise the error propagates.
func (a *ChannelAdapterBase) RunPipeline(ctx context.Context, tc *TurnContext, handler func(context.Context, *TurnContext) error) error {
	err := a.middlewareSet.ReceiveActivity(ctx, tc, handler)
	if err != nil && a.OnTurnError != nil {
		return a.OnTurnError(ctx, tc, err)
	}
	return err
}
