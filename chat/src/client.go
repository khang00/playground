package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client interface {
	send(message []byte) error
	run()
}

type MessageHandler = func(messageType int, message []byte)

type client struct {
	conn    *websocket.Conn
	handler MessageHandler
}

func NewClient(conn *websocket.Conn, handler MessageHandler) Client {
	return &client{
		conn:    conn,
		handler: handler,
	}
}

func (c *client) run() {
	go c.serve()
}

func (c *client) send(message []byte) error {
	err := c.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		return nil
	}

	return nil
}

func (c *client) serve() {
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handler(messageType, message)
	}
}
