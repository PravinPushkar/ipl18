package models

type Team struct {
	TeamId      int    `json:"id"`
	TeamName    string `json:"name"`
	ShortName   string `json:"shortName"`
	PicLocation string `json:"picLocation"`
}

type Teams struct {
	Teams []*Team `json:"teams"`
}

type TeamStats struct {
	TeamId int `json:"teamId"`
	Votes  int `json:"votes"`
}
