package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/auth"
	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

var tokenManager = auth.NewTokenManager(auth.SignMethodSHA512)

var PingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
})

var RegistrationHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&user)

	if _, err := db.DB.Exec("insert into ipluser(firstname, lastname, password, coin, alias, piclocation, inum) values($1, $2, $3, $4, $5, '', $6)",
		user.Firstname,
		user.Lastname,
		user.Password,
		12,
		user.Alias,
		user.INumber); err != nil {
		log.Println("could not insert in db", err.Error())
		util.ErrWriter(w, http.StatusInternalServerError, "could not register new user")
		return
	}
	util.OkWriter(w)
})

var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lm := models.LoginModel{}
	if err := json.NewDecoder(r.Body).Decode(&lm); err != nil {
		log.Println(err.Error())
		util.ErrWriter(w, http.StatusBadRequest, err)
		return
	}

	//check correct password here
	token, err := tokenManager.GetToken(lm.INumber, time.Duration(1))

	if err != nil {
		util.ErrWriter(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	util.StructWriter(w, token)
})
