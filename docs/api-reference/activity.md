# API Reference: activity

Import path: `github.com/microsoft/agents-sdk-go/activity`

The `activity` package implements the Microsoft 365 Agents protocol schema. It has no
external dependencies beyond the Go standard library.

## Activity

`Activity` is the core type that represents every communication between an agent and a
channel.

```go
type Activity struct {
    Type             string
    ID               string
    Timestamp        string
    LocalTimestamp   string
    LocalTimezone    string
    ServiceURL       string
    ChannelID        string
    From             *ChannelAccount
    Conversation     *ConversationAccount
    Recipient        *ChannelAccount
    TextFormat       string
    AttachmentLayout string
    MembersAdded     []*ChannelAccount
    MembersRemoved   []*ChannelAccount
    ReactionsAdded   []*MessageReaction
    ReactionsRemoved []*MessageReaction
    TopicName        string
    HistoryDisclosed bool
    Locale           string
    Text             string
    Speak            string
    InputHint        string
    Summary          string
    SuggestedActions *SuggestedActions
    Attachments      []*Attachment
    Entities         []map[string]interface{}
    ChannelData      interface{}
    Action           string
    ReplyToID        string
    Label            string
    ValueType        string
    Value            interface{}
    Name             string
    RelatesTo        *ConversationReference
    Code             string
    Expiration       string
    Importance       string
    DeliveryMode     string
    ListenFor        []string
    TextHighlights   []*TextHighlight
    SemanticAction   interface{}
    CallerID         string
}
```

All fields map to JSON with `omitempty`. The struct is fully JSON-serializable with
`encoding/json`.

### Activity Methods

| Method | Description |
|---|---|
| `IsType(activityType string) bool` | Case-insensitive type check; also matches `invoke/subtype` prefixes |
| `AsMessageActivity() *Activity` | Returns self if type is `message`, otherwise nil |
| `AsEventActivity() *Activity` | Returns self if type is `event`, otherwise nil |
| `AsInvokeActivity() *Activity` | Returns self if type is `invoke`, otherwise nil |
| `AsTypingActivity() *Activity` | Returns self if type is `typing`, otherwise nil |
| `AsEndOfConversationActivity() *Activity` | Returns self if type is `endOfConversation`, otherwise nil |
| `AsConversationUpdateActivity() *Activity` | Returns self if type is `conversationUpdate`, otherwise nil |
| `AsContactRelationUpdateActivity() *Activity` | Returns self if type is `contactRelationUpdate`, otherwise nil |
| `AsInstallationUpdateActivity() *Activity` | Returns self if type is `installationUpdate`, otherwise nil |
| `AsMessageReactionActivity() *Activity` | Returns self if type is `messageReaction`, otherwise nil |
| `AsMessageUpdateActivity() *Activity` | Returns self if type is `messageUpdate`, otherwise nil |
| `AsMessageDeleteActivity() *Activity` | Returns self if type is `messageDelete`, otherwise nil |
| `AsTraceActivity() *Activity` | Returns self if type is `trace`, otherwise nil |
| `AsHandoffActivity() *Activity` | Returns self if type is `handoff`, otherwise nil |
| `AsSuggestionActivity() *Activity` | Returns self if type is `suggestion`, otherwise nil |
| `GetConversationReference() *ConversationReference` | Builds a reference for proactive messaging |
| `GetReplyConversationReference(reply *ResourceResponse) *ConversationReference` | Reference from a send response |
| `ApplyConversationReference(ref *ConversationReference, isIncoming bool) *Activity` | Applies routing info from a stored reference |
| `CreateReply(text string) *Activity` | Creates a new message activity pre-filled with conversation context |
| `HasContent() bool` | Returns true if the activity has text, summary, attachments, or channel data |
| `IsFromStreamingConnection() bool` | Returns true if ServiceURL does not start with "http" |
| `IsAgenticRequest() bool` | Returns true if the recipient has an agentic identity role |
| `GetAgenticInstanceID() string` | Returns the agentic app ID from the recipient |
| `GetAgenticUser() string` | Returns the agentic user ID from the recipient |
| `GetAgenticTenantID() string` | Returns the tenant ID for an agentic request |

## Constructor Functions

| Function | Description |
|---|---|
| `NewMessageActivity(text string) *Activity` | Message activity with `Type = "message"` and the given text |
| `NewEventActivity(name string, value interface{}) *Activity` | Event activity |
| `NewTypingActivity() *Activity` | Typing indicator activity |
| `NewTraceActivity(name string, value interface{}, valueType, label string) *Activity` | Trace activity for debugging |
| `NewConversationUpdateActivity() *Activity` | Conversation update activity |
| `NewEndOfConversationActivity() *Activity` | End-of-conversation activity |
| `NewInvokeActivity() *Activity` | Invoke activity |
| `NewHandoffActivity() *Activity` | Handoff activity |
| `NewContactRelationUpdateActivity() *Activity` | Contact relation update activity |

## ActivityType Constants

```go
const (
    ActivityTypeMessage              = "message"
    ActivityTypeContactRelationUpdate = "contactRelationUpdate"
    ActivityTypeConversationUpdate   = "conversationUpdate"
    ActivityTypeTyping               = "typing"
    ActivityTypeEndOfConversation    = "endOfConversation"
    ActivityTypeEvent                = "event"
    ActivityTypeInvoke               = "invoke"
    ActivityTypeInvokeResponse       = "invokeResponse"
    ActivityTypeDeleteUserData       = "deleteUserData"
    ActivityTypeMessageUpdate        = "messageUpdate"
    ActivityTypeMessageDelete        = "messageDelete"
    ActivityTypeInstallationUpdate   = "installationUpdate"
    ActivityTypeMessageReaction      = "messageReaction"
    ActivityTypeSuggestion           = "suggestion"
    ActivityTypeTrace                = "trace"
    ActivityTypeHandoff              = "handoff"
)
```

