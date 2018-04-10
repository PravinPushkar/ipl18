package models

import "github.com/gorilla/websocket"

type ConnModel struct {
	INumber string `json"inumber"`
	Conn    *websocket.Conn
}
