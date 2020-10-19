package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GetSingleMeeting returns full-meeting document with given meetingID
func GetSingleMeeting(meetingID string) MeetingResponse {
	meeting := Meeting{}
	err := mCollection.FindOne(context.TODO(), bson.M{"id": meetingID}).Decode(&meeting)
	if err != nil {
		log.Printf("[x] Error: Requested MeetingID - %v not found!\n", meetingID)
		message := MeetingResponse{
			http.StatusNotFound,
			"Meeting not found. Kindly Try Again",
			[]Meeting{},
			time.Now().UTC(),
		}
		return message
	}

	message := MeetingResponse{
		http.StatusOK,
		"Requested Meeting Found",
		[]Meeting{meeting},
		time.Now().UTC(),
	}
	return message
}
