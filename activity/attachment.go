// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

// Attachment represents an attachment within an activity.
type Attachment struct {
	// ContentType is the mimetype/content-type for the file.
	ContentType string `json:"contentType,omitempty"`
	// ContentURL is the content URL.
	ContentURL string `json:"contentUrl,omitempty"`
	// Content is the embedded content.
	Content interface{} `json:"content,omitempty"`
	// Name is the optional name of the attachment.
	Name string `json:"name,omitempty"`
	// ThumbnailURL is an optional thumbnail associated with the attachment.
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
}
