package dao

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

const (
	qSelectPlayers         = "select pid, name, role, tid from player"
	qSelectPlayer          = "select pid, name, role, tid from player where pid=$1"
	qSelectPlayersFromTeam = "select pid, name, role, tid from player where tid=$1"
)

var errPlayerNotFound = fmt.Errorf("player not found")

type PlayerDAO struct{}

func (p PlayerDAO) GetAllPlayers() (*models.PlayersModel, error) {
	log.Println("PlayerDAO: GetAllPlayers")
	rows, err := db.DB.Query(qSelectPlayers)
	if err != nil {
		log.Println("PlayerDAO: ", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	players := []*models.Player{}
	for rows.Next() {
		player := models.Player{}
		err := rows.Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
		if err == sql.ErrNoRows {
			return nil, &errors.DaoError{http.StatusNotFound, err, err}
		} else if err != nil {
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		players = append(players, &player)
	}
	return &models.PlayersModel{players}, nil
}

func (p PlayerDAO) GetPlayerById(pid int) (*models.Player, error) {
	log.Println("PlayerDAO: GetPlayerById ", pid)

	player := models.Player{}
	err := db.DB.QueryRow(qSelectPlayer, pid).Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
	if err == sql.ErrNoRows {
		log.Println("PlayerDAO: GetPlayerById player by id not found", err)
		return nil, &errors.DaoError{http.StatusNotFound, errPlayerNotFound, errPlayerNotFound}
	} else if err != nil {
		log.Println("PlayerDAO: GetPlayerById error getting player by id", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	return &player, nil
}

func (p PlayerDAO) GetAllPlayersByTeam(tid int) (*models.PlayersModel, error) {
	log.Println("PlayerDAO: GetAllPlayersByTeam ", tid)

	rows, err := db.DB.Query(qSelectPlayersFromTeam, tid)
	if err != nil {
		log.Println("PlayerDAO: GetAllPlayersByTeam error getting team players", tid)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	defer rows.Close()
	players := []*models.Player{}
	for rows.Next() {
		player := models.Player{}
		err = rows.Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
		if err != nil {
			log.Println("PlayerDAO: GetAllPlayersByTeam error scanning team players", tid)
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}

		players = append(players, &player)
	}

	return &models.PlayersModel{players}, nil
}
