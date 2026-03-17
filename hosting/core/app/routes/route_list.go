// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package routes

import (
	"context"
	"sort"

	"github.com/ameena3/Agents-for-golang/activity"
)

// RouteList holds all registered routes and selects the best match for an activity.
// Routes are stored in insertion order and sorted on demand so that the highest-priority
// (lowest category+rank+insertionOrder) route is tried first.
type RouteList[StateT any] struct {
	routes  []*Route[StateT]
	counter int // monotonically increasing insertion counter for FIFO tie-breaking
}

// NewRouteList creates an empty RouteList.
func NewRouteList[StateT any]() *RouteList[StateT] {
	return &RouteList[StateT]{}
}

// Add appends a route to the list, assigning its insertion order.
func (rl *RouteList[StateT]) Add(route *Route[StateT]) {
	route.insertionOrder = rl.counter
	rl.counter++
	rl.routes = append(rl.routes, route)
}

// FindRoute returns the highest-priority route whose selector matches the given
// activity, or nil if no route matches.
//
// Routes are evaluated in priority order:
//  1. invoke+agentic routes (category 0)
//  2. invoke-only routes (category 1)
//  3. agentic-only routes (category 2)
//  4. all other routes (category 3)
//
// Within each category, lower rank numbers win; ties are broken by FIFO insertion order.
func (rl *RouteList[StateT]) FindRoute(ctx context.Context, act *activity.Activity) *Route[StateT] {
	// Copy and sort so we evaluate in correct priority order without mutating
	// the underlying slice on every call.
	sorted := make([]*Route[StateT], len(rl.routes))
	copy(sorted, rl.routes)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].less(sorted[j])
	})

	for _, r := range sorted {
		if r.Selector(ctx, act) {
			return r
		}
	}
	return nil
}

// Len returns the number of registered routes.
func (rl *RouteList[StateT]) Len() int {
	return len(rl.routes)
}
