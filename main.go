package main

import (
	"log"
	"net/http"
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
	hub *hub
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveConnect(w, r)
}

func (s *Server) serveConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	s.hub.Register(newClient(conn))
}
