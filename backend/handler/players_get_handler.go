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

type PlayersGetHandler struct {
}

const (
	qSelectPlayers = "select pid, name, role, tid from player"
	qSelectPlayer  = "select pid, name, role, tid from player where pid=$1"
)

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
	if pid, ok := vars["id"]; ok {
		log.Println("PlayersGetHandler: single player query", pid)
		//specific player details
		player := models.Player{}
		err := db.DB.QueryRow(qSelectPlayer, pid).Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
		if err == sql.ErrNoRows {
			errors.ErrWriterPanic(w, http.StatusNotFound, err, errPlayerNotFound, fmt.Sprintf("PlayersGetHandler: player %s not found in db", pid))
		}
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "PlayersGetHandler: could not query db ")
		util.StructWriter(w, &player)
		return
	}

	//all players
	log.Println("PlayersGetHandler: all players query")
	rows, err := db.DB.Query(qSelectPlayers)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "PlayersGetHandler: could not query players")

	players := []models.Player{}
	for rows.Next() {
		player := models.Player{}
		err := rows.Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
		if err == sql.ErrNoRows {
			errors.ErrWriterPanic(w, http.StatusNotFound, err, errPlayerNotFound, "PlayersGetHandler: player not found in db")
		}
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "PlayersGetHandler: could not query db ")
		players = append(players, player)
	}
	util.StructWriter(w, &models.PlayersModel{players})
}
