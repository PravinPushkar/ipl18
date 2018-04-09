package dao

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

type PredictionDAO struct{}

const (
	qSelectPrediction    = "select pid,inumber,mid,vote_team,vote_mom,coinused from prediction where pid=$1"
	qUpdatePrediction    = "update prediction set"
	qInsertNewPrediction = "insert into prediction(inumber,mid,vote_team,vote_mom,coinused) values($1,$2,$3,$4,$5) returning pid"
	qSelectMatchTime     = "select matchdate from match where mid=$1"

	qTeamValidMatch   = "select mid from match where mid=$1 and (tid1=$2 or tid2=$2)"
	qPlayerValidMatch = "select pid from player p where p.pid=$1 and p.tid in (select tid1 from match where mid=$2 union select tid2 from match where mid=$2)"
)

var (
	errPredictionNotFound = fmt.Errorf("could not find prediction with specified id")
	errTeamInvalid        = fmt.Errorf("team not playing in match")
	errPlayerInvalid      = fmt.Errorf("player not playing in match")
)

func (p PredictionDAO) CanMakePrediction(mid int) bool {
	var dt time.Time
	if err := db.DB.QueryRow(qSelectMatchTime, mid).Scan(&dt); err != nil {
		log.Println("PredictionDAO: CanMakePrediction: unable to get match time", err)
	}

	log.Println("PredictionHandler: match time", dt, "current time ", time.Now())
	if dt.Sub(time.Now()).Seconds() < 60*15.0 {
		return false
	}

	return true
}

func (p PredictionDAO) GetPredictionById(pid int) (*models.PredictionsModel, *models.DaoError) {
	log.Println("PredictionDAO: GetPredictionById", pid)

	var voteTeam, voteMom sql.NullInt64
	info := models.PredictionsModel{}
	info.CoinUsed = new(bool)

	err := db.DB.QueryRow(qSelectPrediction, pid).Scan(&info.PredictionId, &info.INumber, &info.MatchId, &voteTeam, &voteMom, info.CoinUsed)
	info.TeamVote = int(voteTeam.Int64)
	info.MoMVote = int(voteMom.Int64)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.DaoError{http.StatusNotFound, err, errPredictionNotFound}
		}
		return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	return &info, nil
}

func (p PredictionDAO) UpdatePredictionById(pid int, info *models.PredictionsModel) *models.DaoError {
	log.Println("PredictionDAO: UpdatePredictionById", pid, info)
	var suffixes []string
	var values []interface{}
	index := 1

	if info.TeamVote != 0 {
		if err := p.checkTeamValidity(info.TeamVote, info.MatchId); err != nil {
			return err
		}

		suffixes = append(suffixes, fmt.Sprintf("vote_team=$%d", index))
		values = append(values, info.TeamVote)
		index += 1
	}

	if info.MoMVote != 0 {
		if err := p.checkPlayerValidity(info.MoMVote, info.MatchId); err != nil {
			return err
		}

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
		log.Println("UpdatePredictionById: ", query, suffixes, values)
		res, err := db.DB.Exec(query, values...)
		if err != nil {
			return &models.DaoError{http.StatusInternalServerError, err, errPredictionNotFound}
		}

		if rowCount, err := res.RowsAffected(); err != nil {
			return &models.DaoError{http.StatusInternalServerError, err, errPredictionNotFound}
		} else if rowCount == 0 {
			return &models.DaoError{http.StatusNotFound, errPredictionNotFound, errPredictionNotFound}
		}
	}
	return nil
}

func (p PredictionDAO) CreateNewPrediction(info *models.PredictionsModel) (*models.GeneralId, *models.DaoError) {
	var tVote, mVote *int
	var coinUsed *bool
	tVote = new(int)
	mVote = new(int)
	coinUsed = new(bool)

	if info.TeamVote != 0 {
		if err := p.checkTeamValidity(info.TeamVote, info.MatchId); err != nil {
			return nil, err
		}
		*tVote = info.TeamVote
	} else {
		tVote = nil
	}

	if info.MoMVote != 0 {
		if err := p.checkPlayerValidity(info.MoMVote, info.MatchId); err != nil {
			return nil, err
		}
		*mVote = info.MoMVote
	} else {
		mVote = nil
	}

	if info.CoinUsed != nil {
		*coinUsed = *info.CoinUsed
	}

	log.Println("PredictionDAO:", qInsertNewPrediction, info)
	if err := db.DB.QueryRow(qInsertNewPrediction, info.INumber, info.MatchId, tVote, mVote, coinUsed).Scan(&info.PredictionId); err != nil {
		return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	log.Println("PredictionDAO: inserted row with id", info.PredictionId)

	return &models.GeneralId{info.PredictionId}, nil
}

func (p PredictionDAO) checkTeamValidity(tid int, mid int) *models.DaoError {
	var midRes int
	err := db.DB.QueryRow(qTeamValidMatch, mid, tid).Scan(&midRes)
	if err != nil {
		if err == sql.ErrNoRows || midRes != mid {
			return &models.DaoError{http.StatusPreconditionFailed, err, errTeamInvalid}
		}
		return &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	return nil
}

func (p PredictionDAO) checkPlayerValidity(pid int, mid int) *models.DaoError {
	var pidRes int
	err := db.DB.QueryRow(qPlayerValidMatch, pid, mid).Scan(&pidRes)
	if err != nil {
		if err == sql.ErrNoRows || pidRes != pid {
			return &models.DaoError{http.StatusPreconditionFailed, err, errPlayerInvalid}
		}
		return &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	return nil
}
