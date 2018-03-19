package handler

import (
	"log"
	"net/http"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/util"

	jwt "github.com/dgrijalva/jwt-go"
)

func sendErrResponse(w http.ResponseWriter, code int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(util.GetJsonErrMessage(code, msg))
}

var PingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
})

var RegistrationHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		sendErrResponse(w, http.StatusBadRequest, err)
		return
	}

	//todo db stuff here
})

var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		sendErrResponse(w, http.StatusBadRequest, err)
		return
	}

	//do db auth here, return on error
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = "Sushant Mahajan"
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix()

	tokenString, err := token.SignedString([]byte("secretsecret"))
	if err != nil {
		log.Println(err.Error())
		sendErrResponse(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	w.Write([]byte(tokenString))
})
