package main

import (
	"errors"
	"sync"
)

/*
Slice that supports lock-free concurrent append and delete
*/

// Double linked node
type node struct {
	next  *node
	prev  *node
	value interface{}
}

const (
	AppendCommand uint8 = iota
	PopLeftCommand
	PopAllCommand
)

// bufferCommand is whatever the channel accepts
type bufferCommand struct {
	command uint8
	value   interface{}
	resCh   chan<- interface{}
}

type Buffer struct {
	// head is where to pop from
	head *node
	// tail is where to append to
	tail *node

	// Protects pop. Only one pop operation can be done at a time
	mu sync.Mutex
	// newCount is number of new appends since last pop
	newCount int

	lastCount int

	commands chan bufferCommand
}

// NewBuffer creates a new buffer
func NewBuffer() *Buffer {
	buf := &Buffer{
		commands: make(chan bufferCommand),
	}
	go buf.run()

	return buf
}

// run should be run in it's own goroutine
func (buf *Buffer) run() {
	for {
		select {
		case c := <-buf.commands:
			switch c.command {
			case AppendCommand:
				buf.appendRight(c.value)
			case PopLeftCommand:
				ch := c.resCh
				val := buf.popLeft()
				ch <- val
			case PopAllCommand:
				buf.mu.Lock()
				// Starts another goroutine to pop all new entries
				// so for loop isn't blocked

				usingCount := buf.newCount
				go func(resCh chan<- interface{}) {
					for i := 0; i < usingCount-buf.lastCount; i++ {
						resCh <- buf.popLeft()
					}

					resCh <- true

					buf.lastCount = usingCount
					buf.mu.Unlock()
				}(c.resCh)
			}
		}
	}
}

// Append adds elements from the right
// Using loks right now. Don't!!!!
// http://www.cs.rochester.edu/~scott/papers/1996_PODC_queues.pdf

// Append adds an element to the end of the buffer
func (buf *Buffer) Append(value interface{}) {
	buf.commands <- bufferCommand{
		command: AppendCommand,
		value:   value,
	}
}

// PopLeft removes elements from the left
func (buf *Buffer) PopLeft() (interface{}, error) {
	resCh := make(chan interface{})
	defer close(resCh)

	buf.commands <- bufferCommand{
		command: PopLeftCommand,
		resCh:   resCh,
	}

	res := <-resCh

	if res == nil {
		return nil, errors.New("Pop on an empty buffer")
	}

	return res, nil
}

// Copy copies the present state of the buffer
// TODO(quarazy): Need to look at convention.
// 1. One way to use this is an interface
// 2. Naming might be different too Clone()
func (buf *Buffer) PopAll() []interface{} {
	resCh := make(chan interface{})
	defer close(resCh)

	buf.commands <- bufferCommand{
		command: PopAllCommand,
		resCh:   resCh,
	}

	var elems []interface{}
	for elem := range resCh {
		// true signifies done
		if elem == true {
			break
		}

		elems = append(elems, elem)
	}

	return elems
}

func (buf *Buffer) appendRight(value interface{}) {
	n := &node{
		prev:  buf.tail,
		value: value,
	}

	if buf.tail != nil {
		buf.tail.next = n
	}

	buf.tail = n

	if buf.head == nil {
		buf.head = buf.tail
	}

	buf.newCount += 1
}

func (buf *Buffer) popLeft() interface{} {
	// Check if there's anything in the queue
	if buf.head == nil {
		return nil
	}

	// Move pointer forward
	value := buf.head.value
	buf.head = buf.head.next

	return value
}

func (buf *Buffer) copy() {

}
