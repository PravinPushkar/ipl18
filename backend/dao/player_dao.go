package dao

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/PravinPushkar/ipl18/backend/cache"
	"github.com/PravinPushkar/ipl18/backend/db"
	"github.com/PravinPushkar/ipl18/backend/errors"
	"github.com/PravinPushkar/ipl18/backend/models"
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

	if cache.Players != nil {
		log.Println("cache hit")
		return cache.Players, nil
	}

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
	data := &models.PlayersModel{players}

	defer cache.Lock.Unlock()
	cache.Lock.Lock()
	if cache.Players == nil {
		cache.Players = data
		for _, pl := range players {
			v := pl
			cache.PlayerIdCache[v.PlayerId] = v
			cache.PlayerNameCache[v.Name] = v
			if _, ok := cache.TeamPlayerCache[v.TeamId]; !ok {
				cache.TeamPlayerCache[v.TeamId] = make(map[int]*models.Player)
			}
			cache.TeamPlayerCache[v.TeamId][v.PlayerId] = v
		}
	}

	return data, nil
}

func (p PlayerDAO) GetPlayerById(pid int) (*models.Player, error) {
	log.Println("PlayerDAO: GetPlayerById ", pid)
	if player, ok := cache.PlayerIdCache[pid]; ok {
		log.Println("cache hit")
		return player, nil
	}

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
	if players, ok := cache.TeamPlayerCache[tid]; ok {
		log.Println("cache hit")
		allPlayers := []*models.Player{}
		for _, pl := range players {
			v := pl
			allPlayers = append(allPlayers, v)
		}
		return &models.PlayersModel{allPlayers}, nil
	}

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
