package dao

import (
	"database/sql"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

type MatchesDAO struct {
}

const (
	qSelectAllMatchesWithPred = "select m.tid1, m.tid2, m.venue, m.matchdate, m.winningteam, m.mid, m.mom, m.star, m.lock, m.status, p.vote_team, p.vote_mom, p.coinused, p.mid, p.pid from match m left outer join prediction p on(m.mid=p.mid) and p.inumber=$1"
	qSelectAllMatches         = "select tid1, tid2, venue, matchdate, winningteam, mid, mom, star, lock, status from match"
	qUpdateMatchResult        = "update match set winningteam=$1,mom=$2,lock=true,status=$3 where mid=$4"
)

func (m MatchesDAO) GetAllMatchesWithPred(inumber string) (*models.Matches, *models.DaoError) {
	log.Println("MatchesDAO: GetAllMatchesWithPred")
	res, err := db.DB.Query(qSelectAllMatchesWithPred, inumber)
	if err != nil {
		log.Println("MatchesDAO:GetAllMatches match info not found")
		return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	defer res.Close()

	var winTeam, mom, voteTeam, voteMom, predMid, predPid sql.NullInt64
	var coinUsed sql.NullBool
	matches := []models.Match{}

	for res.Next() {
		match := models.Match{}
		err = res.Scan(&match.TeamId1, &match.TeamId2, &match.Venue, &match.Date, &winTeam, &match.MatchId, &mom, &match.Star, &match.Lock, &match.Status, &voteTeam, &voteMom, &coinUsed, &predMid, &predPid)
		if err == sql.ErrNoRows {
			log.Println("MatchesDAO: GetAllMatches match info not found")
			return nil, &models.DaoError{http.StatusNotFound, err, err}
		} else if err != nil {
			return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}

		match.MoM = int(mom.Int64)
		match.Result = int(winTeam.Int64)
		pId := int(predPid.Int64)
		if pId != 0 {
			match.Predictions = &models.PredictionsModel{
				TeamVote:     int(voteTeam.Int64),
				MoMVote:      int(voteMom.Int64),
				CoinUsed:     &coinUsed.Bool,
				MatchId:      int(predMid.Int64),
				PredictionId: pId,
				INumber:      inumber,
			}
		}
		matches = append(matches, match)
	}

	return &models.Matches{matches}, nil
}

func (m MatchesDAO) GetAllMatches() (*models.Matches, *models.DaoError) {
	log.Println("MatchesDAO: GetAllMatches")
	log.Println("MatchesDAO: GetAllMatches", qSelectAllMatches)
	res, err := db.DB.Query(qSelectAllMatches)
	if err != nil {
		log.Println("MatchesDAO:GetAllMatches match info not found")
		return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	defer res.Close()

	var winTeam, mom sql.NullInt64
	matches := []models.Match{}

	for res.Next() {
		match := models.Match{}
		err = res.Scan(&match.TeamId1, &match.TeamId2, &match.Venue, &match.Date, &winTeam, &match.MatchId, &mom, &match.Star, &match.Lock, &match.Status)
		if err == sql.ErrNoRows {
			log.Println("MatchesDAO: GetAllMatches match info not found")
			return nil, &models.DaoError{http.StatusNotFound, err, err}
		} else if err != nil {
			return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		matches = append(matches, match)
	}

	log.Println("MatchesDAO: GetAllMatches found ", len(matches), "matches")
	return &models.Matches{matches}, nil
}

func (m MatchesDAO) UpdateResultById(mid, team, mom int, status string) *models.DaoError {
	res, err := db.DB.Exec(qUpdateMatchResult, team, mom, status, mid)
	if err != nil {
		log.Println("MatchesDAO: could not update match result", err)
		return &models.DaoError{http.StatusInternalServerError, err, err}
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println("MatchesDAO: rows not affected by update", err)
		return &models.DaoError{http.StatusInternalServerError, err, err}
	}

	if count != 1 {
		log.Println("MatchesDAO: rows affected", count)
		return &models.DaoError{http.StatusInternalServerError, err, err}
	}
	return nil
}
