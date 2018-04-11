package dao

import (
	"database/sql"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

const (
	qSelectPlayers = "select pid, name, role, tid from player"
)

type PlayerDAO struct{}

func (p PlayerDAO) GetAllPlayers() (*models.PlayersModel, *models.DaoError) {
	log.Println("PlayerDAO: all players query")
	rows, err := db.DB.Query(qSelectPlayers)
	if err != nil {
		log.Println("PlayerDAO: ", err)
		return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	players := []models.Player{}
	for rows.Next() {
		player := models.Player{}
		err := rows.Scan(&player.PlayerId, &player.Name, &player.Role, &player.TeamId)
		if err == sql.ErrNoRows {
			return nil, &models.DaoError{http.StatusNotFound, err, err}
		} else if err != nil {
			return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		players = append(players, player)
	}
	return &models.PlayersModel{players}, nil
}
