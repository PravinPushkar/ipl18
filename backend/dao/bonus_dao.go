package dao

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

type BonusDAO struct{}

const (
	qSelectBonusQuestion       = "SELECT qid, question , answer,relatedEntity,points FROM bonusquestion"
	qInsertIntoBonusPrediction = "INSERT INTO bonusprediction (qid , answer, inumber) VALUES"
)

func (b BonusDAO) GetAllBonusQuestions() (*models.QuestionsModel, error) {
	rows, err := db.DB.Query(qSelectBonusQuestion)
	if err != nil {
		log.Println("BonusDAO:GetAllBonusQuestions: error querying questions", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	questions := []*models.Question{}
	defer rows.Close()
	for rows.Next() {
		question := models.Question{}
		ans := sql.NullString{}
		err := rows.Scan(&question.QuestionID, &question.Question, &ans, &question.RelatedEntity, &question.Points)
		if err != nil {
			log.Println("BonusDAO: GetAllBonusQuestions: error fetching question", err)
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		question.Answer = ""
		if ans.Valid {
			question.Answer = ans.String
		}
		questions = append(questions, &question)
	}
	return &models.QuestionsModel{questions}, nil
}

func (b BonusDAO) InsertPredictions(predictions *models.BonusPredictions) error {
	suffixes := []string{}
	info := []interface{}{}
	j := 1
	for _, data := range predictions.Predictions {
		suffixes = append(suffixes, fmt.Sprintf("($%d,$%d,$%d)", j, j+1, j+2))
		info = append(info, data.QuestionID, data.Answer, data.INumber)
		j += 3
	}

	query := qInsertIntoBonusPrediction + strings.Join(suffixes, ",")
	log.Println("BonusDAO: InsertPredictions:", query, info)

	res, err := db.DB.Exec(query, info...)
	if err != nil {
		log.Println("BonusDAO: error inserting bonus predictions", err)
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	if num, err := res.RowsAffected(); err != nil || int(num) != len(predictions.Predictions) {
		log.Println("BonusPredictionPostHandler: insert into bonus prediction failed", err, num)
		return &errors.DaoError{http.StatusInternalServerError, errors.ErrDBIssue, errors.ErrDBIssue}
	}

	return nil
}
