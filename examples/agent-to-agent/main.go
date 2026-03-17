// agent-to-agent demonstrates how one agent can call another agent
// using the Microsoft 365 Agents SDK proactive messaging and agentic request patterns.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ameena3/Agents-for-golang/activity"
	"github.com/ameena3/Agents-for-golang/hosting/core"
	"github.com/ameena3/Agents-for-golang/hosting/core/app"
	"github.com/ameena3/Agents-for-golang/hosting/nethttp"
)

// OrchestratorState tracks which skill agents we've called.
type OrchestratorState struct {
	LastSkillResult string `json:"lastSkillResult"`
}

func main() {
	// Orchestrator agent: receives messages, delegates to skill agents.
	orchestrator := app.New[OrchestratorState](app.AppOptions[OrchestratorState]{})

	orchestrator.OnMessage("", func(ctx context.Context, tc *core.TurnContext, state OrchestratorState) error {
		text := tc.Activity().Text

		// In a real implementation, this would call a skill agent via the connector.
		// For demo purposes, we simulate a skill response.
		skillResponse := fmt.Sprintf("[Skill processed]: %q", text)
		state.LastSkillResult = skillResponse

		_, err := tc.SendActivity(ctx, &activity.Activity{
			Type: activity.ActivityTypeMessage,
			Text: fmt.Sprintf("Orchestrator received your message.\nSkill says: %s", skillResponse),
		})
		return err
	})

	// Handle agentic requests (agent-to-agent activities).
	orchestrator.OnActivity(activity.ActivityTypeEvent, func(ctx context.Context, tc *core.TurnContext, state OrchestratorState) error {
		if tc.Activity().Name == "agent.response" {
			_, err := tc.SendActivity(ctx, core.Text("Received agent response: "+fmt.Sprintf("%v", tc.Activity().Value)))
			return err
		}
		return nil
	})

	port := 3978
	if p := os.Getenv("PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	log.Printf("Agent-to-agent orchestrator listening on port %d", port)
	if err := nethttp.StartAgentProcess(context.Background(), orchestrator, nethttp.ServerConfig{
		Port:                 port,
		AllowUnauthenticated: true,
	}); err != nil {
		log.Fatal(err)
	}
}
