# Example Walkthrough: Teams Agent

Source: `examples/teams-agent/main.go`

The Teams agent example demonstrates:
- `TeamsActivityHandler` with struct embedding
- `OnTurn` delegation to the Teams pipeline
- Teams-specific handlers: `OnTeamsMembersAdded`, `OnTeamsTaskModuleFetch`, `OnTeamsTaskModuleSubmit`
- Manual `http.HandleFunc` + `nethttp.MessageHandler` setup (alternative to `StartAgentProcess`)

## Full Source

```go
// teams-agent demonstrates Teams-specific activity handling including
// task modules, messaging extensions, and meeting events.
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/ameena3/Agents-for-golang/activity"
    teamstypes "github.com/ameena3/Agents-for-golang/activity/teams"
    "github.com/ameena3/Agents-for-golang/hosting/core"
    "github.com/ameena3/Agents-for-golang/hosting/nethttp"
    hostingteams "github.com/ameena3/Agents-for-golang/hosting/teams"
)

type TeamsAgent struct {
    hostingteams.TeamsActivityHandler
}

func (a *TeamsAgent) OnTurn(ctx context.Context, tc *core.TurnContext) error {
    return a.TeamsActivityHandler.OnTurn(ctx, tc)
}

func (a *TeamsAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error {
    _, err := tc.SendActivity(ctx, core.Text("Teams Echo: "+tc.Activity().Text))
    return err
}

func (a *TeamsAgent) OnTeamsMembersAdded(ctx context.Context, members []*activity.ChannelAccount, tc *core.TurnContext) error {
    for _, member := range members {
        if member.ID != tc.Activity().Recipient.ID {
            if _, err := tc.SendActivity(ctx, core.Text(fmt.Sprintf("Welcome to Teams, %s!", member.Name))); err != nil {
                return err
            }
        }
    }
    return nil
}

func (a *TeamsAgent) OnTeamsTaskModuleFetch(ctx context.Context, req *teamstypes.TaskModuleRequest, tc *core.TurnContext) (*teamstypes.TaskModuleResponse, error) {
    return &teamstypes.TaskModuleResponse{
        Task: &teamstypes.TaskModuleContinueResponse{
            Type: "continue",
            Value: &teamstypes.TaskModuleTaskInfo{
                Title:  "Sample Task",
                Height: 200,
                Width:  400,
                URL:    "https://example.com/taskmodule",
            },
        },
    }, nil
}

func (a *TeamsAgent) OnTeamsTaskModuleSubmit(ctx context.Context, req *teamstypes.TaskModuleRequest, tc *core.TurnContext) (*teamstypes.TaskModuleResponse, error) {
    data, _ := json.Marshal(req.Data)
    _, err := tc.SendActivity(ctx, core.Text("Task submitted: "+string(data)))
    return nil, err
}

func main() {
    agent := &TeamsAgent{}
    adapter := nethttp.NewCloudAdapter(true)

    port := 3978
    if p := os.Getenv("PORT"); p != "" {
        fmt.Sscanf(p, "%d", &port)
    }

    http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, agent))
    log.Printf("Teams agent listening on :%d", port)
    if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
        log.Fatal(err)
    }
}
```

## Annotated Walkthrough

### 1. Struct Embedding

```go
type TeamsAgent struct {
    hostingteams.TeamsActivityHandler
}
```

`TeamsActivityHandler` is embedded (not a field name). This is Go's way of composing
behavior. `TeamsAgent` inherits all the default no-op handlers from
`TeamsActivityHandler`, which in turn inherits from `core.ActivityHandler`.

The embedding order matters for method resolution. When you override `OnMessageActivity`
on `TeamsAgent`, Go's method resolution finds `TeamsAgent.OnMessageActivity` first, then
falls back through the embedding chain.

### 2. The Critical OnTurn Override

```go
func (a *TeamsAgent) OnTurn(ctx context.Context, tc *core.TurnContext) error {
    return a.TeamsActivityHandler.OnTurn(ctx, tc)
}
```

This override is required. Without it, calling `a.OnTurn(ctx, tc)` would dispatch to
`core.ActivityHandler.OnTurn` via the embedding chain. That would route invoke
activities to `core.ActivityHandler.OnInvokeActivity` — the base no-op — instead of
to `TeamsActivityHandler.OnInvokeActivity`, which contains the Teams-specific
`task/fetch`, `task/submit`, etc. routing.

By explicitly delegating to `a.TeamsActivityHandler.OnTurn`, you ensure that the
Teams routing layer is active.

### 3. Message Handling

```go
func (a *TeamsAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error {
    _, err := tc.SendActivity(ctx, core.Text("Teams Echo: "+tc.Activity().Text))
    return err
}
```

