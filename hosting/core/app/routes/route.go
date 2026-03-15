// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package routes provides routing primitives for AgentApplication.
// Routes are matched in priority order: invoke > agentic > regular, then by rank (lower = higher priority).
package routes

import (
	"context"
	"regexp"
	"strings"

	"github.com/microsoft/agents-sdk-go/activity"
)

// Selector is a function that returns true if a route should handle the given activity.
type Selector func(ctx context.Context, act *activity.Activity) bool

// Handler is the function called when a route matches.
// StateT is the application state type threaded through all handlers.
type Handler[StateT any] func(ctx context.Context, tc interface{}, state StateT) error

// RouteRank defines the priority of a route within its category.
// Lower numbers indicate higher priority.
type RouteRank int

const (
	// RouteRankFirst is the highest possible user-assigned rank.
	RouteRankFirst RouteRank = 0
	// RouteRankDefault is the default rank assigned to new routes.
	RouteRankDefault RouteRank = 32767
	// RouteRankLast is the lowest possible user-assigned rank.
	RouteRankLast RouteRank = 65535
)

// routeCategory is the top-level priority bucket.
// Lower values sort before higher values.
type routeCategory int

const (
	routeCategoryInvokeAgentic routeCategory = 0 // invoke AND agentic
	routeCategoryInvoke        routeCategory = 1 // invoke only
	routeCategoryAgentic       routeCategory = 2 // agentic only
	routeCategoryDefault       routeCategory = 3 // everything else
)

// Route represents a single registered handler with its matching selector and priority.
type Route[StateT any] struct {
	Selector Selector
	Handler  Handler[StateT]
	// category is used for priority sorting; lower = higher priority.
	category routeCategory
	// rank is the secondary sort key within a category.
	rank RouteRank
	// insertionOrder is the FIFO tie-breaker among same category+rank routes.
	insertionOrder int
}

// priority returns a comparable triple used for sorting.
func (r *Route[StateT]) priority() [3]int {
	return [3]int{int(r.category), int(r.rank), r.insertionOrder}
}

// less reports whether r should sort before other (higher priority).
func (r *Route[StateT]) less(other *Route[StateT]) bool {
	rp := r.priority()
	op := other.priority()
	for i := range rp {
		if rp[i] != op[i] {
			return rp[i] < op[i]
		}
	}
	return false
}

// newRoute creates a Route with the given selector, handler, and options.
func newRoute[StateT any](
	sel Selector,
	h Handler[StateT],
	isInvoke bool,
	isAgentic bool,
	rank RouteRank,
) *Route[StateT] {
	cat := routeCategoryDefault
	switch {
	case isInvoke && isAgentic:
		cat = routeCategoryInvokeAgentic
	case isInvoke:
		cat = routeCategoryInvoke
	case isAgentic:
		cat = routeCategoryAgentic
	}
	return &Route[StateT]{
		Selector: sel,
		Handler:  h,
		category: cat,
		rank:     rank,
	}
}

// NewMessageRoute creates a route that matches message activities.
// If pattern is empty it matches all messages; otherwise it matches messages
// whose text satisfies regexp.MatchString(pattern, text).
func NewMessageRoute[StateT any](pattern string, handler Handler[StateT]) *Route[StateT] {
	var sel Selector
	if pattern == "" {
		sel = func(_ context.Context, act *activity.Activity) bool {
			return strings.EqualFold(act.Type, activity.ActivityTypeMessage)
		}
	} else {
		re := regexp.MustCompile(pattern)
		sel = func(_ context.Context, act *activity.Activity) bool {
			return strings.EqualFold(act.Type, activity.ActivityTypeMessage) &&
				re.MatchString(act.Text)
		}
	}
	return newRoute(sel, handler, false, false, RouteRankDefault)
}

// NewActivityTypeRoute creates a route that matches a specific activity type.
func NewActivityTypeRoute[StateT any](activityType string, handler Handler[StateT]) *Route[StateT] {
	sel := func(_ context.Context, act *activity.Activity) bool {
		return strings.EqualFold(act.Type, activityType)
	}
	return newRoute(sel, handler, false, false, RouteRankDefault)
}

// NewInvokeRoute creates a route that matches invoke activities with the given name.
// If name is empty, it matches all invoke activities.
func NewInvokeRoute[StateT any](name string, handler Handler[StateT]) *Route[StateT] {
	var sel Selector
	if name == "" {
		sel = func(_ context.Context, act *activity.Activity) bool {
			return strings.EqualFold(act.Type, activity.ActivityTypeInvoke)
		}
	} else {
		sel = func(_ context.Context, act *activity.Activity) bool {
			return strings.EqualFold(act.Type, activity.ActivityTypeInvoke) &&
				act.Name == name
		}
	}
	return newRoute(sel, handler, true, false, RouteRankDefault)
}

// NewConversationUpdateRoute creates a route that matches conversationUpdate activities.
func NewConversationUpdateRoute[StateT any](handler Handler[StateT]) *Route[StateT] {
	sel := func(_ context.Context, act *activity.Activity) bool {
		return strings.EqualFold(act.Type, activity.ActivityTypeConversationUpdate)
	}
	return newRoute(sel, handler, false, false, RouteRankDefault)
}

// NewMembersAddedRoute creates a route that matches conversationUpdate activities
// where the MembersAdded list is non-empty.
func NewMembersAddedRoute[StateT any](handler Handler[StateT]) *Route[StateT] {
	sel := func(_ context.Context, act *activity.Activity) bool {
		return strings.EqualFold(act.Type, activity.ActivityTypeConversationUpdate) &&
			len(act.MembersAdded) > 0
	}
	return newRoute(sel, handler, false, false, RouteRankDefault)
}

// NewAgenticRoute wraps a selector so it only fires when the activity is an
// agentic request (from an agent identity rather than a user).
func NewAgenticRoute[StateT any](sel Selector, handler Handler[StateT], isInvoke bool, rank RouteRank) *Route[StateT] {
	agSel := func(ctx context.Context, act *activity.Activity) bool {
		return act.IsAgenticRequest() && sel(ctx, act)
	}
	return newRoute(agSel, handler, isInvoke, true, rank)
}
