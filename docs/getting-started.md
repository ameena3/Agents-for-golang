# Getting Started

This guide walks you through building your first agent with the Microsoft 365 Agents
SDK for Go, from prerequisites to a running service.

## Prerequisites

1. **Go 1.22 or later.** Verify with:
   ```bash
   go version
   ```

2. **An Azure subscription** with permission to create App Registrations.

3. **Azure Bot registration.** In the Azure Portal:
   - Create a new "Azure Bot" resource (or use an existing one).
   - Note the **Application (client) ID** and **Tenant ID**.
   - Under "Configuration", create a **Client Secret** (or upload a certificate).
   - Set the messaging endpoint to your tunnel URL + `/api/messages`
     (see step 6 for local testing with a dev tunnel).

4. **Optional: M365 Agents Playground** for local testing without Teams:
   - Install from https://github.com/OfficeDev/microsoft-365-agents-toolkit

## Step 1: Create a New Agent Project

```bash
mkdir my-agent && cd my-agent
go mod init example.com/my-agent
go get github.com/ameena3/Agents-for-golang/activity
go get github.com/ameena3/Agents-for-golang/hosting/core
go get github.com/ameena3/Agents-for-golang/hosting/core/app
go get github.com/ameena3/Agents-for-golang/hosting/core/storage
go get github.com/ameena3/Agents-for-golang/hosting/nethttp
```

## Step 2: Implement ActivityHandler (simple inheritance style)

`ActivityHandler` is the simpler of the two programming models. Embed it in a struct
and override only the activity methods you need. The base implementation is a no-op
for all handlers, so you only write the code you need.

```go
// main.go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/ameena3/Agents-for-golang/hosting/core"
    "github.com/ameena3/Agents-for-golang/hosting/nethttp"
)

// MyAgent embeds ActivityHandler to get default no-op behavior for all activity types.
type MyAgent struct {
    core.ActivityHandler
}

// OnMessageActivity is called for every incoming message.
func (a *MyAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error {
    text := tc.Activity().Text
    _, err := tc.SendActivity(ctx, core.Text("You said: "+text))
    return err
}

func main() {
    adapter := nethttp.NewCloudAdapter(true) // true = allow unauthenticated for local dev
    agent := &MyAgent{}

    http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, agent))
    log.Println("Agent listening on :3978")
    log.Fatal(http.ListenAndServe(":3978", nil))
}
```

Run it:
```bash
go run .
```

## Step 3: Implement AgentApplication (modern functional style)

`AgentApplication` is the recommended model for new projects. It uses registered
handler functions and a type-parameterized state struct, avoiding the need to embed
and override.

```go
package main

import (
    "context"
    "log"

    "github.com/ameena3/Agents-for-golang/hosting/core"
    "github.com/ameena3/Agents-for-golang/hosting/core/app"
    "github.com/ameena3/Agents-for-golang/hosting/core/storage"
    "github.com/ameena3/Agents-for-golang/hosting/nethttp"
)

// AppState is persisted across turns for each conversation.
type AppState struct {
    MessageCount int    `json:"messageCount"`
    LastMessage  string `json:"lastMessage"`
}

func main() {
    // In-memory storage for development. Replace with BlobStorage in production.
    store := storage.NewMemoryStorage()

    agent := app.New[AppState](app.AppOptions[AppState]{
        Storage: store,
    })

    // Greet new members when they join the conversation.
    agent.OnMembersAdded(func(ctx context.Context, tc *core.TurnContext, s AppState) error {
        for _, member := range tc.Activity().MembersAdded {
            if member.ID != tc.Activity().Recipient.ID {
                _, err := tc.SendActivity(ctx, core.Text("Welcome! Send me a message."))
                if err != nil {
                    return err
                }
            }
        }
        return nil
    })

    // Handle all incoming messages.
    agent.OnMessage("", func(ctx context.Context, tc *core.TurnContext, s AppState) error {
        s.MessageCount++
        s.LastMessage = tc.Activity().Text
        _, err := tc.SendActivity(ctx, core.Text("You said: "+tc.Activity().Text))
        return err
    })

    // Handle errors.
    agent.OnError(func(ctx context.Context, tc *core.TurnContext, err error) error {
        log.Printf("error in turn: %v", err)
        return nil // error handled; returning nil stops propagation
    })

    log.Fatal(nethttp.StartAgentProcess(context.Background(), agent, nethttp.ServerConfig{
        Port:                 3978,
        AllowUnauthenticated: true,
    }))
}
```

### Pattern matching with OnMessage

`OnMessage` accepts an optional regular expression pattern. Only messages whose text
matches the pattern are dispatched to that handler. Register multiple handlers with
different patterns to build a simple command router:

```go
agent.OnMessage(`^hello`, func(ctx context.Context, tc *core.TurnContext, s AppState) error {
    _, err := tc.SendActivity(ctx, core.Text("Hello back!"))
    return err
})

agent.OnMessage(`^help`, func(ctx context.Context, tc *core.TurnContext, s AppState) error {
    _, err := tc.SendActivity(ctx, core.Text("Available commands: hello, help"))
    return err
})

// Catch-all for unmatched messages.
agent.OnMessage("", func(ctx context.Context, tc *core.TurnContext, s AppState) error {
    _, err := tc.SendActivity(ctx, core.Text("I didn't understand that."))
    return err
})
```

