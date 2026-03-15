// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package storage

import (
	"context"
	"time"

	"github.com/microsoft/agents-sdk-go/activity"
)

// TranscriptInfo holds metadata about a stored conversation transcript.
type TranscriptInfo struct {
	// ChannelID is the channel the conversation took place on.
	ChannelID string
	// ID is the conversation ID.
	ID string
	// Created is the time the transcript was first created.
	Created time.Time
}

// TranscriptLogger writes activities to a persistent transcript store.
type TranscriptLogger interface {
	// LogActivity appends an activity to the transcript for its conversation.
	LogActivity(ctx context.Context, act *activity.Activity) error
}

// TranscriptStore provides read and management access to stored transcripts.
// It extends TranscriptLogger with retrieval and deletion capabilities.
type TranscriptStore interface {
	TranscriptLogger
	// GetTranscriptActivities returns the activities for a conversation,
	// optionally filtered to those at or after startDate.
	// Returns the activities and a continuation token for paged access.
	GetTranscriptActivities(
		ctx context.Context,
		channelID, conversationID string,
		continuationToken string,
		startDate time.Time,
	) ([]*activity.Activity, string, error)
	// ListTranscripts returns the transcript metadata for a channel.
	// Returns the transcripts and a continuation token for paged access.
	ListTranscripts(ctx context.Context, channelID string, continuationToken string) ([]*TranscriptInfo, string, error)
	// DeleteTranscript permanently removes a conversation transcript.
	DeleteTranscript(ctx context.Context, channelID, conversationID string) error
}
