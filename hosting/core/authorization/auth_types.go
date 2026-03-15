// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authorization

// AuthType represents the type of authentication being used for a turn.
type AuthType string

const (
	// AuthTypeConnectorUserAuth is user-to-agent authentication via the Bot Framework connector.
	AuthTypeConnectorUserAuth AuthType = "ConnectorUserAuth"
	// AuthTypeAgenticUserAuth is agent-to-agent authentication with a user context preserved.
	AuthTypeAgenticUserAuth AuthType = "AgenticUserAuth"
	// AuthTypeAnonymous indicates no authentication; the request is unauthenticated.
	AuthTypeAnonymous AuthType = "Anonymous"
)
