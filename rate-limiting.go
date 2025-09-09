package main

/*
	Rate-limiting (anti-flood) por usuario
*/

import (
	"net"
	"sync"
	"time"
)

const (
	RateLimitCount  = 3
	RateLimitWindow = 10
)

// lista circular

type CircularBuffer struct {
	time  [RateLimitCount]time.Time
	index int
	full  bool
}

var userRates = make(map[net.Conn]*CircularBuffer)
var mu sync.RWMutex

func allowCommand(conn net.Conn) bool {
	now := time.Now()
	mu.Lock()
	defer mu.Unlock()

	cb, ok := userRates[conn]

	if !ok {
		cb = &CircularBuffer{}
		userRates[conn] = cb
	}
	if cb.full {
		oldest := cb.time[cb.index]
		if now.Sub(oldest).Seconds() < RateLimitWindow {
			return false
		}
	}
	cb.time[cb.index] = now
	cb.index = (cb.index + 1) % RateLimitCount
	if cb.index == 0 {
		cb.full = true
	}
	return true
}
