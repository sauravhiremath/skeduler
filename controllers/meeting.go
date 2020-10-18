package controllers

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MeetingResponse represents JSON Response structure for Meeting
type MeetingResponse struct {
	Status  int
	Message string
	Data    []Meeting
	Time    time.Time
}

// Meeting represents Meeting collection structure
type Meeting struct {
	ID           string        `json:"id,omitempty"`
	Title        string        `json:"title"`
	Participants []Participant `json:"participants"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	CreatedAt    time.Time     `json:"created_at,omitempty"`
}

// DATABASE INSTANCE
var mCollection *mongo.Collection

var wg sync.WaitGroup

// Collection initialises new collections inside given database
func Collection(c *mongo.Database) {
	mCollection = c.Collection("meetings")
	pCollection = c.Collection("participants")
}

// GetAllMeetings returns all meetings scheduled between given start and end time
func GetAllMeetings(startTime time.Time, endTime time.Time) MeetingResponse {
	meetings := []Meeting{}
	cursor, err := mCollection.Find(context.TODO(), bson.M{
		"start_time": bson.M{
			"$gt": startTime,
			"$lt": endTime,
		}})

	if err != nil {
		log.Printf("[x] Error while getting all meetings, Reason: %v\n", err)
		message := MeetingResponse{
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

	message := MeetingResponse{
		http.StatusOK,
		"List of all Meetings",
		meetings,
		time.Now().UTC(),
	}
	return message
}

// CreateMeeting creates a new meeting entry
func CreateMeeting(meeting Meeting) MeetingResponse {
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

	busyAccounts := checkTimeOverlap(participants, startTime, endTime)

	if len(busyAccounts) > 0 {
		log.Printf("[x] Participants have clashing schedules %v", busyAccounts)
		message := MeetingResponse{
			http.StatusInternalServerError,
			"Some participants have clashing schedules. Kindly Try Again!",
			[]Meeting{meeting},
			time.Now().UTC(),
		}
		return message
	}

	_, err := mCollection.InsertOne(context.TODO(), newMeeting)

	if err != nil {
		log.Printf("Error while inserting new meeting into db, Reason: %v\n", err)
		message := MeetingResponse{
			http.StatusInternalServerError,
			"Something went wrong. Meeting not created. Kindly Try Again",
			[]Meeting{meeting},
			time.Now().UTC(),
		}
		return message
	}

	message := MeetingResponse{
		http.StatusOK,
		"Meeting creating successfully.",
		[]Meeting{},
		time.Now().UTC(),
	}
	return message
}

// GetSingleMeeting returns full-meeting document with given meetingID
func GetSingleMeeting(meetingID string) MeetingResponse {
	meeting := Meeting{}
	err := mCollection.FindOne(context.TODO(), bson.M{"id": meetingID}).Decode(&meeting)
	if err != nil {
		message := MeetingResponse{
			http.StatusInternalServerError,
			"Something went wrong. Meeting not found. Kindly Try Again",
			[]Meeting{meeting},
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

// GetMeetingForParticipant returns all meetings the given participant is inside
func GetMeetingForParticipant(email string) MeetingResponse {
	meetings := []Meeting{}
	cursor, err := mCollection.Find(context.TODO(), bson.M{
		"participants": bson.M{
			"email": bson.M{
				"$all": []string{email},
			},
		}})

	if err != nil {
		log.Println(err)
		message := MeetingResponse{
			http.StatusInternalServerError,
			"Something went wrong. Meeting/s not found for participant. Kindly Try Again",
			meetings,
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

	message := MeetingResponse{
		http.StatusOK,
		"Requested Meeting/s for Participant Found",
		meetings,
		time.Now().UTC(),
	}
	return message

}

func checkTimeOverlap(participants []Participant, startTime time.Time, endTime time.Time) map[string]string {
	var emails []string
	for _, p := range participants {
		emails = append(emails, p.Email)
	}

	overlappingMeetings := make(map[string]string, len(emails))

	for _, email := range emails {
		cursor, err := mCollection.Find(context.TODO(), bson.M{
			"participants.email": bson.M{
				"$all": []string{email},
			},
			"participants.rsvp": "Yes",
		})

		if err != nil {
			log.Println(err)
		}

		// Iterate through the returned cursor.
		for cursor.Next(context.TODO()) {
			var meeting Meeting
			cursor.Decode(&meeting)
			overlap := meeting.StartTime.Before(endTime) && meeting.EndTime.After(startTime)
			if overlap {
				overlappingMeetings[email] = meeting.ID
			}
		}
	}

	return overlappingMeetings
}
