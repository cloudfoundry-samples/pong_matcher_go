package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
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

	http.ListenAndServe(fmt.Sprintf(":%v", getPort()), router)
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

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
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
