package models

type User struct {
	INumber     string `json:"inumber"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Password    string `json:"password"`
	Coins       int    `json:"coins"`
	UID         int    `json:"uid"`
	Alias       string `json:"alias"`
	PicLocation string `json:"pic_loc"`
}
