package models

import "database/sql"

type QuestionModel struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	QuestionId    int            `json:"qid"`
	Question      string         `json:"question"`
	Answer        sql.NullString `json:"answer"`
	RelatedEntity string         `json:"relatedEntity"`
}
