package handler

import (
	"log"
	"net/http"

	"github.com/PravinPushkar/ipl18/backend/dao"
	"github.com/PravinPushkar/ipl18/backend/errors"
	"github.com/PravinPushkar/ipl18/backend/util"
)

// BonusQuestionGetHandler ..
type BonusQuestionGetHandler struct {
	BDao dao.BonusDAO
}

func (q BonusQuestionGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("BonusQuestionGetHandler: new request")

	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	_, err := util.GetValueFromContext(r, "inumber")
	errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "BonusQuestionGetHandler: could not get username from token")

	questions, err := q.BDao.GetAllBonusQuestions()
	errors.ErrAnalyzePanic(w, err, "BonusQuestionGetHandler: error getting questions")

	util.StructWriter(w, questions)
}
