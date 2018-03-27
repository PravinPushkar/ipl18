package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type TeamsGetHandler struct {
}

const (
	qSelectAllTeams        = "select tid, name, shortname, imglocation from team"
	qSelectTeam            = "select tid, name, shortname, imglocation from team where tid=$1"
	qSelectPlayerFromTeam  = "select pid, name, role, tid from player where tid=$1 and pid=$2"
	qSelectPlayersFromTeam = "select pid, name, role, tid from player where tid=$1"
)

var (
	errTeamNotFound     = fmt.Errorf("requested team not found")
	errTeamNotSpecified = fmt.Errorf("team id not specified")
	errPlayerNotFound   = fmt.Errorf("player not found")
)

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
		tid, ok := vars["id"]
		if !ok {
			errors.ErrWriterPanic(w, http.StatusBadGateway, errTeamNotSpecified, errTeamNotSpecified, "TeamsGetHandler: team not specified in request")
		}

		if pid, ok := vars["pid"]; ok {
			//specific player
			player := models.Player{}
			err = db.DB.QueryRow(qSelectPlayerFromTeam, tid, pid).Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
			if err == sql.ErrNoRows {
				errors.ErrWriterPanic(w, http.StatusNotFound, err, errPlayerNotFound, fmt.Sprintf("TeamsGetHandler: player %s not found for team %s", pid, tid))
			}
			errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, fmt.Sprintf("TeamsGetHandler: player %s not found for team %s", pid, tid))
			util.StructWriter(w, &player)
			return
		} else {
			//all players
			rows, err := db.DB.Query(qSelectPlayersFromTeam, tid)
			errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, fmt.Sprintf("TeamsGetHandler: could not query players from team %s", tid))

			defer rows.Close()
			players := []models.Player{}
			for rows.Next() {
				player := models.Player{}
				err = rows.Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
				errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, fmt.Sprintf("TeamsGetHandler: error reading row for %s", tid))
				players = append(players, player)
			}
			util.StructWriter(w, &players)
			return
		}
	}

	//normal team queries
	if tid, ok := vars["id"]; ok {
		log.Println("TeamsGetHandler: request to get team", tid)
		row := db.DB.QueryRow(qSelectTeam, tid)
		t.writeTeam(w, row, tid)
		return
	}

	if len(vars) == 0 {
		log.Println("TeamsGetHandler: request to get all teams")
		rows, err := db.DB.Query(qSelectAllTeams)
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "TeamsGetHandler: could not query db")
		t.writeTeams(w, rows)
		return
	}

	errors.ErrWriter(w, http.StatusBadRequest, "team request is not valid")
}

func (t TeamsGetHandler) writeTeams(w http.ResponseWriter, rows *sql.Rows) {
	teams := []models.Team{}

	defer rows.Close()

	pic := sql.NullString{}
	for rows.Next() {
		team := models.Team{}
		err := rows.Scan(&team.TeamId, &team.TeamName, &team.ShortName, &pic)
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "TeamsGetHandler: could not read from db")
		team.PicLocation = pic.String
		teams = append(teams, team)
	}

	util.StructWriter(w, &models.Teams{teams})
}

func (t TeamsGetHandler) writeTeam(w http.ResponseWriter, row *sql.Row, tid string) {
	team := models.Team{}
	pic := sql.NullString{}
	err := row.Scan(&team.TeamId, &team.TeamName, &team.ShortName, &pic)
	if err == sql.ErrNoRows {
		errors.ErrWriterPanic(w, http.StatusNotFound, err, errTeamNotFound, "TeamsGetHandler: could not find team "+tid)
	}
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "TeamsGetHandler: could not find team "+tid)

	team.PicLocation = pic.String
	util.StructWriter(w, &team)
}
