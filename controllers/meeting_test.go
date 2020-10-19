package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var dummyTime, _ = time.Parse(time.UTC.String(), "2011-08-12T20:17:46.384Z")

func meetingServer(w http.ResponseWriter, r *http.Request) {
	m := MeetingResponse{
		Status:  200,
		Message: "List of all Meetings",
		Data:    "123",
		Time:    dummyTime,
	}
	res, _ := json.Marshal(m)
	w.Write(res)
}

func TestGetSingleMeeting(t *testing.T) {
	t.Run("returns all meetings list", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/meetings/123456789", nil)
		response := httptest.NewRecorder()
		got := MeetingResponse{}
		meetingServer(response, request)

		_ = json.NewDecoder(response.Body).Decode(&got)
		want := MeetingResponse{
			Status:  200,
			Message: "List of all Meetings",
			Data:    "123",
			Time:    dummyTime,
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
