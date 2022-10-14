package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client interface {
	send(msg []byte) error
	run()
}

type ClientMsgHandler = func(msgType int, msg []byte)

type client struct {
	conn    *websocket.Conn
	handler ClientMsgHandler
}

func NewClient(conn *websocket.Conn, handler ClientMsgHandler) Client {
	return &client{
		conn:    conn,
		handler: handler,
	}
}

func (c *client) run() {
	go c.serve()
}

func (c *client) send(msg []byte) error {
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
