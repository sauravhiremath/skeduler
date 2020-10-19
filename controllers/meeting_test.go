package controllers

import "testing"

func TestGetSingleMeeting(t *testing.T) {
	mID = mockDB.meetings[0].ID

	got := GetSingleMeeting(mID)

	if got.Status != 200 {
		t.Errorf("Returned wrong Status Code. Got %q, expected %q", got1.Status, 200)
	}
}
