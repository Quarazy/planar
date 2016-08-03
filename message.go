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
	X    float32     `json:"x"`
	Y    float32     `json:"y"`
	Z    float32     `json:"z"`
}
