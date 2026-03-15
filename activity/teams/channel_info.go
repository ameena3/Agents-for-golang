// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

// ChannelInfo describes a Teams channel.
type ChannelInfo struct {
	// ID is the unique identifier representing a channel.
	ID string `json:"id,omitempty"`
	// Name is the name of the channel.
	Name string `json:"name,omitempty"`
	// Type is the channel type.
	Type string `json:"type,omitempty"`
}

// TeamInfo describes a Teams team.
type TeamInfo struct {
	// ID is the unique identifier representing a team.
	ID string `json:"id,omitempty"`
	// Name is the name of the team.
	Name string `json:"name,omitempty"`
	// AADGroupID is the Azure AD Teams group ID.
	AADGroupID string `json:"aadGroupId,omitempty"`
}

// TenantInfo describes a tenant.
type TenantInfo struct {
	// ID is the unique identifier representing a tenant.
	ID string `json:"id,omitempty"`
}

// NotificationInfo specifies notification settings for a message.
type NotificationInfo struct {
	// Alert indicates whether a notification is to be sent.
	Alert *bool `json:"alert,omitempty"`
	// AlertInMeeting indicates whether notification is to be sent in a meeting context.
	AlertInMeeting *bool `json:"alertInMeeting,omitempty"`
	// ExternalResourceURL is the URL for external resources related to the notification.
	ExternalResourceURL string `json:"externalResourceUrl,omitempty"`
}

// TeamsMeetingInfo describes a Teams meeting.
type TeamsMeetingInfo struct {
	// ID is the unique identifier representing a meeting.
	ID string `json:"id,omitempty"`
}

// TeamsChannelDataSettings represents settings information for a Teams channel data.
type TeamsChannelDataSettings struct {
	// SelectedChannel is information about the selected Teams channel.
	SelectedChannel *ChannelInfo `json:"selectedChannel,omitempty"`
}

// OnBehalfOf specifies the OnBehalfOf entity for meeting notifications.
type OnBehalfOf struct {
	// ItemID is the item id of the OnBehalfOf entity.
	ItemID string `json:"itemId,omitempty"`
	// MentionType is the mention type. Default is "person".
	MentionType string `json:"mentionType,omitempty"`
	// DisplayName is the display name of the OnBehalfOf entity.
	DisplayName string `json:"displayName,omitempty"`
	// MRI is the MRI of the OnBehalfOf entity.
	MRI string `json:"mri,omitempty"`
}
