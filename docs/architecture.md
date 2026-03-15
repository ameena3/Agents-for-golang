# Architecture Overview

## Layered Architecture

The SDK is organized in discrete layers. Each layer depends only on the layers below
it. You can use only the layers you need.

```
+------------------------------------------------------------------+
| Layer 5: Web Adapter                                             |
|   hosting/nethttp  — net/http CloudAdapter, MessageHandler,      |
|                      StartAgentProcess                           |
+------------------------------------------------------------------+
| Layer 4: Platform Extensions & Storage                           |
|   hosting/teams    — TeamsActivityHandler, Teams-specific routing|
|   storage/blob     — Azure Blob Storage backend                  |
|   storage/cosmos   — Azure Cosmos DB backend                     |
+------------------------------------------------------------------+
| Layer 3: Authentication                                          |
|   authentication   — MsalAuth, ConnectionManager                 |
|                      (client secret / certificate / managed id)  |
+------------------------------------------------------------------+
| Layer 2: Core Hosting Engine                                     |
|   hosting/core     — Agent, TurnContext, ActivityHandler,        |
|                      AgentApplication[StateT], Middleware,        |
|                      ChannelAdapter, Storage, connectors         |
+------------------------------------------------------------------+
| Layer 1: Protocol / Schema                                       |
|   activity         — Activity struct and methods                 |
|   activity/card    — Card types                                  |
|   activity/teams   — Teams-specific types                        |
|   activity/config  — AgentConfig, LoadFromEnv                    |
+------------------------------------------------------------------+
```

## Package Descriptions

| Package | Role |
|---|---|
| `activity` | Defines `Activity` and all supporting types. No dependencies outside the standard library. |
| `activity/card` | Card content types: HeroCard, AdaptiveCard, OAuthCard, etc. |
| `activity/teams` | Teams-specific data structures: TeamsChannelData, TaskModule, MessagingExtension, Meeting. |
| `activity/config` | `AgentConfig` and `LoadFromEnv` for reading Bot credentials from environment. |
| `hosting/core` | The runtime engine: `TurnContext`, `ActivityHandler`, `ChannelServiceAdapter`, middleware pipeline, `ChannelAdapter` interface. |
| `hosting/core/app` | `AgentApplication[StateT]` — modern, functional, generic agent framework with route matching. |
| `hosting/core/storage` | `Storage` interface and `MemoryStorage` in-process implementation. |
| `hosting/core/connector` | HTTP connector clients: `ConnectorClient`, `UserTokenClient`, `TeamsConnectorClient`, factory. |
| `hosting/core/authorization` | JWT validation, `ClaimsIdentity`, auth constants. |
| `hosting/nethttp` | `CloudAdapter` for `net/http`, `MessageHandler` `http.HandlerFunc`, `StartAgentProcess`. |
| `hosting/teams` | `TeamsActivityHandler` — extends `ActivityHandler` with Teams invoke and event routing. |
| `authentication` | `MsalAuth` implementing `core.AccessTokenProvider` via MSAL for Go and `azidentity`. |
| `storage/blob` | `BlobStorage` — implements `storage.Storage` using Azure Blob Storage. |
| `storage/cosmos` | `CosmosDBStorage` — implements `storage.Storage` using Azure Cosmos DB. |
| `copilotstudio/client` | `CopilotClient` — Direct-to-Engine client for Copilot Studio agents. |

## Request Flow

```
HTTP POST /api/messages
        |
        v
+------------------+
| net/http server  |  (standard library or any Go HTTP framework)
+------------------+
        |
        v
+------------------+
| nethttp.Message  |  MessageHandler(adapter, agent) http.HandlerFunc
| Handler          |  reads request body, calls adapter.Process()
+------------------+
        |
        v
+---------------------------+
| CloudAdapter.Process()    |  parses JSON → *activity.Activity
|                           |  validates JWT (if not AllowUnauthenticated)
|                           |  creates ClaimsIdentity
+---------------------------+
        |
        v
+-----------------------------------+
| ChannelServiceAdapter             |
| .ProcessActivity()                |  creates TurnContext(adapter, activity, identity)
|                                   |  runs MiddlewareSet pipeline
+-----------------------------------+
        |
        v (middleware chain runs here)
        |
        v
+---------------------+
| Agent.OnTurn()      |  dispatched to your agent (ActivityHandler or AgentApplication)
+---------------------+
        |
        +--- ActivityHandler:
        |        switch activity.Type
        |        → OnMessageActivity / OnConversationUpdateActivity / ...
        |
        +--- AgentApplication[StateT]:
                 1. BeforeTurn hooks
                 2. Load state from Storage
                 3. RouteList.FindRoute() → matched handler
                 4. handler(ctx, tc, appState)
                 5. Save state to Storage
                 6. AfterTurn hooks
        |
        v
tc.SendActivity(ctx, reply)
        |
        v (send handler chain)
        |
        v
ChannelAdapter.SendActivities()
        |
        v
HTTP POST to ServiceURL (channel)
        |
        v
Return 202 Accepted (message) or 200 OK with body (invoke)
```

