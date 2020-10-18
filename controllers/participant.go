package controllers

import (
	"context"
	"log"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ParticipantResponse represents JSON Response structure for Participant
type ParticipantResponse struct {
	Status  int
	Message string
	Data    []Participant
	Time    time.Time
}

// RSVPMessageType represents the allowed RSVP messages for a meeting
type RSVPMessageType int32

const (
	yes         RSVPMessageType = 0
	no          RSVPMessageType = 1
	maybe       RSVPMessageType = 2
	notAnswered RSVPMessageType = 3
)

// Participant represents users
type Participant struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
	RSVP  string `json:"rsvp,omitempty"`
}

// DATABASE INSTANCE
var pCollection *mongo.Collection

func getSingleParticipant(email string, noAccMails chan string) {
	defer wg.Done()

	participant := Participant{}
	err := pCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&participant)
	if err != nil {
		log.Println("[x] Participant not found during meeting creation! Adding to DB...")
		noAccMails <- email
	}
}

func createSingleParticipant(email string) {
	nameRegex := regexp.MustCompile("[^@]+")
	newParticipant := Participant{
		Name:  nameRegex.FindString(email),
		Email: email,
		RSVP:  "Not Answered",
	}

	_, err := pCollection.InsertOne(context.TODO(), newParticipant)

	if err != nil {
		log.Printf("Error while inserting new participant into db, Reason: %v\n", err)
	}
}
