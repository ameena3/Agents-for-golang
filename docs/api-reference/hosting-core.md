# API Reference: hosting/core

Import path: `github.com/ameena3/Agents-for-golang/hosting/core`

The `hosting/core` package is the runtime engine of the SDK. It provides the
`Agent` interface, `TurnContext`, `ActivityHandler`, `ChannelAdapter`, and the
middleware pipeline.

## Interfaces

### Agent

```go
type Agent interface {
    OnTurn(ctx context.Context, turnCtx *TurnContext) error
}
```

Every agent must implement `Agent`. Both `ActivityHandler` (via embedding) and
`AgentApplication` satisfy this interface.

### Middleware

```go
type Middleware interface {
    OnTurn(ctx context.Context, turnCtx *TurnContext, next func(context.Context) error) error
}
```

Implement to intercept every turn. Call `next(ctx)` to continue the pipeline.

### ChannelAdapter

```go
type ChannelAdapter interface {
    SendActivities(ctx context.Context, turnCtx *TurnContext, activities []*activity.Activity) ([]*activity.ResourceResponse, error)
    UpdateActivity(ctx context.Context, turnCtx *TurnContext, act *activity.Activity) (*activity.ResourceResponse, error)
    DeleteActivity(ctx context.Context, turnCtx *TurnContext, activityID string) error
    ContinueConversation(ctx context.Context, ref *activity.ConversationReference, handler func(context.Context, *TurnContext) error) error
    Use(middleware ...Middleware)
}
```

### AccessTokenProvider

```go
type AccessTokenProvider interface {
    GetAccessToken(ctx context.Context, resourceURL string, scopes []string, forceRefresh bool) (string, error)
    AcquireTokenOnBehalfOf(ctx context.Context, scopes []string, userAssertion string) (string, error)
}
```

Implemented by `authentication.MsalAuth`.

## TurnContext

`TurnContext` is created by the adapter for each incoming activity and passed through
the entire pipeline. It is scoped to a single turn; never store or reuse it across
turns.

### Constructor

```go
func NewTurnContext(adapter ChannelAdapter, act *activity.Activity, identity *authorization.ClaimsIdentity) *TurnContext
```

The adapter calls this for you; you rarely need to call it directly.

### Methods

| Method | Description |
|---|---|
| `Activity() *activity.Activity` | Returns the incoming activity for this turn |
| `Adapter() ChannelAdapter` | Returns the channel adapter |
| `Identity() *authorization.ClaimsIdentity` | Returns the claims identity from the auth token |
| `TurnState() map[string]interface{}` | Turn-scoped key/value service cache |
| `Responded() bool` | True if a non-trace activity has been sent this turn |
| `SendActivity(ctx, act) (*ResourceResponse, error)` | Sends one activity; runs through send-activity hooks |
| `SendActivities(ctx, acts) ([]*ResourceResponse, error)` | Sends multiple activities |
| `UpdateActivity(ctx, act) (*ResourceResponse, error)` | Updates an existing activity |
| `DeleteActivity(ctx, activityID) error` | Deletes an activity by ID |
| `OnSendActivities(handler SendActivitiesHandler)` | Registers a hook called before activities are sent |
| `OnUpdateActivity(handler UpdateActivityHandler)` | Registers a hook called before an activity is updated |
| `OnDeleteActivity(handler DeleteActivityHandler)` | Registers a hook called before an activity is deleted |
| `GetConversationReference() *activity.ConversationReference` | Returns a reference for proactive messaging |

### Hook Types

```go
type SendActivitiesHandler func(
    ctx context.Context,
    activities []*activity.Activity,
    next func(context.Context, []*activity.Activity) ([]*activity.ResourceResponse, error),
) ([]*activity.ResourceResponse, error)

type UpdateActivityHandler func(
    ctx context.Context,
    act *activity.Activity,
    next func(context.Context, *activity.Activity) (*activity.ResourceResponse, error),
) (*activity.ResourceResponse, error)

type DeleteActivityHandler func(
    ctx context.Context,
    activityID string,
    next func(context.Context, string) error,
) error
```

### Public Field

```go
BufferedReplies []*activity.Activity
```

When `Activity().DeliveryMode == "expectReplies"`, outbound activities are buffered
here instead of being sent immediately.

## ActivityHandler

`ActivityHandler` is an embeddable base struct that dispatches activities to typed
handler methods. Embed it and override only the methods you need; all methods default
to no-ops.

```go
type MyAgent struct {
    core.ActivityHandler
}
```

### OnTurn (dispatcher)

```go
func (h *ActivityHandler) OnTurn(ctx context.Context, tc *TurnContext) error
```

Dispatches based on `Activity().Type`. This implements the `Agent` interface when
embedded.