Standard message handler. The handler is called for all `message` activities on any
channel, not just Teams. For channel-specific behavior, check `tc.Activity().ChannelID`.

### 4. Teams Members Added

```go
func (a *TeamsAgent) OnTeamsMembersAdded(
    ctx context.Context,
    members []*activity.ChannelAccount,
    tc *core.TurnContext,
) error {
    for _, member := range members {
        if member.ID != tc.Activity().Recipient.ID {
            if _, err := tc.SendActivity(ctx, core.Text(fmt.Sprintf("Welcome to Teams, %s!", member.Name))); err != nil {
                return err
            }
        }
    }
    return nil
}
```

`OnTeamsMembersAdded` is only called for `conversationUpdate` activities on the
`msteams` channel where `MembersAdded` is non-empty.

Contrast with `core.ActivityHandler.OnMembersAdded`, which is called for all channels.
Use `OnTeamsMembersAdded` when you want Teams-only welcome behavior.

The `member.ID != tc.Activity().Recipient.ID` guard prevents sending a welcome to the
bot itself when it is added to the conversation.

### 5. Task Module Fetch

```go
func (a *TeamsAgent) OnTeamsTaskModuleFetch(
    ctx context.Context,
    req *teamstypes.TaskModuleRequest,
    tc *core.TurnContext,
) (*teamstypes.TaskModuleResponse, error) {
    return &teamstypes.TaskModuleResponse{
        Task: &teamstypes.TaskModuleContinueResponse{
            Type: "continue",
            Value: &teamstypes.TaskModuleTaskInfo{
                Title:  "Sample Task",
                Height: 200,
                Width:  400,
                URL:    "https://example.com/taskmodule",
            },
        },
    }, nil
}
```

Called when the user triggers a task module (dialog). The handler returns a
`TaskModuleResponse` with `Task` set to a `TaskModuleContinueResponse`, which opens
a dialog rendered from the given URL.

`TaskModuleTaskInfo` fields:
- `Title`: dialog title bar text
- `Height` / `Width`: dialog dimensions in pixels (or `"small"`, `"medium"`, `"large"`)
- `URL`: URL of the page to render inside the dialog
- `Card`: alternatively, an `Attachment` containing an Adaptive Card to render

### 6. Task Module Submit

```go
func (a *TeamsAgent) OnTeamsTaskModuleSubmit(
    ctx context.Context,
    req *teamstypes.TaskModuleRequest,
    tc *core.TurnContext,
) (*teamstypes.TaskModuleResponse, error) {
    data, _ := json.Marshal(req.Data)
    _, err := tc.SendActivity(ctx, core.Text("Task submitted: "+string(data)))
    return nil, err
}
```

Called when the user submits the task module dialog. `req.Data` contains the form
data submitted by the page. In this example it is marshalled to JSON and echoed back
as a message. Returning `nil` for the response closes the dialog.

### 7. HTTP Setup (Manual)

```go
adapter := nethttp.NewCloudAdapter(true)
http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, agent))
log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
```

This example uses the lower-level manual setup instead of `StartAgentProcess`. This
approach gives you full control over the HTTP server and is suitable when you need to:
- Register additional HTTP routes (health checks, webhooks, etc.)
- Configure TLS
- Use a custom HTTP server with timeouts

`nethttp.MessageHandler` returns a standard `http.HandlerFunc` that:
1. Reads and parses the JSON request body
2. Calls `adapter.Process(ctx, request, agent)`
3. Writes the HTTP response (202 for messages, 200 for invoke responses)

## Running the Example

```bash
cd examples/teams-agent
go run .
```

To test locally with Teams, you need a dev tunnel:

```bash
devtunnel host -a my-tunnel
```

Set the Bot's messaging endpoint to `https://<tunnel-url>/api/messages` and
sideload the Teams app manifest.

## Extending the Example

### Add a Messaging Extension

```go
func (a *TeamsAgent) OnTeamsMessagingExtensionQuery(
    ctx context.Context,
    query *teamstypes.MessagingExtensionQuery,
    tc *core.TurnContext,
) error {
    // query.CommandID identifies which command was used
    // Return results via tc.TurnState[InvokeResponseKey]
    return nil
}
```

### Add Meeting Event Handlers

```go
func (a *TeamsAgent) OnTeamsMeetingStart(
    ctx context.Context,
    meeting *teamstypes.MeetingStartEventDetails,
    tc *core.TurnContext,
) error {
    log.Printf("Meeting started at: %s", meeting.StartTime)
    return nil
}

func (a *TeamsAgent) OnTeamsMeetingEnd(
    ctx context.Context,
    meeting *teamstypes.MeetingEndEventDetails,
    tc *core.TurnContext,
) error {
    log.Printf("Meeting ended at: %s", meeting.EndTime)
    return nil
}
```
