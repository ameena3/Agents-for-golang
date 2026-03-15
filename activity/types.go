// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

// Activity type constants.
const (
	ActivityTypeMessage               = "message"
	ActivityTypeContactRelationUpdate = "contactRelationUpdate"
	ActivityTypeConversationUpdate    = "conversationUpdate"
	ActivityTypeTyping                = "typing"
	ActivityTypeEndOfConversation     = "endOfConversation"
	ActivityTypeEvent                 = "event"
	ActivityTypeInvoke                = "invoke"
	ActivityTypeInvokeResponse        = "invokeResponse"
	ActivityTypeDeleteUserData        = "deleteUserData"
	ActivityTypeMessageUpdate         = "messageUpdate"
	ActivityTypeMessageDelete         = "messageDelete"
	ActivityTypeInstallationUpdate    = "installationUpdate"
	ActivityTypeMessageReaction       = "messageReaction"
	ActivityTypeSuggestion            = "suggestion"
	ActivityTypeTrace                 = "trace"
	ActivityTypeHandoff               = "handoff"
	ActivityTypeCommand               = "command"
	ActivityTypeCommandResult         = "commandResult"
)

// InputHint constants indicate whether the agent is accepting, expecting, or ignoring user input.
const (
	InputHintAcceptingInput = "acceptingInput"
	InputHintIgnoringInput  = "ignoringInput"
	InputHintExpectingInput = "expectingInput"
)

// TextFormatType constants specify the format of text fields.
const (
	TextFormatTypeMarkdown = "markdown"
	TextFormatTypePlain    = "plain"
	TextFormatTypeXML      = "xml"
)

// AttachmentLayoutType constants specify the layout of attachments.
const (
	AttachmentLayoutTypeList     = "list"
	AttachmentLayoutTypeCarousel = "carousel"
)

// DeliveryMode constants specify the delivery mode of an activity.
const (
	DeliveryModeNormal        = "normal"
	DeliveryModeNotification  = "notification"
	DeliveryModeExpectReplies = "expectReplies"
	DeliveryModeEphemeral     = "ephemeral"
	DeliveryModeStream        = "stream"
)

// EndOfConversationCode constants indicate why the conversation ended.
const (
	EndOfConversationCodeUnknown              = "unknown"
	EndOfConversationCodeCompletedSuccessfully = "completedSuccessfully"
	EndOfConversationCodeUserCancelled        = "userCancelled"
	EndOfConversationCodeBotTimedOut          = "botTimedOut"
	EndOfConversationCodeBotIssuedInvalidMessage = "botIssuedInvalidMessage"
	EndOfConversationCodeChannelFailed        = "channelFailed"
)

// RoleType constants specify the role of the entity behind an account.
const (
	RoleTypeUser             = "user"
	RoleTypeAgent            = "bot"
	RoleTypeSkill            = "skill"
	RoleTypeConnectorUser    = "connectoruser"
	RoleTypeAgenticIdentity  = "agenticAppInstance"
	RoleTypeAgenticUser      = "agenticUser"
)

// MessageReactionType constants specify types of reactions to messages.
const (
	MessageReactionTypeReactionsAdded   = "reactionsAdded"
	MessageReactionTypeReactionsRemoved = "reactionsRemoved"
)

// ActionType constants specify the type of card action.
const (
	ActionTypeOpenURL      = "openUrl"
	ActionTypeImBack       = "imBack"
	ActionTypePostBack     = "postBack"
	ActionTypePlayAudio    = "playAudio"
	ActionTypePlayVideo    = "playVideo"
	ActionTypeShowImage    = "showImage"
	ActionTypeDownloadFile = "downloadFile"
	ActionTypeSignin       = "signin"
	ActionTypeCall         = "call"
	ActionTypeMessageBack  = "messageBack"
)

// ActivityImportance constants specify the importance of an activity.
const (
	ActivityImportanceLow    = "low"
	ActivityImportanceNormal = "normal"
	ActivityImportanceHigh   = "high"
)

// CallerIDConstants provides caller ID prefix values.
const (
	CallerIDPublicAzureChannel = "urn:botframework:azure"
	CallerIDUSGovChannel       = "urn:botframework:azureusgov"
	CallerIDSkillPrefix        = "urn:botframework:aadappid:"
)

// SignInConstants provides constants for sign-in flows.
const (
	SignInTokenExchangeOperationName = "signin/tokenExchange"
	SignInVerifyStateOperationName   = "signin/verifyState"
	SignInInvokeActivityName         = "signin/verifyState"
)
