package dao

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/cache"
	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type UserDAO struct{}

const (
	qUpdatePointsIplUser = "update ipluser set points=points+$1 where inumber=$2"
	qSelectAllUsers      = "select firstname,lastname,alias,piclocation,inumber,points from ipluser"
	qSelectLeaders       = "SELECT firstname,lastname,alias,piclocation,inumber,points FROM ipluser where points is not null ORDER BY points DESC"
	qInsertUser          = "insert into ipluser(firstname, lastname, password, coin, alias, inumber) values($1, $2, $3, $4, $5, $6)"
	qVerifyUser          = "select inumber from ipluser where inumber=$1 and password=$2"
)

var errLeaderNotFound = fmt.Errorf("leader not found in db")

func (u UserDAO) UpdateUserPointsByINumber(by int, inumber string) error {
	if res, err := db.DB.Exec(qUpdatePointsIplUser, by, inumber); err != nil {
		log.Println("UserDAO: UpdateUserPointsByINumber error updating points", err)
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	} else {
		if i, err := res.RowsAffected(); err != nil || i != 1 {
			log.Println("UserDAO: UpdateUserPointsByINumber affected rows don't add up", err, i)
			return &errors.DaoError{http.StatusInternalServerError, errors.ErrDBIssue, errors.ErrDBIssue}
		}
	}
	defer cache.Lock.Unlock()
	cache.Lock.Lock()
	if _, ok := cache.UserINumberCache[inumber]; ok {
		cache.UserINumberCache[inumber].Points += by
	}
	return nil
}

func (u UserDAO) GetLeaders() (*models.LeadersModel, error) {
	rows, err := db.DB.Query(qSelectLeaders)
	if err != nil {
		log.Println("UserDAO: GetLeaders error getting leaders", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	leaders := []*models.Leader{}
	defer rows.Close()
	for rows.Next() {
		leader := models.Leader{}
		err := rows.Scan(&leader.Firstname, &leader.Lastname, &leader.Alias, &leader.Piclocation, &leader.INumber, &leader.Points)
		if err == sql.ErrNoRows {
			log.Println("UserDAO: GetLeaders could not find leaders", err)
			return nil, &errors.DaoError{http.StatusNotFound, err, errLeaderNotFound}
		} else if err != nil {
			log.Println("UserDAO: GetLeaders error scanning leaders", err)
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		leaders = append(leaders, &leader)
	}

	return &models.LeadersModel{leaders}, nil
}

func (u UserDAO) InsertUser(user *models.User) error {
	res, err := db.DB.Exec(qInsertUser,
		user.Firstname,
		user.Lastname,
		util.GetHash([]byte(user.Password)),
		12,
		user.Alias,
		user.INumber)

	if err != nil {
		log.Println("UserDAO: InsertUser error inserting new user", err)
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	if num, err := res.RowsAffected(); err != nil || num != 1 {
		log.Println("UserDAO: db not updated inserting new user", err, num)
		return &errors.DaoError{http.StatusInternalServerError, errors.ErrDBIssue, errors.ErrDBIssue}
	}
	defer cache.Lock.Unlock()
	cache.Lock.Lock()
	cache.UserINumberCache[user.INumber] = user

	return nil
}

func (u UserDAO) GetAllUsersBasicInfo() ([]*models.UserBasic, error) {
	log.Println("UserDAO: GetAllUsersBasicInfo")

	if len(cache.UserINumberCache) != 0 {
		log.Println("cache hit")
		users := []*models.UserBasic{}
		for _, v := range cache.UserINumberCache {
			users = append(users, &models.UserBasic{
				v.Firstname + " " + v.Lastname,
				v.INumber,
				v.PicLocation,
			})
		}
		return users, nil
	}

	res, err := db.DB.Query(qSelectAllUsers)
	if err != nil {
		log.Println("UserDAO:GetAllUsersBasicInfo error getting users", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}

	defer res.Close()
	usersB := []*models.UserBasic{}
	users := []*models.User{}
	for res.Next() {
		user := models.User{}
		if err := res.Scan(&user.Firstname, &user.Lastname, &user.Alias, &user.PicLocation, &user.INumber, &user.Points); err != nil {
			log.Println("UserDAO:GetAllUsersBasicInfo error scanning user", err)
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		users = append(users, &user)
		usersB = append(usersB, &models.UserBasic{
			user.Firstname + " " + user.Lastname,
			user.INumber,
			user.PicLocation,
		})
	}

	defer cache.Lock.Unlock()
	cache.Lock.Lock()
	if len(cache.UserINumberCache) == 0 {
		for _, v := range users {
			cache.UserINumberCache[v.INumber] = v
		}
	}

	return usersB, nil
}

func (u UserDAO) VerifyUser(inumber string, pass string) error {
	inum := ""
	err := db.DB.QueryRow(qVerifyUser, inumber, pass).Scan(&inum)
	if err == sql.ErrNoRows {
		log.Println("UserDAO: VerifyUser user not found", err)
		return &errors.DaoError{http.StatusForbidden, errors.ErrUserNotFound, errors.ErrUserNotFound}
	} else if err != nil {
		log.Println("UserDAO: VerifyUser could not verify user", err)
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	return nil
}
