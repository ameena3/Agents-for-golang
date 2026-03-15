// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

import "github.com/microsoft/agents-sdk-go/activity"

// TaskModuleRequestContext contains context for a task module request.
type TaskModuleRequestContext struct {
	// Theme is the current theme.
	Theme string `json:"theme,omitempty"`
}

// TabEntityContext contains the tab entity context.
type TabEntityContext struct {
	// TabEntityID is the tab entity ID.
	TabEntityID string `json:"tabEntityId,omitempty"`
}

// TaskModuleRequest contains the task module invoke request value payload.
type TaskModuleRequest struct {
	// Data contains user input data.
	Data interface{} `json:"data,omitempty"`
	// Context contains the current user context.
	Context *TaskModuleRequestContext `json:"context,omitempty"`
	// TabEntityContext contains the current tab request context.
	TabEntityContext *TabEntityContext `json:"tabContext,omitempty"`
}

// TaskModuleTaskInfo contains information about a task module task.
type TaskModuleTaskInfo struct {
	// Title is the title of the task module.
	Title string `json:"title,omitempty"`
	// Height is the height of the task module.
	Height interface{} `json:"height,omitempty"`
	// Width is the width of the task module.
	Width interface{} `json:"width,omitempty"`
	// URL is the URL of the task module.
	URL string `json:"url,omitempty"`
	// Card is the adaptive card for the task module.
	Card *activity.Attachment `json:"card,omitempty"`
	// FallbackURL is the fallback URL for the task module.
	FallbackURL string `json:"fallbackUrl,omitempty"`
	// CompletionBotID is the completion bot ID.
	CompletionBotID string `json:"completionBotId,omitempty"`
}

// TaskModuleContinueResponse is a response to continue a task module.
type TaskModuleContinueResponse struct {
	// Type is the type of response. Default is 'continue'.
	Type string `json:"type,omitempty"`
	// Value is the task module task info.
	Value *TaskModuleTaskInfo `json:"value,omitempty"`
}

// TaskModuleMessageResponse is a response to display a message in a task module.
type TaskModuleMessageResponse struct {
	// Type is the type of response. Default is 'message'.
	Type string `json:"type,omitempty"`
	// Value is the message to display.
	Value string `json:"value,omitempty"`
}

// TaskModuleResponse is the envelope for a Task Module Response.
type TaskModuleResponse struct {
	// Task is the task module response base.
	Task interface{} `json:"task,omitempty"`
	// CacheInfo contains cache info for this TaskModuleResponse.
	CacheInfo *CacheInfo `json:"cacheInfo,omitempty"`
}
