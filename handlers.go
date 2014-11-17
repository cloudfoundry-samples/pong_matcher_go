package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func AllHandler(w http.ResponseWriter, r *http.Request) {
	deleteAll()
}

type matchRequestPersister func(MatchRequest)

func CreateMatchRequestHandler(persist matchRequestPersister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var matchRequest MatchRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&matchRequest)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		matchRequest.Uuid = mux.Vars(r)["uuid"]

		persist(matchRequest)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}
}

func GetMatchRequestHandler(w http.ResponseWriter, r *http.Request) {
	if found, matchRequest := getMatchRequest(mux.Vars(r)["uuid"]); found {
		js, err := json.Marshal(matchRequest)
		checkErr(err, "Error writing JSON")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	} else {
		w.WriteHeader(404)
	}
}

func MatchHandler(w http.ResponseWriter, r *http.Request) {
	if found, match := getMatch(mux.Vars(r)["uuid"]); found {
		js, err := json.Marshal(match)
		checkErr(err, "Error writing JSON")

		w.WriteHeader(200)
		w.Write(js)
	}
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	var result Result
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&result)
	checkErr(err, "Decoding JSON failed")

	persistResult(result)

	w.WriteHeader(201)
}
