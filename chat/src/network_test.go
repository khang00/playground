package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNetwork_SendToClient(t *testing.T) {
	clientIDs := make([]ClientID, 0)
	n := NewNetwork(func(n Network, localClientID ClientID, msg Message) {
		action := msg.Value("action")
		if action == "send_msg" {
			receiverID := msg.Value("receiver")
			err := n.SendToClient(ClientID(receiverID), msg)
			if err != nil {
				t.Error(err)
			}
		}
	})

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testUpgrade := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := testUpgrade.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
		}

		clientID, err := n.RegisterConn(conn)
		if err != nil {
			t.Error(err)
		}

		clientIDs = append(clientIDs, clientID)
	}))
	defer s.Close()

	c1 := startNewClient(t, s, "1", func(msgType int, msg []byte) {})
	c2InboundChan := make(chan Message)
	_ = startNewClient(t, s, "2", func(msgType int, bytes []byte) {
		inboundMsg, err := NewMsgFromByte(bytes)
		if err != nil {
			t.Error(err)
		}

		c2InboundChan <- inboundMsg
	})

	for len(clientIDs) < 2 {

	}

	outboundMsg := NewMsg().
		Add("action", "send_msg").
		Add("receiver", string(clientIDs[1])).
		Add("content", "hello world!")

	outboundMsgBytes, err := outboundMsg.Bytes()
	if err != nil {
		t.Error(err)
	}

	err = c1.Send(outboundMsgBytes)
	if err != nil {
		t.Error(err)
	}

	result := <-c2InboundChan
	if outboundMsg.Value("content") != result.Value("content") {
		t.Fail()
	}
}

func startNewClient(t *testing.T, s *httptest.Server, clientID ClientID, handler func(msgType int, msg []byte)) Client {
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	ws1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Error(err)
	}
	c1 := NewClient(ws1, clientID, handler)
	c1.Run()

	return c1
}
