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
	qSelectAllMatches = "select m.tid1, m.tid2, m.venue, m.matchdate, m.winningteam, m.mid, m.mom, m.star, m.lock, m.status, p.vote_team, p.vote_mom, p.coinused, p.mid, p.pid from match m left outer join prediction p on(m.mid=p.mid) and p.inumber=$1"
	qSelectMatchById  = qSelectAllMatches + " and m.mid=$2"
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

	inumber, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "MatchesGetHandler: could not get username from token")

	vars := mux.Vars(r)
	if val, ok := vars["id"]; ok {
		//specific match
		log.Println("MatchesGetHandler: request to get latest match by id", val)
		m.handleSpecificMatch(w, r, val, inumber)
	} else {
		//all matches
		log.Println("MatchesGetHandler: request to get all match info")
		m.handleAllMatches(w, r, inumber)
	}
}

func (m MatchesGetHandler) handleSpecificMatch(w http.ResponseWriter, r *http.Request, mid string, inumber string) {
	match := models.Match{}
	var winTeam, mom, voteTeam, voteMom, predMid, predPid sql.NullInt64
	var coinUsed sql.NullBool

	err := db.DB.QueryRow(qSelectMatchById, inumber, mid).Scan(&match.TeamId1, &match.TeamId2, &match.Venue, &match.Date, &winTeam, &match.MatchId, &mom, &match.Star, &match.Lock, &match.Status, &voteTeam, &voteMom, &coinUsed, &predMid, &predPid)

	if err == sql.ErrNoRows {
		errors.ErrWriterPanic(w, http.StatusNotFound, errMatchNotFound, errMatchNotFound, "MatchesGetHandler: match info not found")
	}

	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "MatchesGetHandler: match info not found")
	match.MoM = int(mom.Int64)
	match.Result = int(winTeam.Int64)
	pId := int(predPid.Int64)
	if pId != 0 {
		match.Predictions = &models.PredictionsModel{
			TeamVote:     int(voteTeam.Int64),
			MoMVote:      int(voteMom.Int64),
			CoinUsed:     coinUsed.Bool,
			MatchId:      int(predMid.Int64),
			PredictionId: pId,
		}
	}
	util.StructWriter(w, &match)
}

func (m MatchesGetHandler) handleAllMatches(w http.ResponseWriter, r *http.Request, inumber string) {
	res, err := db.DB.Query(qSelectAllMatches, inumber)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "MatchesGetHandler: match info not found")

	defer res.Close()

	matches := []models.Match{}
	for res.Next() {
		var winTeam, mom, voteTeam, voteMom, predMid, predPid sql.NullInt64
		var coinUsed sql.NullBool
		match := models.Match{}
		err := res.Scan(&match.TeamId1, &match.TeamId2, &match.Venue, &match.Date, &winTeam, &match.MatchId, &mom, &match.Star, &match.Lock, &match.Status, &voteTeam, &voteMom, &coinUsed, &predMid, &predPid)
		if err == sql.ErrNoRows {
			errors.ErrWriterPanic(w, http.StatusNotFound, errMatchNotFound, errMatchNotFound, "MatchesGetHandler: match info not found")
		}
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "MatchesGetHandler: match info not found")
		match.MoM = int(mom.Int64)
		match.Result = int(winTeam.Int64)
		pId := int(predPid.Int64)
		if pId != 0 {
			match.Predictions = &models.PredictionsModel{
				TeamVote:     int(voteTeam.Int64),
				MoMVote:      int(voteMom.Int64),
				CoinUsed:     coinUsed.Bool,
				MatchId:      int(predMid.Int64),
				PredictionId: pId,
			}
		}
		matches = append(matches, match)
	}
	util.StructWriter(w, &models.Matches{matches})
}
