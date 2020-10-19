package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sauravhiremath/skeduler/config"
)

const port = 8080

func handleRequests() {
	router := http.HandlerFunc(Serve)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func main() {
	config.Connect()

	handleRequests()
}
