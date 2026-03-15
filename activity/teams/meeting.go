// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package teams

// UserMeetingDetails contains specific details of a user in a Teams meeting.
type UserMeetingDetails struct {
	// Role is the role of the participant in the current meeting.
	Role string `json:"role,omitempty"`
	// InMeeting indicates whether the participant is in the meeting.
	InMeeting *bool `json:"inMeeting,omitempty"`
}

// TeamsMeetingMember contains data about a meeting participant.
type TeamsMeetingMember struct {
	// User is the channel user data.
	User *TeamsChannelAccount `json:"user,omitempty"`
	// Meeting is the user meeting details.
	Meeting *UserMeetingDetails `json:"meeting,omitempty"`
}

// MeetingStartEventDetails contains specific details of a Teams meeting start event.
type MeetingStartEventDetails struct {
	// StartTime is the timestamp for meeting start, in UTC.
	StartTime string `json:"startTime,omitempty"`
}

// MeetingEndEventDetails contains specific details of a Teams meeting end event.
type MeetingEndEventDetails struct {
	// EndTime is the timestamp for meeting end, in UTC.
	EndTime string `json:"endTime,omitempty"`
}

// MeetingParticipantsEventDetails contains data about the meeting participants.
type MeetingParticipantsEventDetails struct {
	// Members contains the members involved in the meeting event.
	Members []*TeamsMeetingMember `json:"members,omitempty"`
}
