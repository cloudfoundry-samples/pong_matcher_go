package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/guregu/null.v2"
	"net/http"
)

func AllHandler(w http.ResponseWriter, r *http.Request) {
	dbmap := initDb()
	defer dbmap.Db.Close()

	err := dbmap.TruncateTables()
	checkErr(err, "Truncation failed")
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
	dbmap := initDb()
	defer dbmap.Db.Close()

	uuid := mux.Vars(r)["uuid"]

	matchRequest := MatchRequest{}
	err := dbmap.SelectOne(
		&matchRequest,
		"SELECT * FROM match_requests WHERE uuid = ?", uuid,
	)
	if err == nil {
		matchId, err := dbmap.SelectStr(
			`SELECT match_id
					FROM participants
					WHERE match_request_uuid = ?
					AND match_id NOT IN (
						SELECT match_id FROM results
					)`,
			uuid,
		)
		if err == nil && matchId != "" {
			matchRequest.MatchId = null.StringFrom(matchId)
		}

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
	dbmap := initDb()
	defer dbmap.Db.Close()

	matchId := mux.Vars(r)["uuid"]

	var participants []Participant
	_, err := dbmap.Select(
		&participants,
		`SELECT * 
			FROM participants
			WHERE match_id = ?`,
		matchId,
	)
	checkErr(err, "Error getting participants")

	match := Match{
		Id:              matchId,
		MatchRequest1Id: participants[0].MatchRequestUuid,
		MatchRequest2Id: participants[1].MatchRequestUuid,
	}

	js, err := json.Marshal(match)

	w.WriteHeader(200)
	w.Write(js)
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	dbmap := initDb()
	defer dbmap.Db.Close()

	var result Result
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&result)
	checkErr(err, "Decoding JSON failed")

	winningParticipantId, err := dbmap.SelectInt(
		`SELECT id
			FROM participants
			WHERE match_id = :match_id
			AND player_id = :player_id`,
		map[string]interface{}{
			"match_id":  result.MatchId,
			"player_id": result.Winner,
		},
	)
	result.WinningParticipantId = winningParticipantId

	losingParticipantId, err := dbmap.SelectInt(
		`SELECT id
			FROM participants
			WHERE match_id = :match_id
			AND player_id = :player_id`,
		map[string]interface{}{
			"match_id":  result.MatchId,
			"player_id": result.Loser,
		},
	)
	result.LosingParticipantId = losingParticipantId

	dbmap.Insert(&result)

	w.WriteHeader(201)
}
