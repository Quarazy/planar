package main

// Hub maintains the set of active clients
type hub struct {
	// Registered clients
	clients map[*client]bool

	// Register requests from the clients
	register chan *client

	// Unregister requests from the clients
	unregister chan *client

	// Inbound messages from the clients.
	broadcast chan *Message
}

// newHub instantiates a nes hub
func newHub() *hub {
	return &hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		broadcast:  make(chan *Message),
	}
}

// run runs in its own goroutine
func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)

				// Closes the channel
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				// If connection's send buffer is full, then the hub assumes that
				// the client is dead or stuck. In this case, it unregisters the conn
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}

// Register registers a new client to hub
func (h *hub) Register(c *client) error {
	h.register <- c
	return nil
}

// Size returns the number of clients
func (h *hub) Size() int {
	return len(h.clients)
}
