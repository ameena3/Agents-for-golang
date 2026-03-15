// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

import "strings"

// Activity is the basic communication type for the Microsoft 365 Agents protocol.
// It represents a communication between an agent and a channel.
type Activity struct {
	// Type specifies the type of activity (message, invoke, event, etc.)
	Type string `json:"type,omitempty"`
	// ID is the unique identifier for the activity on the channel.
	ID string `json:"id,omitempty"`
	// Timestamp is the datetime when the message was sent, in UTC, expressed in ISO-8601 format.
	Timestamp string `json:"timestamp,omitempty"`
	// LocalTimestamp is the local datetime of the message expressed in ISO-8601 format.
	LocalTimestamp string `json:"localTimestamp,omitempty"`
	// LocalTimezone is the name of the local timezone of the message in IANA format.
	LocalTimezone string `json:"localTimezone,omitempty"`
	// ServiceURL is the URL that specifies the channel's service endpoint.
	ServiceURL string `json:"serviceUrl,omitempty"`
	// ChannelID is the identifier of the channel.
	ChannelID string `json:"channelId,omitempty"`
	// From identifies the sender of the message.
	From *ChannelAccount `json:"from,omitempty"`
	// Conversation identifies the conversation to which the activity belongs.
	Conversation *ConversationAccount `json:"conversation,omitempty"`
	// Recipient identifies the recipient of the message.
	Recipient *ChannelAccount `json:"recipient,omitempty"`
	// TextFormat is the format of text fields. Possible values: 'markdown', 'plain', 'xml'.
	TextFormat string `json:"textFormat,omitempty"`
	// AttachmentLayout is the layout hint for multiple attachments. Possible values: 'list', 'carousel'.
	AttachmentLayout string `json:"attachmentLayout,omitempty"`
	// MembersAdded is the collection of members added to the conversation (conversationUpdate).
	MembersAdded []*ChannelAccount `json:"membersAdded,omitempty"`
	// MembersRemoved is the collection of members removed from the conversation (conversationUpdate).
	MembersRemoved []*ChannelAccount `json:"membersRemoved,omitempty"`
	// ReactionsAdded is the collection of reactions added to the conversation (messageReaction).
	ReactionsAdded []*MessageReaction `json:"reactionsAdded,omitempty"`
	// ReactionsRemoved is the collection of reactions removed from the conversation (messageReaction).
	ReactionsRemoved []*MessageReaction `json:"reactionsRemoved,omitempty"`
	// TopicName is the updated topic name of the conversation (conversationUpdate).
	TopicName string `json:"topicName,omitempty"`
	// HistoryDisclosed indicates whether the prior history of the channel is disclosed (conversationUpdate).
	HistoryDisclosed bool `json:"historyDisclosed,omitempty"`
	// Locale is a locale name for the contents of the text field.
	Locale string `json:"locale,omitempty"`
	// Text is the text content of the message.
	Text string `json:"text,omitempty"`
	// Speak is the text to speak (SSML).
	Speak string `json:"speak,omitempty"`
	// InputHint indicates whether the agent is accepting, expecting, or ignoring user input.
	InputHint string `json:"inputHint,omitempty"`
	// Summary is the text to display if the channel cannot render cards.
	Summary string `json:"summary,omitempty"`
	// SuggestedActions are the suggested actions for the activity.
	SuggestedActions *SuggestedActions `json:"suggestedActions,omitempty"`
	// Attachments are the attachments.
	Attachments []*Attachment `json:"attachments,omitempty"`
	// Entities represents the entities mentioned in the message.
	Entities []map[string]interface{} `json:"entities,omitempty"`
	// ChannelData contains channel-specific content.
	ChannelData interface{} `json:"channelData,omitempty"`
	// Action indicates whether the recipient of a contactRelationUpdate was added or removed.
	Action string `json:"action,omitempty"`
	// ReplyToID contains the ID of the message to which this message is a reply.
	ReplyToID string `json:"replyToId,omitempty"`
	// Label is a descriptive label for the activity.
	Label string `json:"label,omitempty"`
	// ValueType is the type of the activity's value object.
	ValueType string `json:"valueType,omitempty"`
	// Value is a value associated with the activity.
	Value interface{} `json:"value,omitempty"`
	// Name is the name of the operation associated with an invoke or event activity.
	Name string `json:"name,omitempty"`
	// RelatesTo is a reference to another conversation or activity.
	RelatesTo *ConversationReference `json:"relatesTo,omitempty"`
	// Code is the end-of-conversation code indicating why the conversation ended.
	Code string `json:"code,omitempty"`
	// Expiration is the time at which the activity should be considered expired.
	Expiration string `json:"expiration,omitempty"`
	// Importance is the importance level of the activity.
	Importance string `json:"importance,omitempty"`
	// DeliveryMode is a delivery hint to signal alternate delivery paths.
	DeliveryMode string `json:"deliveryMode,omitempty"`
	// ListenFor is a list of phrases and references for speech and language priming systems.
	ListenFor []string `json:"listenFor,omitempty"`
	// TextHighlights is the collection of text fragments to highlight when the activity contains a ReplyToId value.
	TextHighlights []*TextHighlight `json:"textHighlights,omitempty"`
	// SemanticAction is an optional programmatic action accompanying this request.
	SemanticAction interface{} `json:"semanticAction,omitempty"`
	// CallerID is an IRI identifying the caller of an agent.
	CallerID string `json:"callerId,omitempty"`
}

