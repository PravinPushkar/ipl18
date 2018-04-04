package models

// BonusPredictions ..
type BonusPredictions struct {
	Predictions []BonusPrediction `json:"predictions"`
}

// BonusPrediction ..
type BonusPrediction struct {
	QuestionID int    `json:"qid,omitempty"`
	Answer     string `josn:"answer,omitempty"`
	INumber    string `json:"inumber,omitempty"`
}
