// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

// ContentTypeReceiptCard is the content type for a receipt card attachment.
const ContentTypeReceiptCard = "application/vnd.microsoft.card.receipt"

// ReceiptItem represents an item on a receipt card.
type ReceiptItem struct {
	// Title is the title of the card.
	Title string `json:"title,omitempty"`
	// Subtitle appears just below Title field.
	Subtitle string `json:"subtitle,omitempty"`
	// Text field appears just below subtitle.
	Text string `json:"text,omitempty"`
	// Image is the item image.
	Image *CardImage `json:"image,omitempty"`
	// Price is the amount with currency.
	Price string `json:"price,omitempty"`
	// Quantity is the number of items of given kind.
	Quantity string `json:"quantity,omitempty"`
	// Tap is the action to activate when user taps on the item bubble.
	Tap *CardAction `json:"tap,omitempty"`
}

// ReceiptCard represents a receipt card.
type ReceiptCard struct {
	// Title is the title of the card.
	Title string `json:"title,omitempty"`
	// Facts is the array of Fact objects.
	Facts []*Fact `json:"facts,omitempty"`
	// Items is the array of Receipt Items.
	Items []*ReceiptItem `json:"items,omitempty"`
	// Tap is the action to activate when user taps on the card.
	Tap *CardAction `json:"tap,omitempty"`
	// Total is the total amount of money paid (or to be paid).
	Total string `json:"total,omitempty"`
	// Tax is the total amount of tax paid (or to be paid).
	Tax string `json:"tax,omitempty"`
	// VAT is the total amount of VAT paid (or to be paid).
	VAT string `json:"vat,omitempty"`
	// Buttons is the set of actions applicable to the current card.
	Buttons []*CardAction `json:"buttons,omitempty"`
}
