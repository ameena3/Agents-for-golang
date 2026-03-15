// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

// SuggestedActions represents suggested actions that can be performed.
type SuggestedActions struct {
	// To is the list of recipient IDs that the actions should be shown to.
	To []string `json:"to,omitempty"`
	// Actions contains the actions that can be shown to the user.
	Actions []*CardAction `json:"actions,omitempty"`
}

// CardAction represents a clickable action on a card.
type CardAction struct {
	// Type is the type of action implemented by this button.
	Type string `json:"type,omitempty"`
	// Title is the text description which appears on the button.
	Title string `json:"title,omitempty"`
	// Image is the image URL which will appear on the button, next to text label.
	Image string `json:"image,omitempty"`
	// Text is the text for this action.
	Text string `json:"text,omitempty"`
	// DisplayText is optional text to display in the chat feed if the button is clicked.
	DisplayText string `json:"displayText,omitempty"`
	// Value is the supplementary parameter for action.
	Value interface{} `json:"value,omitempty"`
	// ChannelData is the channel-specific data associated with this action.
	ChannelData interface{} `json:"channelData,omitempty"`
	// ImageAltText is the alternate image text to be used in place of the image field.
	ImageAltText string `json:"imageAltText,omitempty"`
}
