package models

type ScraperMatchModel struct {
	MatchNo   int    `json:"matchNo"`
	Team1     string `json:"team1"`
	Team2     string `json:"team2"`
	Winner    string `json:"winner"`
	MoM       string `json:"mom"`
	Abandoned bool   `json:"abandoned"`
}
