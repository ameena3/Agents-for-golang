// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// buildJSONRequest builds an HTTP request with an optional JSON body.
func buildJSONRequest(ctx context.Context, method, rawURL string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("connector: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("connector: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
