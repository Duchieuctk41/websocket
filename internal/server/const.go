package server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4 * 1024,
	WriteBufferSize: 4 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IsUserOnlineRequest struct {
	UserID string `json:"UserID,omitempty"`
}

type IsUserOnlineResponse struct {
	IsOnline bool `json:"isOnline,omitempty"`
}

const (
	// maximum message size allowed from peer.
	maxMessageSize = 4 * 1024

	// time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// writeCh pings to peer with this period. must be less than pongWait
	pingPeriod = (pongWait * 5) / 10

	// time allowed to write a message to the peer
	writeWait = 10 * time.Second
)
