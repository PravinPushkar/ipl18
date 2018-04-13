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
	teams := []models.Team{}

	rows, err := db.DB.Query(qSelectAllTeams)
	if err != nil {
		log.Println("TeamDAO: error querying teams", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	defer rows.Close()

	pic := sql.NullString{}
	for rows.Next() {
		team := models.Team{}
		if err := rows.Scan(&team.TeamId, &team.TeamName, &team.ShortName, &pic); err != nil {
			log.Println("TeamDAO: could not scan teams", err)
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		team.PicLocation = pic.String
		teams = append(teams, team)
	}

	return &models.Teams{teams}, nil
}

func (t TeamDAO) GetTeamById(tid int) (*models.Team, error) {
	log.Println("TeamDAO: GetTeamsById")
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

	return &team, nil
}
