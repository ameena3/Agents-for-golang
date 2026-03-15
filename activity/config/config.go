// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package config provides configuration loading helpers for the
// Microsoft 365 Agents SDK for Go.
package config

import (
	"os"
	"strings"
)

// AgentConfig holds configuration for an agent application.
type AgentConfig struct {
	// ClientID is the Azure AD application (client) ID.
	ClientID string
	// ClientSecret is the Azure AD application client secret.
	ClientSecret string
	// TenantID is the Azure AD tenant ID.
	TenantID string
	// AppID is the Microsoft App ID (same as ClientID for most scenarios).
	AppID string
	// AppPassword is the Microsoft App password (same as ClientSecret for most scenarios).
	AppPassword string
	// ServiceURL is the default service URL for the channel.
	ServiceURL string
	// Connections contains named connection configurations.
	Connections map[string]*ConnectionConfig
}

// ConnectionConfig holds configuration for a named connection (e.g., OAuth).
type ConnectionConfig struct {
	// ServiceURL is the connection service URL.
	ServiceURL string
	// AppID is the app ID for this connection.
	AppID string
	// AppPassword is the app password for this connection.
	AppPassword string
	// TenantID is the tenant ID for this connection.
	TenantID string
	// TokenEndpoint is the token endpoint URL.
	TokenEndpoint string
	// Scopes contains the OAuth scopes.
	Scopes []string
}

// LoadFromEnv loads agent configuration from environment variables.
// It reads standard environment variables used by the Microsoft 365 Agents SDK.
//
// Supported environment variables:
//   - MICROSOFT_APP_ID or AZURE_CLIENT_ID
//   - MICROSOFT_APP_PASSWORD or AZURE_CLIENT_SECRET
//   - MICROSOFT_APP_TENANT_ID or AZURE_TENANT_ID
//   - SERVICE_URL
func LoadFromEnv() *AgentConfig {
	cfg := &AgentConfig{
		Connections: make(map[string]*ConnectionConfig),
	}

	// App ID: try MICROSOFT_APP_ID first, then AZURE_CLIENT_ID
	cfg.AppID = getEnvFallback("MICROSOFT_APP_ID", "AZURE_CLIENT_ID")
	cfg.ClientID = cfg.AppID

	// App Password: try MICROSOFT_APP_PASSWORD first, then AZURE_CLIENT_SECRET
	cfg.AppPassword = getEnvFallback("MICROSOFT_APP_PASSWORD", "AZURE_CLIENT_SECRET")
	cfg.ClientSecret = cfg.AppPassword

	// Tenant ID: try MICROSOFT_APP_TENANT_ID first, then AZURE_TENANT_ID
	cfg.TenantID = getEnvFallback("MICROSOFT_APP_TENANT_ID", "AZURE_TENANT_ID")

	cfg.ServiceURL = os.Getenv("SERVICE_URL")

	return cfg
}

// LoadFromMap loads agent configuration from a map of environment variables,
// following the double-underscore hierarchy convention used by the Python SDK.
// For example, "CONNECTIONS__MyConn__APP_ID" sets Connections["MyConn"].AppID.
func LoadFromMap(envVars map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range envVars {
		levels := strings.Split(key, "__")
		current := result
		for i, level := range levels {
			if i == len(levels)-1 {
				current[level] = value
			} else {
				if _, ok := current[level]; !ok {
					current[level] = make(map[string]interface{})
				}
				if subMap, ok := current[level].(map[string]interface{}); ok {
					current = subMap
				}
			}
		}
	}

	return result
}

// getEnvFallback returns the value of the first non-empty environment variable.
func getEnvFallback(keys ...string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}
