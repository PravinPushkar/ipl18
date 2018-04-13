package handler

import (
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/dao"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/util"
)

// LeadersGetHandler .
type LeadersGetHandler struct {
	UDao dao.UserDAO
}

func (l LeadersGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("LeadersGetHandler:: request to get Leaders")
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "LeadersGetHandler: could not get username from token")

	leaders, err := l.UDao.GetLeaders()
	errors.ErrAnalyzePanic(w, err, "LeadersGetHandler: unable to fetch leaders")

	util.StructWriter(w, &leaders)
}
