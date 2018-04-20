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

	qTeamValidMatch         = "select mid from match where mid=$1 and (tid1=$2 or tid2=$2)"
	qPlayerValidMatch       = "select pid from player p where p.pid=$1 and p.tid in (select tid1 from match where mid=$2 union select tid2 from match where mid=$2)"
	qCoinValidMatch         = "select 1,count(coinused) from prediction where coinused=true and inumber=$1 group by inumber union select 2,mid from match where star=true and mid=$2"
	qSelectAllPredictions   = "select pid,inumber,mid,vote_team,vote_mom,coinused from prediction"
	qInsertPredictionResult = "insert into predictionresult(pid,vote_team_correct,vote_mom_correct,points) values($1,$2,$3,$4)"
	qSelectPredForMatches   = "select concat(u.firstname,' ',u.lastname) as name,p.vote_team,p.vote_mom,p.mid from prediction p,ipluser u where p.inumber=u.inumber and p.mid=$1"
)

var (
	errPredictionNotFound = fmt.Errorf("could not find prediction with specified id")
	errTeamInvalid        = fmt.Errorf("team not playing in match")
	errPlayerInvalid      = fmt.Errorf("player not playing in match")
	errCannotUseCoin      = fmt.Errorf("cannot use coin in match")
	errCoinQuotaOver      = fmt.Errorf("all coins used up")
	errNoVote             = fmt.Errorf("vote for team or mom")
)

func (p PredictionDAO) GetAllPredictions() ([]*models.PredictionsModel, error) {
	log.Println("PredictionDAO: GetAllPredictions")
	res, err := db.DB.Query(qSelectAllPredictions)
	if err != nil {
		log.Println("PredictionDAO: GetAllPredictions match info not found")
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	defer res.Close()

	var voteTeam, voteMom sql.NullInt64
	var coinUsed sql.NullBool
	predictions := []*models.PredictionsModel{}

	for res.Next() {
		pred := models.PredictionsModel{}
		err = res.Scan(&pred.PredictionId, &pred.INumber, &pred.MatchId, &voteTeam, &voteMom, &coinUsed)
		if err == sql.ErrNoRows {
			log.Println("PredictionDAO: GetAllPredictions match info not found")
			return nil, &errors.DaoError{http.StatusNotFound, err, err}
		} else if err != nil {
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}

		pred.MoMVote = int(voteMom.Int64)
		pred.TeamVote = int(voteTeam.Int64)

		predictions = append(predictions, &pred)
	}
	log.Println("PredictionDAO: GetAllPredictions found", len(predictions), "predictions")
	return predictions, nil
}

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

func (p PredictionDAO) GetPredictionById(pid int) (*models.PredictionsModel, error) {
	log.Println("PredictionDAO: GetPredictionById", pid)

	var voteTeam, voteMom sql.NullInt64
	info := models.PredictionsModel{}
	info.CoinUsed = new(bool)

	err := db.DB.QueryRow(qSelectPrediction, pid).Scan(&info.PredictionId, &info.INumber, &info.MatchId, &voteTeam, &voteMom, info.CoinUsed)
	info.TeamVote = int(voteTeam.Int64)
	info.MoMVote = int(voteMom.Int64)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.DaoError{http.StatusNotFound, err, errPredictionNotFound}
		}
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	return &info, nil
}

func (p PredictionDAO) GetPredictionsForMatch(mid int) (*models.Predictions, error) {
	log.Println("PredictionDAO: GetPredictionsForMatch", mid)
	var voteTeam, voteMom sql.NullInt64

	res, err := db.DB.Query(qSelectPredForMatches, mid)
	if err != nil {
		log.Println("PredictionDAO: GetPredictionsForMatch error querying preds", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	defer res.Close()

	preds := []*models.PredictionStatModel{}
	for res.Next() {
		info := models.PredictionStatModel{}
		err := res.Scan(&info.Name, &voteTeam, &voteMom, &info.MatchId)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, &errors.DaoError{http.StatusNotFound, err, errPredictionNotFound}
			}
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		info.TeamVote = int(voteTeam.Int64)
		info.MoMVote = int(voteMom.Int64)

		preds = append(preds, &info)
	}
	return &models.Predictions{preds}, nil
}

func (p PredictionDAO) UpdatePredictionById(pid int, info *models.PredictionsModel) error {
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
		if *info.CoinUsed {
			if err := p.checkCoinValidity(info.INumber, info.MatchId); err != nil {
				return err
			}
		}

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
			return &errors.DaoError{http.StatusInternalServerError, err, errPredictionNotFound}
		}

		if rowCount, err := res.RowsAffected(); err != nil {
			return &errors.DaoError{http.StatusInternalServerError, err, errPredictionNotFound}
		} else if rowCount == 0 {
			return &errors.DaoError{http.StatusNotFound, errPredictionNotFound, errPredictionNotFound}
		}
	}
	return nil
}

