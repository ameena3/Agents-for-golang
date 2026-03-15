// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

// ContentTypeHeroCard is the content type for a hero card attachment.
const ContentTypeHeroCard = "application/vnd.microsoft.card.hero"

// HeroCard is a card with a single, large image.
type HeroCard struct {
	// Title is the title of the card.
	Title string `json:"title,omitempty"`
	// Subtitle is the subtitle of the card.
	Subtitle string `json:"subtitle,omitempty"`
	// Text is the text for the card.
	Text string `json:"text,omitempty"`
	// Images is an array of images for the card.
	Images []*CardImage `json:"images,omitempty"`
	// Buttons is the set of actions applicable to the current card.
	Buttons []*CardAction `json:"buttons,omitempty"`
	// Tap is the action to activate when user taps on the card itself.
	Tap *CardAction `json:"tap,omitempty"`
}
