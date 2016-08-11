package main

import "time"

// tickRate is the number of times a second the server will run a simulation of the game
// During each tick, server processes user input, runs simulation, takes a snapshot and broadcasts these back to the client

// TODO(quarazy): Clients should be able to request a tickRate to accomodate
// lower bandwidth
const tickRate = 25

// interpolation sends updates to the client at a set interval
type interpolation struct {
	clock *GameClock
	hub   *hub

	// Buffer that holds all actions for player
	buffer *Buffer

	actions <-chan *Message
}

func newInterpolation(hub *hub) *interpolation {
	inter := &interpolation{
		hub:     hub,
		clock:   &GameClock{},
		actions: hub.actions,
		buffer:  NewBuffer(),
	}

	go inter.run()
	return inter
}

// run is expected to run in it's own goroutine
func (i *interpolation) run() {
	i.clock.Start(time.Millisecond * 1000 / tickRate)

	for {
		select {
		// FIXME(quarazy): This feels wrong. Shouldn't be passing channels
		case message := <-i.actions: // Add all messages to buffer
			i.buffer.Append(ToCommand(*message))
		case <-i.clock.timeChannel:
			i.clock.GetDelta()
			// Send snapshot to current
			for _, player := range i.hub.players {
				i.hub.broadcast <- &Message{
					Id:   player.Id,
					X:    player.Location.X,
					Y:    player.Location.Y,
					Z:    player.Location.Z,
					Step: i.clock.Tick,
					Msec: float64(time.Now().UnixNano()) / float64(time.Millisecond),
				}
			}

			for _, v := range i.buffer.PopAll() {
				com, _ := v.(Command)
				i.hub.players[com.ObjectId].Update(com, com.Delta)
			}
		}
	}
}
