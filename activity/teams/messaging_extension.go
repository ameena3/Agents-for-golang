// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

import "github.com/microsoft/agents-sdk-go/activity"

// MessagingExtensionParameter represents a messaging extension query parameter.
type MessagingExtensionParameter struct {
	// Name is the name of the parameter.
	Name string `json:"name,omitempty"`
	// Value is the value of the parameter.
	Value interface{} `json:"value,omitempty"`
}

// MessagingExtensionQueryOptions represents messaging extension query options.
type MessagingExtensionQueryOptions struct {
	// Skip is the number of entities to skip.
	Skip *int `json:"skip,omitempty"`
	// Count is the number of entities to fetch.
	Count *int `json:"count,omitempty"`
}

// MessagingExtensionQuery represents a messaging extension query.
type MessagingExtensionQuery struct {
	// CommandID is the ID of the command assigned by the Bot.
	CommandID string `json:"commandId,omitempty"`
	// Parameters contains the query parameters.
	Parameters []*MessagingExtensionParameter `json:"parameters,omitempty"`
	// QueryOptions contains query options.
	QueryOptions *MessagingExtensionQueryOptions `json:"queryOptions,omitempty"`
	// State is the state parameter passed back to the bot after authentication/configuration flow.
	State string `json:"state,omitempty"`
}

// MessagingExtensionSuggestedAction represents suggested actions for a messaging extension result.
type MessagingExtensionSuggestedAction struct {
	// Actions contains the suggested actions.
	Actions []*activity.CardAction `json:"actions,omitempty"`
}

// MessagingExtensionAttachment is an attachment for a messaging extension result.
type MessagingExtensionAttachment struct {
	// ContentType is the mimetype/content-type for the file.
	ContentType string `json:"contentType,omitempty"`
	// ContentURL is the content URL.
	ContentURL string `json:"contentUrl,omitempty"`
	// Content is the embedded content.
	Content interface{} `json:"content,omitempty"`
	// Name is the optional name of the attachment.
	Name string `json:"name,omitempty"`
	// ThumbnailURL is the optional thumbnail associated with the attachment.
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	// Preview is the preview attachment.
	Preview *activity.Attachment `json:"preview,omitempty"`
}

// MessagingExtensionResult represents a messaging extension result.
type MessagingExtensionResult struct {
	// AttachmentLayout is the hint for how to deal with multiple attachments.
	AttachmentLayout string `json:"attachmentLayout,omitempty"`
	// Type is the type of the result.
	Type string `json:"type,omitempty"`
	// Attachments contains the attachments (only when type is result).
	Attachments []*MessagingExtensionAttachment `json:"attachments,omitempty"`
	// SuggestedActions contains suggested actions for the result.
	SuggestedActions *MessagingExtensionSuggestedAction `json:"suggestedActions,omitempty"`
	// Text contains the text (only when type is message).
	Text string `json:"text,omitempty"`
	// ActivityPreview is the message activity to preview (only when type is botMessagePreview).
	ActivityPreview *activity.Activity `json:"activityPreview,omitempty"`
}

// CacheInfo contains cache information for responses.
type CacheInfo struct {
	// CacheType is the type of cache.
	CacheType string `json:"cacheType,omitempty"`
	// CacheDuration is the cache duration in seconds.
	CacheDuration *int `json:"cacheDuration,omitempty"`
}

// MessagingExtensionResponse represents a messaging extension response.
type MessagingExtensionResponse struct {
	// ComposeExtension is the compose extension result.
	ComposeExtension *MessagingExtensionResult `json:"composeExtension,omitempty"`
	// CacheInfo contains cache info for this response.
	CacheInfo *CacheInfo `json:"cacheInfo,omitempty"`
}
