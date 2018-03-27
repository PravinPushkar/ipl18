package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type UserPutHandler struct {
}

const (
	maxMemory = 1024 * 1024 * 2
)

var (
	errINumberDiff  = fmt.Errorf("token info mismatch")
	errSaveImage    = fmt.Errorf("could not save image")
	errInvalidField = fmt.Errorf("invalid key in form")
)

func (p UserPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("UserPutHandler: new update user profile request")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked")
		}
	}()

	inumber, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "UserPutHandler: could not get username from token")

	pathVar := mux.Vars(r)
	if pathVar["inumber"] != inumber {
		errors.ErrWriterPanic(w, http.StatusForbidden, errINumberDiff, errors.ErrTokenInfoMismatch, fmt.Sprintf("UserPutHandler: token info and path var mismatch %s-%s", pathVar["inumber"], inumber))
	}
	err = p.parseAndUpdate(r, inumber)
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, errors.ErrParseRequest, "UserPutHandler: error parsing form data")
	util.OkWriter(w)
}

func (p UserPutHandler) parseAndUpdate(r *http.Request, inumber string) error {
	log.Println("UserPutHandler: parsing request")
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return err
	}

	query := "update ipluser set inumber=$1"
	values := []interface{}{inumber}

	i := 2
	if location, err := p.handleImage(r, inumber); err != nil {
		return err
	} else {
		query += fmt.Sprintf(",piclocation=$%d", i)
		values = append(values, location)
	}

	for k, _ := range r.Form {
		val := r.Form.Get(k)
		log.Println("UserPutHandler: found field ", k)
		i++
		switch k {

		case "alias":
			query += fmt.Sprintf(",alias=$%d", i)
			values = append(values, val)

		case "password":
			query += fmt.Sprintf(",password=$%d", i)
			values = append(values, util.GetHash([]byte(val)))

		default:
			return errInvalidField
		}
	}

	if len(values) > 1 {
		i++
		query += fmt.Sprintf(" where inumber=$%d", i)
		values = append(values, inumber)
		log.Println(query, values)
		_, err := db.DB.Exec(query, values...)
		return err
	}
	return nil
}

func (p UserPutHandler) handleImage(r *http.Request, inumber string) (string, error) {
	file, handle, err := r.FormFile("image")
	if err != nil {
		log.Println("UserPutHandler: error getting file handle", err.Error())
		return "", err
	}

	defer file.Close()
	piclocation := fmt.Sprintf("./static/assets/img/users/%s_%d_%s", inumber, time.Now().Unix(), handle.Filename)
	if f, err := os.OpenFile(piclocation, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		log.Println("UserPutHandler: error opening new file for writing", err.Error())
		return "", err
	} else {
		defer f.Close()
		_, err := io.Copy(f, file)
		if err != nil {
			log.Println("UserPutHandler: could not create file in fs", err.Error())
			return "", err
		}
		return piclocation, nil
	}
}