### Handler Methods (override to implement)

| Method | Activity Type | Notes |
|---|---|---|
| `OnMessageActivity(ctx, tc) error` | `message` | — |
| `OnMessageUpdateActivity(ctx, tc) error` | `messageUpdate` | — |
| `OnMessageDeleteActivity(ctx, tc) error` | `messageDelete` | — |
| `OnConversationUpdateActivity(ctx, tc) error` | `conversationUpdate` | Calls `OnMembersAdded` / `OnMembersRemoved` by default |
| `OnMembersAdded(ctx, members, tc) error` | — | Called from `OnConversationUpdateActivity` |
| `OnMembersRemoved(ctx, members, tc) error` | — | Called from `OnConversationUpdateActivity` |
| `OnMessageReactionActivity(ctx, tc) error` | `messageReaction` | Calls `OnReactionsAdded` / `OnReactionsRemoved` |
| `OnReactionsAdded(ctx, reactions, tc) error` | — | — |
| `OnReactionsRemoved(ctx, reactions, tc) error` | — | — |
| `OnEventActivity(ctx, tc) error` | `event` | Calls `OnTokenResponseEvent` for `tokens/response` |
| `OnTokenResponseEvent(ctx, tc) error` | — | OAuth token response events |
| `OnEvent(ctx, tc) error` | — | All other events |
| `OnInvokeActivity(ctx, tc) error` | `invoke` | — |
| `OnTypingActivity(ctx, tc) error` | `typing` | — |
| `OnEndOfConversationActivity(ctx, tc) error` | `endOfConversation` | — |
| `OnInstallationUpdateActivity(ctx, tc) error` | `installationUpdate` | — |
| `OnUnrecognizedActivityType(ctx, tc) error` | any unrecognized type | — |

## Message Factory Functions

Convenience constructors in `hosting/core` for building common activity types:

| Function | Description |
|---|---|
| `Text(text string) *activity.Activity` | Message activity with text |
| `TextWithAttachment(text string, att *activity.Attachment) *activity.Activity` | Message with text and one attachment |
| `Attachments(atts ...*activity.Attachment) *activity.Activity` | Message with only attachments |
| `ContentURL(url, contentType, name string) *activity.Activity` | Message with a ContentURL attachment |
| `SuggestedActions(text string, actions ...*activity.CardAction) *activity.Activity` | Message with suggested actions |
| `TypingActivity() *activity.Activity` | Typing indicator |
| `EndOfConversation() *activity.Activity` | End-of-conversation signal |
| `Event(name string, value interface{}) *activity.Activity` | Event activity |
| `Carousel(atts ...*activity.Attachment) *activity.Activity` | Message with carousel attachment layout |

## Card Factory Functions

```go
import "github.com/ameena3/Agents-for-golang/hosting/core"

// Returns an *activity.Attachment suitable for use in an activity
core.HeroCard(card *card.HeroCard) *activity.Attachment
core.ThumbnailCard(card *card.ThumbnailCard) *activity.Attachment
core.OAuthCard(card *card.OAuthCard) *activity.Attachment
core.SigninCard(card *card.SigninCard) *activity.Attachment
core.AdaptiveCard(content interface{}) *activity.Attachment
```

## AgentApplication[StateT]

Import path: `github.com/ameena3/Agents-for-golang/hosting/core/app`

`AgentApplication` is the modern, functional-style agent framework. It implements
`core.Agent` and uses Go generics for type-safe application state.

### Constructor

```go
func New[StateT any](opts AppOptions[StateT]) *AgentApplication[StateT]
```

### AppOptions

```go
type AppOptions[StateT any] struct {
    // Storage backend for conversation, user, and temp state.
    // If nil, state is not persisted.
    Storage storage.Storage

    // TurnStateFactory creates the initial AppState for each turn.
    // If nil, the zero value of StateT is used.
    TurnStateFactory func() StateT

    // LongRunningMessages enables 202 Accepted responses for long-running turns.
    LongRunningMessages bool

    // BotAppID is the App ID for the bot (used for some auth scenarios).
    BotAppID string

    // StartTypingTimer sends a typing indicator while the handler runs.
    StartTypingTimer bool

    // RemoveRecipientMention strips @mention text from messages in Teams.
    RemoveRecipientMention bool

    // NormalizeMentions normalizes mention entities in the activity.
    NormalizeMentions bool
}
```

### Route Registration Methods

All registration methods return `*AgentApplication[StateT]` for chaining.

