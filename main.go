package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/guregu/null.v2"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.
		HandleFunc("/all", AllHandler).
		Methods("DELETE")
	router.
		HandleFunc("/match_requests/{uuid}", CreateMatchRequestHandler).
		Methods("PUT")
	router.
		HandleFunc("/match_requests/{uuid}", GetMatchRequestHandler).
		Methods("GET")
	router.
		HandleFunc("/matches/{uuid}", MatchHandler).
		Methods("GET")
	router.
		HandleFunc("/results", ResultsHandler).
		Methods("POST")

	http.ListenAndServe(fmt.Sprintf(":%v", GetPort()), router)
}

type MatchRequest struct {
	Id          int64       `json:"-" db:"id"`
	Uuid        string      `json:"id" db:"uuid"`
	RequesterId string      `json:"player" db:"requester_id"`
	MatchId     null.String `json:"match_id" db:"-"`
}

type Participant struct {
	Id               int64  `db:"id"`
	MatchId          string `db:"match_id"`
	MatchRequestUuid string `db:"match_request_uuid"`
	PlayerId         string `db:"player_id"`
	OpponentId       string `db:"opponent_id"`
}

type Match struct {
	Id              string `json:"id"`
	MatchRequest1Id string `json:"match_request_1_id"`
	MatchRequest2Id string `json:"match_request_2_id"`
}

type Result struct {
	Id                   int64  `json:"-" db:"id"`
	MatchId              string `json:"match_id" db:"match_id"`
	Winner               string `json:"winner" db:"winner"`
	Loser                string `json:"loser" db:"loser"`
	WinningParticipantId int64  `db:"winning_participant_id"`
	LosingParticipantId  int64  `db:"losing_participant_id"`
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

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

func suitableOpponentMatchRequests(dbmap *gorp.DbMap, requesterId string) []MatchRequest {
	var matchRequests []MatchRequest
	_, err := dbmap.Select(
		&matchRequests,
		`SELECT *
		FROM match_requests
		WHERE requester_id <> :requester_id
		AND uuid NOT IN (
			SELECT match_request_uuid
			FROM participants 
		)
		AND requester_id NOT IN (
			SELECT opponent_id
			FROM participants
			WHERE player_id = :requester_id
		)
		LIMIT 1`,
		map[string]interface{}{"requester_id": requesterId},
	)
	if err != nil {
		checkErr(err, "Error selecting match request")
	}
	return matchRequests
}

func recordMatch(dbmap *gorp.DbMap, openMatchRequest MatchRequest, newMatchRequest MatchRequest) {
	matchIdUuid, err := uuid.NewV4()
	checkErr(err, "Couldn't generate UUID")
	matchId := fmt.Sprintf("%v", matchIdUuid)

	participant1 := Participant{
		MatchId:          matchId,
		MatchRequestUuid: openMatchRequest.Uuid,
		PlayerId:         openMatchRequest.RequesterId,
		OpponentId:       newMatchRequest.RequesterId,
	}
	participant2 := Participant{
		MatchId:          matchId,
		MatchRequestUuid: newMatchRequest.Uuid,
		PlayerId:         newMatchRequest.RequesterId,
		OpponentId:       openMatchRequest.RequesterId,
	}
	err = dbmap.Insert(&participant1, &participant2)
	checkErr(err, "Couldn't insert participants")
}

func initDb() *gorp.DbMap {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "mysql2://gopong:gopong@127.0.0.1:3306/pong_matcher_go_development?reconnect=true"
	}

	url, err := url.Parse(databaseUrl)
	checkErr(err, "Error parsing DATABASE_URL")

	formattedUrl := fmt.Sprintf(
		"%v@tcp(%v)%v",
		url.User,
		url.Host,
		url.Path,
	)

	db, err := sql.Open("mysql", formattedUrl)
	checkErr(err, "failed to establish database connection")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(MatchRequest{}, "match_requests").SetKeys(true, "Id")
	participants := dbmap.AddTableWithName(Participant{}, "participants").SetKeys(true, "Id")
	participants.ColMap("match_request_uuid").SetUnique(true)
	dbmap.AddTableWithName(Result{}, "results").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
