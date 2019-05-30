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
	qSelectTeam     = "select tid, name, shortname, imglocation from team where tid=$1"
	qSelectAllTeams = "select tid, name, shortname, imglocation from team"
)

var (
	errTeamNotFound     = fmt.Errorf("requested team not found")
	errTeamNotSpecified = fmt.Errorf("team id not specified")
)

type TeamDAO struct{}

func (t TeamDAO) GetAllTeams() (*models.Teams, error) {
	log.Println("TeamDAO: GetAllTeams ")
	//check cache or fill
	if cache.Teams != nil {
		log.Println("cache hit")
		return cache.Teams, nil
	}

	rows, err := db.DB.Query(qSelectAllTeams)
	if err != nil {
		log.Println("TeamDAO: error querying teams", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	defer rows.Close()

	teams := []*models.Team{}
	pic := sql.NullString{}
	for rows.Next() {
		team := models.Team{}
		if err := rows.Scan(&team.TeamId, &team.TeamName, &team.ShortName, &pic); err != nil {
			log.Println("TeamDAO: could not scan teams", err)
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		team.PicLocation = pic.String
		teams = append(teams, &team)
	}

	data := &models.Teams{teams}

	defer cache.Lock.Unlock()
	cache.Lock.Lock()
	if cache.Teams == nil {
		cache.Teams = data
		for _, v := range data.Teams {
			cache.TeamIdCache[v.TeamId] = v
			cache.TeamSNameCache[v.ShortName] = v
			cache.TeamNameCache[v.TeamName] = v
		}
	}
	return data, nil
}

func (t TeamDAO) GetTeamById(tid int) (*models.Team, error) {
	log.Println("TeamDAO: GetTeamsById")
	if team, ok := cache.TeamIdCache[tid]; ok {
		log.Println("cache hit")
		return team, nil
	}

	team := models.Team{}
	pic := sql.NullString{}

	err := db.DB.QueryRow(qSelectTeam, tid).Scan(&team.TeamId, &team.TeamName, &team.ShortName, &pic)
	if err == sql.ErrNoRows {
		log.Println("TeamDAO: GetTeamById could not find team ", tid)
		return nil, &errors.DaoError{http.StatusNotFound, err, errTeamNotFound}
	} else if err != nil {
		log.Println("TeamDAO: GetTeamById could not query team", tid)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	team.PicLocation = pic.String
	data := &team

	defer cache.Lock.Unlock()
	cache.Lock.Lock()
	if team, ok := cache.TeamIdCache[tid]; !ok {
		cache.TeamIdCache[tid] = data
		cache.TeamSNameCache[team.ShortName] = data
		cache.TeamNameCache[team.TeamName] = data
	}

	return data, nil
}