// MessageReaction represents a reaction to a message.
type MessageReaction struct {
	// Type is the message reaction type.
	Type string `json:"type,omitempty"`
}

// TextHighlight represents a highlighted text fragment.
type TextHighlight struct {
	// Text is the text to highlight.
	Text string `json:"text,omitempty"`
	// Occurrence specifies which occurrence of the text to highlight.
	Occurrence int `json:"occurrence,omitempty"`
}

// NewMessageActivity creates a new message activity with the given text.
func NewMessageActivity(text string) *Activity {
	return &Activity{
		Type: ActivityTypeMessage,
		Text: text,
	}
}

// NewEventActivity creates a new event activity with the given name and value.
func NewEventActivity(name string, value interface{}) *Activity {
	return &Activity{
		Type:  ActivityTypeEvent,
		Name:  name,
		Value: value,
	}
}

// NewTypingActivity creates a new typing activity.
func NewTypingActivity() *Activity {
	return &Activity{Type: ActivityTypeTyping}
}

// NewTraceActivity creates a trace activity with the given name, value, value type, and label.
func NewTraceActivity(name string, value interface{}, valueType, label string) *Activity {
	if valueType == "" && value != nil {
		valueType = "object"
	}
	return &Activity{
		Type:      ActivityTypeTrace,
		Name:      name,
		Value:     value,
		ValueType: valueType,
		Label:     label,
	}
}

// NewConversationUpdateActivity creates a new conversationUpdate activity.
func NewConversationUpdateActivity() *Activity {
	return &Activity{Type: ActivityTypeConversationUpdate}
}

// NewEndOfConversationActivity creates a new endOfConversation activity.
func NewEndOfConversationActivity() *Activity {
	return &Activity{Type: ActivityTypeEndOfConversation}
}

// NewInvokeActivity creates a new invoke activity.
func NewInvokeActivity() *Activity {
	return &Activity{Type: ActivityTypeInvoke}
}

// NewHandoffActivity creates a new handoff activity.
func NewHandoffActivity() *Activity {
	return &Activity{Type: ActivityTypeHandoff}
}

// NewContactRelationUpdateActivity creates a new contactRelationUpdate activity.
func NewContactRelationUpdateActivity() *Activity {
	return &Activity{Type: ActivityTypeContactRelationUpdate}
}

// IsType returns true if the activity has the given type.
func (a *Activity) IsType(activityType string) bool {
	if a.Type == "" {
		return false
	}
	t := strings.ToLower(a.Type)
	at := strings.ToLower(activityType)
	if t == at {
		return true
	}
	// also handle sub-type prefix matching (e.g., "invoke/something")
	if strings.HasPrefix(t, at+"/") {
		return true
	}
	return false
}

