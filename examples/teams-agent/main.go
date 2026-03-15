// teams-agent demonstrates Teams-specific activity handling including
// task modules, messaging extensions, and meeting events.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/microsoft/agents-sdk-go/activity"
	teamstypes "github.com/microsoft/agents-sdk-go/activity/teams"
	"github.com/microsoft/agents-sdk-go/hosting/core"
	"github.com/microsoft/agents-sdk-go/hosting/nethttp"
	hostingteams "github.com/microsoft/agents-sdk-go/hosting/teams"
)

// TeamsAgent handles Teams-specific activities.
type TeamsAgent struct {
	hostingteams.TeamsActivityHandler
}

// OnTurn dispatches to the Teams-aware handler pipeline.
func (a *TeamsAgent) OnTurn(ctx context.Context, tc *core.TurnContext) error {
	return a.TeamsActivityHandler.OnTurn(ctx, tc)
}

func (a *TeamsAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error {
	_, err := tc.SendActivity(ctx, core.Text("Teams Echo: "+tc.Activity().Text))
	return err
}

func (a *TeamsAgent) OnTeamsMembersAdded(ctx context.Context, members []*activity.ChannelAccount, tc *core.TurnContext) error {
	for _, member := range members {
		if member.ID != tc.Activity().Recipient.ID {
			if _, err := tc.SendActivity(ctx, core.Text(fmt.Sprintf("Welcome to Teams, %s!", member.Name))); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *TeamsAgent) OnTeamsTaskModuleFetch(ctx context.Context, req *teamstypes.TaskModuleRequest, tc *core.TurnContext) (*teamstypes.TaskModuleResponse, error) {
	return &teamstypes.TaskModuleResponse{
		Task: &teamstypes.TaskModuleContinueResponse{
			Type: "continue",
			Value: &teamstypes.TaskModuleTaskInfo{
				Title:  "Sample Task",
				Height: 200,
				Width:  400,
				URL:    "https://example.com/taskmodule",
			},
		},
	}, nil
}

func (a *TeamsAgent) OnTeamsTaskModuleSubmit(ctx context.Context, req *teamstypes.TaskModuleRequest, tc *core.TurnContext) (*teamstypes.TaskModuleResponse, error) {
	data, _ := json.Marshal(req.Data)
	_, err := tc.SendActivity(ctx, core.Text("Task submitted: "+string(data)))
	return nil, err
}

func main() {
	agent := &TeamsAgent{}
	adapter := nethttp.NewCloudAdapter(true) // allowUnauthenticated for local testing

	port := 3978
	if p := os.Getenv("PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, agent))
	log.Printf("Teams agent listening on :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
