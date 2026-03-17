// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package nethttp

import (
	"net/http"

	"github.com/ameena3/Agents-for-golang/hosting/core"
)

// MessageHandler returns an http.HandlerFunc that processes bot messages.
// Mount it at /api/messages:
//
//	http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, agent))
func MessageHandler(adapter *CloudAdapter, agent core.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := adapter.Process(r.Context(), r, w, agent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
