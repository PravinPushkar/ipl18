package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend/models"
	"github.wdf.sap.corp/I334816/ipl18/backend/service"

	"github.wdf.sap.corp/I334816/ipl18/backend/errors"

	"github.com/gorilla/websocket"
)

// FeedsSocketHandler ...
type FeedsSocketHandler struct {
	SockMgr service.SocketHandler
}

var upgrader = websocket.Upgrader{}

var errSameOriginNotFound = fmt.Errorf("Different Origin")

func (f FeedsSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panicked", r)
		}
	}()

	f.checkSameOrigin(w, r)

	log.Println("buzz_websocket_handler : starting connection with client")

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade : ", err)
		return
	}
	// defer c.Close()
	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read: ", err)
	}
	log.Println("recv:", mt, string(message))
	feedsAuthModel := models.FeedsAuthModel{}
	if err := json.NewDecoder(bytes.NewReader(message)).Decode(&feedsAuthModel); err != nil {
		errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrParseContext, "Parse Error")
	} else {
		claims, err := tokenManager.GetClaims(feedsAuthModel.Authorization)
		errors.ErrWriterPanic(w, http.StatusForbidden, err, errors.ErrTokenInfoMismatch, "token not valid")
		log.Println("FeedsSocketHandler: adding connection to manager", c.RemoteAddr())
		inumber, _ := claims["inumber"]
		f.SockMgr.Add(&models.ConnModel{inumber.(string), c})
	}
}

func (f FeedsSocketHandler) checkSameOrigin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		errors.ErrWriterPanic(w, http.StatusForbidden, errSameOriginNotFound, errSameOriginNotFound, "Origin not allowed")
		return
	}
}
