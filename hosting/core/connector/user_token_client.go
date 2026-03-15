// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// TokenResponse holds an OAuth token response.
type TokenResponse struct {
	ChannelID      string `json:"channelId,omitempty"`
	ConnectionName string `json:"connectionName,omitempty"`
	Token          string `json:"token,omitempty"`
	Expiration     string `json:"expiration,omitempty"`
}

// UserTokenClient manages user OAuth tokens.
type UserTokenClient struct {
	serviceURL string
	httpClient *http.Client
}

// NewUserTokenClient creates a new UserTokenClient.
func NewUserTokenClient(serviceURL string) *UserTokenClient {
	return &UserTokenClient{
		serviceURL: serviceURL,
		httpClient: &http.Client{},
	}
}

// doRequest is a helper that performs an HTTP request and optionally decodes the response.
func (c *UserTokenClient) doRequest(ctx context.Context, method, rawURL string, body interface{}, result interface{}) error {
	req, err := buildJSONRequest(ctx, method, rawURL, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("user_token_client: execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user_token_client: unexpected status %d: %s", resp.StatusCode, string(responseBody))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("user_token_client: decode response: %w", err)
		}
	}
	return nil
}

// GetToken retrieves an OAuth token for a user.
func (c *UserTokenClient) GetToken(ctx context.Context, userID, connectionName, channelID, magicCode string) (*TokenResponse, error) {
	params := url.Values{}
	params.Set("userId", userID)
	params.Set("connectionName", connectionName)
	if channelID != "" {
		params.Set("channelId", channelID)
	}
	if magicCode != "" {
		params.Set("code", magicCode)
	}
	rawURL := fmt.Sprintf("%s/api/usertoken/GetToken?%s", c.serviceURL, params.Encode())

	var result TokenResponse
	if err := c.doRequest(ctx, http.MethodGet, rawURL, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SignOutUser signs a user out of a connection.
func (c *UserTokenClient) SignOutUser(ctx context.Context, userID, connectionName, channelID string) error {
	params := url.Values{}
	params.Set("userId", userID)
	if connectionName != "" {
		params.Set("connectionName", connectionName)
	}
	if channelID != "" {
		params.Set("channelId", channelID)
	}
	rawURL := fmt.Sprintf("%s/api/usertoken/SignOut?%s", c.serviceURL, params.Encode())
	return c.doRequest(ctx, http.MethodDelete, rawURL, nil, nil)
}

// GetSignInURL returns the sign-in URL for OAuth.
func (c *UserTokenClient) GetSignInURL(ctx context.Context, connectionName, state string) (string, error) {
	params := url.Values{}
	params.Set("connectionName", connectionName)
	if state != "" {
		params.Set("state", state)
	}
	rawURL := fmt.Sprintf("%s/api/botsignin/GetSignInUrl?%s", c.serviceURL, params.Encode())

	var signInURL string
	req, err := buildJSONRequest(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("user_token_client: execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("user_token_client: unexpected status %d: %s", resp.StatusCode, string(responseBody))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("user_token_client: read response body: %w", err)
	}
	signInURL = string(bodyBytes)
	return signInURL, nil
}

// ExchangeToken exchanges a token for another token.
func (c *UserTokenClient) ExchangeToken(ctx context.Context, userID, connectionName, channelID string, body interface{}) (*TokenResponse, error) {
	params := url.Values{}
	params.Set("userId", userID)
	params.Set("connectionName", connectionName)
	if channelID != "" {
		params.Set("channelId", channelID)
	}
	rawURL := fmt.Sprintf("%s/api/usertoken/exchange?%s", c.serviceURL, params.Encode())

	var result TokenResponse
	if err := c.doRequest(ctx, http.MethodPost, rawURL, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
