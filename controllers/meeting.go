package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
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
	ID           string    `json:"id,omitempty"`
	Title        string    `json:"title"`
	Participants []string  `json:"participants"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
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

	_, err := mCollection.InsertOne(context.TODO(), newMeeting)

	newAccountsCreated := _addNewParticipants(newMeeting.Participants, newMeeting)
	log.Println(newAccountsCreated)

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
		"Meeting creating successfully." +
			"New Accounts created for: " +
			strings.Join(newAccountsCreated, ", "),
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
	meeting := Meeting{}
	err := mCollection.FindOne(context.TODO(), bson.M{
		"participants": bson.M{
			"$all": []string{email},
		}}).Decode(&meeting)
	if err != nil {
		log.Println(err)
		message := MeetingResponse{
			http.StatusInternalServerError,
			"Something went wrong. Meeting/s not found for participant. Kindly Try Again",
			[]Meeting{meeting},
			time.Now().UTC(),
		}
		return message
	}

	message := MeetingResponse{
		http.StatusOK,
		"Requested Meeting/s for Participant Found",
		[]Meeting{meeting},
		time.Now().UTC(),
	}
	return message

}

// func (doc *Meeting) participantPool() chan string {
// 	notAccEmails := make(chan string, len(doc.Participants))

// 	for _, email := range doc.Participants {
// 		go getSingleParticipant(email, notAccEmails)
// 	}
// 	defer close(notAccEmails)

// 	return notAccEmails
// }

func _addNewParticipants(emails []string, doc Meeting) []string {
	var newAccounts []string
	notAccEmails := make(chan string, len(emails))

	for _, email := range emails {
		wg.Add(1)
		go getSingleParticipant(email, notAccEmails)
	}
	wg.Wait()
	close(notAccEmails)

	for email := range notAccEmails {
		newAccounts = append(newAccounts, email)
		createSingleParticipant(email)
	}

	return newAccounts
}
