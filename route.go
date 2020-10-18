package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var routes = []route{
	newRoute("GET", "/test", getTest),
	newRoute("GET", "/meetings", handleMeeting),
	newRoute("GET", "/meetings/([^/]+)", getMeetingWithID),
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
	log.Println(fields)
	return fields[index]
}

func createMeeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "createMeeting\n")
}

func handleMeeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "handleMeeting\n")
	paramParticipant := r.URL.Query().Get("participant")
	paramStart := r.URL.Query().Get("start")
	paramEnd := r.URL.Query().Get("end")

	if paramParticipant != "" {
		listMeetingsGivenParticipant(w, r)
	} else if paramStart != "" && paramEnd != "" {
		listMeetingsGivenTime(w, r)
	}
}

func listMeetingsGivenTime(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "listMeetingsWithTime %s\n", r.URL.Query())
}

func listMeetingsGivenParticipant(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "listMeetingsWithTime %s\n", r.URL.Query())
}

func getMeetingWithID(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "getMeetingWithId %s\n", slug)
}

func getTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API test\n")
}
