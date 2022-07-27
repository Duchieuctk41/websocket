package internal

import (
	"context"
	"go_chat_tutorial/internal/common"
)

// Client presents a websocket client
// it includes a socket connection, and user information (if authenticated)
type Client interface {
	Reader()
	Writer()
	Send(m common.Message) error
	Close()
}

type Hub interface {
	Run(ctx context.Context) error
	Register(c Client)
	Unregister(c Client)
	BroadCastCh(m common.Message)
}
