// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams_test

import (
	"context"
	"testing"

	"github.com/microsoft/agents-sdk-go/activity"
	"github.com/microsoft/agents-sdk-go/hosting/core"
	"github.com/microsoft/agents-sdk-go/hosting/core/authorization"
	"github.com/microsoft/agents-sdk-go/hosting/teams"
)

// --- stub adapter ---

type stubAdapter struct{}

func (s *stubAdapter) SendActivities(_ context.Context, _ *core.TurnContext, acts []*activity.Activity) ([]*activity.ResourceResponse, error) {
	resp := make([]*activity.ResourceResponse, len(acts))
	for i := range acts {
		resp[i] = &activity.ResourceResponse{ID: "stub-resp"}
	}
	return resp, nil
}

func (s *stubAdapter) UpdateActivity(_ context.Context, _ *core.TurnContext, act *activity.Activity) (*activity.ResourceResponse, error) {
	return &activity.ResourceResponse{ID: act.ID}, nil
}

func (s *stubAdapter) DeleteActivity(_ context.Context, _ *core.TurnContext, _ string) error {
	return nil
}

func (s *stubAdapter) ContinueConversation(_ context.Context, _ *activity.ConversationReference, _ func(context.Context, *core.TurnContext) error) error {
	return nil
}

func (s *stubAdapter) Use(_ ...core.Middleware) {}

// --- helpers ---

func newTurnContext(actType, channelID, name string) *core.TurnContext {
	act := &activity.Activity{
		Type:      actType,
		ChannelID: channelID,
		Name:      name,
		ID:        "test-activity-id",
		Conversation: &activity.ConversationAccount{
			ID: "conv-1",
		},
		From:      &activity.ChannelAccount{ID: "user-1", Name: "User"},
		Recipient: &activity.ChannelAccount{ID: "agent-1", Name: "Agent"},
	}
	identity := authorization.NewClaimsIdentity(true, "Bearer", map[string]string{"aud": "app-id"})
	return core.NewTurnContext(&stubAdapter{}, act, identity)
}

// --- TeamsActivityHandler tests ---

func TestTeamsActivityHandler_Creation(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	if handler == nil {
		t.Fatal("expected non-nil TeamsActivityHandler")
	}
}

func TestTeamsActivityHandler_OnTurn_MessageActivity(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	// Should not panic and return nil (default no-op)
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_InvokeActivity_TaskFetch(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeInvoke, "msteams", "task/fetch")
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for task/fetch invoke: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_InvokeActivity_TaskSubmit(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeInvoke, "msteams", "task/submit")
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for task/submit invoke: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_InvokeActivity_ConfigFetch(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeInvoke, "msteams", "config/fetch")
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for config/fetch invoke: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_InvokeActivity_ConfigSubmit(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeInvoke, "msteams", "config/submit")
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for config/submit invoke: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_InvokeActivity_ComposeExtensionQuery(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeInvoke, "msteams", "composeExtension/query")
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for composeExtension/query: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_InvokeActivity_EmptyName_Teams(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeInvoke, "msteams", "")
	// Empty name on msteams routes to OnTeamsCardActionInvoke (no-op by default)
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for empty-name Teams invoke: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_InvokeActivity_Unknown(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeInvoke, "msteams", "some/unknown/invoke")
	// Falls back to base ActivityHandler.OnInvokeActivity (no-op)
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for unknown invoke: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_ConversationUpdate_MembersAdded(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeConversationUpdate, "msteams", "")
	tc.Activity().MembersAdded = []*activity.ChannelAccount{
		{ID: "new-member", Name: "New Member"},
	}
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for membersAdded: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_ConversationUpdate_MembersRemoved(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeConversationUpdate, "msteams", "")
	tc.Activity().MembersRemoved = []*activity.ChannelAccount{
		{ID: "old-member", Name: "Old Member"},
	}
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for membersRemoved: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_ConversationUpdate_ChannelCreated(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeConversationUpdate, "msteams", "")
	tc.Activity().ChannelData = map[string]interface{}{
		"eventType": "channelCreated",
		"channel":   map[string]interface{}{"id": "ch-1", "name": "General"},
		"team":      map[string]interface{}{"id": "team-1", "name": "Team"},
	}
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for channelCreated: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_ConversationUpdate_ChannelDeleted(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeConversationUpdate, "msteams", "")
	tc.Activity().ChannelData = map[string]interface{}{
		"eventType": "channelDeleted",
		"channel":   map[string]interface{}{"id": "ch-1", "name": "General"},
		"team":      map[string]interface{}{"id": "team-1", "name": "Team"},
	}
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for channelDeleted: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_ConversationUpdate_TeamRenamed(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeConversationUpdate, "msteams", "")
	tc.Activity().ChannelData = map[string]interface{}{
		"eventType": "teamRenamed",
		"team":      map[string]interface{}{"id": "team-1", "name": "NewName"},
	}
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for teamRenamed: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_EventActivity_MeetingStart(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeEvent, "msteams", "application/vnd.microsoft.meetingStart")
	tc.Activity().Value = map[string]interface{}{"startTime": "2024-01-01T10:00:00Z"}
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for meetingStart: %v", err)
	}
}

