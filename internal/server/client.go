package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go_chat_tutorial/internal"
	"go_chat_tutorial/internal/common"
	"log"
)

// wsClientImpl is a connection manager
type WsClientImpl struct {
	id   uuid.UUID // id is client ID
	Conn *websocket.Conn
	hub  internal.Hub // references to the hub
}

// NewWsClient ...
func NewWsClient(h internal.Hub, conn *websocket.Conn) internal.Client {
	connID := uuid.New()
	return &WsClientImpl{
		id:   connID,
		Conn: conn,
		hub:  h,
	}
}

func (c *WsClientImpl) Reader() {
	defer func() {
		c.hub.Unregister(c)
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := common.Message{Type: messageType, Body: string(p)}
		c.hub.BroadCastCh(message)
		fmt.Printf("Message Received: %+v\n", message)
	}
}

func (c *WsClientImpl) Send(m common.Message) error {
	if err := c.Conn.WriteJSON(m); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Close closes the Redis subscriber
// then client connection
func (c *WsClientImpl) Close() {
	if err := c.Conn.Close(); err != nil {
		fmt.Errorf("error while closing client connection: %v", err)
	}
}

func (c *WsClientImpl) Writer() {
}
