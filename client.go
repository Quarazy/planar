package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type client struct {
	// Unique identifier for client
	id uint64

	// The websocket connection
	ws *websocket.Conn

	// Messsages to be sent
	send chan *Message

	// Reference to the hub
	hub *hub
}

// newClient instantiates a new client
func newClient(id uint64, ws *websocket.Conn, hub *hub) *client {
	return &client{
		id:   id,
		ws:   ws,
		send: make(chan *Message, 256),
		hub:  hub,
	}
}

func (c *client) reader() {
	for {
		message := &Message{}
		if err := c.ws.ReadJSON(message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Close error: %v\n", err)
			}
			log.Printf("Error when reading message: %s\n", err)
			break
		}
		c.hub.actions <- message
	}
}

// writer should run in it's own go process
func (c *client) writer() {
	for message := range c.send {
		if err := c.ws.WriteJSON(message); err != nil {
			continue
		}
	}

	c.ws.Close()
}
