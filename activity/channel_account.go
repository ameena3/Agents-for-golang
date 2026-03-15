// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

// ChannelAccount contains channel account information needed to route a message.
type ChannelAccount struct {
	// ID is the channel id for the user or agent on this channel.
	ID string `json:"id,omitempty"`
	// Name is the display friendly name.
	Name string `json:"name,omitempty"`
	// AADObjectID is this account's object ID within Azure Active Directory.
	AADObjectID string `json:"aadObjectId,omitempty"`
	// Role is the role of the entity behind the account.
	Role string `json:"role,omitempty"`
	// AgenticUserID is the agentic user ID for agentic requests.
	AgenticUserID string `json:"agenticUserId,omitempty"`
	// AgenticAppID is the agentic application instance ID.
	AgenticAppID string `json:"agenticAppId,omitempty"`
	// TenantID is the tenant ID for this account.
	TenantID string `json:"tenantId,omitempty"`
	// Properties contains additional channel-specific properties.
	Properties map[string]interface{} `json:"properties,omitempty"`
}
