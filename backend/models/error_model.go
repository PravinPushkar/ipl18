package models

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DaoError struct {
	Code      int
	ActualErr error
	UserErr   error
}
