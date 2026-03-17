# API Reference: hosting/teams

Import path: `github.com/ameena3/Agents-for-golang/hosting/teams`

The `hosting/teams` package extends `ActivityHandler` with Teams-specific invoke and
event routing. Use it when your agent targets Microsoft Teams and needs to handle
task modules, messaging extensions, meeting events, or card actions.

## TeamsActivityHandler

`TeamsActivityHandler` embeds `core.ActivityHandler` and overrides `OnInvokeActivity`,
`OnConversationUpdateActivity`, and `OnEventActivity` to dispatch Teams-specific
events before falling back to the base handler.

```go
type TeamsActivityHandler struct {
    core.ActivityHandler
}
```

### Setup

Embed `TeamsActivityHandler` in your agent struct and override `OnTurn` to delegate
to `TeamsActivityHandler.OnTurn`. This activates the Teams routing pipeline.

```go
type MyAgent struct {
    teams.TeamsActivityHandler
}

func (a *MyAgent) OnTurn(ctx context.Context, tc *core.TurnContext) error {
    return a.TeamsActivityHandler.OnTurn(ctx, tc)
}
```

Without the `OnTurn` override, calls to your struct go directly to
`core.ActivityHandler.OnTurn`, bypassing the Teams-specific dispatch.

## Handler Methods

### Message & Conversation

| Method | Signature | When Called |
|---|---|---|
| `OnMessageActivity` | `(ctx, tc) error` | Inherited from `core.ActivityHandler` — override as normal |
| `OnTeamsMembersAdded` | `(ctx, members []*activity.ChannelAccount, tc) error` | conversationUpdate with members added on msteams channel |
| `OnTeamsMembersRemoved` | `(ctx, members []*activity.ChannelAccount, tc) error` | conversationUpdate with members removed on msteams channel |
| `OnTeamsChannelCreated` | `(ctx, channelInfo *teams.ChannelInfo, teamInfo *teams.TeamInfo, tc) error` | channelCreated event |
| `OnTeamsChannelDeleted` | `(ctx, channelInfo *teams.ChannelInfo, teamInfo *teams.TeamInfo, tc) error` | channelDeleted event |
| `OnTeamsTeamRenamed` | `(ctx, teamInfo *teams.TeamInfo, tc) error` | teamRenamed event |

### Task Modules

Task modules are dialogs launched from Teams cards or command boxes.

| Method | Signature | When Called |
|---|---|---|
| `OnTeamsTaskModuleFetch` | `(ctx, req *teams.TaskModuleRequest, tc) (*teams.TaskModuleResponse, error)` | `invoke` with name `task/fetch` |
| `OnTeamsTaskModuleSubmit` | `(ctx, req *teams.TaskModuleRequest, tc) (*teams.TaskModuleResponse, error)` | `invoke` with name `task/submit` |

Return a `TaskModuleResponse` with a `TaskModuleContinueResponse` to open a dialog,
or a `TaskModuleMessageResponse` to show a completion message.

### Messaging Extensions

| Method | Signature | When Called |
|---|---|---|
| `OnTeamsMessagingExtensionQuery` | `(ctx, query *teams.MessagingExtensionQuery, tc) error` | `composeExtension/query` |
| `OnTeamsMessagingExtensionSelectItem` | `(ctx, query interface{}, tc) error` | `composeExtension/selectItem` |
| `OnTeamsMessagingExtensionSubmitAction` | `(ctx, action interface{}, tc) error` | `composeExtension/submitAction` |
| `OnTeamsMessagingExtensionFetchTask` | `(ctx, action interface{}, tc) error` | `composeExtension/fetchTask` |

### Config

| Method | Signature | When Called |
|---|---|---|
| `OnTeamsConfigFetch` | `(ctx, query interface{}, tc) (*teams.ConfigResponse, error)` | `config/fetch` |
| `OnTeamsConfigSubmit` | `(ctx, query interface{}, tc) (*teams.ConfigResponse, error)` | `config/submit` |

### Meeting Events

| Method | Signature | When Called |
|---|---|---|
| `OnTeamsMeetingStart` | `(ctx, meeting *teams.MeetingStartEventDetails, tc) error` | Event: `application/vnd.microsoft.meetingStart` |
| `OnTeamsMeetingEnd` | `(ctx, meeting *teams.MeetingEndEventDetails, tc) error` | Event: `application/vnd.microsoft.meetingEnd` |

### Card Actions

