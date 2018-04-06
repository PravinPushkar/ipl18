package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/auth"
	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

var (
	tokenManager = auth.NewTokenManager(auth.SignMethodSHA512)
)

var PingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	util.StructWriter(w, models.PingModel{true})
})

var RegistrationHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("RegistrationHandler: new user registration request ")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked")
		}
	}()

	user := models.User{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrParseRequest, "RegistrationHandler: could not parse user information")

	_, err = db.DB.Exec("insert into ipluser(firstname, lastname, password, coin, alias, inumber) values($1, $2, $3, $4, $5, $6)",
		user.Firstname,
		user.Lastname,
		util.GetHash([]byte(user.Password)),
		12,
		user.Alias,
		user.INumber)

	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "RegistrationHandler: could not register new user")
	util.OkWriter(w)
})

var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler: new user login request ")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	defer r.Body.Close()
	lm := models.LoginModel{}

	err := json.NewDecoder(r.Body).Decode(&lm)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrParseRequest, "LoginHandler: could not parse request body")

	inumber := ""
	hashPass := util.GetHash([]byte(lm.Password))
	log.Println("got details:", lm.INumber, hashPass)

	err = db.DB.QueryRow("select inumber from ipluser where inumber=$1 and password=$2", lm.INumber, hashPass).Scan(&inumber)
	if err == sql.ErrNoRows {
		errors.ErrWriterPanic(w, http.StatusForbidden, err, errUserNotFound, "UserGetHandler: user not found in db")
	}
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "UserGetHandler: could not query db")

	token, err := tokenManager.GetToken(lm.INumber, time.Duration(1))
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrGettingToken, "LoginHandler: could not get new token")

	util.StructWriter(w, token)
})
