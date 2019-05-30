package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PravinPushkar/ipl18/backend/dao"
	"github.com/PravinPushkar/ipl18/backend/models"

	"github.com/PravinPushkar/ipl18/backend/errors"
	"github.com/PravinPushkar/ipl18/backend/util"
)

// BonusPredictionPostHandler ...
type BonusPredictionPostHandler struct {
	BDao dao.BonusDAO
}

var (
	errNotAllAnswered = fmt.Errorf("all questions were not answered")
	errNoAnswer       = fmt.Errorf("answer not provided")
)

func (bpred BonusPredictionPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("BonusPredictionPostHandler : new request")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	inumber, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "BonusPredictionPostHandler: could not get username from token")

	defer r.Body.Close()
	bonusPredictions := models.BonusPredictions{}
	err = json.NewDecoder(r.Body).Decode(&bonusPredictions)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrEncodingResponse, "BonusPredictionPostHandler: could not decode request body")

	if len(bonusPredictions.Predictions) != 8 {
		errors.ErrWriterPanic(w, http.StatusBadRequest, errNotAllAnswered, errNotAllAnswered, "BonusPredictionHandler: all questions not answered")
	}

	for _, data := range bonusPredictions.Predictions {
		if inumber != data.INumber {
			errors.ErrWriterPanic(w, http.StatusForbidden, errINumberDiff, errors.ErrTokenInfoMismatch, "BonusPredictionPostHandler: token info and payload mismatch")
		}
		if data.Answer == "" {
			errors.ErrWriterPanic(w, http.StatusBadRequest, errNoAnswer, errNoAnswer, fmt.Sprintf("BonusPredictionPostHandler: answer to question not provided %v", data))
		}
	}

	errors.ErrAnalyzePanic(w, bpred.BDao.InsertPredictions(&bonusPredictions), "BonusPredictionPostHandler: unable to insert predictions")
	util.OkWriter(w)
}
