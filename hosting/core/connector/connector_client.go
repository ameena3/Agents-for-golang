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

	"github.com/microsoft/agents-sdk-go/activity"
)

// TokenProvider is a local interface to avoid circular imports.
// Callers may pass any implementation that provides GetAccessToken.
type TokenProvider interface {
	// GetAccessToken returns an access token for the given resource URL.
	GetAccessToken(ctx context.Context, resource string) (string, error)
}

// ConnectorClient sends activities to the Bot Framework channel service.
type ConnectorClient struct {
	serviceURL    string
	httpClient    *http.Client
	tokenProvider TokenProvider
}

// NewConnectorClient creates a new ConnectorClient for the given service URL.
func NewConnectorClient(serviceURL string, tokenProvider TokenProvider) *ConnectorClient {
	return &ConnectorClient{
		serviceURL:    serviceURL,
		httpClient:    &http.Client{},
		tokenProvider: tokenProvider,
	}
}

// sendRequest is a helper that marshals the body, sends an HTTP request, and
// optionally decodes the response into result.
func (c *ConnectorClient) sendRequest(ctx context.Context, method, url string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("connector: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("connector: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if c.tokenProvider != nil {
		token, err := c.tokenProvider.GetAccessToken(ctx, c.serviceURL)
		if err != nil {
			return fmt.Errorf("connector: get access token: %w", err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("connector: execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("connector: unexpected status %d: %s", resp.StatusCode, string(responseBody))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("connector: decode response: %w", err)
		}
	}
	return nil
}

// SendToConversation sends an activity to a conversation.
func (c *ConnectorClient) SendToConversation(ctx context.Context, conversationID string, act *activity.Activity) (*activity.ResourceResponse, error) {
	url := fmt.Sprintf("%s/v3/conversations/%s/activities", c.serviceURL, conversationID)
	var result activity.ResourceResponse
	if err := c.sendRequest(ctx, http.MethodPost, url, act, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ReplyToActivity sends a reply to a specific activity.
func (c *ConnectorClient) ReplyToActivity(ctx context.Context, conversationID, activityID string, act *activity.Activity) (*activity.ResourceResponse, error) {
	url := fmt.Sprintf("%s/v3/conversations/%s/activities/%s", c.serviceURL, conversationID, activityID)
	var result activity.ResourceResponse
	if err := c.sendRequest(ctx, http.MethodPost, url, act, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateActivity updates an existing activity.
func (c *ConnectorClient) UpdateActivity(ctx context.Context, conversationID, activityID string, act *activity.Activity) (*activity.ResourceResponse, error) {
	url := fmt.Sprintf("%s/v3/conversations/%s/activities/%s", c.serviceURL, conversationID, activityID)
	var result activity.ResourceResponse
	if err := c.sendRequest(ctx, http.MethodPut, url, act, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteActivity deletes an activity.
func (c *ConnectorClient) DeleteActivity(ctx context.Context, conversationID, activityID string) error {
	url := fmt.Sprintf("%s/v3/conversations/%s/activities/%s", c.serviceURL, conversationID, activityID)
	return c.sendRequest(ctx, http.MethodDelete, url, nil, nil)
}
