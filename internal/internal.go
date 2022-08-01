package internal

import (
	"context"
	"github.com/google/uuid"
	"go_chat_tutorial/internal/common"
	pb "go_chat_tutorial/proto/ws-gateway"
)

// Client presents a websocket client
// it includes a socket connection, and user information (if authenticated)
type Client interface {
	ReadPump()
	WritePump()
	Send(m common.Message)
	Close()
	GetID() uuid.UUID
	GetUserID() string
	SetUserID(id string)
}

type Hub interface {
	Run(ctx context.Context) error
	Register(c Client)
	Unregister(c Client)
	ProcessMessage(m *pb.Command, client Client)
	NotifyOnline(c Client)
}

// Processor processes a message
type Processor interface {
	Handle(cmd *pb.Command, client Client) error
}
