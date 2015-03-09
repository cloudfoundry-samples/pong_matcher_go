package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-samples/pong_matcher_go/io"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "mysql2://gopong:gopong@127.0.0.1:3306/pong_matcher_go_development?reconnect=true"
	}

	io.InitDb(io.MigratedDbMap(databaseUrl, "db/migrations"))
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
