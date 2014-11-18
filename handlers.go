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
type matchRequestRetriever func(string) (bool, MatchRequest)
type matchRetriever func(string) (bool, Match)
type resultPersister func(Result)

func CreateMatchRequestHandler(persist matchRequestPersister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var matchRequest MatchRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&matchRequest); err == nil {
			matchRequest.Uuid = mux.Vars(r)["uuid"]
			persist(matchRequest)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
		}
	}
}

func GetMatchRequestHandler(retrieve matchRequestRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if found, matchRequest := retrieve(mux.Vars(r)["uuid"]); found {
			js, err := json.Marshal(matchRequest)
			checkErr(err, "Error writing JSON")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(js)
		} else {
			w.WriteHeader(404)
		}
	}
}

func MatchHandler(retrieve matchRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if found, match := retrieve(mux.Vars(r)["uuid"]); found {
			js, err := json.Marshal(match)
			checkErr(err, "Error writing JSON")
			w.WriteHeader(200)
			w.Write(js)
		}
	}
}

func ResultsHandler(persist resultPersister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var result Result
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&result); err == nil {
			persist(result)
			w.WriteHeader(201)
		} else {
			w.WriteHeader(400)
		}
	}
}
