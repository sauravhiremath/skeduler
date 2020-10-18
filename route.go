package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var routes = []route{
	newRoute("POST", "/meetings", createMeeting),
	newRoute("GET", "/meetings/([^/]+)", getMeetingWithID),
	newRoute("GET", "/meetings?start=([^/]+)&end=([^/]+)", listMeetingsGivenTime),
	newRoute("GET", "/meetings?participant=([^/]+)", listMeetingsGivenParticipant),
	newRoute("GET", "/test", getTest),
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
	return fields[index]
}

func createMeeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "createMeeting\n")
}

func listMeetingsGivenTime(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "listMeetingsWithTime %s\n", slug)
}

func listMeetingsGivenParticipant(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "listMeetingsWithTime %s\n", slug)
}

func getMeetingWithID(w http.ResponseWriter, r *http.Request) {
	slug := getField(r, 0)
	fmt.Fprintf(w, "getMeetingWithId %s\n", slug)
}

func getTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API test\n")
}
