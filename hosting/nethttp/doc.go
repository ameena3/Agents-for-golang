// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package nethttp provides a net/http adapter for the Microsoft 365 Agents SDK.
// It replaces the Python aiohttp and FastAPI adapters with a single standard-library
// HTTP implementation suitable for any Go HTTP framework or the standard net/http server.
//
// Basic usage:
//
//	adapter := nethttp.NewCloudAdapter(true) // true = allow unauthenticated (local testing)
//	app := app.New[MyState](app.AppOptions[MyState]{})
//	app.OnMessage("", func(ctx context.Context, tc *core.TurnContext, s MyState) error {
//	    return tc.SendActivity(ctx, core.Text("Hello!"))
//	})
//	http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, app))
//	http.ListenAndServe(":3978", nil)
package nethttp
