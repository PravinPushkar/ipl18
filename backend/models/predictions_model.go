package models

type PredictionsModel struct {
	PredictionId int    `json:"predId,omitempty"`
	MatchId      int    `json:"mid,omitempty"`
	TeamVote     int    `json:"teamVote,omitempty"`
	MoMVote      int    `json:"momVote,omitempty"`
	CoinUsed     *bool  `json:"coinUsed,omitempty"`
	INumber      string `json:"inumber,omitempty"`
}
