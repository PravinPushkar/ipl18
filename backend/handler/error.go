package handler

import (
	"errors"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

var NotImplementedErr = errors.New("Not Implemented")

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	util.ErrWriter(w, http.StatusNotFound, NotImplementedErr)
})
