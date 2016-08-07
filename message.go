package main

type MessageType uint8

const (
	Position MessageType = iota
	User
	PlayerJoin
	PlayerLeave
)

type Message struct {
	Type MessageType `json:"type"`
	Id   uint64      `json:"id"`
	Msec float64     `json:"msec"`
	X    float32     `json:"x"`
	Y    float32     `json:"y"`
	Z    float32     `json:"z"`
	Step uint16      `json:"step"`
}

// Half-Life style user input
type UserCommand struct {
	lerpMsec uint16
	msec     byte
}
