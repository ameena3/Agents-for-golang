// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TeamsConnectorClient provides Teams-specific API operations.
// It embeds ConnectorClient for standard channel operations.
type TeamsConnectorClient struct {
	ConnectorClient
}

// NewTeamsConnectorClient creates a new TeamsConnectorClient.
func NewTeamsConnectorClient(serviceURL string, tokenProvider TokenProvider) *TeamsConnectorClient {
	return &TeamsConnectorClient{
		ConnectorClient: ConnectorClient{
			serviceURL:    serviceURL,
			httpClient:    &http.Client{},
			tokenProvider: tokenProvider,
		},
	}
}

// teamsGet is a helper for Teams GET endpoints that return JSON.
func (c *TeamsConnectorClient) teamsGet(ctx context.Context, path string, result interface{}) error {
	rawURL := fmt.Sprintf("%s%s", c.serviceURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("teams_connector: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if c.tokenProvider != nil {
		token, err := c.tokenProvider.GetAccessToken(ctx, c.serviceURL)
		if err != nil {
			return fmt.Errorf("teams_connector: get access token: %w", err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("teams_connector: execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("teams_connector: unexpected status %d: %s", resp.StatusCode, string(responseBody))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("teams_connector: decode response: %w", err)
		}
	}
	return nil
}

// GetTeamDetails retrieves team details.
func (c *TeamsConnectorClient) GetTeamDetails(ctx context.Context, teamID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/v3/teams/%s", teamID)
	var result map[string]interface{}
	if err := c.teamsGet(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetChannels retrieves the list of channels in a team.
func (c *TeamsConnectorClient) GetChannels(ctx context.Context, teamID string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/v3/teams/%s/conversations", teamID)
	var wrapper struct {
		Conversations []map[string]interface{} `json:"conversations"`
	}
	if err := c.teamsGet(ctx, path, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.Conversations, nil
}

// GetMembers retrieves the members of a conversation.
func (c *TeamsConnectorClient) GetMembers(ctx context.Context, conversationID string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/v3/conversations/%s/members", conversationID)
	var result []map[string]interface{}
	if err := c.teamsGet(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}
