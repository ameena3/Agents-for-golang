// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import "github.com/microsoft/agents-sdk-go/activity"

// Text creates a simple text message activity with acceptingInput as the input hint.
func Text(text string) *activity.Activity {
	return &activity.Activity{
		Type:      activity.ActivityTypeMessage,
		Text:      text,
		InputHint: activity.InputHintAcceptingInput,
	}
}

// TextWithAttachment creates a message with text and a single attachment.
func TextWithAttachment(text string, attachment *activity.Attachment) *activity.Activity {
	act := Text(text)
	if attachment != nil {
		act.Attachments = []*activity.Attachment{attachment}
	}
	return act
}

// Attachments creates a message activity that contains only attachments (no text body).
func Attachments(attachments ...*activity.Attachment) *activity.Activity {
	list := make([]*activity.Attachment, 0, len(attachments))
	for _, a := range attachments {
		if a != nil {
			list = append(list, a)
		}
	}
	return &activity.Activity{
		Type:        activity.ActivityTypeMessage,
		Attachments: list,
		InputHint:   activity.InputHintAcceptingInput,
	}
}

// ContentURL creates a message activity with a single attachment whose content is
// hosted at the given URL. The contentType parameter specifies the MIME type and
// name is the optional display name.
func ContentURL(url, contentType, name string) *activity.Activity {
	return Attachments(&activity.Attachment{
		ContentType: contentType,
		ContentURL:  url,
		Name:        name,
	})
}

// SuggestedActions creates a message with suggested action buttons displayed
// alongside the text.
func SuggestedActions(text string, actions ...*activity.CardAction) *activity.Activity {
	act := Text(text)
	if len(actions) > 0 {
		act.SuggestedActions = &activity.SuggestedActions{
			Actions: actions,
		}
	}
	return act
}

// TypingActivity creates a typing indicator activity.
func TypingActivity() *activity.Activity {
	return &activity.Activity{Type: activity.ActivityTypeTyping}
}

// EndOfConversation creates an end-of-conversation activity with the given code.
func EndOfConversation(code string) *activity.Activity {
	return &activity.Activity{
		Type: activity.ActivityTypeEndOfConversation,
		Code: code,
	}
}

// Event creates an event activity with the given name and value.
func Event(name string, value interface{}) *activity.Activity {
	return &activity.Activity{
		Type:  activity.ActivityTypeEvent,
		Name:  name,
		Value: value,
	}
}

// Carousel creates a message activity that displays attachments in carousel layout.
func Carousel(attachments ...*activity.Attachment) *activity.Activity {
	list := make([]*activity.Attachment, 0, len(attachments))
	for _, a := range attachments {
		if a != nil {
			list = append(list, a)
		}
	}
	return &activity.Activity{
		Type:             activity.ActivityTypeMessage,
		AttachmentLayout: activity.AttachmentLayoutTypeCarousel,
		Attachments:      list,
		InputHint:        activity.InputHintAcceptingInput,
	}
}
