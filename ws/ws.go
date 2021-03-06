package ws

import (
	"time"

	"encoding/json"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/mateuszdyminski/logag/model"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection id.
	ID string

	// The websocket connection.
	Ws *websocket.Conn

	// The websocket filter.
	Filter *model.Filter

	// Buffered channel of outbound messages.
	Send chan *model.Log
}

// write writes a message with the given message type and payload.
func (c *Connection) Write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.CloseMessage, []byte{})
				return
			}

			json, err := json.Marshal(message)
			if err != nil {
				logrus.Errorf("Can't marshal log for ws cients. Err: %v", err)
				c.Write(websocket.CloseMessage, []byte{})
				return
			}

			sendMsg := false
			if c.Filter == nil {
				sendMsg = true
			} else {
				if c.Filter.Level == "" {
					if len(c.Filter.Keywords) > 0 {
						for _, k := range c.Filter.Keywords {
							if strings.Contains(message.Msg, k) {
								sendMsg = true
								break
							}
						}
					} else {
						sendMsg = true
					}
				} else if c.Filter.Level == message.Level {
					if len(c.Filter.Keywords) > 0 {
						for _, k := range c.Filter.Keywords {
							if strings.Contains(message.Msg, k) {
								sendMsg = true
								break
							}
						}
					} else {
						sendMsg = true
					}
				}
			}

			if sendMsg {
				if err := c.Write(websocket.TextMessage, json); err != nil {
					return
				}
			}
		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
