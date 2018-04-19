package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type PredictionHandler struct {
	PDao predictor
}

type predictor interface {
	CanMakePrediction(int) bool
	GetPredictionById(int) (*models.PredictionsModel, error)
	CreateNewPrediction(*models.PredictionsModel) (*models.GeneralId, error)
	UpdatePredictionById(int, *models.PredictionsModel) error
	GetPredictionsForMatches() (*models.Predictions, error)
}

var (
	errInvalidPredId       = fmt.Errorf("prediction id not valid")
	errTimeToPredictPassed = fmt.Errorf("cannot predict after 15 minutes to game")
)

func (p PredictionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("PredictionHandler: new request")

	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	inumber, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "PredictionHandler: could not get username from token")

	switch r.Method {
	case http.MethodPost:
		p.handlePost(w, r, inumber)
	case http.MethodPut:
		p.handlePut(w, r, inumber)
	case http.MethodGet:
		p.handleGet(w, r, inumber)
	}
}

func (p PredictionHandler) parseBody(w http.ResponseWriter, r *http.Request, inumber string) *models.PredictionsModel {
	defer r.Body.Close()
	info := models.PredictionsModel{}
	err := json.NewDecoder(r.Body).Decode(&info)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrParseRequest, "PredictionHandler: cannot parse")
	if info.INumber != inumber {
		errors.ErrWriterPanic(w, http.StatusForbidden, errors.ErrTokenInfoMismatch, errors.ErrTokenInfoMismatch, "PredictionHandler: invalid inumber")
	}
	return &info
}

func (p PredictionHandler) handlePost(w http.ResponseWriter, r *http.Request, inumber string) {
	info := p.parseBody(w, r, inumber)
	if !p.PDao.CanMakePrediction(info.MatchId) {
		errors.ErrWriterPanic(w, http.StatusPreconditionFailed, errTimeToPredictPassed, errTimeToPredictPassed, "PredictionHandler: time has passed")
	}

	if gid, err := p.PDao.CreateNewPrediction(info); err != nil {
		errors.ErrAnalyzePanic(w, err, "PredictionHandler:")
	} else {
		util.StructWriter(w, gid)
	}
}

func (p PredictionHandler) handlePut(w http.ResponseWriter, r *http.Request, inumber string) {
	info := p.parseBody(w, r, inumber)
	if !p.PDao.CanMakePrediction(info.MatchId) {
		errors.ErrWriterPanic(w, http.StatusPreconditionFailed, errTimeToPredictPassed, errTimeToPredictPassed, "PredictionHandler: time has passed")
	}

	vars := mux.Vars(r)
	if pStr, ok := vars["id"]; ok {
		pid, err := strconv.Atoi(pStr)
		errors.ErrWriterPanic(w, http.StatusBadRequest, err, errInvalidPredId, "PredictionHandler: invalid prediction id in put")

		errDao := p.PDao.UpdatePredictionById(pid, info)
		errors.ErrAnalyzePanic(w, errDao, "PredictionHandler:")
	}

	util.OkWriter(w)
}

func (p PredictionHandler) handleGet(w http.ResponseWriter, r *http.Request, inumber string) {
	if strings.Contains(r.URL.Path, "userStats") {
		fmt.Println("XXX in stats")
		info, err := p.PDao.GetPredictionsForMatches()
		errors.ErrAnalyzePanic(w, err, "PredictionHandler:")
		util.StructWriter(w, info)
		return
	}

	vars := mux.Vars(r)
	if pStr, ok := vars["id"]; ok {
		pid, err := strconv.Atoi(pStr)
		errors.ErrWriterPanic(w, http.StatusBadRequest, err, errInvalidPredId, "PredictionHandler: invalid prediction id in get")

		info, err := p.PDao.GetPredictionById(pid)
		errors.ErrAnalyzePanic(w, err, "PredictionHandler:")
		util.StructWriter(w, info)
	}
}
