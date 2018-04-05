package models

//ProfileViewModel .
type ProfileViewModel struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Coin        int    `json:"coin"`
	Alias       string `json:"alias"`
	PicLocation string `json:"picLocation"`
	INumber     string `json:"inumber"`
	Points      int    `json:"points"`
}
