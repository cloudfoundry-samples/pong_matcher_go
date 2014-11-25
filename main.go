package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.com/camelpunch/pong_matcher_go/io"
)

func main() {
	io.InitDb(io.MigratedDbMap())
	defer io.CloseDb()

	router := mux.NewRouter()

	router.
		HandleFunc("/all", AllHandler(io.DeleteAll)).
		Methods("DELETE")
	router.
		HandleFunc("/match_requests/{uuid}", CreateMatchRequestHandler(io.PersistMatchRequest)).
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
