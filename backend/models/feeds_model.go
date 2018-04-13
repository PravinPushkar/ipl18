package models

import (
	"time"
)

// FeedsAuthModel ..
type FeedsAuthModel struct {
	Authorization string `json:"Authorization"`
}

// FeedsMessageModel ..
type FeedsMessageModel struct {
	INumber     string    `json:"inumber"`
	Name        string    `json:"name"`
	Message     string    `json:"message"`
	Date        time.Time `json:"date"`
	PicLocation string    `json:"picLocation"`
}
