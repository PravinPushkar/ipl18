package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

type PredictPutHandler struct{}

func (p PredictPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("MatchesGetHandler: new request")

	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "MatchesGetHandler: could not get username from token")

	vars := mux.Vars(r)
	if _, ok := vars["id"]; ok {
		//parse body
		defer r.Body.Close()
		json.NewDecoder(r.Body).Decode(nil)
	}
}
