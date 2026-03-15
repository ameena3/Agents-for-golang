// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

// TeamsChannelAccount details user Azure Active Directory details for Teams channels.
type TeamsChannelAccount struct {
	// ID is the channel id for the user or bot on this channel.
	ID string `json:"id,omitempty"`
	// Name is the display friendly name.
	Name string `json:"name,omitempty"`
	// GivenName is the given name part of the user name.
	GivenName string `json:"givenName,omitempty"`
	// Surname is the surname part of the user name.
	Surname string `json:"surname,omitempty"`
	// Email is the email ID of the user.
	Email string `json:"email,omitempty"`
	// UserPrincipalName is the unique user principal name.
	UserPrincipalName string `json:"userPrincipalName,omitempty"`
	// TenantID is the tenant ID of the user.
	TenantID string `json:"tenantId,omitempty"`
	// UserRole is the user role of the user.
	UserRole string `json:"userRole,omitempty"`
	// Properties contains additional channel-specific properties.
	Properties map[string]interface{} `json:"properties,omitempty"`
}
