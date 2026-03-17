// echo-agent is a simple Microsoft 365 agent that echoes back any message.
// Run it with:
//
//	go run ./examples/echo-agent/
//
// Then use the M365 Agents Playground or Azure Bot Service to send messages.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ameena3/Agents-for-golang/activity/config"
	"github.com/ameena3/Agents-for-golang/hosting/core"
	"github.com/ameena3/Agents-for-golang/hosting/core/app"
	"github.com/ameena3/Agents-for-golang/hosting/core/storage"
	"github.com/ameena3/Agents-for-golang/hosting/nethttp"
)

// AppState holds per-conversation state for this agent.
type AppState struct {
	MessageCount int `json:"messageCount"`
}

func main() {
	// Load configuration from environment variables.
	cfg := config.LoadFromEnv()
	_ = cfg

	// Set up in-memory storage (replace with BlobStorage for production).
	store := storage.NewMemoryStorage()

	// Create the agent application.
	agentApp := app.New[AppState](app.AppOptions[AppState]{
		Storage: store,
	})

	// Handle members joining the conversation.
	agentApp.OnMembersAdded(func(ctx context.Context, tc *core.TurnContext, state AppState) error {
		for _, member := range tc.Activity().MembersAdded {
			if member.ID != tc.Activity().Recipient.ID {
				_, err := tc.SendActivity(ctx, core.Text("Hello! I'm the Echo Agent. Send me a message and I'll echo it back."))
				return err
			}
		}
		return nil
	})

	// Echo all messages back.
	agentApp.OnMessage("", func(ctx context.Context, tc *core.TurnContext, state AppState) error {
		state.MessageCount++
		text := tc.Activity().Text
		reply := fmt.Sprintf("Echo (#%d): %s", state.MessageCount, text)
		_, err := tc.SendActivity(ctx, core.Text(reply))
		return err
	})

	// Determine port.
	port := 3978
	if p := os.Getenv("PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	log.Printf("Echo agent starting on port %d...", port)
	if err := nethttp.StartAgentProcess(context.Background(), agentApp, nethttp.ServerConfig{
		Port:                 port,
		AllowUnauthenticated: true, // Set to false in production with real auth.
	}); err != nil {
		log.Fatal(err)
	}
}
