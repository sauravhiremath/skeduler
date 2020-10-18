package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MsgResponse represents JSON Response structure
type MsgResponse struct {
	Status  int
	Message string
	Data    []Meeting
	Time    time.Time
}

// Meeting represents Meeting collection structure
type Meeting struct {
	ID           string    `json:"id,omitempty"`
	Title        string    `json:"title"`
	Participants []string  `json:"participants"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// DATABASE INSTANCE
var collection *mongo.Collection

// MeetingCollection initialises new collection
func MeetingCollection(c *mongo.Database) {
	collection = c.Collection("meetings")
}

// GetAllMeetings returns all meetings scheduled between given start and end time
func GetAllMeetings(startTime time.Time, endTime time.Time) MsgResponse {
	meetings := []Meeting{}
	cursor, err := collection.Find(context.TODO(), bson.M{
		"start_time": bson.M{
			"$gt": startTime,
			"$lt": endTime,
		}})

	if err != nil {
		log.Printf("[x] Error while getting all meetings, Reason: %v\n", err)
		message := MsgResponse{
			http.StatusInternalServerError,
			"Something went wrong",
			[]Meeting{},
			time.Now().UTC(),
		}
		return message
	}

	// Iterate through the returned cursor.
	for cursor.Next(context.TODO()) {
		var meeting Meeting
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}

	message := MsgResponse{
		http.StatusOK,
		"List of all Meetings",
		meetings,
		time.Now().UTC(),
	}
	return message
}

// CreateMeeting creates a new meeting entry
func CreateMeeting(meeting Meeting) MsgResponse {
	id := UUID()
	title := meeting.Title
	participants := meeting.Participants
	startTime := meeting.StartTime
	endTime := meeting.EndTime

	newMeeting := Meeting{
		ID:           id,
		Title:        title,
		Participants: participants,
		StartTime:    startTime,
		EndTime:      endTime,
		CreatedAt:    time.Now().UTC(),
	}

	_, err := collection.InsertOne(context.TODO(), newMeeting)

	if err != nil {
		log.Printf("Error while inserting new meeting into db, Reason: %v\n", err)
		message := MsgResponse{
			http.StatusInternalServerError,
			"Something went wrong. Meeting not created. Kindly Try Again",
			[]Meeting{meeting},
			time.Now().UTC(),
		}
		return message
	}

	message := MsgResponse{
		http.StatusOK,
		"Meeting creating successfully",
		[]Meeting{},
		time.Now().UTC(),
	}
	return message
}

// GetSingleMeeting returns full-meeting document with given meetingID
func GetSingleMeeting(meetingID string) MsgResponse {
	meeting := Meeting{}
	err := collection.FindOne(context.TODO(), bson.M{"id": meetingID}).Decode(&meeting)
	if err != nil {
		message := MsgResponse{
			http.StatusInternalServerError,
			"Something went wrong. Meeting not found. Kindly Try Again",
			[]Meeting{meeting},
			time.Now().UTC(),
		}
		return message
	}

	message := MsgResponse{
		http.StatusOK,
		"Requested Meeting Found",
		[]Meeting{meeting},
		time.Now().UTC(),
	}
	return message
}

// GetMeetingForParticipant returns all meetings the given participant is inside
func GetMeetingForParticipant(email string) MsgResponse {
	meeting := Meeting{}
	err := collection.FindOne(context.TODO(), bson.M{
		"participants": bson.M{
			"$all": []string{email},
		}}).Decode(&meeting)
	if err != nil {
		log.Println(err)
		message := MsgResponse{
			http.StatusInternalServerError,
			"Something went wrong. Meeting/s not found for participant. Kindly Try Again",
			[]Meeting{meeting},
			time.Now().UTC(),
		}
		return message
	}

	message := MsgResponse{
		http.StatusOK,
		"Requested Meeting/s for Participant Found",
		[]Meeting{meeting},
		time.Now().UTC(),
	}
	return message

}
