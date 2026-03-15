# Changelog

## [Unreleased] - Go SDK

### Added
- Complete Go 1.24 implementation of the Microsoft 365 Agents SDK
- `activity` package: full Activity protocol with 200+ types, card types, Teams types
- `hosting/core` package: Agent interface, TurnContext, ActivityHandler, AgentApplication[StateT],
  middleware pipeline, state management (ConversationState, UserState, TempState),
  connector clients, MemoryStorage
- `hosting/nethttp` package: net/http adapter replacing Python aiohttp/FastAPI adapters
- `hosting/teams` package: TeamsActivityHandler with 20+ Teams-specific handlers, TeamsInfo
- `authentication` package: MSAL-based OAuth (client secret, certificate, managed identity)
- `storage/blob` package: Azure Blob Storage backend
- `storage/cosmos` package: Azure Cosmos DB backend
- `copilotstudio/client` package: Direct-to-Engine client for Copilot Studio agents
- Complete documentation in `docs/`
- Example applications in `examples/`

### Changed
- Converted from Python to Go
- Replaced Python async/await with Go context + goroutines
- Replaced Pydantic models with Go structs + JSON tags
- Replaced Python Protocols with Go interfaces
- Replaced Python generics (TypeVar) with Go generics (1.18+)
- Replaced aiohttp/FastAPI with standard net/http

### Removed
- All Python source files and packages
- Python-specific configuration (pytest.ini, flake8, pre-commit)
