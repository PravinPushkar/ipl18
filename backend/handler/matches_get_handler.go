package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

var (
	errMatchNotFound  = fmt.Errorf("info for match not found")
	errInvalidMatchId = fmt.Errorf("match id not valid")
)

type MatchesGetHandler struct {
	MDao dao.MatchesDAO
}

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
		if mid, err := strconv.Atoi(val); err != nil {
			errors.ErrWriterPanic(w, http.StatusBadRequest, err, errInvalidMatchId, "MatchesGetHandler: match id not valid")
		} else {
			//specific match stats
			if strings.Contains(r.URL.Path, "/stats") {
				log.Println("MatchesGetHandler: request to get match stats by id", val)
				m.handleMatchStats(w, r, mid)
				return
			}
			//specific match
			log.Println("MatchesGetHandler: request to get match by id", val)
			m.handleSpecificMatch(w, r, mid, inumber)
		}
	} else {
		//all matches
		log.Println("MatchesGetHandler: request to get all match info")
		m.handleAllMatches(w, r, inumber)
	}
}

func (m MatchesGetHandler) handleSpecificMatch(w http.ResponseWriter, r *http.Request, mid int, inumber string) {
	if match, err := m.MDao.GetMatchWithPredById(mid, inumber); err != nil {
		errors.ErrAnalyzePanic(w, err, "MatchesGetHandler: error getting match by id")
	} else {
		util.StructWriter(w, match)
	}
}

func (m MatchesGetHandler) handleAllMatches(w http.ResponseWriter, r *http.Request, inumber string) {
	if matches, err := m.MDao.GetAllMatchesWithPred(inumber); err != nil {
		errors.ErrAnalyzePanic(w, err, "MatchesGetHandler: unable to get matches")
	} else {
		util.StructWriter(w, matches)
	}
}

func (m MatchesGetHandler) handleMatchStats(w http.ResponseWriter, r *http.Request, mid int) {
	if stats, err := m.MDao.GetMatchStatsById(mid); err != nil {
		errors.ErrAnalyzePanic(w, err, "MatchesGetHandler: unable to get matches")
	} else {
		util.StructWriter(w, stats)
	}
}
