package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendMessage(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testUpgrade := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := testUpgrade.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
		}

		inboundChan1 := make(chan []byte, 0)
		inboundChan2 := make(chan []byte, 0)
		c1 := NewClient(conn, func(messageType int, message []byte) {
			inboundChan1 <- message
		})
		c1.run()

		c2 := NewClient(conn, func(messageType int, message []byte) {
			inboundChan2 <- message
		})
		c2.run()

		msg1 := []byte("hello from c2 to c1")
		err = c2.send(msg1)
		if err != nil {
			t.Error(err)
		}

		msg1Result := <-inboundChan1
		if string(msg1) != string(msg1Result) {
			t.Fail()
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Error(err)
	}
}