| Method | Signature | When Called |
|---|---|---|
| `OnTeamsCardActionInvoke` | `(ctx, tc) error` | `invoke` with no name on msteams channel, or `actionableMessage/executeAction` |

## Routing Logic

`TeamsActivityHandler.OnInvokeActivity` dispatches by `Activity().Name`:

| Invoke Name | Handler Called |
|---|---|
| `task/fetch` | `OnTeamsTaskModuleFetch` |
| `task/submit` | `OnTeamsTaskModuleSubmit` |
| `config/fetch` | `OnTeamsConfigFetch` |
| `config/submit` | `OnTeamsConfigSubmit` |
| `composeExtension/query` | `OnTeamsMessagingExtensionQuery` |
| `composeExtension/selectItem` | `OnTeamsMessagingExtensionSelectItem` |
| `composeExtension/submitAction` | `OnTeamsMessagingExtensionSubmitAction` |
| `composeExtension/fetchTask` | `OnTeamsMessagingExtensionFetchTask` |
| `actionableMessage/executeAction` | `OnTeamsCardActionInvoke` |
| `""` on msteams channel | `OnTeamsCardActionInvoke` |
| anything else | `core.ActivityHandler.OnInvokeActivity` (base) |

`OnConversationUpdateActivity` routes by channel ID and event type:

- If `ChannelID == "msteams"` and `MembersAdded` is non-empty → `OnTeamsMembersAdded`
- If `ChannelID == "msteams"` and `MembersRemoved` is non-empty → `OnTeamsMembersRemoved`
- If `eventType == "channelCreated"` → `OnTeamsChannelCreated`
- If `eventType == "channelDeleted"` → `OnTeamsChannelDeleted`
- If `eventType == "teamRenamed"` → `OnTeamsTeamRenamed`
- Otherwise → `core.ActivityHandler.OnConversationUpdateActivity`

`OnEventActivity` routes by activity name when `ChannelID == "msteams"`:

- `application/vnd.microsoft.meetingStart` → `OnTeamsMeetingStart`
- `application/vnd.microsoft.meetingEnd` → `OnTeamsMeetingEnd`
- Otherwise → `core.ActivityHandler.OnEventActivity`

## TeamsInfo

Import path: `github.com/ameena3/Agents-for-golang/hosting/teams`

`TeamsInfo` provides helper functions for querying Teams-specific information via
the connector. It requires a `ConnectorClient` from the turn context's turn state.

```go
// Get details about the team (name, ID, etc.)
func GetTeamDetails(ctx context.Context, tc *core.TurnContext, teamID string) (*teams.TeamInfo, error)

// Get all channels in a team
func GetChannels(ctx context.Context, tc *core.TurnContext, teamID string) ([]*teams.ChannelInfo, error)

// Get members of a channel or conversation
func GetMembers(ctx context.Context, tc *core.TurnContext) ([]*teams.TeamsChannelAccount, error)
```

## Example: Teams Agent with Task Module

```go
type MyTeamsAgent struct {
    hostingteams.TeamsActivityHandler
}

func (a *MyTeamsAgent) OnTurn(ctx context.Context, tc *core.TurnContext) error {
    return a.TeamsActivityHandler.OnTurn(ctx, tc)
}

func (a *MyTeamsAgent) OnMessageActivity(ctx context.Context, tc *core.TurnContext) error {
    _, err := tc.SendActivity(ctx, core.Text("Hello from Teams!"))
    return err
}

func (a *MyTeamsAgent) OnTeamsTaskModuleFetch(
    ctx context.Context,
    req *teamstypes.TaskModuleRequest,
    tc *core.TurnContext,
) (*teamstypes.TaskModuleResponse, error) {
    return &teamstypes.TaskModuleResponse{
        Task: &teamstypes.TaskModuleContinueResponse{
            Type: "continue",
            Value: &teamstypes.TaskModuleTaskInfo{
                Title:  "My Dialog",
                Height: 300,
                Width:  500,
                URL:    "https://myapp.example.com/taskmodule",
            },
        },
    }, nil
}

func (a *MyTeamsAgent) OnTeamsTaskModuleSubmit(
    ctx context.Context,
    req *teamstypes.TaskModuleRequest,
    tc *core.TurnContext,
) (*teamstypes.TaskModuleResponse, error) {
    // req.Data contains the submitted form data
    _, err := tc.SendActivity(ctx, core.Text("Task submitted!"))
    return nil, err
}
```
