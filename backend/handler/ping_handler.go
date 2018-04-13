package handler

import (
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

var PingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	util.StructWriter(w, models.PingModel{true})
})
