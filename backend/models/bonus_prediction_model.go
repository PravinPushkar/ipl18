package models

// BonusPrediction ..
type BonusPrediction struct {
	QuestionID int    `json:"qid,omitempty"`
	Answer     string `josn:"answer,omitempty"`
	INumber    string `json:"inumber,omitempty"`
}
