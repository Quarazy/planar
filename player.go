package main

type player struct {
	Id       uint64
	Velocity float64
	Location Location
}

// Move updates the player on the new command
func (p *player) Update(c Command, delta float64) {
	p.Location.X += c.X * p.Velocity * c.Delta
	p.Location.Y += c.Y * p.Velocity * c.Delta
	p.Location.Z += c.Z * p.Velocity * c.Delta
}

// NewPlayer creates a new player struct
func NewPlayer(id uint64, velocity float64) *player {
	return &player{
		Id:       id,
		Velocity: velocity,
	}
}