func TestTeamsActivityHandler_OnTurn_EventActivity_MeetingEnd(t *testing.T) {
	handler := &teams.TeamsActivityHandler{}
	tc := newTurnContext(activity.ActivityTypeEvent, "msteams", "application/vnd.microsoft.meetingEnd")
	tc.Activity().Value = map[string]interface{}{"endTime": "2024-01-01T11:00:00Z"}
	err := handler.OnTurn(context.Background(), tc)
	if err != nil {
		t.Fatalf("unexpected error for meetingEnd: %v", err)
	}
}

// --- GetPagedMembers tests (no HTTP calls needed) ---

func TestGetPagedMembers_ErrorWhenNoConnectorClient(t *testing.T) {
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	// TurnState has no ConnectorClient — should return an error.
	_, _, err := teams.GetPagedMembers(context.Background(), tc, 10, "")
	if err == nil {
		t.Fatal("expected error when ConnectorClient is missing from TurnState")
	}
}

func TestGetMembers_ErrorWhenNoConnectorClient(t *testing.T) {
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	_, err := teams.GetMembers(context.Background(), tc)
	if err == nil {
		t.Fatal("expected error when ConnectorClient is missing from TurnState")
	}
}

func TestGetMember_ErrorWhenUserIDEmpty(t *testing.T) {
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	_, err := teams.GetMember(context.Background(), tc, "")
	if err == nil {
		t.Fatal("expected error when userID is empty")
	}
}

func TestGetTeamDetails_ErrorWhenTeamIDEmpty(t *testing.T) {
	// Activity has no channel data, so teamID extraction will fail.
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	_, err := teams.GetTeamDetails(context.Background(), tc, "")
	if err == nil {
		t.Fatal("expected error when teamID is empty and not in channel data")
	}
}

func TestGetTeamChannels_ErrorWhenTeamIDEmpty(t *testing.T) {
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	_, err := teams.GetTeamChannels(context.Background(), tc, "")
	if err == nil {
		t.Fatal("expected error when teamID is empty and not in channel data")
	}
}

// --- SendNotificationToUser tests ---

func TestSendNotificationToUser_ErrorWhenUserIDEmpty(t *testing.T) {
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	err := teams.SendNotificationToUser(context.Background(), tc, "", map[string]interface{}{"text": "hello"})
	if err == nil {
		t.Fatal("expected error when userID is empty")
	}
}

func TestSendNotificationToUser_ErrorWhenNotificationNil(t *testing.T) {
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	err := teams.SendNotificationToUser(context.Background(), tc, "user-1", nil)
	if err == nil {
		t.Fatal("expected error when notification is nil")
	}
}

func TestSendNotificationToUser_SuccessWithValidArgs(t *testing.T) {
	tc := newTurnContext(activity.ActivityTypeMessage, "msteams", "")
	// Has a valid stubAdapter that returns stub responses.
	err := teams.SendNotificationToUser(context.Background(), tc, "user-1", map[string]interface{}{
		"text": "Hello!",
		"type": "message",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
