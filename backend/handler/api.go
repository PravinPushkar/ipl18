package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type ProfileHandler struct {
}

const (
	qFetchUserDetails = "select firstname, lastname, password, coin, alias, piclocation, inumber from ipluser where inumber=$1"
)

func (p ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	claims, err := tokenManager.GetClaims(token)
	if err != nil {
		log.Println("ProfileHandler: error parsing token ", err.Error())
		util.ErrWriter(w, http.StatusForbidden, "could not parse token")
		return
	}

	inumber, ok := claims["inumber"].(string)
	if !ok {
		log.Println("ProfileHandler: claims not correct")
		util.ErrWriter(w, http.StatusForbidden, "token claims not valid")
		return
	}

	log.Printf("ProfileHandler: got claim %s", inumber)

	var fname, lname, pass, alias, piclocation, inum string
	var coin int

	if err := db.DB.QueryRow(qFetchUserDetails, inumber).Scan(&fname, &lname, &pass, &coin, &alias, &piclocation, &inum); err == sql.ErrNoRows {
		log.Println("ProfileHandler: user not found in db", err.Error())
		util.ErrWriter(w, http.StatusForbidden, "user not found")
		return
	} else if err != nil {
		log.Println("ProfileHandler: user not found in db", err.Error())
		util.ErrWriter(w, http.StatusForbidden, "could not query user in db")
		return
	}

	b := new(bytes.Buffer)
	profModel := models.ProfileViewModel{
		fname, lname, coin, alias, piclocation, inum,
	}
	if err := json.NewEncoder(b).Encode(&profModel); err != nil {
		log.Println("ProfileHandler: error encoding user data ", err.Error())
		util.ErrWriter(w, http.StatusInternalServerError, "error sending data")
		return
	}
	w.Write(b.Bytes())
}
