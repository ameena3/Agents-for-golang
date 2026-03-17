// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package app provides the AgentApplication — the modern, decorator-style
// agent programming model for the Microsoft 365 Agents SDK for Go.
package app

import (
	"github.com/ameena3/Agents-for-golang/hosting/core/storage"
)

// AppOptions configures an AgentApplication.
// All fields are optional; sensible defaults are used when fields are zero.
type AppOptions[StateT any] struct {
	// Storage is the backend used to persist ConversationState and UserState.
	// If nil, state changes are only held in memory for the duration of the turn.
	Storage storage.Storage

	// LongRunningMessages enables long-running message support.
	// When true, the incoming request is immediately converted to a proactive
	// conversation so the handler can run beyond the channel timeout.
	// Requires a configured adapter and BotAppID.
	LongRunningMessages bool

	// BotAppID is the application (client) ID of the bot.
	// Required when LongRunningMessages is true.
	BotAppID string

	// StartTypingTimer enables the automatic typing indicator.
	// When true, a typing activity is sent at the start of each message turn.
	StartTypingTimer bool

	// RemoveRecipientMention strips the bot @mention from incoming message text.
	RemoveRecipientMention bool

	// NormalizeMentions normalizes all @mention text in incoming activities.
	NormalizeMentions bool

	// TurnStateFactory is an optional factory called once per turn to create
	// the initial AppStateT value. If nil the zero value of AppStateT is used.
	TurnStateFactory func() StateT
}