## Key Interfaces

### Agent

Every agent must implement `core.Agent`:

```go
type Agent interface {
    OnTurn(ctx context.Context, turnCtx *TurnContext) error
}
```

Both `ActivityHandler` (via embedding) and `AgentApplication` implement `Agent`.

### Middleware

Cross-cutting concerns implement `core.Middleware`:

```go
type Middleware interface {
    OnTurn(ctx context.Context, turnCtx *TurnContext, next func(context.Context) error) error
}
```

Call `next` to continue the pipeline. Return an error to abort.

Register middleware on an adapter:

```go
adapter.Use(myMiddleware1, myMiddleware2)
```

### Storage

All persistence backends implement `storage.Storage`:

```go
type Storage interface {
    Read(ctx context.Context, keys []string) (map[string]StoreItem, error)
    Write(ctx context.Context, changes map[string]StoreItem) error
    Delete(ctx context.Context, keys []string) error
}
```

Implementations: `MemoryStorage`, `BlobStorage`, `CosmosDBStorage`.

### ChannelAdapter

The adapter layer that sends and receives activities:

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

For acquiring OAuth tokens:

```go
type AccessTokenProvider interface {
    GetAccessToken(ctx context.Context, resourceURL string, scopes []string, forceRefresh bool) (string, error)
    AcquireTokenOnBehalfOf(ctx context.Context, scopes []string, userAssertion string) (string, error)
}
```

`MsalAuth` implements this interface.

## Design Patterns

### Interface-Oriented Design

All major abstractions are Go interfaces, enabling substitution and testing with mocks.
Unlike the Python SDK, which uses abstract base classes and Protocol types, the Go SDK
relies on Go's structural typing: any struct that has the required methods satisfies
the interface.

### Struct Embedding (Inheritance Substitute)

The Go SDK uses struct embedding to approximate the Python inheritance model:

```go
// Go: embed ActivityHandler, override methods
type MyAgent struct {
    core.ActivityHandler
}

func (a *MyAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error { ... }
```

This is equivalent to Python's:

```python
class MyAgent(ActivityHandler):
    async def on_message_activity(self, context: TurnContext): ...
```

### Generics for Type-Safe State

`AgentApplication[StateT]` and `TurnState[AppStateT]` use Go generics to carry
application state through the pipeline without `interface{}` casts:

```go
// State type is declared once at construction
agent := app.New[MyState](app.AppOptions[MyState]{Storage: store})

// All handlers receive MyState directly
agent.OnMessage("", func(ctx context.Context, tc *core.TurnContext, s MyState) error {
    s.Count++ // strongly typed
    return nil
})
```

### Middleware Pipeline

The middleware pipeline follows a chain-of-responsibility pattern identical to the
Python SDK, adapted to Go:

```go
// Python
async def on_turn(self, context, next):
    # before
    await next()
    # after

// Go
func (m *MyMiddleware) OnTurn(ctx context.Context, tc *core.TurnContext, next func(context.Context) error) error {
    // before
    if err := next(ctx); err != nil {
        return err
    }
    // after
    return nil
}
```

### Context Instead of Async/Await

The Python SDK uses `async/await` throughout. The Go SDK uses `context.Context` as
the first parameter of every function that performs I/O, and returns `error` as the
last return value. Cancellation and deadlines propagate through the context.

## Python vs Go Comparison

| Python SDK | Go SDK |
|---|---|
| `async def on_message_activity(self, context)` | `func (a *MyAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error` |
| `await context.send_activity(...)` | `tc.SendActivity(ctx, ...)` |
| `class MyAgent(ActivityHandler):` | `type MyAgent struct { core.ActivityHandler }` |
| `app = AgentApplication[TurnState]()` | `agent := app.New[MyState](app.AppOptions[MyState]{})` |
| `@app.message()` decorator | `agent.OnMessage("", handler)` method call |
| Pydantic models with validators | Plain Go structs with `json` tags |
| `asyncio` event loop | Standard goroutines + `context.Context` |
| Python `exceptions` | Go `error` return values |
| `pip install microsoft-agents-*` | `go get github.com/microsoft/agents-sdk-go/...` |
| Multiple packages in `libraries/` | Multiple packages in one module |

## State Management

`AgentApplication` manages three state scopes automatically when `Storage` is
configured:

| Scope | Key Pattern | Lifetime |
|---|---|---|
| Conversation | `conv/{channelID}/{conversationID}` | Per conversation |
| User | `user/{channelID}/{userID}` | Per user across conversations |
| Temp | In-memory only | Current turn only |

State is loaded at the start of `OnTurn` and saved after the route handler returns.
The application-level typed state (`AppState`) is accessible as the third parameter
of every handler function.
