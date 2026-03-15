// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

// ConfigResponseBase is the base type for configuration responses.
type ConfigResponseBase struct {
	// ResponseType is the type of config response.
	ResponseType string `json:"responseType,omitempty"`
}

// ConfigResponse is the envelope for a Config Response Payload.
type ConfigResponse struct {
	ConfigResponseBase
	// Config is the response to the config message. Possible values: 'auth', 'task'.
	Config interface{} `json:"config,omitempty"`
	// CacheInfo contains the response cache info.
	CacheInfo *CacheInfo `json:"cacheInfo,omitempty"`
}
