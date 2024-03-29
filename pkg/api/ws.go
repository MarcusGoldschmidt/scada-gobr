package api

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