## Other Constants

### InputHint

```go
const (
    InputHintAcceptingInput  = "acceptingInput"
    InputHintIgnoringInput   = "ignoringInput"
    InputHintExpectingInput  = "expectingInput"
)
```

### TextFormatType

```go
const (
    TextFormatTypeMarkdown = "markdown"
    TextFormatTypePlain    = "plain"
    TextFormatTypeXML      = "xml"
)
```

### AttachmentLayoutType

```go
const (
    AttachmentLayoutTypeList     = "list"
    AttachmentLayoutTypeCarousel = "carousel"
)
```

### DeliveryMode

```go
const (
    DeliveryModeNormal        = "normal"
    DeliveryModeExpectReplies = "expectReplies"
    DeliveryModeNotification  = "notification"
)
```

### RoleType

```go
const (
    RoleTypeUser            = "user"
    RoleTypeBot             = "bot"
    RoleTypeSkill           = "skill"
    RoleTypeAgenticIdentity = "AgenticIdentity"
    RoleTypeAgenticUser     = "AgenticUser"
)
```

## Supporting Types

### ChannelAccount

```go
type ChannelAccount struct {
    ID            string
    Name          string
    Role          string
    AadObjectID   string
    AgenticAppID  string
    AgenticUserID string
    TenantID      string
}
```

### ConversationAccount

```go
type ConversationAccount struct {
    IsGroup          bool
    ConversationType string
    TenantID         string
    ID               string
    Name             string
    AadObjectID      string
    Role             string
}
```

### ConversationReference

```go
type ConversationReference struct {
    ActivityID   string
    User         *ChannelAccount
    Agent        *ChannelAccount
    Conversation *ConversationAccount
    ChannelID    string
    Locale       string
    ServiceURL   string
}
```

Method: `GetContinuationActivity() *Activity` — creates an activity pre-filled from
the reference, suitable for `ContinueConversation`.

### Attachment

```go
type Attachment struct {
    ContentType  string
    ContentURL   string
    Content      interface{}
    Name         string
    ThumbnailURL string
}
```

### ResourceResponse

```go
type ResourceResponse struct {
    ID string
}
```

### SuggestedActions

```go
type SuggestedActions struct {
    To      []string
    Actions []*CardAction
}
```

### CardAction

```go
type CardAction struct {
    Type        string
    Title       string
    Image       string
    Text        string
    DisplayText string
    Value       interface{}
    ChannelData interface{}
    ImageAltText string
}
```

## Card Types (activity/card)

Import path: `github.com/microsoft/agents-sdk-go/activity/card`

| Type | ContentType Constant |
|---|---|
| `HeroCard` | `ContentTypeHeroCard` |
| `ThumbnailCard` | `ContentTypeThumbnailCard` |
| `OAuthCard` | `ContentTypeOAuthCard` |
| `SigninCard` | `ContentTypeSigninCard` |
| `ReceiptCard` | `ContentTypeReceiptCard` |
| `AudioCard` | `ContentTypeAudioCard` |
| `VideoCard` | `ContentTypeVideoCard` |
| `AnimationCard` | `ContentTypeAnimationCard` |
| `AdaptiveCardInvokeValue` / `AdaptiveCardInvokeResponse` | — |

To attach a card to a message, set `Attachment.ContentType` to the appropriate
constant and `Attachment.Content` to the card struct.

Example:

```go
import (
    "github.com/microsoft/agents-sdk-go/activity"
    "github.com/microsoft/agents-sdk-go/activity/card"
)

hero := &card.HeroCard{
    Title:    "Hello",
    Subtitle: "This is a hero card",
    Text:     "Some body text",
    Buttons: []*activity.CardAction{
        {Type: "openUrl", Title: "Learn more", Value: "https://example.com"},
    },
}

act := &activity.Activity{
    Type: activity.ActivityTypeMessage,
    Attachments: []*activity.Attachment{
        {ContentType: card.ContentTypeHeroCard, Content: hero},
    },
}
```

## Teams Types (activity/teams)

Import path: `github.com/microsoft/agents-sdk-go/activity/teams`

Key types used by `hosting/teams`:

| Type | Description |
|---|---|
| `TeamsChannelData` | ChannelData for Teams activities |
| `TeamsChannelAccount` | Extended ChannelAccount with Teams-specific fields |
| `ChannelInfo` | Teams channel ID and name |
| `TeamInfo` | Team ID and name |
| `TaskModuleRequest` | Request payload for task/fetch and task/submit |
| `TaskModuleResponse` | Response containing a task module continuation or message |
| `TaskModuleContinueResponse` | Continue response with TaskModuleTaskInfo |
| `TaskModuleTaskInfo` | Title, height, width, URL, card |
| `MessagingExtensionQuery` | Query for composeExtension/query invokes |
| `MessagingExtensionResponse` | Response for messaging extension queries |
| `MeetingStartEventDetails` | Meeting start time |
| `MeetingEndEventDetails` | Meeting end time |
| `ConfigResponse` | Response for config/fetch and config/submit |
