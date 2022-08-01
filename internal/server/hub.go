package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"go_chat_tutorial/internal"
	"go_chat_tutorial/internal/common"
	"go_chat_tutorial/internal/processor"
	pb "go_chat_tutorial/proto/ws-gateway"
	"net/http"
	"time"
)

type zeroSt struct{}

// hubImpl is a centralized messages to manage clients, send messages
type hubImpl struct {
	registerCh      chan internal.Client // register requests from the clients
	unregisterCh    chan internal.Client // unregister requests from the clients
	config          *Config
	log             *logrus.Entry
	nowFunc         func() time.Time
	httpSever       *http.Server
	clients         map[internal.Client]zeroSt
	processors      []internal.Processor
	onlineCh        chan internal.Client
	clientsByUserID map[string]map[internal.Client]bool
}

// NewHub ...
func NewHub(hubCfg *Config, l *logrus.Entry) internal.Hub {
	return &hubImpl{
		registerCh:   make(chan internal.Client),
		unregisterCh: make(chan internal.Client),
		config:       hubCfg,
		log:          l,
		nowFunc:      time.Now,
		clients:      make(map[internal.Client]zeroSt),
		processors:   []internal.Processor{
			//processor.NewAuth(hubCfg.AuthUrl),
		},
		onlineCh:        make(chan internal.Client),
		clientsByUserID: make(map[string]map[internal.Client]bool),
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
	log.Debugf(tag+" new client connected `%s` (host: `%s`)", client.GetID(), r.Host)

	// allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	h.Register(client)

	go client.WritePump()
	go client.ReadPump()

}

func (h *hubImpl) Register(c internal.Client) {
	h.registerCh <- c
}

func (h *hubImpl) Unregister(c internal.Client) {
	h.unregisterCh <- c
}

// Run ...
func (h *hubImpl) Run(ctx context.Context) error {
	// log 2
	go h.serveHTTP() // khởi tạo websocket với ReadPump và WritePump

	h.processors = append(
		h.processors,
		processor.NewAuth(h.config.AuthUrl), // check auth
	)
	for {
		select {
		case client := <-h.registerCh:
			h.clients[client] = zeroSt{}
		case client := <-h.unregisterCh:
			client.Close()
			delete(h.clients, client)
		case client := <-h.onlineCh:
			userID := client.GetUserID()
			if _, ok := h.clientsByUserID[userID]; !ok {
				h.clientsByUserID[userID] = make(map[internal.Client]bool)
			}
			h.clientsByUserID[userID][client] = true

			// unregister case, we need to remove client from clients list
			// and also the reference in user clients list
		}
	}
}

//
//// HandleUserRegisterEvent will handle the Join event for New socket users
//func HandleUserRegisterEvent(hub *hubImpl, client internal.Client) {
//	hub.clients[client] = zeroSt{}
//	handleSocketPayloadEvents(client, common.Message{
//		EventName:    "join",
//		EventPayload: client.GetUserID(),
//	})
//}
//
//func handleSocketPayloadEvents(client internal.Client, socketEventPayload common.Message) {
//	//var socketEventResponse common.Message
//	//switch socketEventPayload.EventName {
//	//case "join":
//	//	log.Printf("Join Event triggered")
//	BroadcastSocketEventToAllClient(client.hub, common.Message{
//		EventName: socketEventPayload.EventName,
//		EventPayload: JoinDisconnectPayload{
//			UserID: client.userID,
//			Users:  getAllConnectedUsers(client.hub),
//		},
//	})

//case "disconnect":
//	log.Printf("Disconnect Event triggered")
//	BroadcastSocketEventToAllClient(client.hub, SocketEventStruct{
//		EventName: socketEventPayload.EventName,
//		EventPayload: JoinDisconnectPayload{
//			UserID: client.userID,
//			Users:  getAllConnectedUsers(client.hub),
//		},
//	})
//
//case "message":
//	log.Printf("Message Event triggered")
//	selectedUserID := socketEventPayload.EventPayload.(map[string]interface{})["userID"].(string)
//	socketEventResponse.EventName = "message response"
//	socketEventResponse.EventPayload = map[string]interface{}{
//		"username": getUsernameByUserID(client.hub, selectedUserID),
//		"message":  socketEventPayload.EventPayload.(map[string]interface{})["message"],
//		"userID":   selectedUserID,
//	}
//	EmitToSpecificClient(client.hub, socketEventResponse, selectedUserID)
//}
//}
//
//// BroadcastSocketEventToAllClient will emit the socket events to all socket users
//func BroadcastSocketEventToAllClient(hub *Hub, payload common.Message) {
//	for client := range hub.clients {
//		select {
//		case client.send <- payload:
//		default:
//			close(client.send)
//			delete(hub.clients, client)
//		}
//	}
//}

func (h *hubImpl) ProcessMessage(cmd *pb.Command, c internal.Client) {
	var (
		err error
		log = logrus.WithField("func", "hubImpl.handleCommand")
	)

	for _, p := range h.processors {
		if err = p.Handle(cmd, c); err != nil {
			log.WithError(err).Debug("process message error")
		}
	}

	keys := make([]internal.Client, len(h.clients))

	i := 0
	for k := range h.clients {
		keys[i] = k
		i++
	}
	for _, client := range keys {
		if client.GetUserID() == c.GetUserID() {
			client.Send(common.Message{
				EventName: "message response",
				EventPayload: map[string]interface{}{
					"username": "chwa co ten",
					"message":  "message ne",
					"userID":   c.GetUserID(),
				}})
		}
	}
}

func (h *hubImpl) NotifyOnline(c internal.Client) {
	h.onlineCh <- c
}
