# Microsoft 365 Agents SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/ameena3/Agents-for-golang.svg)](https://pkg.go.dev/github.com/ameena3/Agents-for-golang)
[![Build Status](https://img.shields.io/github/actions/workflow/status/microsoft/agents-sdk-go/ci.yml)](https://github.com/ameena3/Agents-for-golang/actions)

A Go SDK for building enterprise-grade conversational agents for Microsoft 365,
Teams, Copilot Studio, and other channels.

## Overview

The Microsoft 365 Agents SDK for Go provides everything needed to build, host, and
deploy conversational agents on Microsoft platforms. It follows a layered architecture
adapted to Go idioms: interfaces instead of abstract classes, struct embedding instead
of inheritance, `context.Context` for request lifetimes, and standard `error` returns.

## Requirements

- Go 1.22 or later (1.24+ recommended)
- An Azure Bot registration (App ID + secret, certificate, or Managed Identity)

## Quick Start

### Option 1: AgentApplication (recommended, functional style)

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

type MyState struct {
    Count int `json:"count"`
}

func main() {
    store := storage.NewMemoryStorage()

    agent := app.New[MyState](app.AppOptions[MyState]{Storage: store})

    agent.OnMembersAdded(func(ctx context.Context, tc *core.TurnContext, s MyState) error {
        _, err := tc.SendActivity(ctx, core.Text("Hello! Send me a message."))
        return err
    })

    agent.OnMessage("", func(ctx context.Context, tc *core.TurnContext, s MyState) error {
        s.Count++
        _, err := tc.SendActivity(ctx, core.Text("You said: "+tc.Activity().Text))
        return err
    })

    log.Fatal(nethttp.StartAgentProcess(context.Background(), agent, nethttp.ServerConfig{
        Port:                 3978,
        AllowUnauthenticated: true, // false in production
    }))
}
```

### Option 2: ActivityHandler (inheritance style)

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/ameena3/Agents-for-golang/hosting/core"
    "github.com/ameena3/Agents-for-golang/hosting/nethttp"
)

type EchoAgent struct {
    core.ActivityHandler
}

func (a *EchoAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error {
    _, err := tc.SendActivity(ctx, core.Text("Echo: "+tc.Activity().Text))
    return err
}

func main() {
    adapter := nethttp.NewCloudAdapter(true)
    agent := &EchoAgent{}
    http.HandleFunc("/api/messages", nethttp.MessageHandler(adapter, agent))
    log.Fatal(http.ListenAndServe(":3978", nil))
}
```

## Installation

Install individual packages as needed:

```bash
# Core packages (required)
go get github.com/ameena3/Agents-for-golang/activity
go get github.com/ameena3/Agents-for-golang/hosting/core
go get github.com/ameena3/Agents-for-golang/hosting/nethttp

# Teams extension
go get github.com/ameena3/Agents-for-golang/hosting/teams

# Authentication
go get github.com/ameena3/Agents-for-golang/authentication

# Storage backends
go get github.com/ameena3/Agents-for-golang/storage/blob
go get github.com/ameena3/Agents-for-golang/storage/cosmos

# Copilot Studio client
go get github.com/ameena3/Agents-for-golang/copilotstudio/client
```

## Packages

| Import Path | Description |
|---|---|
| `github.com/ameena3/Agents-for-golang/activity` | Activity protocol types, constructors, and constants |
| `github.com/ameena3/Agents-for-golang/activity/card` | Card types: HeroCard, AdaptiveCard, OAuthCard, and more |
| `github.com/ameena3/Agents-for-golang/activity/teams` | Teams-specific activity types and data structures |
| `github.com/ameena3/Agents-for-golang/hosting/core` | Core agent runtime: Agent, TurnContext, ActivityHandler, middleware |
| `github.com/ameena3/Agents-for-golang/hosting/core/app` | AgentApplication — modern functional-style agent framework |
| `github.com/ameena3/Agents-for-golang/hosting/core/storage` | Storage interface and in-memory implementation |
| `github.com/ameena3/Agents-for-golang/hosting/nethttp` | net/http adapter: CloudAdapter, MessageHandler, StartAgentProcess |
| `github.com/ameena3/Agents-for-golang/hosting/teams` | TeamsActivityHandler with Teams-specific routing |
| `github.com/ameena3/Agents-for-golang/authentication` | MSAL-based OAuth: client secret, certificate, managed identity |
| `github.com/ameena3/Agents-for-golang/storage/blob` | Azure Blob Storage state persistence |
| `github.com/ameena3/Agents-for-golang/storage/cosmos` | Azure Cosmos DB state persistence |
| `github.com/ameena3/Agents-for-golang/copilotstudio/client` | Direct-to-Engine client for Copilot Studio agents |

## Examples

Ready-to-run examples are in the `examples/` directory:

| Example | Description |
|---|---|
| `examples/echo-agent/` | Minimal echo agent using AgentApplication with in-memory state |
| `examples/teams-agent/` | Teams agent with task modules and messaging extensions |
| `examples/agent-to-agent/` | Orchestrator agent that delegates to skill agents |

Run an example:

```bash
go run ./examples/echo-agent/
```

Then connect using the [M365 Agents Playground](https://github.com/OfficeDev/microsoft-365-agents-toolkit)
or Azure Bot Service.

## Documentation

Full documentation is in the [`docs/`](./docs/) directory:

- [Getting Started](./docs/getting-started.md)
- [Architecture Overview](./docs/architecture.md)
- [API Reference — activity](./docs/api-reference/activity.md)
- [API Reference — hosting/core](./docs/api-reference/hosting-core.md)
- [API Reference — hosting/teams](./docs/api-reference/hosting-teams.md)
- [API Reference — authentication](./docs/api-reference/authentication.md)
- [API Reference — storage/blob](./docs/api-reference/storage-blob.md)
- [API Reference — storage/cosmos](./docs/api-reference/storage-cosmos.md)
- [API Reference — copilotstudio/client](./docs/api-reference/copilotstudio.md)

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to
agree to a Contributor License Agreement (CLA). For details, visit
https://cla.opensource.microsoft.com.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
See the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any questions.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.

## Trademarks

This project may contain trademarks or logos for projects, products, or services.
Authorized use of Microsoft trademarks or logos is subject to and must follow
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Any use of third-party trademarks or logos are subject to those third-party's policies.
