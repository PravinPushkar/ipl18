package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"

	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

// BonusPredictionPostHandler ...
type BonusPredictionPostHandler struct{}

const (
	qInsertIntoBonusPrediction = "INSERT INTO bonusprediction (qid , answer, inumber) VALUES"
)

var (
	errNotAllAnswered = fmt.Errorf("all questions were not answered")
	errNoAnswer       = fmt.Errorf("answer not provided")
)

func (bpred BonusPredictionPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("BonusPredictionPostHandler : insert bonus prediction")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked")
		}
	}()

	inumber, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "BonusPredictionPostHandler: could not get username from token")

	defer r.Body.Close()
	bonusPredictions := models.BonusPredictions{}
	err = json.NewDecoder(r.Body).Decode(&bonusPredictions)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrEncodingResponse, "BonusPredictionPostHandler: could not decode request body")

	suffixes := []string{}
	info := []interface{}{}
	j := 1
	for _, data := range bonusPredictions.Predictions {
		if inumber != data.INumber {
			errors.ErrWriterPanic(w, http.StatusForbidden, errINumberDiff, errors.ErrTokenInfoMismatch, "BonusPredictionPostHandler: token info and payload mismatch")
		}
		if data.Answer == "" {
			errors.ErrWriterPanic(w, http.StatusBadRequest, errNoAnswer, errNoAnswer, fmt.Sprintf("BonusPredictionPostHandler: answer to question not provided %v", data))
		}
		suffixes = append(suffixes, fmt.Sprintf("($%d,$%d,$%d)", j, j+1, j+2))
		info = append(info, data.QuestionID, data.Answer, data.INumber)
		j += 3
	}

	if len(suffixes) != 8 {
		errors.ErrWriterPanic(w, http.StatusBadRequest, errNotAllAnswered, errNotAllAnswered, "BonusPredictionHandler: all questions not answered")
	}

	query := qInsertIntoBonusPrediction + strings.Join(suffixes, ",")
	log.Println("BonusPredictionHandler:", query, info)

	_, err = db.DB.Query(query, info...)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "BonusPredictionPostHandler: insert into bonus prediction failed")
	util.OkWriter(w)
}
