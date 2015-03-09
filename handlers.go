package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/cloudfoundry-samples/pong_matcher_go/domain"
)

type MatchRequestPersister func(domain.MatchRequest) error
type MatchRequestRetriever func(string) (bool, domain.MatchRequest, error)
type MatchRetriever func(string) (bool, domain.Match)
type ResultPersister func(domain.Result) error
type Wiper func() error

func AllHandler(wipe Wiper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := wipe(); err != nil {
			w.WriteHeader(500)
		}
	}
}

func CreateMatchRequestHandler(persist MatchRequestPersister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var matchRequest domain.MatchRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&matchRequest); err != nil {
			w.WriteHeader(400)
			return
		}

		matchRequest.Uuid = mux.Vars(r)["uuid"]

		if err := persist(matchRequest); err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
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
			if err != nil {
				log.Fatalln("Error writing JSON:", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(js)
		} else {
			w.WriteHeader(404)
		}
	}
}

func MatchHandler(retrieve MatchRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if found, match := retrieve(mux.Vars(r)["uuid"]); found {
			js, err := json.Marshal(match)
			if err != nil {
				log.Fatalln("Error writing JSON:", err)
			}
			w.WriteHeader(200)
			w.Write(js)
		} else {
			w.WriteHeader(404)
		}
	}
}

func ResultsHandler(persist ResultPersister) http.HandlerFunc {
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
