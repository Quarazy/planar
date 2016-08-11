package main

type CommandCode uint16

const (
	Move CommandCode = iota
	Shoot
	Jump
)

// Command is an action sent from the client
type Command struct {
	ObjectId   uint64
	Code       CommandCode
	Repeatable bool
	X, Y, Z    float64
	Delta      float64
}

// Converts a Message type to a Command
// TODO(quarazy): This should be simple. Find the naming convention for type conversion methods
func ToCommand(msg Message) Command {
	direction := Direction{}

	if ButtonW&msg.Buttons == ButtonW {
		direction.Y = 1
	} else if ButtonS&msg.Buttons == ButtonS {
		direction.Y = -1
	}

	if ButtonA&msg.Buttons == ButtonA {
		direction.X = -1
	} else if ButtonD&msg.Buttons == ButtonD {
		direction.X = 1
	}

	return Command{
		ObjectId: msg.Id,
		Code:     Move,
		Delta:    msg.Msec,
		X:        direction.X,
		Y:        direction.Y,
		Z:        direction.Z,
	}
}
