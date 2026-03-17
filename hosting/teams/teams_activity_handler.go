// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

import (
	"context"

	"github.com/ameena3/Agents-for-golang/activity"
	"github.com/ameena3/Agents-for-golang/activity/teams"
	"github.com/ameena3/Agents-for-golang/hosting/core"
)

// TeamsActivityHandler extends ActivityHandler with Teams-specific invoke routing.
// Embed this in your Teams agent struct and override the OnTeams* methods you need.
//
// Example:
//
//	type MyTeamsAgent struct {
//	    teams.TeamsActivityHandler
//	}
//
//	func (a *MyTeamsAgent) OnTeamsTaskModuleFetch(ctx context.Context, request *teams.TaskModuleRequest, tc *core.TurnContext) (*teams.TaskModuleResponse, error) {
//	    // handle task/fetch
//	    return &teams.TaskModuleResponse{}, nil
//	}
type TeamsActivityHandler struct {
	core.ActivityHandler
}

// OnTurn overrides the base OnTurn to add Teams-specific routing before
// falling back to the base ActivityHandler.
func (h *TeamsActivityHandler) OnTurn(ctx context.Context, tc *core.TurnContext) error {
	return h.ActivityHandler.OnTurn(ctx, tc)
}

// OnInvokeActivity overrides the base invoke handler to route Teams-specific invoke names.
// It inspects the activity Name and dispatches to the appropriate Teams handler.
// Unrecognized invoke names fall back to the base ActivityHandler.OnInvokeActivity.
func (h *TeamsActivityHandler) OnInvokeActivity(ctx context.Context, tc *core.TurnContext) error {
	act := tc.Activity()

	// If there's no name and we're in Teams, route to card action invoke.
	if act.Name == "" && act.ChannelID == "msteams" {
		return h.OnTeamsCardActionInvoke(ctx, tc)
	}

	switch act.Name {
	case "task/fetch":
		var req teams.TaskModuleRequest
		if act.Value != nil {
			if m, ok := act.Value.(map[string]interface{}); ok {
				if data, hasData := m["data"]; hasData {
					req.Data = data
				}
			}
		}
		_, err := h.OnTeamsTaskModuleFetch(ctx, &req, tc)
		return err

	case "task/submit":
		var req teams.TaskModuleRequest
		if act.Value != nil {
			if m, ok := act.Value.(map[string]interface{}); ok {
				if data, hasData := m["data"]; hasData {
					req.Data = data
				}
			}
		}
		_, err := h.OnTeamsTaskModuleSubmit(ctx, &req, tc)
		return err

	case "config/fetch":
		_, err := h.OnTeamsConfigFetch(ctx, act.Value, tc)
		return err

	case "config/submit":
		_, err := h.OnTeamsConfigSubmit(ctx, act.Value, tc)
		return err

	case "composeExtension/query":
		var query teams.MessagingExtensionQuery
		if act.Value != nil {
			if m, ok := act.Value.(map[string]interface{}); ok {
				if cmdID, ok := m["commandId"].(string); ok {
					query.CommandID = cmdID
				}
			}
		}
		return h.OnTeamsMessagingExtensionQuery(ctx, &query, tc)

	case "composeExtension/selectItem":
		return h.OnTeamsMessagingExtensionSelectItem(ctx, act.Value, tc)

	case "composeExtension/submitAction":
		return h.OnTeamsMessagingExtensionSubmitAction(ctx, act.Value, tc)

	case "composeExtension/fetchTask":
		return h.OnTeamsMessagingExtensionFetchTask(ctx, act.Value, tc)

	case "actionableMessage/executeAction":
		return h.OnTeamsCardActionInvoke(ctx, tc)

	default:
		return h.ActivityHandler.OnInvokeActivity(ctx, tc)
	}
}

