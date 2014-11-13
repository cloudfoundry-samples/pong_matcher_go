package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func AllHandler(w http.ResponseWriter, r *http.Request) {
	deleteAll()
}

func CreateMatchRequestHandler(w http.ResponseWriter, r *http.Request) {
	dbmap := initDb()
	defer dbmap.Db.Close()

	uuid := mux.Vars(r)["uuid"]

	var matchRequest MatchRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&matchRequest)
	checkErr(err, "Decoding JSON failed")
	matchRequest.Uuid = uuid

	err = dbmap.Insert(&matchRequest)
	checkErr(err, "Creation of MatchRequest failed")

	openMatchRequests := suitableOpponentMatchRequests(dbmap, matchRequest.RequesterId)
	if len(openMatchRequests) > 0 {
		recordMatch(dbmap, openMatchRequests[0], matchRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
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
