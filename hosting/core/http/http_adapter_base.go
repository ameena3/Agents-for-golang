// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ameena3/Agents-for-golang/activity"
	"github.com/ameena3/Agents-for-golang/hosting/core"
	"github.com/ameena3/Agents-for-golang/hosting/core/authorization"
)

// HttpAdapterBase provides common logic for processing HTTP requests
// that carry Bot Framework activity payloads. Web-framework adapters
// embed or call this base to avoid duplicating parsing and dispatch.
type HttpAdapterBase struct {
	adapter              *core.ChannelServiceAdapter
	allowUnauthenticated bool
}

// NewHttpAdapterBase creates a new HttpAdapterBase.
func NewHttpAdapterBase(adapter *core.ChannelServiceAdapter, allowUnauthenticated bool) *HttpAdapterBase {
	return &HttpAdapterBase{
		adapter:              adapter,
		allowUnauthenticated: allowUnauthenticated,
	}
}

// Process parses the incoming HTTP request, constructs a ClaimsIdentity,
// and dispatches the activity through the adapter pipeline.
// It writes an appropriate HTTP response (200/202/400/401/405/500).
func (h *HttpAdapterBase) Process(ctx context.Context, req HttpRequestProtocol, agent core.Agent, w http.ResponseWriter) {
	if req.Method() != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body())
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var act activity.Activity
	if err := json.Unmarshal(body, &act); err != nil {
		http.Error(w, fmt.Sprintf("invalid activity JSON: %v", err), http.StatusBadRequest)
		return
	}

	var identity *authorization.ClaimsIdentity
	if h.allowUnauthenticated {
		identity = authorization.NewClaimsIdentity(false, "Anonymous", nil)
	} else {
		identity = req.GetClaimsIdentity(ctx)
		if identity == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Process the activity through the pipeline.
	var tc *core.TurnContext
	pipelineErr := h.adapter.ProcessActivity(ctx, identity, &act, func(pCtx context.Context, turnCtx *core.TurnContext) error {
		tc = turnCtx
		return agent.OnTurn(pCtx, turnCtx)
	})

	if pipelineErr != nil {
		http.Error(w, fmt.Sprintf("error processing activity: %v", pipelineErr), http.StatusInternalServerError)
		return
	}

	statusCode := http.StatusAccepted
	if tc != nil {
		statusCode = h.adapter.GetHTTPStatusCode(tc)
	}

	if statusCode == http.StatusOK && tc != nil && len(tc.BufferedReplies) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(tc.BufferedReplies)
		return
	}

	w.WriteHeader(statusCode)
}
