package dao

import (
	"log"
	"net/http"
	"time"

	"github.wdf.sap.corp/I334816/ipl18/backend/db"
	"github.wdf.sap.corp/I334816/ipl18/backend/errors"
	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

type ChatDAO struct{}

const (
	qInsertChat = "INSERT INTO chats(data,inumber,chatdate) VALUES($1,$2,$3)"
	qSelectChat = "SELECT data,inumber,chatdate FROM chats ORDER BY cid DESC LIMIT $1"
)

func (cd ChatDAO) InsertChat(message []byte, inumber string, date time.Time) error {
	resp, err := db.DB.Exec(qInsertChat, message, inumber, date)
	if err != nil {
		log.Println("ChatDAO : InsertChat", err)
		return &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	num, err := resp.RowsAffected()
	if err != nil || num != 1 {
		log.Println("ChatDAO : InsertChat rows affected", num, err)
		return &errors.DaoError{http.StatusInternalServerError, errors.ErrDBIssue, errors.ErrDBIssue}
	}
	log.Println("ChatDAO: InsertChat: inserted")

	return nil
}

func (cd ChatDAO) GetRecentChats(top int) ([]*models.FeedsMessageModel, error) {
	resp, err := db.DB.Query(qSelectChat, top)
	if err != nil {
		log.Println("ChatDAO : GetRecentChats", err)
		return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
	}
	defer resp.Close()
	chats := []*models.FeedsMessageModel{}
	for resp.Next() {
		chat := models.FeedsMessageModel{}
		data := []byte{}
		if err := resp.Scan(&data, &chat.INumber, &chat.Date); err != nil {
			log.Println("ChatDAO : GetRecentChats", err)
			return nil, &errors.DaoError{http.StatusInternalServerError, err, errors.ErrDBIssue}
		}
		chat.Message = string(data)
		chats = append(chats, &chat)
	}
	return chats, nil
}