// AsMessageActivity returns this activity if it is a message activity, otherwise nil.
func (a *Activity) AsMessageActivity() *Activity {
	if a.IsType(ActivityTypeMessage) {
		return a
	}
	return nil
}

// AsEventActivity returns this activity if it is an event activity, otherwise nil.
func (a *Activity) AsEventActivity() *Activity {
	if a.IsType(ActivityTypeEvent) {
		return a
	}
	return nil
}

// AsInvokeActivity returns this activity if it is an invoke activity, otherwise nil.
func (a *Activity) AsInvokeActivity() *Activity {
	if a.IsType(ActivityTypeInvoke) {
		return a
	}
	return nil
}

// AsTypingActivity returns this activity if it is a typing activity, otherwise nil.
func (a *Activity) AsTypingActivity() *Activity {
	if a.IsType(ActivityTypeTyping) {
		return a
	}
	return nil
}

// AsEndOfConversationActivity returns this activity if it is an endOfConversation activity, otherwise nil.
func (a *Activity) AsEndOfConversationActivity() *Activity {
	if a.IsType(ActivityTypeEndOfConversation) {
		return a
	}
	return nil
}

// AsConversationUpdateActivity returns this activity if it is a conversationUpdate activity, otherwise nil.
func (a *Activity) AsConversationUpdateActivity() *Activity {
	if a.IsType(ActivityTypeConversationUpdate) {
		return a
	}
	return nil
}

// AsContactRelationUpdateActivity returns this activity if it is a contactRelationUpdate activity, otherwise nil.
func (a *Activity) AsContactRelationUpdateActivity() *Activity {
	if a.IsType(ActivityTypeContactRelationUpdate) {
		return a
	}
	return nil
}

// AsInstallationUpdateActivity returns this activity if it is an installationUpdate activity, otherwise nil.
func (a *Activity) AsInstallationUpdateActivity() *Activity {
	if a.IsType(ActivityTypeInstallationUpdate) {
		return a
	}
	return nil
}

// AsMessageReactionActivity returns this activity if it is a messageReaction activity, otherwise nil.
func (a *Activity) AsMessageReactionActivity() *Activity {
	if a.IsType(ActivityTypeMessageReaction) {
		return a
	}
	return nil
}

// AsMessageUpdateActivity returns this activity if it is a messageUpdate activity, otherwise nil.
func (a *Activity) AsMessageUpdateActivity() *Activity {
	if a.IsType(ActivityTypeMessageUpdate) {
		return a
	}
	return nil
}

// AsMessageDeleteActivity returns this activity if it is a messageDelete activity, otherwise nil.
func (a *Activity) AsMessageDeleteActivity() *Activity {
	if a.IsType(ActivityTypeMessageDelete) {
		return a
	}
	return nil
}

// AsTraceActivity returns this activity if it is a trace activity, otherwise nil.
func (a *Activity) AsTraceActivity() *Activity {
	if a.IsType(ActivityTypeTrace) {
		return a
	}
	return nil
}

// AsHandoffActivity returns this activity if it is a handoff activity, otherwise nil.
func (a *Activity) AsHandoffActivity() *Activity {
	if a.IsType(ActivityTypeHandoff) {
		return a
	}
	return nil
}

// AsSuggestionActivity returns this activity if it is a suggestion activity, otherwise nil.
func (a *Activity) AsSuggestionActivity() *Activity {
	if a.IsType(ActivityTypeSuggestion) {
		return a
	}
	return nil
}

// GetConversationReference creates a ConversationReference based on this activity.
func (a *Activity) GetConversationReference() *ConversationReference {
	ref := &ConversationReference{
		ChannelID:    a.ChannelID,
		ServiceURL:   a.ServiceURL,
		Conversation: a.Conversation,
		Locale:       a.Locale,
	}
	if a.From != nil {
		fromCopy := *a.From
		ref.User = &fromCopy
	}
	if a.Recipient != nil {
		recipientCopy := *a.Recipient
		ref.Agent = &recipientCopy
	}
	// Don't set ActivityID for conversationUpdate on directline/webchat
	if a.Type != ActivityTypeConversationUpdate ||
		(a.ChannelID != "directline" && a.ChannelID != "webchat") {
		ref.ActivityID = a.ID
	}
	return ref
}

