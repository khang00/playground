package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type NetMsgHandler = func(net Network, localClientID ClientID, message Message)
type Network interface {
	SendToClient(clientID ClientID, message Message) error
	RegisterConn(conn *websocket.Conn) (ClientID, error)
}

type network struct {
	msgHandler NetMsgHandler
	clients    map[ClientID]Client
}

func NewNetwork(handler NetMsgHandler) Network {
	return &network{
		msgHandler: handler,
		clients:    make(map[ClientID]Client, 0),
	}
}

func (n *network) RegisterConn(conn *websocket.Conn) (ClientID, error) {
	clientID := ClientID(uuid.New().String())
	c := NewClient(conn, clientID, func(msgType int, msgBytes []byte) {
		msg, err := NewMsgFromByte(msgBytes)
		if err != nil {
			return
		}

		n.msgHandler(n, clientID, msg)
	})

	n.clients[clientID] = c
	c.Run()

	return clientID, nil
}

func (n *network) SendToClient(clientID ClientID, message Message) error {
	c, ok := n.clients[clientID]
	if !ok {
		return errors.New("error: no client with ID " + string(clientID))
	}

	bytes, err := message.Bytes()
	if err != nil {
		return err
	}

	err = c.Send(bytes)
	if err != nil {
		return err
	}

	return nil
}
