package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go_chat_tutorial/internal"
	"go_chat_tutorial/internal/common"
	"go_chat_tutorial/pkg/logger"
	pb "go_chat_tutorial/proto/ws-gateway"
	"time"
)

// wsClientImpl is a connection manager
type WsClientImpl struct {
	id      uuid.UUID // id is client ID
	UserID  string    // userID stores ID of user if authenticated
	Conn    *websocket.Conn
	hub     internal.Hub // references to the hub
	writeCh chan common.Message
	log     *logrus.Entry
	exit    chan bool
	authCh  chan string
	nowFunc func() time.Time
}

// NewWsClient ...
func NewWsClient(h internal.Hub, conn *websocket.Conn) internal.Client {
	connID := uuid.New()
	return &WsClientImpl{
		id:      connID,
		Conn:    conn,
		hub:     h,
		writeCh: make(chan common.Message, 256),
		log:     logger.WithField("connID", connID.String()),
		exit:    make(chan bool),
		authCh:  make(chan string, 4),
	}
}

func (c *WsClientImpl) pongHandler(_ string) error {
	const tag = "[wsClientImpl.pongHandler] "
	log := c.log.WithField("userID", c.UserID)

	err := c.Conn.SetReadDeadline(c.nowFunc().Add(pongWait))
	if err != nil {
		log.Errorf(tag+"error while setting read deadline: %v", err)
	}
	return err
}

func (c *WsClientImpl) ReadPump() {
	const tag = "[wsClientImpl.readPump]"
	log := c.log.WithField("tag", tag)

	defer func() {
		log.Debugf(tag+" close client (user: %q, conn: %q)", c.UserID, c.id.String())
		c.hub.Unregister(c)
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Errorf(tag+" error while set read deadline: %v", err)
		return
	}
	//c.Conn.SetPongHandler(c.pongHandler)
	expectedCloses := []int{
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
		websocket.CloseNoStatusReceived,
	}
	var (
		payload []byte
		cmd     pb.Command
		msgType int
		//err     error
	)

	for {
		msgType, payload, err = c.Conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, expectedCloses...) {
				log.Infof(tag+" connection is closed (expected): %v", err)
			} else {
				log.Warnf(tag+" connection is closed (unknown reason): %v", err)
			}
			return
		}

		if msgType != websocket.TextMessage {
			log.Warnf(tag+" unknown message type: %v", msgType)
			continue
		}
		messageID := uuid.New().String()
		log.WithField("payload", string(payload)).
			WithField("messageID", messageID).
			Debugf(tag+" processing message %q", messageID)

		authBody := pb.Command_Auth{
			Auth: &pb.AuthBody{Token: string(payload)},
		}
		cmd = pb.Command{Body: &authBody}
		//err = json.Unmarshal(payload, &cmd)
		//if err != nil {
		//	log.WithField("messageID", messageID).Errorf(tag+" unable to unmarshal message: %v", err)
		//	continue
		//}
		//if cmd.GetId() == "" {
		//	cmd.Id = messageID
		//}

		c.hub.ProcessMessage(&cmd, c)
	}
}

//func (c *WsClientImpl) ReadPump() {
//	const tag = "[wsClientImpl.readPump]"
//	log := c.log.WithField("tag", tag)
//
//	defer func() {
//		log.Debugf(tag+" close client (user: %q, conn: %q)", c.userID, c.id.String())
//		c.hub.Unregister(c)
//	}()
//
//	c.Conn.SetReadLimit(maxMessageSize)
//	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
//	if err != nil {
//		log.Errorf(tag+" error while set read deadline: %v", err)
//		return
//	}
//	c.Conn.SetPongHandler(c.pongHandler)
//	expectedCloses := []int{
//		websocket.CloseNormalClosure,
//		websocket.CloseGoingAway,
//		websocket.CloseAbnormalClosure,
//		websocket.CloseNoStatusReceived,
//	}
//
//	for {
//		mmsgType, payload, err = c.conn.ReadMessage()if err != nil {
//			if !websocket.IsUnexpectedCloseError(err, expectedCloses...) {
//				log.Infof(tag+" connection is closed (expected): %v", err)
//			} else {
//				log.Warnf(tag+" connection is closed (unknown reason): %v", err)
//			}
//			return
//		}
//
//		if msgType != websocket.TextMessage {
//			log.Warnf(tag+" unknown message type: %v", msgType)
//			continue
//		}
//		messageID := uuid.New().String()
//		log.WithField("payload", string(payload)).
//			WithField("messageID", messageID).
//			Debugf(tag+" processing message %q", messageID)
//		p := bytes.NewReader(payload)
//		err = jsonpb.Unmarshal(p, &cmd)
//		if err != nil {
//			log.WithField("messageID", messageID).Errorf(tag+" unable to unmarshal message: %v", err)
//			continue
//		}
//		if cmd.GetId() == "" {
//			cmd.Id = messageID
//		}
//
//		c.hub.ProcessMessage(&cmd, c)
//	}
//	var (
//		payload []byte
//		cmd     pb.Command
//		msgType int
//	)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//
//		authBody := pb.Command_Auth{
//			Auth: &pb.AuthBody{Token: string(payload)},
//		}
//
//		cmd := pb.Command{Body: &authBody}
//		c.hub.ProcessMessage(&cmd, c)
//	}
//}

func (c *WsClientImpl) Send(m common.Message) {
	c.writeCh <- m
}

// Close closes the Redis subscriber
// then client connection
func (c *WsClientImpl) Close() {
	if err := c.Conn.Close(); err != nil {
		fmt.Errorf("error while closing client connection: %v", err)
	}
}

func (c *WsClientImpl) WritePump() {
	tag := "[wsClientImpl.writePump]"
	log := c.log.WithField("tag", tag)
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.writeCh:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Errorf(tag+" error while write close message: %v", err)
				}
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Errorf(tag+" failed to take next writer: %v", err)
				return
			}

			byt, err := json.Marshal(message)
			if err != nil {
				continue
			}
			if _, err = w.Write(byt); err != nil {
				log.Error(tag+" failed to write message: %v", err)
				continue
			}

			// add queued chat messages to the current websocket message
			n := len(c.writeCh)
			for i := 0; i < n; i++ {
				if byt, err = json.Marshal(<-c.writeCh); err != nil {
					log.Errorf(tag+" failed marshaler message: %v", err)
					continue
				}
				if _, err = w.Write(byt); err != nil {
					log.Errorf(tag+" failed to write message: %v", err)
				}
			}

			log.Debugf(tag+"send `%d` message(s) successfully", n+1)

			if err = w.Close(); err != nil {
				log.Errorf(tag+"failed to close writer: %v", err)
				return
			}

		case <-ticker.C:
			if err := c.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Errorf(tag+" failed to send ping message: %v", err)
				return
			}

		case <-c.exit:
			log.Infof(tag + "received close signal")
			return

		// after authorized, update logger for better logging
		case userID, ok := <-c.authCh:
			if ok {
				log = log.WithField("userID", userID)
			}
		}
	}
}

func (c *WsClientImpl) GetID() uuid.UUID {
	return c.id
}

// GetUserID ...
func (c *WsClientImpl) GetUserID() string {
	return c.UserID
}

// 	SetUserID sets userID to this client
// then unblock subscribe channel to start subscribing
func (c *WsClientImpl) SetUserID(id string) {
	if c.UserID != id {
		c.UserID = id
		c.authCh <- id
		c.hub.NotifyOnline(c)
	}
}
