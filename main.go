package main

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var dbmap *gorp.DbMap

func main() {
	dbmap = initDb()
	defer dbmap.Db.Close()

	router := mux.NewRouter()

	router.
		HandleFunc("/all", AllHandler).
		Methods("DELETE")
	router.
		HandleFunc(
			"/match_requests/{uuid}",
			CreateMatchRequestHandler(persistMatchRequest),
		).
		Methods("PUT")
	router.
		HandleFunc("/match_requests/{uuid}", GetMatchRequestHandler(getMatchRequest)).
		Methods("GET")
	router.
		HandleFunc("/matches/{uuid}", MatchHandler(getMatch)).
		Methods("GET")
	router.
		HandleFunc("/results", ResultsHandler(persistResult)).
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
