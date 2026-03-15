// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import "context"

// MiddlewareSet manages a collection of Middleware and runs them as a pipeline.
// The set itself satisfies the Middleware interface so it can be nested.
type MiddlewareSet struct {
	middleware []Middleware
}

// NewMiddlewareSet creates an empty MiddlewareSet.
func NewMiddlewareSet() *MiddlewareSet {
	return &MiddlewareSet{}
}

// Use adds one or more middleware to the set.
func (ms *MiddlewareSet) Use(middleware ...Middleware) {
	ms.middleware = append(ms.middleware, middleware...)
}

// ReceiveActivity runs the middleware pipeline and then calls the handler.
// Each middleware receives the context, turn context, and a next function that
// continues the pipeline. The handler is invoked after all middleware have run.
func (ms *MiddlewareSet) ReceiveActivity(ctx context.Context, tc *TurnContext, handler func(context.Context, *TurnContext) error) error {
	return ms.receiveInternal(ctx, tc, handler, 0)
}

// OnTurn implements the Middleware interface so a MiddlewareSet can be nested
// inside another MiddlewareSet or ChannelAdapterBase.
func (ms *MiddlewareSet) OnTurn(ctx context.Context, tc *TurnContext, next func(context.Context) error) error {
	err := ms.receiveInternal(ctx, tc, func(innerCtx context.Context, innerTc *TurnContext) error {
		return next(innerCtx)
	}, 0)
	return err
}

// receiveInternal is the recursive helper that walks through middleware.
func (ms *MiddlewareSet) receiveInternal(ctx context.Context, tc *TurnContext, handler func(context.Context, *TurnContext) error, index int) error {
	if index == len(ms.middleware) {
		if handler != nil {
			return handler(ctx, tc)
		}
		return nil
	}
	mw := ms.middleware[index]
	return mw.OnTurn(ctx, tc, func(nextCtx context.Context) error {
		return ms.receiveInternal(nextCtx, tc, handler, index+1)
	})
}
