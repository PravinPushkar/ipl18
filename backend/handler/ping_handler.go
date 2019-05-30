package handler

import (
	"net/http"

	"github.com/PravinPushkar/ipl18/backend/models"
	"github.com/PravinPushkar/ipl18/backend/util"
)

var PingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	util.StructWriter(w, models.PingModel{true})
})