// OnConversationUpdateActivity overrides the base to add Teams-specific routing
// for Teams channel events (channelCreated, channelDeleted, teamRenamed, etc.)
// and members added/removed.
func (h *TeamsActivityHandler) OnConversationUpdateActivity(ctx context.Context, tc *core.TurnContext) error {
	act := tc.Activity()

	if act.ChannelID == "msteams" {
		// Get channel data event type from the raw channel data map.
		var eventType string
		var channelInfo *teams.ChannelInfo
		var teamInfo *teams.TeamInfo
		if m, ok := act.ChannelData.(map[string]interface{}); ok {
			if et, ok := m["eventType"].(string); ok {
				eventType = et
			}
			if ch, ok := m["channel"].(map[string]interface{}); ok {
				channelInfo = &teams.ChannelInfo{}
				if id, ok := ch["id"].(string); ok {
					channelInfo.ID = id
				}
				if name, ok := ch["name"].(string); ok {
					channelInfo.Name = name
				}
			}
			if t, ok := m["team"].(map[string]interface{}); ok {
				teamInfo = &teams.TeamInfo{}
				if id, ok := t["id"].(string); ok {
					teamInfo.ID = id
				}
				if name, ok := t["name"].(string); ok {
					teamInfo.Name = name
				}
			}
		}

		if len(act.MembersAdded) > 0 {
			return h.OnTeamsMembersAdded(ctx, act.MembersAdded, tc)
		}

		if len(act.MembersRemoved) > 0 {
			return h.OnTeamsMembersRemoved(ctx, act.MembersRemoved, tc)
		}

		switch eventType {
		case "channelCreated":
			return h.OnTeamsChannelCreated(ctx, channelInfo, teamInfo, tc)
		case "channelDeleted":
			return h.OnTeamsChannelDeleted(ctx, channelInfo, teamInfo, tc)
		case "teamRenamed":
			return h.OnTeamsTeamRenamed(ctx, teamInfo, tc)
		}
	}

	return h.ActivityHandler.OnConversationUpdateActivity(ctx, tc)
}

// OnEventActivity overrides the base to route Teams-specific event names
// (meetingStart, meetingEnd) before falling back to the base handler.
func (h *TeamsActivityHandler) OnEventActivity(ctx context.Context, tc *core.TurnContext) error {
	act := tc.Activity()

	if act.ChannelID == "msteams" {
		switch act.Name {
		case "application/vnd.microsoft.meetingStart":
			var meeting teams.MeetingStartEventDetails
			if m, ok := act.Value.(map[string]interface{}); ok {
				if st, ok := m["startTime"].(string); ok {
					meeting.StartTime = st
				}
			}
			return h.OnTeamsMeetingStart(ctx, &meeting, tc)

		case "application/vnd.microsoft.meetingEnd":
			var meeting teams.MeetingEndEventDetails
			if m, ok := act.Value.(map[string]interface{}); ok {
				if et, ok := m["endTime"].(string); ok {
					meeting.EndTime = et
				}
			}
			return h.OnTeamsMeetingEnd(ctx, &meeting, tc)
		}
	}

	return h.ActivityHandler.OnEventActivity(ctx, tc)
}

// --- Teams Card Actions ---

// OnTeamsCardActionInvoke handles Teams card action invoke activities.
// Default implementation is a no-op; override to process card actions.
func (h *TeamsActivityHandler) OnTeamsCardActionInvoke(ctx context.Context, tc *core.TurnContext) error {
	return nil
}

// --- Messaging Extension ---

// OnTeamsMessagingExtensionQuery handles messaging extension query.
// Default implementation is a no-op; override to return search results.
func (h *TeamsActivityHandler) OnTeamsMessagingExtensionQuery(ctx context.Context, query *teams.MessagingExtensionQuery, tc *core.TurnContext) error {
	return nil
}

// OnTeamsMessagingExtensionSelectItem handles messaging extension item selection.
// Default implementation is a no-op; override to handle the selected item.
func (h *TeamsActivityHandler) OnTeamsMessagingExtensionSelectItem(ctx context.Context, query interface{}, tc *core.TurnContext) error {
	return nil
}

// OnTeamsMessagingExtensionSubmitAction handles messaging extension submit action.
// Default implementation is a no-op; override to handle the submitted action.
func (h *TeamsActivityHandler) OnTeamsMessagingExtensionSubmitAction(ctx context.Context, action interface{}, tc *core.TurnContext) error {
	return nil
}

