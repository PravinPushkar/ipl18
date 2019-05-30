package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/mux"

	"github.com/PravinPushkar/ipl18/backend/cache"
	"github.com/PravinPushkar/ipl18/backend/db"
	"github.com/PravinPushkar/ipl18/backend/errors"
	"github.com/PravinPushkar/ipl18/backend/util"
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
	errAliasInvalid = fmt.Errorf("alias cannot be more than 10 chars and should be alphanumeric")
	errPassInvalid  = fmt.Errorf("password should be 8-20 chars and alphanumeric with allowed special chars [-_.$%^&*!@]")
	errImageInvaid  = fmt.Errorf("image can be only png/jpg/jpeg and size 4 MB max")

	rePass, _  = regexp.Compile("^[-_.$%^&*!@0-9a-zA-Z]{8,20}$")
	reAlias, _ = regexp.Compile("^[0-9a-zA-Z]{1,10}$")
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
	errors.ErrWriterPanic(w, http.StatusBadRequest, err, err, "UserPutHandler: error parsing form data")
	util.OkWriter(w)
}

func (p UserPutHandler) parseAndUpdate(r *http.Request, inumber string) error {
	log.Println("UserPutHandler: parsing request")
	err := r.ParseMultipartForm(maxMemory)
	updates := [4]bool{}
	query := "update ipluser set inumber=$1"
	values := []interface{}{inumber}
	i := 2
	if err != nil {
		log.Println("ProfileDAO: ", err)
	} else {
		if location, err := p.handleImage(r, inumber); err != nil {
			return err
		} else if location != "" {
			query += fmt.Sprintf(",piclocation=$%d", i)
			i++
			values = append(values, location)
			updates[0] = true
		}
	}

	for k, _ := range r.Form {
		val := r.Form.Get(k)
		log.Println("UserPutHandler: found field ", k)
		switch k {
		case "alias":
			if !reAlias.MatchString(val) {
				return errAliasInvalid
			}
			query += fmt.Sprintf(",alias=$%d", i)
			values = append(values, val)
			updates[1] = true
			i++

		case "password":
			if !rePass.MatchString(val) {
				return errPassInvalid
			}
			query += fmt.Sprintf(",password=$%d", i)
			values = append(values, util.GetHash([]byte(val)))
			updates[2] = true
			i++

		default:
			return errInvalidField
		}
	}

	if len(values) > 1 {
		query += fmt.Sprintf(" where inumber=$%d", i)
		values = append(values, inumber)
		log.Println(query, values)
		_, err := db.DB.Exec(query, values...)
		if err == nil {
			//update cache
			log.Println("ProfilePutHandler: updating cache")
			defer cache.Lock.Unlock()
			cache.Lock.Lock()
			if user, ok := cache.UserINumberCache[inumber]; ok {
				index := 0
				if updates[0] {
					user.PicLocation = values[index].(string)
					index += 1
				}
				if updates[1] {
					user.Alias = values[index].(string)
					index += 1
				}
			}
		}
		return err
	}

	return nil
}

func (p UserPutHandler) handleImage(r *http.Request, inumber string) (string, error) {
	file, handle, err := r.FormFile("image")
	if err != nil {
		log.Println("UserPutHandler: error getting file handle", err.Error())
		return "", nil
	}

	defer file.Close()
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", err
	}
	imgType := http.DetectContentType(buf)
	if !(imgType == "image/png" || imgType == "image/jpg" || imgType == "image/jpeg") {
		return "", errImageInvaid
	}
	if r.ContentLength > 1024*1024*5 {
		return "", errImageInvaid
	}

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
