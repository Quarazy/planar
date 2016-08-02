package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type client struct {
	// The websocket connection
	ws *websocket.Conn

	// Messsages to be sent
	send chan []byte
}

// newClient instantiates a new client
func newClient(ws *websocket.Conn) *client {
	return &client{
		ws:   ws,
		send: make(chan []byte, 256),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
