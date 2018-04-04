package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

const (
	qSelectAllMatches  = "select tid1, tid2, venue, matchdate, winningteam, mid, mom, star, lock, status from match"
	qSelectMatchById   = qSelectAllMatches + " where mid=$1"
	qSelectLatestMatch = qSelectAllMatches + " order by mid desc limit 1 offset 0"
)

var (
	errMatchNotFound = fmt.Errorf("info for match not found")
)

type MatchesGetHandler struct{}

func (m MatchesGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("MatchesGetHandler: new request")

	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "MatchesGetHandler: could not get username from token")

	if r.URL.Query().Get("q") == "latest" {
		//get latest match stuff
		log.Println("MatchesGetHandler: request to get latest match")
		m.handleLatestMatch(w, r)
		return
	}

	vars := mux.Vars(r)
	if val, ok := vars["id"]; ok {
		//specific match
		log.Println("MatchesGetHandler: request to get latest match by id", val)
		m.handleSpecificMatch(w, r, val)
	} else {
		//all matches
		log.Println("MatchesGetHandler: request to get all match info")
		m.handleAllMatches(w, r)
	}
}

func (m MatchesGetHandler) handleLatestMatch(w http.ResponseWriter, r *http.Request) {
	match := models.Match{}
	var winTeam, mom sql.NullInt64
	err := db.DB.QueryRow(qSelectLatestMatch).Scan(&match.TeamId1, &match.TeamId2, &match.Venue, &match.Date, &winTeam, &match.MatchId, &mom, &match.Star, &match.Lock, &match.Status)

	if err == sql.ErrNoRows {
		errors.ErrWriterPanic(w, http.StatusNotFound, errMatchNotFound, errMatchNotFound, "MatchesGetHandler: match info not found")
	}

	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "MatchesGetHandler: match info not found")
	match.MoM = int(mom.Int64)
	match.Result = int(winTeam.Int64)
	util.StructWriter(w, &match)
}

func (m MatchesGetHandler) handleSpecificMatch(w http.ResponseWriter, r *http.Request, id string) {
	match := models.Match{}
	var winTeam, mom sql.NullInt64
	err := db.DB.QueryRow(qSelectMatchById, id).Scan(&match.TeamId1, &match.TeamId2, &match.Venue, &match.Date, &winTeam, &match.MatchId, &mom, &match.Star, &match.Lock, &match.Status)

	if err == sql.ErrNoRows {
		errors.ErrWriterPanic(w, http.StatusNotFound, errMatchNotFound, errMatchNotFound, "MatchesGetHandler: match info not found")
	}

	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "MatchesGetHandler: match info not found")
	match.MoM = int(mom.Int64)
	match.Result = int(winTeam.Int64)
	util.StructWriter(w, &match)
}

func (m MatchesGetHandler) handleAllMatches(w http.ResponseWriter, r *http.Request) {
	res, err := db.DB.Query(qSelectAllMatches)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "MatchesGetHandler: match info not found")

	defer res.Close()

	matches := []models.Match{}
	for res.Next() {
		var winTeam, mom sql.NullInt64
		match := models.Match{}
		err := res.Scan(&match.TeamId1, &match.TeamId2, &match.Venue, &match.Date, &winTeam, &match.MatchId, &mom, &match.Star, &match.Lock, &match.Status)
		if err == sql.ErrNoRows {
			errors.ErrWriterPanic(w, http.StatusNotFound, errMatchNotFound, errMatchNotFound, "MatchesGetHandler: match info not found")
		}
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "MatchesGetHandler: match info not found")
		match.MoM = int(mom.Int64)
		match.Result = int(winTeam.Int64)
		matches = append(matches, match)
	}
	util.StructWriter(w, &models.Matches{matches})
}
