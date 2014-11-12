package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/guregu/null.v2"
	"log"
	"net/http"
	"strings"
)

type MatchRequest struct {
	Id          int64       `json:"-" db:"id"`
	Uuid        string      `json:"id" db:"uuid"`
	RequesterId string      `json:"player" db:"requester_id"`
	MatchId     null.String `json:"match_id" db:"-"`
}

type Participant struct {
	Id             int64  `db:"id"`
	MatchId        string `db:"match_id"`
	MatchRequestId string `db:"match_request_id"`
	PlayerId       string `db:"player_id"`
	OpponentId     string `db:"opponent_id"`
}

type Match struct {
	Id              string
	MatchRequest1Id string
	MatchRequest2Id string
}

func main() {
	dbmap := initDb()
	defer dbmap.Db.Close()

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		err := dbmap.TruncateTables()
		checkErr(err, "Truncation failed")
	})

	http.HandleFunc("/match_requests/", func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		uuid := urlParts[len(urlParts)-1]

		switch r.Method {

		case "PUT":
			var matchRequest MatchRequest

			decoder := json.NewDecoder(r.Body)

			err := decoder.Decode(&matchRequest)
			checkErr(err, "Decoding JSON failed")

			matchRequest.Uuid = uuid

			fmt.Printf("decoded from PUT: %+v\n", matchRequest)

			err = dbmap.Insert(&matchRequest)
			checkErr(err, "Creation of MatchRequest failed")

			openMatchRequests := suitableOpponentMatchRequests(dbmap, matchRequest.RequesterId)
			if len(openMatchRequests) > 0 {
				fmt.Printf("Found open match request: %v\n", openMatchRequests[0])
				recordMatch(dbmap, openMatchRequests[0], matchRequest)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)

		case "GET":
			matchRequest := MatchRequest{}
			err := dbmap.SelectOne(
				&matchRequest,
				"SELECT * FROM match_requests WHERE uuid = ?", uuid,
			)
			if err == nil {
				matchId, err := dbmap.SelectStr(
					`SELECT match_id
					FROM participants
					WHERE match_request_id = ?`,
					uuid,
				)
				if err == nil && matchId != "" {
					matchRequest.MatchId = null.StringFrom(matchId)
				}

				fmt.Printf("reading: %+v\n", matchRequest)

				js, err := json.Marshal(matchRequest)
				checkErr(err, "Error writing JSON")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(js)
			} else {
				w.WriteHeader(404)
			}
		}
	})

	http.HandleFunc("/matches/", func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		matchId := urlParts[len(urlParts)-1]

		var participants []Participant
		_, err := dbmap.Select(
			&participants,
			`SELECT * 
			FROM participants
			WHERE match_id = ?`,
			matchId,
		)
		checkErr(err, "Error getting participants")

		fmt.Printf("Participants: %v\n", participants)

		match := Match{
			Id:              matchId,
			MatchRequest1Id: participants[0].MatchRequestId,
			MatchRequest2Id: participants[1].MatchRequestId,
		}

		js, err := json.Marshal(match)

		w.WriteHeader(200)
		w.Write(js)
	})

	http.ListenAndServe(":3000", nil)
}

func suitableOpponentMatchRequests(dbmap *gorp.DbMap, requesterId string) []MatchRequest {
	var matchRequests []MatchRequest
	_, err := dbmap.Select(
		&matchRequests,
		`SELECT *
		FROM match_requests
		WHERE requester_id <> :requester_id
		AND id NOT IN (
			SELECT match_request_id
			FROM participants 
			WHERE opponent_id = :requester_id
			OR player_id = :requester_id
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
		MatchId:        matchId,
		MatchRequestId: openMatchRequest.Uuid,
		PlayerId:       openMatchRequest.RequesterId,
		OpponentId:     newMatchRequest.RequesterId,
	}
	participant2 := Participant{
		MatchId:        matchId,
		MatchRequestId: newMatchRequest.Uuid,
		PlayerId:       newMatchRequest.RequesterId,
		OpponentId:     openMatchRequest.RequesterId,
	}
	dbmap.Insert(&participant1, &participant2)
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "gopong:gopong@/pong_matcher_go_development?charset=utf8&parseTime=True")
	checkErr(err, "failed to establish database connection")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(MatchRequest{}, "match_requests").SetKeys(true, "Id")
	dbmap.AddTableWithName(Participant{}, "participants").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
