// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

// ConversationAccount represents the identity of the conversation within a channel.
type ConversationAccount struct {
	// IsGroup indicates whether the conversation contains more than two participants.
	IsGroup *bool `json:"isGroup,omitempty"`
	// ConversationType indicates the type of the conversation in channels that distinguish between types.
	ConversationType string `json:"conversationType,omitempty"`
	// ID is the channel id for the conversation.
	ID string `json:"id,omitempty"`
	// Name is the display friendly name of the conversation.
	Name string `json:"name,omitempty"`
	// AADObjectID is this conversation's object ID within Azure Active Directory.
	AADObjectID string `json:"aadObjectId,omitempty"`
	// Role is the role of the entity behind the account.
	Role string `json:"role,omitempty"`
	// TenantID is the conversation's tenant ID.
	TenantID string `json:"tenantId,omitempty"`
	// Properties contains additional conversation properties.
	Properties interface{} `json:"properties,omitempty"`
}
