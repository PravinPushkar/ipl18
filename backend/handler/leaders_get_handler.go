package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

// LeadersGetHandler .
type LeadersGetHandler struct{}

const (
	qSelectLeaders = "SELECT firstname,lastname,alias,piclocation,inumber,points FROM ipluser where points is not null ORDER BY points DESC"
)

var errLeaderNotFound = fmt.Errorf("leader not found in db")

func (l LeadersGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("LeadersGetHandler:: request to get Leaders")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "LeadersGetHandler: could not get username from token")

	rows, err := db.DB.Query(qSelectLeaders)
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "LeadersGetHandler: could not fetch leaderboard")

	leaders := []models.Leader{}
	defer rows.Close()
	for rows.Next() {
		leader := models.Leader{}
		err := rows.Scan(&leader.Firstname, &leader.Lastname, &leader.Alias, &leader.Piclocation, &leader.INumber, &leader.Points)
		if err == sql.ErrNoRows {
			errors.ErrWriterPanic(w, http.StatusNotFound, err, errLeaderNotFound, "LeadersGetHandler: leader not found")
		}
		errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "LeadersGetHandler: db issue in leaders query")
		leaders = append(leaders, leader)
	}

	util.StructWriter(w, &leaders)
}
