package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/cache"
	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type PlayersGetHandler struct {
	PDao dao.PlayerDAO
}

func (p PlayersGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("PlayersGetHandler: new request")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "PlayersGetHandler: could not get username from token")

	vars := mux.Vars(r)
	if pidS, ok := vars["id"]; ok {
		pid, err := strconv.Atoi(pidS)
		errors.ErrAnalyzePanic(w, err, "PlayersGetHandler: pid is not valid")

		if player, ok := cache.PlayerIdCache[pid]; ok {
			util.StructWriter(w, player)
			return
		}

		player, err := p.PDao.GetPlayerById(pid)
		errors.ErrAnalyzePanic(w, err, "PlayersGetHandler: unable to get player by id")

		util.StructWriter(w, player)
		return
	}

	//all players
	log.Println("PlayersGetHandler: all players query")
	players, err := func() (*models.PlayersModel, error) {
		players := []*models.Player{}
		if len(cache.PlayerIdCache) != 0 {
			for _, player := range cache.PlayerIdCache {
				players = append(players, player)
			}
			return &models.PlayersModel{players}, nil
		} else {
			return p.PDao.GetAllPlayers()
		}
	}()
	errors.ErrAnalyzePanic(w, err, "TeamsGetHandler: unable to get all teams")
	util.StructWriter(w, players)
}
