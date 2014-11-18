package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"pong_matcher_go/io"
	"log"
	"net/http"
	"os"
)

func main() {
	io.InitDb()
	defer io.CloseDb()

	router := mux.NewRouter()

	router.
		HandleFunc("/all", AllHandler(io.DeleteAll)).
		Methods("DELETE")
	router.
		HandleFunc(
		"/match_requests/{uuid}",
		CreateMatchRequestHandler(io.PersistMatchRequest),
	).
		Methods("PUT")
	router.
		HandleFunc("/match_requests/{uuid}", GetMatchRequestHandler(io.GetMatchRequest)).
		Methods("GET")
	router.
		HandleFunc("/matches/{uuid}", MatchHandler(io.GetMatch)).
		Methods("GET")
	router.
		HandleFunc("/results", ResultsHandler(io.PersistResult)).
		Methods("POST")

	if err := http.ListenAndServe(fmt.Sprintf(":%v", getPort()), router); err != nil {
		log.Fatalln(err)
	}
}

func getPort() string {
	if configuredPort := os.Getenv("PORT"); configuredPort == "" {
		return "3000"
	} else {
		return configuredPort
	}
}
