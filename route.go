package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"skeduler/controllers"
	"strings"
	"time"
)

var routes = []route{
	newRoute("GET", "/test", getTest),
	newRoute("GET", "/meeting/([^/]+)", getMeetingWithID),
	newRoute("GET", "/meetings", getMeetingHandler),
	newRoute("POST", "/meetings", createMeeting),
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

var malformedRequest, _ = json.Marshal(controllers.MalformedRequest{
	Status:  http.StatusBadRequest,
	Message: "Malformed Request, check your query again",
	Time:    time.Now().UTC(),
})

// Serve serves the endpoints for the allowed routes
func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			w.Header().Set("Content-Type", "application/json")
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func getTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API test\n")
}

func createMeeting(w http.ResponseWriter, r *http.Request) {
	var body controllers.Meeting
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Can't read body. Try again!", http.StatusBadRequest)
		return
	}
	// TODO: Use for unit tests later
	// fmt.Fprint(w, body.Participants[0], "\n")
	// fmt.Fprint(w, body.StartTime, "\n")
	// fmt.Fprint(w, body.EndTime, "\n")
	// fmt.Fprint(w, body.Title, "\n")

	message := controllers.CreateMeeting(body)
	response, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	w.Write(response)

}

func getMeetingHandler(w http.ResponseWriter, r *http.Request) {
	paramParticipantEmail := r.URL.Query().Get("participant")
	paramStart := r.URL.Query().Get("start")
	paramEnd := r.URL.Query().Get("end")

	if paramParticipantEmail != "" {
		meetingsGivenParticipant(w, paramParticipantEmail)
		return
	}
	if paramStart != "" && paramEnd != "" {
		meetingsGivenTime(w, paramStart, paramEnd)
		return
	}

	w.Write(malformedRequest)
}

func getMeetingWithID(w http.ResponseWriter, r *http.Request) {
	meetingID := getField(r, 0)
	if meetingID == "" {
		w.Write(malformedRequest)
		return
	}

	message := controllers.GetSingleMeeting(meetingID)
	response, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	w.Write(response)
}

func meetingsGivenTime(w http.ResponseWriter, start string, end string) {
	startTime, err1 := time.Parse(time.RFC3339, start)
	endTime, err2 := time.Parse(time.RFC3339, end)

	if err1 != nil || err2 != nil {
		w.Write(malformedRequest)
		return
	}

	message := controllers.GetAllMeetings(startTime, endTime)
	response, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	w.Write(response)
}

func meetingsGivenParticipant(w http.ResponseWriter, email string) {
	message := controllers.GetMeetingForParticipant(email)
	response, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	w.Write(response)
}
