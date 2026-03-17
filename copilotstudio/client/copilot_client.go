// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ameena3/Agents-for-golang/activity"
)

// CopilotClient communicates with a Copilot Studio agent via Direct-to-Engine protocol.
type CopilotClient struct {
	settings       *ConnectionSettings
	tokenProvider  TokenProvider
	httpClient     *http.Client
	conversationID string
	token          string
}

// NewCopilotClient creates a new CopilotClient.
func NewCopilotClient(_ context.Context, settings *ConnectionSettings, tokenProvider TokenProvider) (*CopilotClient, error) {
	if settings == nil {
		return nil, fmt.Errorf("copilotstudio: settings are required")
	}
	if settings.Cloud == "" {
		settings.Cloud = PowerPlatformCloudPublic
	}
	if settings.AgentType == "" {
		settings.AgentType = AgentTypePublished
	}
	return &CopilotClient{
		settings:      settings,
		tokenProvider: tokenProvider,
		httpClient:    &http.Client{},
	}, nil
}

// StartConversation initializes a new conversation with the Copilot Studio agent.
func (c *CopilotClient) StartConversation(ctx context.Context) (*StartConversationResponse, error) {
	url := c.settings.GetEndpointURL() + "/conversations"
	resp, err := c.doPost(ctx, url, struct{}{})
	if err != nil {
		return nil, fmt.Errorf("copilotstudio: starting conversation: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("copilotstudio: reading start response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("copilotstudio: server error %d: %s", resp.StatusCode, string(body))
	}

	var result StartConversationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("copilotstudio: decoding start response: %w", err)
	}
	c.conversationID = result.ConversationID
	c.token = result.Token
	return &result, nil
}

// ExecuteTurn sends an activity and retrieves the bot's response activities.
func (c *CopilotClient) ExecuteTurn(ctx context.Context, act *activity.Activity) (*ExecuteTurnResponse, error) {
	if c.conversationID == "" {
		return nil, fmt.Errorf("copilotstudio: no active conversation; call StartConversation first")
	}
	url := fmt.Sprintf("%s/conversations/%s/activities",
		c.settings.GetEndpointURL(), c.conversationID)

	resp, err := c.doPost(ctx, url, &ExecuteTurnRequest{Activity: act})
	if err != nil {
		return nil, fmt.Errorf("copilotstudio: executing turn: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("copilotstudio: reading turn response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("copilotstudio: server error %d: %s", resp.StatusCode, string(body))
	}

	var result ExecuteTurnResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("copilotstudio: decoding turn response: %w", err)
	}
	return &result, nil
}

// SendText is a convenience method that sends a plain-text message.
func (c *CopilotClient) SendText(ctx context.Context, text string) (*ExecuteTurnResponse, error) {
	return c.ExecuteTurn(ctx, &activity.Activity{
		Type: activity.ActivityTypeMessage,
		Text: text,
	})
}

func (c *CopilotClient) doPost(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Prefer token provider over conversation token
	if c.tokenProvider != nil {
		tok, err := c.tokenProvider.GetAccessToken(ctx, "https://api.powerplatform.com")
		if err == nil && tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
	} else if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return c.httpClient.Do(req)
}