| Method | Description |
|---|---|
| `OnMessage(pattern string, handler HandlerFunc[StateT]) *AgentApplication[StateT]` | Handle message activities. Empty `pattern` matches all messages; otherwise matches as a Go regexp. |
| `OnActivity(activityType string, handler HandlerFunc[StateT]) *AgentApplication[StateT]` | Handle activities of a specific type. |
| `OnInvoke(name string, handler HandlerFunc[StateT]) *AgentApplication[StateT]` | Handle invoke activities by name. Empty name matches all invoke activities. |
| `OnConversationUpdate(handler HandlerFunc[StateT]) *AgentApplication[StateT]` | Handle all conversationUpdate activities. |
| `OnMembersAdded(handler HandlerFunc[StateT]) *AgentApplication[StateT]` | Handle conversationUpdate where `MembersAdded` is non-empty. |

### Lifecycle Hooks

| Method | Description |
|---|---|
| `BeforeTurn(handler func(context.Context, *core.TurnContext) error) *AgentApplication[StateT]` | Run before every turn, before route matching. Abort turn by returning an error. |
| `AfterTurn(handler func(context.Context, *core.TurnContext) error) *AgentApplication[StateT]` | Run after every turn, after state is saved. |
| `OnError(handler func(context.Context, *core.TurnContext, error) error) *AgentApplication[StateT]` | Called when any pipeline step returns an error. Return nil to consume the error. |

### HandlerFunc Type

```go
type HandlerFunc[StateT any] func(ctx context.Context, tc *core.TurnContext, appState StateT) error
```

### OnTurn Pipeline Order

1. BeforeTurn hooks (abort on error)
2. Load state from Storage
3. `RouteList.FindRoute()` — first matching route wins
4. Execute route handler
5. Save state to Storage
6. AfterTurn hooks

### Route Priority

Routes are matched in this priority order:
1. Invoke + agentic routes (highest priority)
2. Invoke routes
3. Agentic routes
4. All other routes (message, activity type, conversation update, etc.)

Within the same priority level, routes match in registration order.

## Storage (hosting/core/storage)

Import path: `github.com/ameena3/Agents-for-golang/hosting/core/storage`

### Storage Interface

```go
type Storage interface {
    Read(ctx context.Context, keys []string) (map[string]StoreItem, error)
    Write(ctx context.Context, changes map[string]StoreItem) error
    Delete(ctx context.Context, keys []string) error
}

type StoreItem interface{} // any JSON-serializable value
```

### MemoryStorage

```go
func NewMemoryStorage() *MemoryStorage
```

Thread-safe in-memory implementation. Suitable for development, testing, and
single-instance deployments that do not require persistence across restarts.

## TurnState (hosting/core/app/state)

Import path: `github.com/ameena3/Agents-for-golang/hosting/core/app/state`

`TurnState[AppStateT]` manages the three state scopes for a turn. `AgentApplication`
creates and manages this automatically when `Storage` is configured; you do not need
to use it directly in most cases.

```go
type TurnState[AppStateT any] struct {
    App AppStateT // application-specific state, threaded through handlers
}

func NewTurnState[AppStateT any](store storage.Storage) *TurnState[AppStateT]
func (ts *TurnState[AppStateT]) Load(ctx context.Context, channelID, convID, userID string) error
func (ts *TurnState[AppStateT]) Save(ctx context.Context, channelID, convID, userID string) error
func (ts *TurnState[AppStateT]) Conversation() *ConversationState
func (ts *TurnState[AppStateT]) User() *UserState
func (ts *TurnState[AppStateT]) Temp() *TempState
```

### ConversationState

Persisted per conversation across turns.

```go
func (s *ConversationState) Get(key string) (interface{}, bool)
func (s *ConversationState) Set(key string, value interface{})
func (s *ConversationState) Delete(key string)
func (s *ConversationState) Load(ctx context.Context, store storage.Storage, channelID, convID string) error
func (s *ConversationState) Save(ctx context.Context, store storage.Storage, channelID, convID string) error
func (s *ConversationState) Clear()
```

### UserState

Persisted per user across conversations.

```go
func (s *UserState) Get(key string) (interface{}, bool)
func (s *UserState) Set(key string, value interface{})
func (s *UserState) Delete(key string)
func (s *UserState) Load(ctx context.Context, store storage.Storage, channelID, userID string) error
func (s *UserState) Save(ctx context.Context, store storage.Storage, channelID, userID string) error
func (s *UserState) Clear()
```

### TempState

Ephemeral state, exists only for the current turn.

```go
func (s *TempState) Get(key string) (interface{}, bool)
func (s *TempState) Set(key string, value interface{})
func (s *TempState) Delete(key string)
func (s *TempState) Keys() []string
```

## Constants

```go
const InvokeResponseKey = "TurnContext.InvokeResponse"
```

Used to store an invoke response in `TurnState` so the adapter can return it as
the HTTP response body.
