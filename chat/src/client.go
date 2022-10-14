package main

import (
	"github.com/gorilla/websocket"
	"log"
)

/*
Client manages a web socket connection to another peer
*/
type Client interface {
	Send(msg []byte) error
	Run()
	ID() ClientID
}

/*
ClientID is the local id of this client.
*/
type ClientID string
type ClientMsgHandler = func(msgType int, msg []byte)

type client struct {
	conn     *websocket.Conn
	handler  ClientMsgHandler
	clientID ClientID
}

func NewClient(conn *websocket.Conn, clientID ClientID, handler ClientMsgHandler) Client {
	return &client{
		conn:     conn,
		handler:  handler,
		clientID: clientID,
	}
}

func (c *client) ID() ClientID {
	return c.clientID
}

func (c *client) Run() {
	go c.serve()
}

func (c *client) Send(msg []byte) error {
	err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		return nil
	}

	return nil
}

func (c *client) serve() {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handler(msgType, msg)
	}
}
