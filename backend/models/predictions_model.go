package models

import "fmt"

type PredictionsModel struct {
	PredictionId int    `json:"predId,omitempty"`
	MatchId      int    `json:"mid,omitempty"`
	TeamVote     int    `json:"teamVote,omitempty"`
	MoMVote      int    `json:"momVote,omitempty"`
	CoinUsed     *bool  `json:"coinUsed,omitempty"`
	INumber      string `json:"inumber,omitempty"`
}

func (p *PredictionsModel) String() string {
	coinUsed := false
	if p.CoinUsed != nil {
		coinUsed = *p.CoinUsed
	}
	return fmt.Sprintln(p.PredictionId, p.MatchId, p.TeamVote, p.MoMVote, coinUsed, p.INumber)
}
