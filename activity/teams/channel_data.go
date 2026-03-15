// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

// TeamsChannelData contains channel data specific to messages received in Microsoft Teams.
type TeamsChannelData struct {
	// Channel is information about the channel in which the message was sent.
	Channel *ChannelInfo `json:"channel,omitempty"`
	// EventType is the type of event.
	EventType string `json:"eventType,omitempty"`
	// Team is information about the team in which the message was sent.
	Team *TeamInfo `json:"team,omitempty"`
	// Notification contains notification settings for the message.
	Notification *NotificationInfo `json:"notification,omitempty"`
	// Tenant is information about the tenant in which the message was sent.
	Tenant *TenantInfo `json:"tenant,omitempty"`
	// Meeting is information about the meeting in which the message was sent.
	Meeting *TeamsMeetingInfo `json:"meeting,omitempty"`
	// Settings contains information about the settings.
	Settings *TeamsChannelDataSettings `json:"settings,omitempty"`
	// OnBehalfOf is the OnBehalfOf list for user attribution.
	OnBehalfOf []*OnBehalfOf `json:"onBehalfOf,omitempty"`
}
