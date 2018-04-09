package models

import "time"

type Match struct {
	TeamId1     int               `json:"teamId1"`
	TeamId2     int               `json:"teamId2"`
	Venue       string            `json:"venue"`
	Date        time.Time         `json:"date"`
	Status      string            `json:"status"`
	Result      int               `json:"winner"`
	MatchId     int               `json:"id"`
	MoM         int               `json:"mom"`
	Star        bool              `json:"star"`
	Lock        bool              `json:"lock"`
	Predictions *PredictionsModel `json:"predictions"`
}

type Matches struct {
	Matches []Match `json:"matches"`
}

type MatchStats struct {
	TeamsStats   *[]TeamStats   `json:"teamStats,omitempty"`
	PlayersStats *[]PlayerStats `json:"playerStats,omitempty"`
}
