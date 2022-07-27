package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_chat_tutorial/internal"
	"go_chat_tutorial/internal/common"
	"net/http"
	"time"
)

// hubImpl is a centralized messages to manage clients, send messages
type hubImpl struct {
	registerCh   chan internal.Client // register requests from the clients
	unregisterCh chan internal.Client // unregister requests from the clients
	config       *Config
	log          *logrus.Entry
	nowFunc      func() time.Time
	httpSever    *http.Server
	clients      map[internal.Client]bool
	broadCastCh  chan common.Message
}

// NewHub ...
func NewHub(hubCfg *Config, l *logrus.Entry) internal.Hub {
	return &hubImpl{
		registerCh:   make(chan internal.Client),
		unregisterCh: make(chan internal.Client),
		config:       hubCfg,
		log:          l,
		nowFunc:      time.Now,
		clients:      make(map[internal.Client]bool),
		broadCastCh:  make(chan common.Message),
	}
}

// wsHandler ...
func (h *hubImpl) wsHandler(w http.ResponseWriter, r *http.Request) {
	tag := "[hubImpl.wsHandler]"
	log := h.log.WithField("tag", tag)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf(tag+" failed to upgrade connection: %v", err)
		return
	}

	client := NewWsClient(h, conn)
	//log.Debugf(tag+" new client connected `%s` (host: `%s`)", client.GetID(), r.Host)

	h.registerCh <- client

	// allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.Reader()
}

func (h *hubImpl) Register(c internal.Client) {
	h.registerCh <- c
}

func (h *hubImpl) Unregister(c internal.Client) {
	h.unregisterCh <- c
}

func (h *hubImpl) BroadCastCh(m common.Message) {
	h.broadCastCh <- m
}

// Run ...
func (h *hubImpl) Run(ctx context.Context) error {
	go h.serveHTTP() // khởi tạo websocket với Reader và Writer

	for {
		select {
		case client := <-h.registerCh:
			h.clients[client] = true
		case client := <-h.unregisterCh:
			delete(h.clients, client)
			client.Close()
		case message := <-h.broadCastCh:
			for client := range h.clients {
				if err := client.Send(message); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
