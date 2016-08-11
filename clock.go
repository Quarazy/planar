package main

import "time"

type GameClock struct {
	started     bool
	Tick        uint64
	time        time.Time
	timeChannel <-chan time.Time
}

// GetDelta returns the time that has elapsed since the last call
// Not thread safe
func (c *GameClock) GetDelta() float64 {
	now := time.Now()
	c.Tick += 1

	if !c.started {
		c.started = true
		c.time = now
		return 0
	}

	elapsed := float64(now.Sub(c.time) / time.Second)
	c.time = now
	return elapsed
}

func (c *GameClock) Start(d time.Duration) {
	c.timeChannel = time.NewTicker(d).C
}
