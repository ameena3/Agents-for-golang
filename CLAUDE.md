# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Repository Overview

This is the **Microsoft 365 Agents SDK for Go**, a framework for building enterprise-grade
conversational agents for M365, Teams, Copilot Studio, and other platforms.

**Module**: `github.com/microsoft/agents-sdk-go`

## Development Setup

### Requirements
- Go 1.26+

### Setup
```bash
go mod download
```

### Testing
```bash
go test ./...
go test ./activity/...
go test ./hosting/core/...
```

### Building
```bash
go build ./...
```

### Linting
```bash
go vet ./...
```

## Architecture

Layer 5: hosting/nethttp (HTTP adapter)
Layer 4: hosting/teams, storage/blob, storage/cosmos, authentication
Layer 3: copilotstudio/client
Layer 2: hosting/core (Agent, TurnContext, ActivityHandler, AgentApplication, State)
Layer 1: activity (Activity protocol types)

## Key Patterns

- Go interfaces for all major abstractions (no inheritance)
- Go generics for AgentApplication[StateT]
- context.Context for cancellation/timeouts
- error returns (no exceptions)
- JSON struct tags for serialization

## Common Pitfalls

1. Always pass context.Context as the first parameter
2. Check errors from all SDK calls
3. State is auto-saved only if Storage is configured in AppOptions
4. TurnContext is scoped to a single turn
5. Middleware runs in registration order, reverses on return
