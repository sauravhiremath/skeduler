package controllers

import (
	"testing"
	"time"
)

// TODO: Add mock mongoDB calls, currently used DB integration
type MockDB struct {
	meetings     []Meeting
	participants []Participant
}

var mockDB = new(MockDB)

func TestCreateMeeting(t *testing.T) {
	startTime, _ := time.Parse(time.UTC.String(), "2019-10-18T10:44:56Z")
	endTime, _ := time.Parse(time.UTC.String(), "2019-10-18T12:51:56Z")
	m1 := Meeting{
		Title: "Sample Meeting 1",
		Participants: []Participant{
			{
				Name:  "abc",
				Email: "abc@def.com",
				RSVP:  "Yes",
			},
		},
		StartTime: startTime,
		EndTime:   endTime,
	}
	m2 := Meeting{
		Title: "Sample Meeting 2",
		Participants: []Participant{
			{
				Name:  "xyz",
				Email: "xyz@def.com",
				RSVP:  "Yes",
			},
		},
		StartTime: startTime,
		EndTime:   endTime,
	}

	got1 := CreateMeeting(m1)
	got2 := CreateMeeting(m2)

	if got1.Status != 200 || got2.Status != 200 {
		t.Errorf("Returned wrong Status Code. Got %q, expected %q", got1.Status, 200)
	}

	mockDB.meetings = append(mockDB.meetings, got1.Data[0])
	mockDB.meetings = append(mockDB.meetings, got2.Data[0])
}

func TestGetAllMeetings(t *testing.T) {
	startTime, _ := time.Parse(time.UTC.String(), "2019-10-18T10:24:56Z")
	endTime, _ := time.Parse(time.UTC.String(), "2019-10-18T12:52:56Z")

	got := GetAllMeetings(startTime, endTime)

	if got.Status != 200 {
		t.Errorf("Returned wrong Status Code. Got %q, expected %q", got.Status, 200)
	}
}

func TestGetMeetingForParticipant(t *testing.T) {
	email := "abc@def.com"

	got := GetMeetingForParticipant(email)

	if got.Status != 200 {
		t.Errorf("Returned wrong Status Code. Got %q, expected %q", got.Status, 200)
	}
}
