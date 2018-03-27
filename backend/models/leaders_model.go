package models

type LeadersModel struct {
	Leaders []Leader `json:"leaders"`
}

type Leader struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	INumber   string `json:"inumber"`
	Point     int    `json:"point"`
}
