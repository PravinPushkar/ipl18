package service

import (
	"log"
	"time"

	"github.com/PravinPushkar/ipl18/backend/cache"
	"github.com/PravinPushkar/ipl18/backend/dao"
	"github.com/PravinPushkar/ipl18/backend/models"
)

//WebSocketManager and its dependencies
type WebSocketManager struct {
	conns map[*models.ConnModel]bool
	rch   chan *models.FeedsMessageModel
	cDao  dao.ChatDAO
	uDao  dao.UserDAO
}

//SocketHandler - interface exposed by the service
type SocketHandler interface {
	Add(conn *models.ConnModel)
}

//NewWebSocketManager - Constructor for service
func NewWebSocketManager() *WebSocketManager {
	ws := WebSocketManager{
		conns: make(map[*models.ConnModel]bool),
		rch:   make(chan *models.FeedsMessageModel),
		cDao:  dao.ChatDAO{},
		uDao:  dao.UserDAO{},
	}
	ws.uDao.GetAllUsersBasicInfo()
	return &ws
}

//Start - start monitoring channels
func (ws *WebSocketManager) Start() {
	go ws.monitorChan()
}

//Add new connection to pool
func (ws *WebSocketManager) Add(conn *models.ConnModel) {
	ws.writeInitialFeed(conn)
	go ws.poll(conn)
	ws.conns[conn] = true
}

func (ws *WebSocketManager) writeInitialFeed(c *models.ConnModel) {
	if chats, err := ws.cDao.GetRecentChats(100); err != nil {
		log.Println("WebSocketManager: error getting recent chats", err)
	} else {
		log.Println("WebSocketManager: writing initial chat history to", c.Conn.RemoteAddr())
		for _, chat := range chats {
			user := cache.UserINumberCache[chat.INumber]
			chat.Name = user.Firstname + " " + user.Lastname
			chat.PicLocation = user.PicLocation
			c.Conn.WriteJSON(chat)
		}
	}
}

func (ws *WebSocketManager) monitorChan() {
	log.Println("WebSocketManager: monitorChan starting to monitor channel")
	for {
		select {
		case msg := <-ws.rch:
			log.Println("WebSocketManager: monitorChan new msg", msg)
			for c := range ws.conns {
				log.Println("WebSocketManager: monitorChan writing to", c.Conn.RemoteAddr())
				c.Conn.WriteJSON(msg)
			}
		}
	}
}

func (ws *WebSocketManager) poll(c *models.ConnModel) {
	log.Println("WebSocketManager: poll starting poll for", c.Conn.RemoteAddr())
	defer c.Conn.Close()
	for {
		mt, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocketManager: poll read error", err)
			delete(ws.conns, c)
			break
		}
		msg := string(message)
		log.Println("WebSocketManager: poll new message", mt, msg)
		date := time.Now()

		user := cache.UserINumberCache[c.INumber]
		ws.rch <- &models.FeedsMessageModel{c.INumber, user.Firstname + " " + user.Lastname, msg, date, user.PicLocation}
		ws.cDao.InsertChat(message, c.INumber, date)
	}
}
