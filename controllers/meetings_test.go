package controllers

import (
	"testing"
	"time"
)

type MockDB struct {
	meetings     []Meeting
	participants []Participant
}

var mockDB = new(MockDB)

func TestCreateMeeting(t *testing.T) {
	startTime, _ := time.Parse(time.UTC.String(), "2019-10-18T10:44:56Z")
	endTime, _ := time.Parse(time.UTC.String(), "2019-10-18T12:51:56Z")
	m := Meeting{
		Title: "Sample Meeting",
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

	got := CreateMeeting(m)

	if got.Status != 200 {
		t.Errorf("Got %q, expected %q", got.Status, 200)
	}

	mockDB.meetings = append(mockDB.meetings, got.Data[0])
}
