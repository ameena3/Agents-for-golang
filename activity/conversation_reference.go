// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

// ConversationReference contains an object relating to a particular point in a conversation.
type ConversationReference struct {
	// ActivityID is the ID of the activity to refer to.
	ActivityID string `json:"activityId,omitempty"`
	// User is the user participating in this conversation.
	User *ChannelAccount `json:"user,omitempty"`
	// Agent (serialized as "bot") is the agent participating in this conversation.
	Agent *ChannelAccount `json:"bot,omitempty"`
	// Conversation is the conversation reference.
	Conversation *ConversationAccount `json:"conversation,omitempty"`
	// ChannelID is the channel ID.
	ChannelID string `json:"channelId,omitempty"`
	// Locale is a locale name for the contents of the text field.
	Locale string `json:"locale,omitempty"`
	// ServiceURL is the service endpoint where operations concerning the referenced conversation may be performed.
	ServiceURL string `json:"serviceUrl,omitempty"`
}

// GetContinuationActivity creates an event activity for continuing a conversation from this reference.
func (r *ConversationReference) GetContinuationActivity() *Activity {
	return &Activity{
		Type:         ActivityTypeEvent,
		Name:         "ContinueConversation",
		ChannelID:    r.ChannelID,
		ServiceURL:   r.ServiceURL,
		Conversation: r.Conversation,
		Recipient:    r.Agent,
		From:         r.User,
		RelatesTo:    r,
	}
}
