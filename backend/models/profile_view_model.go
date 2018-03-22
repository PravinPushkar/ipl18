package models

type ProfileViewModel struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Coin        int    `json:"coin"`
	Alias       string `json:"alias"`
	PicLocation string `json:"piclocation"`
	INumber     string `json:"inumber"`
}