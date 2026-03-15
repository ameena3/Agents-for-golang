# Microsoft 365 Agents SDK for Go — Documentation

This directory contains the full documentation for the Go SDK.

## Contents

### Guides

| Document | Description |
|---|---|
| [Getting Started](./getting-started.md) | Step-by-step guide from prerequisites to deployment |
| [Architecture Overview](./architecture.md) | Layered architecture, request flow, and design patterns |

### API Reference

| Document | Package |
|---|---|
| [activity](./api-reference/activity.md) | `github.com/microsoft/agents-sdk-go/activity` |
| [hosting/core](./api-reference/hosting-core.md) | `github.com/microsoft/agents-sdk-go/hosting/core` and `hosting/core/app` |
| [hosting/teams](./api-reference/hosting-teams.md) | `github.com/microsoft/agents-sdk-go/hosting/teams` |
| [authentication](./api-reference/authentication.md) | `github.com/microsoft/agents-sdk-go/authentication` |
| [storage/blob](./api-reference/storage-blob.md) | `github.com/microsoft/agents-sdk-go/storage/blob` |
| [storage/cosmos](./api-reference/storage-cosmos.md) | `github.com/microsoft/agents-sdk-go/storage/cosmos` |
| [copilotstudio/client](./api-reference/copilotstudio.md) | `github.com/microsoft/agents-sdk-go/copilotstudio/client` |

### Examples

| Document | Example |
|---|---|
| [Echo Agent](./examples/echo-agent.md) | Minimal echo agent — annotated walkthrough |
| [Teams Agent](./examples/teams-agent.md) | Teams agent with task modules and messaging extensions |

## Module Path

```
github.com/microsoft/agents-sdk-go
```

## Go Version

Go 1.22 or later. Go 1.24+ recommended.
