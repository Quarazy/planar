package main

type MessageType uint8

const (
	Position MessageType = iota
	User
	PlayerJoin
	PlayerLeave
)

const (
	ButtonW uint16 = 1 << iota
	ButtonA
	ButtonS
	ButtonD
	ButtonR
	ButtonSpace
	ButtonControl
)

type Vector3D struct {
	X, Y, Z float64
}

// Location is Vector3D representing the location in a 3D world
type Location Vector3D

// Direction is Vector3D representing the direction in a 3D world
type Direction Vector3D

func NewVector3D(X, Y, Z float64) Vector3D {
	return Vector3D{
		X: X,
		Y: Y,
		Z: Z,
	}
}

func NewLocation(X, Y, Z float64) Location {
	return Location(NewVector3D(X, Y, Z))
}

func NewDirection(X, Y, Z float64) Direction {
	return Direction(NewVector3D(X, Y, Z))
}

type Message struct {
	Type MessageType `json:"type"`
	Id   uint64      `json:"id"`
	// Duration in milliseconds of command
	Msec float64 `json:"msec"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
	Step uint16  `json:"step"`
	// Buttons is a bit field of the buttons being held down
	Buttons uint16 `json:"buttons"`
}

// Half-Life style user input
type UserCommand struct {
	lerpMsec uint16
	msec     byte
}
