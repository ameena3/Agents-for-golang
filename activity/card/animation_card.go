// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

// ContentTypeAnimationCard is the content type for an animation card attachment.
const ContentTypeAnimationCard = "application/vnd.microsoft.card.animation"

// AnimationCard represents an animation card (e.g., gif or short video clip).
type AnimationCard struct {
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
