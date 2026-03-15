// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package app

import (
	"context"
	"fmt"
	"log"

	"github.com/microsoft/agents-sdk-go/activity"
	"github.com/microsoft/agents-sdk-go/hosting/core"
	"github.com/microsoft/agents-sdk-go/hosting/core/app/routes"
	"github.com/microsoft/agents-sdk-go/hosting/core/app/state"
)

// HandlerFunc is the signature for every route handler registered with AgentApplication.
// ctx carries request-scoped values; tc is the current turn context; appState is the
// caller-supplied application state for this turn.
type HandlerFunc[StateT any] func(ctx context.Context, tc *core.TurnContext, appState StateT) error

// AgentApplication is the modern, functional-style agent framework.
// It replaces inheritance-based ActivityHandler with registered handler functions.
//
// Usage:
//
//	app := New[MyState](AppOptions[MyState]{Storage: myStorage})
//
//	app.OnMessage("", func(ctx context.Context, tc *core.TurnContext, s MyState) error {
//	    return tc.SendActivity(ctx, activity.NewMessageActivity("hello"))
//	})
//
// StateT is the type of the application-specific state threaded through all handlers.
//
// The routing design: RouteList[StateT] stores routes whose Handler has signature
// func(ctx context.Context, tcIface interface{}, appState StateT) error.
// The tcIface parameter carries *core.TurnContext at runtime; it is interface{} because
// the routes package must remain independent of the core package.
type AgentApplication[StateT any] struct {
	options     AppOptions[StateT]
	routeList   *routes.RouteList[StateT]
	before      []func(context.Context, *core.TurnContext) error
	after       []func(context.Context, *core.TurnContext) error
	errHandlers []func(context.Context, *core.TurnContext, error) error
}

// New creates a new AgentApplication with the given options.
func New[StateT any](opts AppOptions[StateT]) *AgentApplication[StateT] {
	return &AgentApplication[StateT]{
		options:   opts,
		routeList: routes.NewRouteList[StateT](),
	}
}

// OnTurn implements the core.Agent interface.
// It executes the full per-turn pipeline:
//  1. Run BeforeTurn hooks (abort turn on error).
//  2. Load state from storage (if configured).
//  3. Find and execute the matching route handler.
//  4. Save state to storage.
//  5. Run AfterTurn hooks.
//
// If any step returns an error and error handlers are registered, they are called
// instead of returning the error directly. If no error handlers are registered the
// error is returned to the caller.
func (a *AgentApplication[StateT]) OnTurn(ctx context.Context, tc *core.TurnContext) error {
	// --- Before hooks ---
	for _, hook := range a.before {
		if err := hook(ctx, tc); err != nil {
			return a.handleError(ctx, tc, fmt.Errorf("before turn hook: %w", err))
		}
	}

	// --- Build application state ---
	var appState StateT
	if a.options.TurnStateFactory != nil {
		appState = a.options.TurnStateFactory()
	}

	// --- Load persistent state ---
	var turnState *state.TurnState[StateT]
	if a.options.Storage != nil {
		turnState = state.NewTurnState[StateT](a.options.Storage)
		channelID, convID, userID := extractStateKeys(tc)
		if err := turnState.Load(ctx, channelID, convID, userID); err != nil {
			return a.handleError(ctx, tc, fmt.Errorf("load state: %w", err))
		}
		// App field carries the per-turn typed state.
		turnState.App = appState
		appState = turnState.App
	}

	// --- Find and execute matching route ---
	route := a.routeList.FindRoute(ctx, tc.Activity())
	if route == nil {
		act := tc.Activity()
		log.Printf("AgentApplication: no route matched activity type=%q name=%q text=%q",
			act.Type, act.Name, act.Text)
	} else {
		// route.Handler has signature: func(ctx, tcIface interface{}, appState StateT) error
		// We pass tc as interface{} so the routes package stays decoupled from core.
		if err := route.Handler(ctx, tc, appState); err != nil {
			return a.handleError(ctx, tc, fmt.Errorf("route handler: %w", err))
		}
	}

	// --- Save persistent state ---
	if a.options.Storage != nil && turnState != nil {
		channelID, convID, userID := extractStateKeys(tc)
		if err := turnState.Save(ctx, channelID, convID, userID); err != nil {
			return a.handleError(ctx, tc, fmt.Errorf("save state: %w", err))
		}
	}

	// --- After hooks ---
	for _, hook := range a.after {
		if err := hook(ctx, tc); err != nil {
			return a.handleError(ctx, tc, fmt.Errorf("after turn hook: %w", err))
		}
	}

	return nil
}

