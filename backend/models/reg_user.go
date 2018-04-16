package models

type User struct {
	INumber     string `json:"inumber"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Password    string `json:"password"`
	Coins       int    `json:"coins"`
	UID         int    `json:"id"`
	Alias       string `json:"alias"`
	PicLocation string `json:"picLocation"`
	Points      int    `json:"points,omitempty"`
}

type UserBasic struct {
	INumber     string `json:"inumber"`
	Name        string `json:"name"`
	PicLocation string `json:"picLocation"`
}
