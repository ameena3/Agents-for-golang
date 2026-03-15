// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package client provides a Direct-to-Engine client for Microsoft Copilot Studio agents.
// It enables Go applications to communicate directly with Copilot Studio
// (Power Virtual Agents) agents without going through the Bot Framework connector.
//
// Basic usage:
//
//	settings := &client.ConnectionSettings{
//	    EnvironmentID: os.Getenv("COPILOT_ENVIRONMENT_ID"),
//	    BotID:         os.Getenv("COPILOT_BOT_ID"),
//	    Cloud:         client.PowerPlatformCloudPublic,
//	}
//	c, err := client.NewCopilotClient(context.Background(), settings, tokenProvider)
//	resp, err := c.StartConversation(ctx)
//	turnResp, err := c.SendText(ctx, "Hello!")
package client
