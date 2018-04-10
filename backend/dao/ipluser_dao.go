package dao

import (
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

type UserDAO struct{}

const (
	qUpdatePointsIplUser = "update ipluser set points=points+$1 where inumber=$2"
	qSelectAllUsers      = "select concat(firstname,' ',lastname) as name,inumber from ipluser"
)

func (u UserDAO) UpdateUserPointsByINumber(by int, inumber string) *models.DaoError {
	if res, err := db.DB.Exec(qUpdatePointsIplUser, by, inumber); err != nil {
		log.Println("UserDAO: UpdateUserPointsByINumber error updating points", err)
		return &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	} else {
		if i, err := res.RowsAffected(); err != nil || i != 1 {
			log.Println("UserDAO: UpdateUserPointsByINumber affected rows don't add up", err, i)
			return &models.DaoError{http.StatusInternalServerError, errors.ErrDBIssue, errors.ErrDBIssue}
		}
	}
	return nil
}

func (u UserDAO) GetAllUsersBasicInfo() ([]*models.UserBasic, *models.DaoError) {
	log.Println("UserDAO: GetAllUsersBasicInfo")
	res, err := db.DB.Query(qSelectAllUsers)
	if err != nil {
		log.Println("UserDAO:GetAllUsersBasicInfo error getting users", err)
		return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	defer res.Close()
	users := []*models.UserBasic{}
	for res.Next() {
		user := models.UserBasic{}
		if err := res.Scan(&user.Name, &user.INumber); err != nil {
			log.Println("UserDAO:GetAllUsersBasicInfo error scanning user", err)
			return nil, &models.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		users = append(users, &user)
	}
	return users, nil
}