## Step 4: Register Routes and Handle Teams Messages

For Teams-specific scenarios, use `TeamsActivityHandler`:

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/ameena3/Agents-for-golang/activity"
    teamstypes "github.com/ameena3/Agents-for-golang/activity/teams"
    "github.com/ameena3/Agents-for-golang/hosting/core"
    "github.com/ameena3/Agents-for-golang/hosting/nethttp"
    hostingteams "github.com/ameena3/Agents-for-golang/hosting/teams"
)

type TeamsAgent struct {
    hostingteams.TeamsActivityHandler
}

// OnTurn must delegate to TeamsActivityHandler.OnTurn to activate Teams routing.
func (a *TeamsAgent) OnTurn(ctx context.Context, tc *core.TurnContext) error {
    return a.TeamsActivityHandler.OnTurn(ctx, tc)
}

func (a *TeamsAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error {
    _, err := tc.SendActivity(ctx, core.Text("Teams Echo: "+tc.Activity().Text))
    return err
}

func (a *TeamsAgent) OnTeamsMembersAdded(
    ctx context.Context,
    members []*activity.ChannelAccount,
    tc *core.TurnContext,
) error {
    for _, m := range members {
        if m.ID != tc.Activity().Recipient.ID {
            if _, err := tc.SendActivity(ctx, core.Text("Welcome, "+m.Name+"!")); err != nil {
                return err
            }
        }
    }
    return nil
}

func (a *TeamsAgent) OnTeamsTaskModuleFetch(
    ctx context.Context,
    req *teamstypes.TaskModuleRequest,
    tc *core.TurnContext,
) (*teamstypes.TaskModuleResponse, error) {
    return &teamstypes.TaskModuleResponse{
        Task: &teamstypes.TaskModuleContinueResponse{
            Type: "continue",
            Value: &teamstypes.TaskModuleTaskInfo{
                Title:  "My Task",
                Height: 200,
                Width:  400,
                URL:    "https://example.com/task",
            },
        },
    }, nil
}

func main() {
    adapter := nethttp.NewCloudAdapter(true)
    agent := &TeamsAgent{}
    http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, agent))
    log.Fatal(http.ListenAndServe(":3978", nil))
}
```

## Step 5: Configure Authentication

For production, replace `AllowUnauthenticated: true` with real MSAL authentication.

```bash
go get github.com/ameena3/Agents-for-golang/authentication
```

```go
import "github.com/ameena3/Agents-for-golang/authentication"

auth, err := authentication.NewMsalAuth(authentication.Config{
    TenantID:     "your-tenant-id",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})
if err != nil {
    log.Fatal(err)
}
// auth implements core.AccessTokenProvider
```

You can also load credentials from standard environment variables:

```bash
export AZURE_TENANT_ID=your-tenant-id
export AZURE_CLIENT_ID=your-client-id
export AZURE_CLIENT_SECRET=your-client-secret
```

```go
auth, err := authentication.NewMsalAuthFromEnv()
```

Then set `AllowUnauthenticated: false` in `ServerConfig` and pass `auth` where the
adapter or connector factory needs a token provider.

See [authentication API reference](./api-reference/authentication.md) for certificate
and managed identity options.

## Step 6: Local Testing with Dev Tunnels

To receive messages from the Azure Bot Service locally, expose your local port via a
dev tunnel:

```bash
# Install (Windows)
winget install Microsoft.devtunnel

# Authenticate
devtunnel user login

# Create and host a persistent tunnel on port 3978
devtunnel create my-agent -a
devtunnel port create -p 3978 my-agent
devtunnel host -a my-agent
```

Set the Bot's messaging endpoint to `https://<tunnel-url>/api/messages` in the Azure
Portal.

## Step 7: Deploy to Azure

### Azure Container Apps (recommended)

1. Build and push a Docker image:
   ```bash
   docker build -t myregistry.azurecr.io/my-agent:latest .
   docker push myregistry.azurecr.io/my-agent:latest
   ```

2. Deploy:
   ```bash
   az containerapp create \
     --name my-agent \
     --resource-group my-rg \
     --image myregistry.azurecr.io/my-agent:latest \
     --target-port 3978 \
     --ingress external \
     --env-vars AZURE_TENANT_ID=... AZURE_CLIENT_ID=... AZURE_CLIENT_SECRET=...
   ```

3. Update the Bot's messaging endpoint to the Container App URL + `/api/messages`.

### Environment Variables

The SDK recognizes these standard variables when using `authentication.NewMsalAuthFromEnv()`
and `config.LoadFromEnv()`:

| Variable | Description |
|---|---|
| `AZURE_TENANT_ID` | Azure AD tenant ID |
| `AZURE_CLIENT_ID` | App registration client ID |
| `AZURE_CLIENT_SECRET` | Client secret (or use certificate/managed identity) |
| `AZURE_CERTIFICATE_PATH` | Path to PEM/PFX certificate file |
| `AZURE_CERTIFICATE_PASSWORD` | Certificate password (PFX only) |
| `AZURE_USE_MANAGED_IDENTITY` | Set to `true` for Managed Identity |
| `AZURE_USER_ASSIGNED_CLIENT_ID` | Client ID for user-assigned managed identity |
| `PORT` | HTTP listen port (default: 3978) |
