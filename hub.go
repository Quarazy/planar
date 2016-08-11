package main

// Hub maintains the set of active clients
type hub struct {
	// Registered clients
	clients map[*client]bool

	// Registered players
	// TODO(quarazy): Consider combining player and client
	players map[uint64]*player

	// Register requests from the clients
	register chan *client

	// Unregister requests from the clients
	unregister chan *client

	// Outbound broadcast messages to the clients.
	broadcast chan *Message

	// Inbound messages from the client
	actions chan *Message
}

// newHub instantiates a new hub
func newHub() *hub {
	hub := &hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		players:    make(map[uint64]*player),
		broadcast:  make(chan *Message),
		actions:    make(chan *Message),
	}

	go hub.run()
	return hub
}

// run runs in its own goroutine
func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.players[client.id] = NewPlayer(client.id, 4.0)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.remove(client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.remove(client)
				}
			}
		}
	}
}

func (h *hub) remove(client *client) {
	delete(h.clients, client)
	// Closes the channel
	close(client.send)
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

// ConnectedPlayers returns the ids of all currently connected players
// TODO(quarazy): Should check to see best practice on how to check state like this
// Can either do this or constantly update another structure when players join/leave
func (h *hub) ConnectedPlayers() []uint64 {
	keys := make([]uint64, len(h.clients))

	i := 0
	for client := range h.clients {
		keys[i] = client.id
		i++
	}

	return keys
}
