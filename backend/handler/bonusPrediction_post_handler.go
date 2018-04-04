package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"

	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

// BonusPredictionPostHandler ...
type BonusPredictionPostHandler struct{}

const (
	qInsertIntoBonusPrediction = "INSERT INTO bonusprediction (qid , answer, inumber) VALUES ($1, $2, $3)"
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

	decoder := json.NewDecoder(r.Body)
	bonusPredictions := []models.BonusPrediction{}
	err = decoder.Decode(&bonusPredictions)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrEncodingResponse, "BonusPredictionPostHandler: could not decode request body")

	for _, data := range bonusPredictions {
		if inumber != data.INumber {
			errors.ErrWriterPanic(w, http.StatusForbidden, errINumberDiff, errors.ErrTokenInfoMismatch, "BonusPredictionPostHandler: token info and payload mismatch")
		}
		_, err = db.DB.Query(qInsertIntoBonusPrediction, data.QuestionID, data.Answer, data.INumber)
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "BonusPredictionPostHandler: insert into bonus prediction failed")
	}
	util.OkWriter(w)
}
