// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

import (
	"context"
	"fmt"

	"github.com/microsoft/agents-sdk-go/activity"
	"github.com/microsoft/agents-sdk-go/hosting/core"
	"github.com/microsoft/agents-sdk-go/hosting/core/connector"
)

// TeamsInfo provides utility methods for interacting with the Teams APIs.
// All methods are package-level functions that take a TurnContext for obtaining
// the Teams connector client stored in turn state.
type TeamsInfo struct{}

// getTeamsConnectorClient retrieves a TeamsConnectorClient from the TurnContext's turn state.
// The connector client is expected to be stored under the "ConnectorClient" key.
func getTeamsConnectorClient(tc *core.TurnContext) (*connector.TeamsConnectorClient, error) {
	ts := tc.TurnState()
	if ts == nil {
		return nil, fmt.Errorf("teams_info: turn state is nil")
	}
	raw, ok := ts["ConnectorClient"]
	if !ok || raw == nil {
		return nil, fmt.Errorf("teams_info: TeamsConnectorClient is not available in the turn state")
	}
	client, ok := raw.(*connector.TeamsConnectorClient)
	if !ok {
		return nil, fmt.Errorf("teams_info: ConnectorClient in turn state is not a *TeamsConnectorClient")
	}
	return client, nil
}

// GetTeamDetails retrieves team details for the given team ID.
// If teamID is empty, it attempts to extract the team ID from the activity's channel data.
func GetTeamDetails(ctx context.Context, tc *core.TurnContext, teamID string) (map[string]interface{}, error) {
	if teamID == "" {
		act := tc.Activity()
		if m, ok := act.ChannelData.(map[string]interface{}); ok {
			if t, ok := m["team"].(map[string]interface{}); ok {
				if id, ok := t["id"].(string); ok {
					teamID = id
				}
			}
		}
	}
	if teamID == "" {
		return nil, fmt.Errorf("teams_info: teamID is required")
	}
	client, err := getTeamsConnectorClient(tc)
	if err != nil {
		return nil, err
	}
	return client.GetTeamDetails(ctx, teamID)
}

// GetTeamChannels retrieves all channels in a team.
// If teamID is empty, it attempts to extract the team ID from the activity's channel data.
func GetTeamChannels(ctx context.Context, tc *core.TurnContext, teamID string) ([]map[string]interface{}, error) {
	if teamID == "" {
		act := tc.Activity()
		if m, ok := act.ChannelData.(map[string]interface{}); ok {
			if t, ok := m["team"].(map[string]interface{}); ok {
				if id, ok := t["id"].(string); ok {
					teamID = id
				}
			}
		}
	}
	if teamID == "" {
		return nil, fmt.Errorf("teams_info: teamID is required")
	}
	client, err := getTeamsConnectorClient(tc)
	if err != nil {
		return nil, err
	}
	return client.GetChannels(ctx, teamID)
}

// GetMember retrieves a single member of a conversation by user ID.
func GetMember(ctx context.Context, tc *core.TurnContext, userID string) (map[string]interface{}, error) {
	if userID == "" {
		return nil, fmt.Errorf("teams_info: userID is required")
	}

	act := tc.Activity()
	conversationID := ""
	if act.Conversation != nil {
		conversationID = act.Conversation.ID
	}
	if conversationID == "" {
		return nil, fmt.Errorf("teams_info: conversation ID is required")
	}

	client, err := getTeamsConnectorClient(tc)
	if err != nil {
		return nil, err
	}

	members, err := client.GetMembers(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		if id, ok := member["id"].(string); ok && id == userID {
			return member, nil
		}
	}
	return nil, fmt.Errorf("teams_info: member %q not found in conversation %q", userID, conversationID)
}

// GetMembers retrieves all members of the current conversation.
func GetMembers(ctx context.Context, tc *core.TurnContext) ([]map[string]interface{}, error) {
	act := tc.Activity()
	conversationID := ""
	if act.Conversation != nil {
		conversationID = act.Conversation.ID
	}
	if conversationID == "" {
		return nil, fmt.Errorf("teams_info: conversation ID is required")
	}

	client, err := getTeamsConnectorClient(tc)
	if err != nil {
		return nil, err
	}
	return client.GetMembers(ctx, conversationID)
}

// GetPagedMembers retrieves members in pages.
// Returns the members slice, a continuation token (empty string if no more pages), and an error.
// When pageSize <= 0 or pageSize >= total members, all members are returned in one page.
// The continuationToken is treated as a start index encoded as a decimal integer string;
// on the first call pass an empty string.
func GetPagedMembers(ctx context.Context, tc *core.TurnContext, pageSize int, continuationToken string) ([]map[string]interface{}, string, error) {
	all, err := GetMembers(ctx, tc)
	if err != nil {
		return nil, "", err
	}

	// If no page size requested or fits in one page, return everything.
	if pageSize <= 0 || pageSize >= len(all) {
		return all, "", nil
	}

	// Parse the start offset from continuationToken.
	start := 0
	if continuationToken != "" {
		n, parseErr := fmt.Sscanf(continuationToken, "%d", &start)
		if parseErr != nil || n != 1 || start < 0 || start >= len(all) {
			start = 0
		}
	}

	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}

	nextToken := ""
	if end < len(all) {
		nextToken = fmt.Sprintf("%d", end)
	}

	return all[start:end], nextToken, nil
}

// SendNotificationToUser sends a notification activity to a specific user in the current conversation.
// The activity Recipient is set to the given userID before sending.
func SendNotificationToUser(ctx context.Context, tc *core.TurnContext, userID string, notification map[string]interface{}) error {
	if userID == "" {
		return fmt.Errorf("teams_info: userID is required")
	}
	if notification == nil {
		return fmt.Errorf("teams_info: notification is required")
	}

	// Build a message activity from the notification map.
	act := &activity.Activity{}
	if text, ok := notification["text"].(string); ok {
		act.Text = text
	}
	if t, ok := notification["type"].(string); ok {
		act.Type = t
	}
	if act.Type == "" {
		act.Type = activity.ActivityTypeMessage
	}
	act.Recipient = &activity.ChannelAccount{ID: userID}

	_, err := tc.SendActivity(ctx, act)
	return err
}
