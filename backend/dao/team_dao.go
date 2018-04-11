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
	qSelectAllTeams = "select tid, name, shortname, imglocation from team"
)

type TeamDAO struct{}

func (t TeamDAO) GetAllTeams() (*models.Teams, *models.DaoError) {
	teams := []models.Team{}

	rows, err := db.DB.Query(qSelectAllTeams)
	if err != nil {
		log.Println("TeamDAO: error querying teams", err)
		return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	defer rows.Close()

	pic := sql.NullString{}
	for rows.Next() {
		team := models.Team{}
		if err := rows.Scan(&team.TeamId, &team.TeamName, &team.ShortName, &pic); err != nil {
			log.Println("TeamDAO: could not scan teams", err)
			return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		team.PicLocation = pic.String
		teams = append(teams, team)
	}

	return &models.Teams{teams}, nil
}
