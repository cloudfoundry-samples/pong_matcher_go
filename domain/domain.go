package domain

import (
	"gopkg.in/guregu/null.v2"
)

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
}
