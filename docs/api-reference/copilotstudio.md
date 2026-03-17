# API Reference: copilotstudio/client

Import path: `github.com/ameena3/Agents-for-golang/copilotstudio/client`

The `copilotstudio/client` package provides a Direct-to-Engine HTTP client for
communicating with agents created in Microsoft Copilot Studio.

## CopilotClient

```go
type CopilotClient struct { /* unexported */ }
```

`CopilotClient` communicates with a Copilot Studio agent via the Direct-to-Engine
protocol: start a conversation, then send activities and receive responses.

### Constructor

```go
func NewCopilotClient(
    ctx context.Context,
    settings *ConnectionSettings,
    tokenProvider TokenProvider,
) (*CopilotClient, error)
```

- `settings` is required. If `Cloud` or `AgentType` are empty, they default to
  `PowerPlatformCloudPublic` and `AgentTypePublished` respectively.
- `tokenProvider` is optional. If provided, it supplies a Bearer token for each
  request. If nil, the conversation token from `StartConversation` is used.

### Methods

```go
// StartConversation initializes a new conversation with the Copilot Studio agent.
// Must be called before ExecuteTurn.
func (c *CopilotClient) StartConversation(ctx context.Context) (*StartConversationResponse, error)

// ExecuteTurn sends an activity and returns the agent's response activities.
func (c *CopilotClient) ExecuteTurn(ctx context.Context, act *activity.Activity) (*ExecuteTurnResponse, error)

// SendText is a convenience wrapper that sends a plain-text message.
func (c *CopilotClient) SendText(ctx context.Context, text string) (*ExecuteTurnResponse, error)
```

### Response Types

```go
type StartConversationResponse struct {
    ConversationID string `json:"conversationId"`
    Token          string `json:"token"`
    ExpiresIn      int    `json:"expires_in"`
}

type ExecuteTurnRequest struct {
    Activity *activity.Activity `json:"activity"`
}

type ExecuteTurnResponse struct {
    Activities []*activity.Activity `json:"activities"`
    Watermark  string               `json:"watermark"`
}
```

## ConnectionSettings

```go
type ConnectionSettings struct {
    // EnvironmentID is the Power Platform environment ID. Required.
    EnvironmentID string

    // BotID is the bot identifier in Copilot Studio. Required.
    BotID string

    // BotName is the display name (optional).
    BotName string

    // Cloud is the Power Platform cloud region. Defaults to PowerPlatformCloudPublic.
    Cloud PowerPlatformCloud

    // AgentType is Published or Preview. Defaults to AgentTypePublished.
    AgentType AgentType

    // CustomEndpoint overrides the computed Direct-to-Engine URL.
    CustomEndpoint string
}
```

`GetEndpointURL()` method computes the full Direct-to-Engine base URL from `Cloud`
and `EnvironmentID`, then appends the bot path. Override with `CustomEndpoint` for
non-standard deployments.

## Constants

### PowerPlatformCloud

```go
type PowerPlatformCloud string

const (
    PowerPlatformCloudPublic PowerPlatformCloud = "Public"
    // Additional cloud constants defined in power_platform.go
)
```

### AgentType

```go
type AgentType string

const (
    AgentTypePublished AgentType = "Published"
    AgentTypePreview   AgentType = "Preview"
)
```

## TokenProvider Interface

```go
type TokenProvider interface {
    GetAccessToken(ctx context.Context, resourceURL string) (string, error)
}
```

Pass an implementation (e.g., an `*authentication.MsalAuth` wrapped in a thin
adapter) to provide per-request Bearer tokens for the Power Platform API.

## Usage Example

```go
import (
    "context"
    "fmt"
    "log"

    copilot "github.com/ameena3/Agents-for-golang/copilotstudio/client"
)

func main() {
    settings := &copilot.ConnectionSettings{
        EnvironmentID: "your-environment-id",
        BotID:         "your-bot-id",
    }

    client, err := copilot.NewCopilotClient(context.Background(), settings, nil)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Start the conversation
    conv, err := client.StartConversation(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Conversation started:", conv.ConversationID)

    // Send a message and read the response
    resp, err := client.SendText(ctx, "Hello!")
    if err != nil {
        log.Fatal(err)
    }

    for _, act := range resp.Activities {
        if act.Type == "message" {
            fmt.Println("Bot says:", act.Text)
        }
    }
}
```

## Notes

- `StartConversation` must be called before `ExecuteTurn` or `SendText`. The
  conversation ID returned is stored internally and used for subsequent requests.
- Each `CopilotClient` instance represents a single conversation. Create a new
  instance for each new conversation.
- If `tokenProvider` is set, its token is used in preference to the conversation
  token from `StartConversation`.
- The client uses the standard `net/http.Client` with no custom timeout; callers
  should set a deadline on the context or configure a custom `http.Client` by
  wrapping `CopilotClient`.
