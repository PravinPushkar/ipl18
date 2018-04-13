package models

// QuestionsModel .
type QuestionsModel struct {
	Questions []*Question `json:"questions"`
}

// Question ...
type Question struct {
	QuestionID    int    `json:"qid"`
	Question      string `json:"question"`
	Answer        string `json:"answer"`
	RelatedEntity string `json:"relatedEntity"`
	Points        int    `json:"points"`
}
