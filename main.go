package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.
		HandleFunc("/all", AllHandler).
		Methods("DELETE")
	router.
		HandleFunc("/match_requests/{uuid}", CreateMatchRequestHandler).
		Methods("PUT")
	router.
		HandleFunc("/match_requests/{uuid}", GetMatchRequestHandler).
		Methods("GET")
	router.
		HandleFunc("/matches/{uuid}", MatchHandler).
		Methods("GET")
	router.
		HandleFunc("/results", ResultsHandler).
		Methods("POST")

	err := http.ListenAndServe(fmt.Sprintf(":%v", getPort()), router)
	checkErr(err, "Error starting server")
}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		return "3000"
	} else {
		return configuredPort
	}
}
