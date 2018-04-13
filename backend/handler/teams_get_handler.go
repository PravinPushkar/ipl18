package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/cache"
	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type TeamsGetHandler struct {
	PDao dao.PlayerDAO
	TDao dao.TeamDAO
}

var errTeamNotSpecified = fmt.Errorf("team id not specified")

func (t TeamsGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("TeamsGetHandler: new request")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked")
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "TeamsGetHandler: could not get username from token")

	vars := mux.Vars(r)
	//player specific query
	if strings.Contains(r.URL.Path, "/players") {
		tidS, ok := vars["id"]
		if !ok {
			errors.ErrWriterPanic(w, http.StatusBadGateway, errTeamNotSpecified, errTeamNotSpecified, "TeamsGetHandler: team not specified in request")
		}
		tid, err := strconv.Atoi(tidS)
		errors.ErrAnalyzePanic(w, err, "TeamsGetHandler: team id not valid")

		if pidS, ok := vars["pid"]; ok {
			//specific player
			pid, err := strconv.Atoi(pidS)
			errors.ErrAnalyzePanic(w, err, "TeamsGetHandler: player id not valid")

			if player, ok := cache.PlayerIdCache[pid]; ok {
				util.StructWriter(w, player)
				return
			}
			player, err := t.PDao.GetPlayerById(pid)
			errors.ErrAnalyzePanic(w, err, "TeamsGetHandler: error getting player")
			util.StructWriter(w, player)
			return
		} else {
			//all players of team
			players, err := func() (*models.PlayersModel, error) {
				players := []*models.Player{}
				if playerMap, ok := cache.TeamPlayerCache[tid]; ok {
					log.Println("TeamsGetHandler: Found in cache")
					for _, player := range playerMap {
						players = append(players, player)
					}
					return &models.PlayersModel{players}, nil
				} else {
					return t.PDao.GetAllPlayersByTeam(tid)
				}
			}()

			errors.ErrAnalyzePanic(w, err, "TeamsGetHandler: error getting team players")
			util.StructWriter(w, players)
			return
		}
	}

	//normal team queries
	if tidS, ok := vars["id"]; ok {
		log.Println("TeamsGetHandler: request to get team", tidS)
		tid, err := strconv.Atoi(tidS)
		errors.ErrAnalyzePanic(w, err, "TeamsGetHandler: team id not valid")

		if team, ok := cache.TeamIdCache[tid]; ok {
			log.Println("TeamsGetHandler: Found in cache")
			util.StructWriter(w, team)
			return
		}
		team, err := t.TDao.GetTeamById(tid)
		errors.ErrAnalyzePanic(w, err, "TeamDAO: error getting team by id")
		util.StructWriter(w, team)
		return
	}

	if len(vars) == 0 {
		log.Println("TeamsGetHandler: request to get all teams")
		teams, err := func() (*models.Teams, error) {
			teams := []*models.Team{}
			if len(cache.TeamIdCache) != 0 {
				for _, team := range cache.TeamIdCache {
					teams = append(teams, team)
				}
				return &models.Teams{teams}, nil
			} else {
				return t.TDao.GetAllTeams()
			}
		}()
		errors.ErrAnalyzePanic(w, err, "TeamsGetHandler: unable to get all teams")
		util.StructWriter(w, teams)
		return
	}

	errors.ErrWriter(w, http.StatusBadRequest, "team request is not valid")
}
