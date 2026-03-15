// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package client

import "fmt"

// ConnectionSettings holds the configuration for connecting to a Copilot Studio agent.
type ConnectionSettings struct {
	// EnvironmentID is the Power Platform environment ID.
	EnvironmentID string
	// BotID is the bot identifier in Copilot Studio.
	BotID string
	// BotName is the display name of the bot (optional).
	BotName string
	// Cloud is the Power Platform cloud. Defaults to Public.
	Cloud PowerPlatformCloud
	// AgentType is the agent type (Published or Preview). Defaults to Published.
	AgentType AgentType
	// CustomEndpoint overrides the default Direct-to-Engine endpoint URL.
	CustomEndpoint string
}

// GetEndpointURL returns the Direct-to-Engine conversations endpoint URL.
func (s *ConnectionSettings) GetEndpointURL() string {
	if s.CustomEndpoint != "" {
		return s.CustomEndpoint
	}
	env := &PowerPlatformEnvironment{EnvironmentID: s.EnvironmentID, Cloud: s.Cloud}
	base := env.GetDirectToEngineURL()
	return fmt.Sprintf("%s/powervirtualagents/environments/%s/bots/%s/directline",
		base, s.EnvironmentID, s.BotID)
}
