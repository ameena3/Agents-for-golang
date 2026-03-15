// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

import (
	"encoding/json"
	"testing"
)

func TestTeamsChannelAccountFields(t *testing.T) {
	acct := &TeamsChannelAccount{
		ID:                "user-id-1",
		Name:              "Jane Doe",
		GivenName:         "Jane",
		Surname:           "Doe",
		Email:             "jane@example.com",
		UserPrincipalName: "jane@contoso.com",
		TenantID:          "tenant-123",
		UserRole:          "user",
	}
	if acct.ID != "user-id-1" {
		t.Errorf("ID: got %q, want %q", acct.ID, "user-id-1")
	}
	if acct.GivenName != "Jane" {
		t.Errorf("GivenName: got %q, want %q", acct.GivenName, "Jane")
	}
	if acct.Surname != "Doe" {
		t.Errorf("Surname: got %q, want %q", acct.Surname, "Doe")
	}
	if acct.TenantID != "tenant-123" {
		t.Errorf("TenantID: got %q, want %q", acct.TenantID, "tenant-123")
	}
}

func TestTeamsChannelAccountJSONRoundTrip(t *testing.T) {
	original := &TeamsChannelAccount{
		ID:                "id-abc",
		Name:              "Test User",
		GivenName:         "Test",
		Surname:           "User",
		Email:             "test@example.com",
		UserPrincipalName: "test@contoso.onmicrosoft.com",
		TenantID:          "tid-xyz",
		UserRole:          "owner",
		Properties:        map[string]interface{}{"extra": "value"},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got TeamsChannelAccount
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.ID != original.ID {
		t.Errorf("ID: got %q, want %q", got.ID, original.ID)
	}
	if got.Email != original.Email {
		t.Errorf("Email: got %q, want %q", got.Email, original.Email)
	}
	if got.UserRole != original.UserRole {
		t.Errorf("UserRole: got %q, want %q", got.UserRole, original.UserRole)
	}
	if got.Properties == nil {
		t.Error("Properties should not be nil after round-trip")
	}
}

func TestTeamsChannelAccountJSONFieldNames(t *testing.T) {
	acct := &TeamsChannelAccount{
		ID:                "id-1",
		Name:              "Tester",
		GivenName:         "T",
		Surname:           "S",
		Email:             "t@example.com",
		UserPrincipalName: "t@upn.com",
		TenantID:          "tid",
		UserRole:          "user",
	}
	data, err := json.Marshal(acct)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	expectedKeys := []string{"id", "name", "givenName", "surname", "email", "userPrincipalName", "tenantId", "userRole"}
	for _, key := range expectedKeys {
		if _, ok := m[key]; !ok {
			t.Errorf("expected JSON field %q to be present", key)
		}
	}
}

func TestTeamsChannelDataJSONRoundTrip(t *testing.T) {
	alert := true
	cd := &TeamsChannelData{
		Channel:   &ChannelInfo{ID: "chan-1", Name: "General", Type: "standard"},
		EventType: "teamMemberAdded",
		Team:      &TeamInfo{ID: "team-1", Name: "My Team", AADGroupID: "aad-group-1"},
		Notification: &NotificationInfo{
			Alert:               &alert,
			ExternalResourceURL: "https://example.com",
		},
		Tenant:  &TenantInfo{ID: "tenant-1"},
		Meeting: &TeamsMeetingInfo{ID: "meeting-1"},
		Settings: &TeamsChannelDataSettings{
			SelectedChannel: &ChannelInfo{ID: "sel-chan", Name: "Selected"},
		},
		OnBehalfOf: []*OnBehalfOf{
			{ItemID: "item-1", MentionType: "person", DisplayName: "Bot", MRI: "mri-1"},
		},
	}

	data, err := json.Marshal(cd)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got TeamsChannelData
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.EventType != "teamMemberAdded" {
		t.Errorf("EventType: got %q, want %q", got.EventType, "teamMemberAdded")
	}
	if got.Channel == nil || got.Channel.ID != "chan-1" {
		t.Errorf("Channel.ID: got %v", got.Channel)
	}
	if got.Team == nil || got.Team.Name != "My Team" {
		t.Errorf("Team.Name: got %v", got.Team)
	}
	if got.Tenant == nil || got.Tenant.ID != "tenant-1" {
		t.Errorf("Tenant.ID: got %v", got.Tenant)
	}
	if got.Meeting == nil || got.Meeting.ID != "meeting-1" {
		t.Errorf("Meeting.ID: got %v", got.Meeting)
	}
	if got.Settings == nil || got.Settings.SelectedChannel == nil || got.Settings.SelectedChannel.ID != "sel-chan" {
		t.Errorf("Settings.SelectedChannel: got %v", got.Settings)
	}
	if len(got.OnBehalfOf) != 1 || got.OnBehalfOf[0].DisplayName != "Bot" {
		t.Errorf("OnBehalfOf: got %v", got.OnBehalfOf)
	}
}

func TestTeamsChannelDataOmitsEmptyFields(t *testing.T) {
	cd := &TeamsChannelData{EventType: "channelCreated"}
	data, err := json.Marshal(cd)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	for _, key := range []string{"channel", "team", "notification", "tenant", "meeting", "settings", "onBehalfOf"} {
		if _, ok := m[key]; ok {
			t.Errorf("expected %q to be omitted but was present", key)
		}
	}
}

func TestChannelInfoRoundTrip(t *testing.T) {
	info := &ChannelInfo{ID: "c-1", Name: "Lobby", Type: "general"}
	data, _ := json.Marshal(info)
	var got ChannelInfo
	json.Unmarshal(data, &got)
	if got.ID != "c-1" || got.Name != "Lobby" || got.Type != "general" {
		t.Errorf("ChannelInfo round-trip failed: %+v", got)
	}
}

func TestTeamInfoAADGroupID(t *testing.T) {
	team := &TeamInfo{ID: "t-1", Name: "Engineering", AADGroupID: "aad-123"}
	data, _ := json.Marshal(team)
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	if m["aadGroupId"] != "aad-123" {
		t.Errorf("expected aadGroupId to be %q, got %v", "aad-123", m["aadGroupId"])
	}
}

func TestNotificationInfoAlert(t *testing.T) {
	tr := true
	n := &NotificationInfo{Alert: &tr, ExternalResourceURL: "https://example.com/notify"}
	data, _ := json.Marshal(n)
	var got NotificationInfo
	json.Unmarshal(data, &got)
	if got.Alert == nil || *got.Alert != true {
		t.Errorf("Alert not round-tripped correctly: %v", got.Alert)
	}
	if got.ExternalResourceURL != "https://example.com/notify" {
		t.Errorf("ExternalResourceURL: got %q", got.ExternalResourceURL)
	}
}
