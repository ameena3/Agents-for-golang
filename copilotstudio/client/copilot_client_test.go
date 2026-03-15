// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package client

import (
	"context"
	"testing"
)

func TestNewCopilotClient_NilSettings(t *testing.T) {
	_, err := NewCopilotClient(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error for nil settings")
	}
}

func TestNewCopilotClient_DefaultsApplied(t *testing.T) {
	c, err := NewCopilotClient(context.Background(), &ConnectionSettings{
		EnvironmentID: "env123",
		BotID:         "bot123",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.settings.Cloud != PowerPlatformCloudPublic {
		t.Errorf("expected default cloud Public, got %v", c.settings.Cloud)
	}
	if c.settings.AgentType != AgentTypePublished {
		t.Errorf("expected default AgentType Published, got %v", c.settings.AgentType)
	}
}

func TestGetEndpointURL_CustomOverride(t *testing.T) {
	s := &ConnectionSettings{
		CustomEndpoint: "https://my.custom.endpoint",
	}
	if got := s.GetEndpointURL(); got != "https://my.custom.endpoint" {
		t.Errorf("expected custom endpoint, got %q", got)
	}
}

func TestGetEndpointURL_Public(t *testing.T) {
	s := &ConnectionSettings{
		Cloud:         PowerPlatformCloudPublic,
		EnvironmentID: "env1",
		BotID:         "bot1",
	}
	url := s.GetEndpointURL()
	if url == "" {
		t.Fatal("expected non-empty URL")
	}
}

func TestExecuteTurn_NoConversation(t *testing.T) {
	c, _ := NewCopilotClient(context.Background(), &ConnectionSettings{BotID: "b"}, nil)
	_, err := c.SendText(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected error when no conversation started")
	}
}