// addRoute is the internal helper for registering a route. The handler wrapper
// converts the routes.Handler[StateT] convention (tcIface interface{}) back to
// *core.TurnContext before calling the user-supplied HandlerFunc.
func (a *AgentApplication[StateT]) addRoute(route *routes.Route[StateT]) {
	a.routeList.Add(route)
}

// makeHandler wraps a HandlerFunc[StateT] into the routes.Handler[StateT] type.
// The routes package passes the turn context as interface{}; we type-assert it back
// to *core.TurnContext here so user-facing handlers always receive the concrete type.
func makeHandler[StateT any](h HandlerFunc[StateT]) routes.Handler[StateT] {
	return func(ctx context.Context, tcIface interface{}, appState StateT) error {
		tc, ok := tcIface.(*core.TurnContext)
		if !ok {
			return fmt.Errorf("AgentApplication: unexpected turn context type %T", tcIface)
		}
		return h(ctx, tc, appState)
	}
}

// OnMessage registers a handler for message activities.
// If pattern is non-empty, only messages whose text matches the regular expression
// are handled. An empty pattern matches all messages.
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) OnMessage(pattern string, handler HandlerFunc[StateT]) *AgentApplication[StateT] {
	route := routes.NewMessageRoute(pattern, makeHandler(handler))
	a.addRoute(route)
	return a
}

// OnActivity registers a handler for a specific activity type.
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) OnActivity(activityType string, handler HandlerFunc[StateT]) *AgentApplication[StateT] {
	route := routes.NewActivityTypeRoute(activityType, makeHandler(handler))
	a.addRoute(route)
	return a
}

// OnInvoke registers a handler for invoke activities with the given name.
// An empty name matches all invoke activities.
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) OnInvoke(name string, handler HandlerFunc[StateT]) *AgentApplication[StateT] {
	route := routes.NewInvokeRoute(name, makeHandler(handler))
	a.addRoute(route)
	return a
}

// OnConversationUpdate registers a handler for conversationUpdate activities.
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) OnConversationUpdate(handler HandlerFunc[StateT]) *AgentApplication[StateT] {
	route := routes.NewConversationUpdateRoute(makeHandler(handler))
	a.addRoute(route)
	return a
}

// OnMembersAdded registers a handler for conversationUpdate activities where
// the MembersAdded list is non-empty (i.e., someone joined the conversation).
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) OnMembersAdded(handler HandlerFunc[StateT]) *AgentApplication[StateT] {
	route := routes.NewMembersAddedRoute(makeHandler(handler))
	a.addRoute(route)
	return a
}

// OnError registers an error handler. Handlers are called in registration order
// when any pipeline step returns an error. If a handler returns nil the error is
// considered handled; if it returns a non-nil error that error is propagated.
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) OnError(handler func(context.Context, *core.TurnContext, error) error) *AgentApplication[StateT] {
	a.errHandlers = append(a.errHandlers, handler)
	return a
}

// BeforeTurn registers a hook called before every turn, prior to route matching.
// If the hook returns an error the turn is aborted and error handlers are invoked.
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) BeforeTurn(handler func(context.Context, *core.TurnContext) error) *AgentApplication[StateT] {
	a.before = append(a.before, handler)
	return a
}

// AfterTurn registers a hook called after every turn, after state has been saved.
// If the hook returns an error, error handlers are invoked.
//
// Returns the AgentApplication for method chaining.
func (a *AgentApplication[StateT]) AfterTurn(handler func(context.Context, *core.TurnContext) error) *AgentApplication[StateT] {
	a.after = append(a.after, handler)
	return a
}

// --- compile-time interface assertion ---

var _ core.Agent = (*AgentApplication[struct{}])(nil)

// --- helpers ---

// handleError calls all registered error handlers in order.
// If no error handlers are registered, the original error is returned unchanged.
func (a *AgentApplication[StateT]) handleError(ctx context.Context, tc *core.TurnContext, err error) error {
	if len(a.errHandlers) == 0 {
		return err
	}
	current := err
	for _, h := range a.errHandlers {
		if herr := h(ctx, tc, current); herr != nil {
			current = herr
		} else {
			// handler consumed the error
			current = nil
			break
		}
	}
	return current
}

// extractStateKeys extracts the channel ID, conversation ID, and user ID from
// a TurnContext's Activity for use as storage keys.
func extractStateKeys(tc *core.TurnContext) (channelID, conversationID, userID string) {
	if tc == nil {
		return "", "", ""
	}
	act := tc.Activity()
	if act == nil {
		return "", "", ""
	}
	channelID = act.ChannelID
	if act.Conversation != nil {
		conversationID = act.Conversation.ID
	}
	if act.From != nil {
		userID = act.From.ID
	}
	return
}

// Ensure activity package is referenced (avoids unused import if only activity constants used).
var _ = activity.ActivityTypeMessage
