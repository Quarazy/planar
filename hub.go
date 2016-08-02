package main

// Hub maintains the set of active clients
type hub struct {
	// Registered clients
	clients map[*client]bool

	// Register requests from the clients
	register chan *client

	// Unregister requests from the clients
	unregister chan *client
}

// newHub instantiates a nes hub
func newHub() *hub {
	return &hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
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
