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

	steps []map[uint16]*Message
}

// newHub instantiates a new hub
func newHub() *hub {
	hub := &hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		broadcast:  make(chan *Message),
		steps:      make([]map[uint16]*Message, 2),
	}

	// initialize steps map
	for i := 0; i < 2; i++ {
		hub.steps[i] = make(map[uint16]*Message)
	}

	return hub
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
			// Send move to another channel if I ever use multiple go routines
			// Right now there's only one go routine, so this logic will be safe
			h.steps[message.Id][message.Step] = message

			if allStepsReceived(message.Step, h.steps) {
				for _, v := range h.steps {
					for client := range h.clients {
						select {
						case client.send <- v[message.Step]:
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

// allStepsReceived checks to see if all steps
// have been received
func allStepsReceived(step uint16, steps []map[uint16]*Message) bool {
	for _, v := range steps {
		if _, ok := v[step]; !ok {
			return false
		}
	}

	return true
}
