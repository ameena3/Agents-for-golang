# Changelog

## [Unreleased]

### Added
- Complete Go 1.24 implementation of the Microsoft 365 Agents SDK
- `activity` package: full Activity protocol with 200+ types, card types, Teams types
- `hosting/core` package: Agent interface, TurnContext, ActivityHandler, AgentApplication[StateT],
  middleware pipeline, state management (ConversationState, UserState, TempState),
  connector clients, MemoryStorage
- `hosting/nethttp` package: standard net/http adapter for serving agent requests
- `hosting/teams` package: TeamsActivityHandler with 20+ Teams-specific handlers, TeamsInfo
- `authentication` package: MSAL-based OAuth (client secret, certificate, managed identity)
- `storage/blob` package: Azure Blob Storage backend
- `storage/cosmos` package: Azure Cosmos DB backend
- `copilotstudio/client` package: Direct-to-Engine client for Copilot Studio agents
- `hosting/core/http` package: HttpRequestProtocol and HttpAdapterBase abstractions
- `hosting/core/app/oauth` package: OAuthFlow, SignInState, AuthHandler
- Complete documentation in `docs/`
- Example applications in `examples/`
- CI workflows for Go (GitHub Actions, Azure DevOps)
- Dev container configuration for Go development
