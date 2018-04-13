package models

// LeadersModel .
type LeadersModel struct {
	Leaders []*Leader `json:"leaders"`
}

// Leader .
type Leader struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	INumber     string `json:"inumber"`
	Points      int    `json:"points"`
	Alias       string `json:"alias"`
	Piclocation string `json:"piclocation"`
}
