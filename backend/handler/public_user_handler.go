package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type PublicUserHandler struct {
	UDao dao.UserDAO
}

func (p PublicUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("PublicUserHandler: new request")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	if strings.Contains(r.URL.Path, "/register") {
		p.handleRegistration(w, r)
	} else if strings.Contains(r.URL.Path, "/login") {
		p.handleLogin(w, r)
	}
}

func (p PublicUserHandler) handleRegistration(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrParseRequest, "RegistrationHandler: could not parse user information")

	errors.ErrAnalyzePanic(w, p.UDao.InsertUser(&user), "RegistrationHandler: could not register new user")
	util.OkWriter(w)
}

func (p PublicUserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lm := models.LoginModel{}

	err := json.NewDecoder(r.Body).Decode(&lm)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrParseRequest, "LoginHandler: could not parse request body")

	hashPass := util.GetHash([]byte(lm.Password))
	log.Println("got details:", lm.INumber, hashPass)

	err = p.UDao.VerifyUser(lm.INumber, hashPass)
	errors.ErrAnalyzePanic(w, err, "LoginHandler: user not valid")

	token, err := tokenManager.GetToken(lm.INumber, time.Duration(1))
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrGettingToken, "LoginHandler: could not get new token")

	util.StructWriter(w, token)
}
