package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type BonusQuestionGetHandler struct{}

const (
	qSelectBonusQuestion = "SELECT qid, question , answer,relatedEntity FROM bonusquestion"
)

func (q BonusQuestionGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("BonusQuestionGetHandler: new request")

	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "BonusQuestionGetHandler: could not get username from token")

	rows, err := db.DB.Query(qSelectBonusQuestion)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "BonusQuestionGetHandler: could not fetch bonus questions")

	questions := []models.Question{}
	defer rows.Close()
	for rows.Next() {
		question := models.Question{}
		ans := sql.NullString{}
		err := rows.Scan(&question.QuestionID, &question.Question, &ans, &question.RelatedEntity)
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "BonusQuestionGetHandler : db issue in get question query")
		question.Answer = ""
		if ans.Valid {
			question.Answer = ans.String
		}
		questions = append(questions, question)
	}
	util.StructWriter(w, &models.QuestionsModel{Questions: questions})
}
