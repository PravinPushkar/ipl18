package models

import "database/sql"

type Question struct {
	QuestionId    int            `json:"qid"`
	Question      string         `json:"question"`
	Answer        sql.NullString `json:"answer"`
	RelatedEntity string         `json:"relatedEntity"`
}
