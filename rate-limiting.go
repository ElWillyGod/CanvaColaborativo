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
	RateLimitWindow = 10 * time.Second
)

// lista circular

type CircularBuffer struct {
	time  [RateLimitCount]time.Time
	index int
	full  bool
	mutex sync.Mutex
}

var userRates = make(map[net.Conn]*CircularBuffer)
var mu sync.RWMutex

func allowCommand(conn net.Conn) bool {
	mu.Lock()
	limiter, exists := userRates[conn]
	if !exists {
		limiter = &CircularBuffer{}
		userRates[conn] = limiter
	}
	mu.Unlock()

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	now := time.Now()

	if limiter.full {
		oldestRequest := limiter.time[limiter.index]
		if now.Sub(oldestRequest) < RateLimitWindow {
			return false
		}
	}

	limiter.time[limiter.index] = now
	limiter.index = (limiter.index + 1) % RateLimitCount
	if !limiter.full && limiter.index == 0 {
		limiter.full = true
	}

	return true
}

func removeLimiter(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	delete(userRates, conn)
}
