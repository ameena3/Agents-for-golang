// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

// MediaURL represents a media URL.
type MediaURL struct {
	// URL is the URL for the media.
	URL string `json:"url,omitempty"`
	// Profile is an optional profile hint to differentiate multiple MediaURL objects.
	Profile string `json:"profile,omitempty"`
}

// ThumbnailURL represents a thumbnail URL.
type ThumbnailURL struct {
	// URL is the URL pointing to the thumbnail.
	URL string `json:"url,omitempty"`
	// Alt is the HTML alt text to include on this thumbnail image.
	Alt string `json:"alt,omitempty"`
}

// CardImage represents an image on a card.
type CardImage struct {
	// URL is the URL thumbnail image for major content property.
	URL string `json:"url,omitempty"`
	// Alt is image description intended for screen readers.
	Alt string `json:"alt,omitempty"`
	// Tap is the action assigned to this attachment.
	Tap *CardAction `json:"tap,omitempty"`
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

// Fact represents a key-value pair on a receipt card.
type Fact struct {
	// Key is the key for this fact.
	Key string `json:"key,omitempty"`
	// Value is the value for this fact.
	Value string `json:"value,omitempty"`
}

// MediaCard is the base type for media cards (AudioCard, VideoCard, AnimationCard).
type MediaCard struct {
	// Title is the title of this card.
	Title string `json:"title,omitempty"`
	// Subtitle is the subtitle of this card.
	Subtitle string `json:"subtitle,omitempty"`
	// Text is the text of this card.
	Text string `json:"text,omitempty"`
	// Image is the thumbnail placeholder.
	Image *ThumbnailURL `json:"image,omitempty"`
	// Media contains media URLs for this card.
	Media []*MediaURL `json:"media,omitempty"`
	// Buttons contains actions on this card.
	Buttons []*CardAction `json:"buttons,omitempty"`
	// Shareable indicates whether this content may be shared with others.
	Shareable *bool `json:"shareable,omitempty"`
	// Autoloop indicates whether the client should loop playback at end of content.
	Autoloop *bool `json:"autoloop,omitempty"`
	// Autostart indicates whether the client should automatically start playback.
	Autostart *bool `json:"autostart,omitempty"`
	// Aspect is the aspect ratio of thumbnail/media placeholder.
	Aspect string `json:"aspect,omitempty"`
	// Duration describes the length of the media content in ISO 8601 Duration format.
	Duration string `json:"duration,omitempty"`
	// Value is a supplementary parameter for this card.
	Value interface{} `json:"value,omitempty"`
}
