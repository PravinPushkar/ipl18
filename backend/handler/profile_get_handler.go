package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type UserGetHandler struct {
}

const (
	qFetchUserDetails = "select firstname, lastname, coin, alias, piclocation, inumber from ipluser where inumber=$1"
)

var (
	errUserNotFound = fmt.Errorf("user not found")
)

func (p UserGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("UserGetHandler: request to get profile")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked")
		}
	}()

	inumber, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "UserGetHandler: could not get username from token")

	pathVar := mux.Vars(r)
	if pathVar["inumber"] != inumber {
		errors.ErrWriterPanic(w, http.StatusForbidden, errors.ErrTokenInfoMismatch, errors.ErrTokenInfoMismatch, fmt.Sprintf("UserPutHandler: token info and path var mismatch %s-%s", pathVar["inumber"], inumber))
	}

	var pic sql.NullString
	info := models.ProfileViewModel{}

	if err := db.DB.QueryRow(qFetchUserDetails, inumber).Scan(&info.Firstname, &info.Lastname, &info.Coin, &info.Alias, &pic, &info.INumber); err == sql.ErrNoRows {
		errors.ErrWriterPanic(w, http.StatusForbidden, err, errUserNotFound, "UserGetHandler: user not found in db")
	}
	errors.ErrWriterPanic(w, http.StatusInternalServerError, err, errors.ErrDBIssue, "UserGetHandler: could not query db")

	info.PicLocation = pic.String
	util.StructWriter(w, &info)
}
