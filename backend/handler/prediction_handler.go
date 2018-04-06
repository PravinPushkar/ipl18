package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type PredictionHandler struct{}

const (
	qInsertNewPrediction = "insert into prediction(inumber,mid,vote_team,vote_mom,coinused) values($1,$2,$3,$4,$5) returning pid"
	qUpdatePrediction    = "update prediction set"
	qSelectPrediction    = "select pid,inumber,mid,vote_team,vote_mom,coinused from prediction where pid=$1"
	qUpdateProfile       = "update ipluser set coin=coin%s where inumber=$1"
	qSelectMatchTime     = "select matchdate from match where mid=$1"
)

var (
	errInvalidPredId       = fmt.Errorf("prediction id not valid")
	errPredictionNotFound  = fmt.Errorf("could not find prediction with specified id")
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

func (p PredictionHandler) checkTime(w http.ResponseWriter, mid int) {
	log.Println("PredictionHandler: checking match time")

	var dt time.Time
	err := db.DB.QueryRow(qSelectMatchTime, mid).Scan(&dt)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "PredictionHandler: could not fetch match time info")

	log.Println("PredictionHandler: match time", dt, "current time ", time.Now())
	if dt.Sub(time.Now()).Seconds() < 60*15.0 {
		errors.ErrWriterPanic(w, http.StatusForbidden, errTimeToPredictPassed, errTimeToPredictPassed, "PredictionHandler: cannot allow prediction, time passed")
	}
}

func (p PredictionHandler) handlePost(w http.ResponseWriter, r *http.Request, inumber string) {
	info := p.parseBody(w, r, inumber)
	p.checkTime(w, info.MatchId)

	var tVote, mVote *int
	var coinUsed *bool
	tVote = new(int)
	mVote = new(int)
	coinUsed = new(bool)

	if info.TeamVote != 0 {
		*tVote = info.TeamVote
	} else {
		tVote = nil
	}

	if info.MoMVote != 0 {
		*mVote = info.MoMVote
	} else {
		mVote = nil
	}

	if info.CoinUsed != nil {
		*coinUsed = *info.CoinUsed
	}

	log.Println("PredictionHandler:", qInsertNewPrediction, info)
	err := db.DB.QueryRow(qInsertNewPrediction, info.INumber, info.MatchId, tVote, mVote, coinUsed).Scan(&info.PredictionId)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "PredictionHandler: could not add new prediction (post)")

	log.Println("PredictionHandler: inserted row with id", info.PredictionId)
	util.StructWriter(w, &models.GeneralId{info.PredictionId})
}

func (p PredictionHandler) handlePut(w http.ResponseWriter, r *http.Request, inumber string) {
	info := p.parseBody(w, r, inumber)
	p.checkTime(w, info.MatchId)

	var suffixes []string
	var values []interface{}
	index := 1

	vars := mux.Vars(r)
	if p, ok := vars["id"]; ok {
		pid, err := strconv.Atoi(p)
		errors.ErrWriterPanic(w, http.StatusBadRequest, err, errInvalidPredId, "PredictionHandler: invalid prediction id in put")

		if info.TeamVote != 0 {
			suffixes = append(suffixes, fmt.Sprintf("vote_team=$%d", index))
			values = append(values, info.TeamVote)
			index += 1
		}

		if info.MoMVote != 0 {
			suffixes = append(suffixes, fmt.Sprintf("vote_mom=$%d", index))
			values = append(values, info.MoMVote)
			index += 1
		}

		if info.CoinUsed != nil {
			suffixes = append(suffixes, fmt.Sprintf("coinused=$%d", index))
			values = append(values, *info.CoinUsed)
			index += 1
		}

		query := qUpdatePrediction
		if len(suffixes) != 0 {
			query = fmt.Sprintf("%s %s where pid=%d", query, strings.Join(suffixes, ","), pid)
			log.Println("PredictionHandler: ", query, suffixes, values)
			res, err := db.DB.Exec(query, values...)
			errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "PredictionHandler: could not handle update prediction (put)")

			rowCount, err := res.RowsAffected()
			if rowCount == 0 || err != nil {
				errors.ErrWriterPanic(w, http.StatusNotFound, errPredictionNotFound, errPredictionNotFound, "PredictionHandler: could not handle update prediction (put)")
			}
		}
	}

	util.OkWriter(w)
}

func (p PredictionHandler) handleGet(w http.ResponseWriter, r *http.Request, inumber string) {
	vars := mux.Vars(r)
	if p, ok := vars["id"]; ok {
		pid, err := strconv.Atoi(p)
		errors.ErrWriterPanic(w, http.StatusBadRequest, err, errInvalidPredId, "PredictionHandler: invalid prediction id in get")

		var voteTeam, voteMom sql.NullInt64
		info := models.PredictionsModel{}
		info.CoinUsed = new(bool)

		err = db.DB.QueryRow(qSelectPrediction, pid).Scan(&info.PredictionId, &info.INumber, &info.MatchId, &voteTeam, &voteMom, info.CoinUsed)
		if err == sql.ErrNoRows {
			errors.ErrWriterPanic(w, http.StatusNotFound, err, errPredictionNotFound, "PredictionHandler: could not find prediction")
		}
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "PredictionHandler: error selecting prediction")
		info.TeamVote = int(voteTeam.Int64)
		info.MoMVote = int(voteMom.Int64)

		util.StructWriter(w, &info)
	}
}
