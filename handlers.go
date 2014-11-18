package main

import (
	"pong_matcher_go/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type matchRequestPersister func(domain.MatchRequest) error
type MatchRequestRetriever func(string) (bool, domain.MatchRequest, error)
type matchRetriever func(string) (bool, domain.Match)
type resultPersister func(domain.Result) error
type wiper func() error

func AllHandler(wipe wiper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := wipe(); err != nil {
			w.WriteHeader(500)
		}
	}
}

func CreateMatchRequestHandler(persist matchRequestPersister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var matchRequest domain.MatchRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&matchRequest); err == nil {
			matchRequest.Uuid = mux.Vars(r)["uuid"]
			if err = persist(matchRequest); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
			} else {
				w.WriteHeader(500)
			}
		} else {
			w.WriteHeader(400)
		}
	}
}

func GetMatchRequestHandler(retrieve MatchRequestRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		found, matchRequest, err := retrieve(mux.Vars(r)["uuid"])

		if err != nil {
			w.WriteHeader(500)
			return
		}

		if found {
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
		var result domain.Result
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&result); err != nil {
			w.WriteHeader(400)
			return
		}
		if err := persist(result); err != nil {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(201)
		}
	}
}