func (p PredictionDAO) CreateNewPrediction(info *models.PredictionsModel) (*models.GeneralId, error) {
	var tVote, mVote *int
	var coinUsed *bool
	tVote = new(int)
	mVote = new(int)
	coinUsed = new(bool)

	if info.TeamVote == 0 && info.MoMVote == 0 {
		log.Println("PredictionDAO: CreateNewPrediction empty prediction")
		return nil, &errors.DaoError{http.StatusInternalServerError, errNoVote, errNoVote}
	}

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
		if *info.CoinUsed {
			if err := p.checkCoinValidity(info.INumber, info.MatchId); err != nil {
				return nil, err
			}
		}

		*coinUsed = *info.CoinUsed
	}

	log.Println("PredictionDAO:", qInsertNewPrediction, info)
	if err := db.DB.QueryRow(qInsertNewPrediction, info.INumber, info.MatchId, tVote, mVote, coinUsed).Scan(&info.PredictionId); err != nil {
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	log.Println("PredictionDAO: inserted row with id", info.PredictionId)

	return &models.GeneralId{info.PredictionId}, nil
}

func (p PredictionDAO) checkTeamValidity(tid int, mid int) error {
	var midRes int
	err := db.DB.QueryRow(qTeamValidMatch, mid, tid).Scan(&midRes)
	if err != nil {
		if err == sql.ErrNoRows || midRes != mid {
			return &errors.DaoError{http.StatusPreconditionFailed, err, errTeamInvalid}
		}
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	return nil
}

func (p PredictionDAO) checkPlayerValidity(pid int, mid int) error {
	var pidRes int
	err := db.DB.QueryRow(qPlayerValidMatch, pid, mid).Scan(&pidRes)
	if err != nil {
		if err == sql.ErrNoRows || pidRes != pid {
			return &errors.DaoError{http.StatusPreconditionFailed, err, errPlayerInvalid}
		}
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	return nil
}

func (p PredictionDAO) checkCoinValidity(inumber string, mid int) error {
	var vType, vVal int
	log.Println("PredictionDAO:", qCoinValidMatch, inumber, mid)

	if res, err := db.DB.Query(qCoinValidMatch, inumber, mid); err != nil {
		log.Println("error executing query to validate coin")
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	} else {
		defer res.Close()
		results := map[int]int{}

		for res.Next() {
			err := res.Scan(&vType, &vVal)
			if err != nil {
				return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
			}
			results[vType] = vVal
		}

		//no error
		status := 0
		//some coins used
		if val, ok := results[1]; ok {
			//all coins used error
			if val == 12 {
				status = 1
			}
		}
		if _, ok := results[2]; !ok {
			//match is not star match error
			status = 2
		}

		switch status {
		case 1:
			log.Println("PredictionDAO:", inumber, mid, "all coins used")
			return &errors.DaoError{http.StatusPreconditionFailed, errCoinQuotaOver, errCoinQuotaOver}
		case 2:
			log.Println("PredictionDAO:", inumber, mid, "not star match")
			return &errors.DaoError{http.StatusPreconditionFailed, errCannotUseCoin, errCannotUseCoin}
		}
	}
	return nil
}

func (p PredictionDAO) WritePredictionResult(pid, team, mom, points int) error {
	log.Println("PredictionDAO: UpdatePredictionResult")
	if _, err := db.DB.Exec(qInsertPredictionResult, pid, team, mom, points); err != nil {
		log.Println("PredictionDAO: failed to update prediction result", err)
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	return nil
}