// GetReplyConversationReference creates a ConversationReference from this activity's conversation info
// and the ResourceResponse returned from sending an activity.
func (a *Activity) GetReplyConversationReference(reply *ResourceResponse) *ConversationReference {
	ref := a.GetConversationReference()
	if reply != nil {
		ref.ActivityID = reply.ID
	}
	return ref
}

// ApplyConversationReference updates this activity with delivery information from a ConversationReference.
// If isIncoming is true, the agent is treated as the recipient; otherwise the agent is the sender.
func (a *Activity) ApplyConversationReference(reference *ConversationReference, isIncoming bool) *Activity {
	a.ChannelID = reference.ChannelID
	a.ServiceURL = reference.ServiceURL
	a.Conversation = reference.Conversation
	if reference.Locale != "" {
		a.Locale = reference.Locale
	}
	if isIncoming {
		a.From = reference.User
		a.Recipient = reference.Agent
		if reference.ActivityID != "" {
			a.ID = reference.ActivityID
		}
	} else {
		a.From = reference.Agent
		a.Recipient = reference.User
		if reference.ActivityID != "" {
			a.ReplyToID = reference.ActivityID
		}
	}
	return a
}

// CreateReply creates a new message activity as a response to this activity.
func (a *Activity) CreateReply(text string) *Activity {
	reply := &Activity{
		Type:       ActivityTypeMessage,
		ServiceURL: a.ServiceURL,
		ChannelID:  a.ChannelID,
		Text:       text,
		Locale:     a.Locale,
	}
	if a.Recipient != nil {
		from := &ChannelAccount{ID: a.Recipient.ID, Name: a.Recipient.Name}
		reply.From = from
	}
	if a.From != nil {
		recipient := &ChannelAccount{ID: a.From.ID, Name: a.From.Name}
		reply.Recipient = recipient
	}
	if a.Conversation != nil {
		conv := *a.Conversation
		reply.Conversation = &conv
	}
	if a.Type != ActivityTypeConversationUpdate ||
		(a.ChannelID != "directline" && a.ChannelID != "webchat") {
		reply.ReplyToID = a.ID
	}
	return reply
}

// HasContent returns true if the activity has any content to send.
func (a *Activity) HasContent() bool {
	if strings.TrimSpace(a.Text) != "" {
		return true
	}
	if strings.TrimSpace(a.Summary) != "" {
		return true
	}
	if len(a.Attachments) > 0 {
		return true
	}
	if a.ChannelData != nil {
		return true
	}
	return false
}

// IsFromStreamingConnection returns true if the activity was sent via a streaming connection.
func (a *Activity) IsFromStreamingConnection() bool {
	if a.ServiceURL == "" {
		return false
	}
	return !strings.HasPrefix(strings.ToLower(a.ServiceURL), "http")
}

// IsAgenticRequest returns true if this activity is from an agentic identity.
func (a *Activity) IsAgenticRequest() bool {
	if a.Recipient == nil {
		return false
	}
	role := a.Recipient.Role
	return role == RoleTypeAgenticIdentity || role == RoleTypeAgenticUser
}

// GetAgenticInstanceID returns the agent instance ID if this is an agentic request.
func (a *Activity) GetAgenticInstanceID() string {
	if !a.IsAgenticRequest() || a.Recipient == nil {
		return ""
	}
	return a.Recipient.AgenticAppID
}

// GetAgenticUser returns the agentic user ID if this is an agentic request.
func (a *Activity) GetAgenticUser() string {
	if !a.IsAgenticRequest() || a.Recipient == nil {
		return ""
	}
	return a.Recipient.AgenticUserID
}

// GetAgenticTenantID returns the agentic tenant ID if this is an agentic request.
func (a *Activity) GetAgenticTenantID() string {
	if !a.IsAgenticRequest() {
		return ""
	}
	if a.Recipient != nil && a.Recipient.TenantID != "" {
		return a.Recipient.TenantID
	}
	if a.Conversation != nil && a.Conversation.TenantID != "" {
		return a.Conversation.TenantID
	}
	return ""
}
