package controllers

import "time"

// MeetingResponse represents JSON Response structure for Meeting
type MeetingResponse struct {
	Status  int
	Message string
	Data    interface{}
	Time    time.Time
}

// ParticipantResponse represents JSON Response structure for Participant
type ParticipantResponse struct {
	Status  int
	Message string
	Data    []Participant
	Time    time.Time
}