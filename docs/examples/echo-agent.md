# Example Walkthrough: Echo Agent

Source: `examples/echo-agent/main.go`

The echo agent is the simplest complete agent example. It demonstrates:
- `AgentApplication[StateT]` with typed per-conversation state
- `OnMembersAdded` for welcome messages
- `OnMessage` for echoing user input
- `storage.MemoryStorage` as the state backend
- `nethttp.StartAgentProcess` for one-line server startup

## Full Source

```go
// echo-agent is a simple Microsoft 365 agent that echoes back any message.
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/ameena3/Agents-for-golang/activity/config"
    "github.com/ameena3/Agents-for-golang/hosting/core"
    "github.com/ameena3/Agents-for-golang/hosting/core/app"
    "github.com/ameena3/Agents-for-golang/hosting/core/storage"
    "github.com/ameena3/Agents-for-golang/hosting/nethttp"
)

// AppState holds per-conversation state for this agent.
type AppState struct {
    MessageCount int `json:"messageCount"`
}

func main() {
    cfg := config.LoadFromEnv()
    _ = cfg

    store := storage.NewMemoryStorage()

    agentApp := app.New[AppState](app.AppOptions[AppState]{
        Storage: store,
    })

    agentApp.OnMembersAdded(func(ctx context.Context, tc *core.TurnContext, state AppState) error {
        for _, member := range tc.Activity().MembersAdded {
            if member.ID != tc.Activity().Recipient.ID {
                _, err := tc.SendActivity(ctx, core.Text("Hello! I'm the Echo Agent. Send me a message and I'll echo it back."))
                return err
            }
        }
        return nil
    })

    agentApp.OnMessage("", func(ctx context.Context, tc *core.TurnContext, state AppState) error {
        state.MessageCount++
        text := tc.Activity().Text
        reply := fmt.Sprintf("Echo (#%d): %s", state.MessageCount, text)
        _, err := tc.SendActivity(ctx, core.Text(reply))
        return err
    })

    port := 3978
    if p := os.Getenv("PORT"); p != "" {
        fmt.Sscanf(p, "%d", &port)
    }

    log.Printf("Echo agent starting on port %d...", port)
    if err := nethttp.StartAgentProcess(context.Background(), agentApp, nethttp.ServerConfig{
        Port:                 port,
        AllowUnauthenticated: true,
    }); err != nil {
        log.Fatal(err)
    }
}
```

## Annotated Walkthrough

### 1. Application State

```go
type AppState struct {
    MessageCount int `json:"messageCount"`
}
```

`AppState` is the strongly-typed state struct for this agent. It is parameterized into
`AgentApplication[AppState]`. The `json` tag is required for serialization to storage.

When `Storage` is configured, `AgentApplication` automatically:
- Loads `AppState` from storage at the start of each turn (keyed by channel + conversation)
- Passes the loaded state as the third argument to every handler
- Saves the state back to storage after the handler returns

Note: in this example, `state.MessageCount++` increments the local copy. The modified
value is saved back to storage automatically. On the next turn, the handler receives
the incremented count. This gives you per-conversation turn counters without any
explicit persistence code.

### 2. In-Memory Storage

```go
store := storage.NewMemoryStorage()
```

`MemoryStorage` is a goroutine-safe in-memory map. State is lost when the process
restarts. It is suitable for:
- Local development and testing
- Single-instance deployments where persistence across restarts is not required

For production, replace with `blob.NewBlobStorage(...)` or
`cosmos.NewCosmosDBStorage(...)`.

### 3. Creating the AgentApplication

```go
agentApp := app.New[AppState](app.AppOptions[AppState]{
    Storage: store,
})
```

`app.New[AppState]` creates a new `AgentApplication` parameterized on `AppState`.
The `AppOptions` struct configures the storage backend. All other options use their
zero-value defaults.

### 4. Handling New Members

```go
agentApp.OnMembersAdded(func(ctx context.Context, tc *core.TurnContext, state AppState) error {
    for _, member := range tc.Activity().MembersAdded {
        if member.ID != tc.Activity().Recipient.ID {
            _, err := tc.SendActivity(ctx, core.Text("Hello! ..."))
            return err
        }
    }
    return nil
})
```

`OnMembersAdded` registers a handler for `conversationUpdate` activities where
`MembersAdded` is non-empty. This fires when a user (or the bot itself) joins a
conversation.

The `member.ID != tc.Activity().Recipient.ID` check skips the bot's own join event.
Without this check, the agent would send a welcome message to itself.

`core.Text(...)` is a convenience function that creates a `*activity.Activity` of
type `message` with the given text string.

### 5. Echoing Messages

```go
agentApp.OnMessage("", func(ctx context.Context, tc *core.TurnContext, state AppState) error {
    state.MessageCount++
    text := tc.Activity().Text
    reply := fmt.Sprintf("Echo (#%d): %s", state.MessageCount, text)
    _, err := tc.SendActivity(ctx, core.Text(reply))
    return err
})
```

`OnMessage("")` with an empty pattern matches all incoming message activities.
The handler:
1. Increments `state.MessageCount` (the mutated value is saved by the framework)
2. Reads the incoming message text from `tc.Activity().Text`
3. Builds a reply string that includes the turn counter
4. Calls `tc.SendActivity` to send the reply back to the channel

`tc.SendActivity` returns `(*activity.ResourceResponse, error)`. The resource
response contains the ID of the sent activity, which is discarded here with `_`.

### 6. Server Startup

```go
nethttp.StartAgentProcess(context.Background(), agentApp, nethttp.ServerConfig{
    Port:                 port,
    AllowUnauthenticated: true,
})
```

`StartAgentProcess` is a one-line server starter. It:
- Creates a `CloudAdapter`
- Registers `POST /api/messages` using `MessageHandler`
- Calls `http.ListenAndServe`

`AllowUnauthenticated: true` skips JWT validation, which is required for local testing
without real Azure Bot credentials. Set this to `false` and provide a real credential
for production deployments.

## Running the Example

```bash
cd examples/echo-agent
go run .
```

Then connect with the M365 Agents Playground pointing to `http://localhost:3978/api/messages`.

## Production Adaptations

1. Replace `storage.NewMemoryStorage()` with `blob.NewBlobStorage(...)`:
   ```go
   store, _ := blob.NewBlobStorage(ctx, blob.Config{
       ConnectionString: os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
       ContainerName:    "agent-state",
   })
   ```

2. Set `AllowUnauthenticated: false` and configure your Bot registration credentials
   via environment variables (`AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`).

3. Add error handling:
   ```go
   agentApp.OnError(func(ctx context.Context, tc *core.TurnContext, err error) error {
       log.Printf("turn error: %v", err)
       return nil
   })
   ```
