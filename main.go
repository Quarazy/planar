package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	hub := newHub()
	go hub.run()

	s := &Server{hub: hub}
	http.Handle("/ws", s)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Fatal(http.ListenAndServe("localhost:9000", nil))
}

// A server keeps track of all client connections
type Server struct {
	hub     *hub
	counter uint64
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveConnect(w, r)
}

func (s *Server) serveConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := newClient(s.counter, conn, s.hub)
	defer func() {
		s.hub.unregister <- client
		client.ws.Close()
	}()

	s.hub.Register(client)
	go client.writer()

	// Broadcast to let all users know about new user
	s.hub.broadcast <- &Message{
		Type: PlayerJoin,
		Id:   s.counter,
	}

	// Send world information back to user
	// TODO(quarazy): Need to broadcast location of all existing players
	client.send <- &Message{
		Type: User,
		Id:   s.counter,
	}

	atomic.AddUint64(&s.counter, 1)
	client.reader()
}