// OnTeamsMessagingExtensionFetchTask handles messaging extension fetch task.
// Default implementation is a no-op; override to return a task module response.
func (h *TeamsActivityHandler) OnTeamsMessagingExtensionFetchTask(ctx context.Context, action interface{}, tc *core.TurnContext) error {
	return nil
}

// --- Task Module ---

// OnTeamsTaskModuleFetch handles task module fetch.
// Default implementation returns nil; override to return a TaskModuleResponse.
func (h *TeamsActivityHandler) OnTeamsTaskModuleFetch(ctx context.Context, request *teams.TaskModuleRequest, tc *core.TurnContext) (*teams.TaskModuleResponse, error) {
	return nil, nil
}

// OnTeamsTaskModuleSubmit handles task module submit.
// Default implementation returns nil; override to process the submitted data.
func (h *TeamsActivityHandler) OnTeamsTaskModuleSubmit(ctx context.Context, request *teams.TaskModuleRequest, tc *core.TurnContext) (*teams.TaskModuleResponse, error) {
	return nil, nil
}

// --- Config ---

// OnTeamsConfigFetch handles Teams config fetch.
// Default implementation returns nil; override to return a ConfigResponse.
func (h *TeamsActivityHandler) OnTeamsConfigFetch(ctx context.Context, query interface{}, tc *core.TurnContext) (*teams.ConfigResponse, error) {
	return nil, nil
}

// OnTeamsConfigSubmit handles Teams config submit.
// Default implementation returns nil; override to process config data.
func (h *TeamsActivityHandler) OnTeamsConfigSubmit(ctx context.Context, query interface{}, tc *core.TurnContext) (*teams.ConfigResponse, error) {
	return nil, nil
}

// --- Meeting Events ---

// OnTeamsMeetingStart handles meeting start events.
// Default implementation is a no-op; override to respond to meeting starts.
func (h *TeamsActivityHandler) OnTeamsMeetingStart(ctx context.Context, meeting *teams.MeetingStartEventDetails, tc *core.TurnContext) error {
	return nil
}

// OnTeamsMeetingEnd handles meeting end events.
// Default implementation is a no-op; override to respond to meeting ends.
func (h *TeamsActivityHandler) OnTeamsMeetingEnd(ctx context.Context, meeting *teams.MeetingEndEventDetails, tc *core.TurnContext) error {
	return nil
}

// --- Conversation Update ---

// OnTeamsMembersAdded handles Teams members added (extends base OnMembersAdded).
// Default implementation falls back to the base ActivityHandler.OnMembersAdded.
func (h *TeamsActivityHandler) OnTeamsMembersAdded(ctx context.Context, membersAdded []*activity.ChannelAccount, tc *core.TurnContext) error {
	return h.ActivityHandler.OnMembersAdded(ctx, membersAdded, tc)
}

// OnTeamsMembersRemoved handles Teams members removed.
// Default implementation falls back to the base ActivityHandler.OnMembersRemoved.
func (h *TeamsActivityHandler) OnTeamsMembersRemoved(ctx context.Context, membersRemoved []*activity.ChannelAccount, tc *core.TurnContext) error {
	return h.ActivityHandler.OnMembersRemoved(ctx, membersRemoved, tc)
}

// OnTeamsChannelCreated handles channel created events.
// Default implementation is a no-op; override to respond to channel creation.
func (h *TeamsActivityHandler) OnTeamsChannelCreated(ctx context.Context, channelInfo *teams.ChannelInfo, teamInfo *teams.TeamInfo, tc *core.TurnContext) error {
	return nil
}

// OnTeamsChannelDeleted handles channel deleted events.
// Default implementation is a no-op; override to respond to channel deletion.
func (h *TeamsActivityHandler) OnTeamsChannelDeleted(ctx context.Context, channelInfo *teams.ChannelInfo, teamInfo *teams.TeamInfo, tc *core.TurnContext) error {
	return nil
}

// OnTeamsTeamRenamed handles team renamed events.
// Default implementation is a no-op; override to respond to team renames.
func (h *TeamsActivityHandler) OnTeamsTeamRenamed(ctx context.Context, teamInfo *teams.TeamInfo, tc *core.TurnContext) error {
	return nil
}
