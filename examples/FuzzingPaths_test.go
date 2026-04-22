package main

import (
	"sync"
	"testing"
	"time"
)

func TestFuzzingPaths(t *testing.T) {
	ch := make(chan int)
	var mu sync.Mutex

	// Goroutine 1: channel
	go func() {
		time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)

		select {
		case ch <- 1: // normal send
		default: // potentially no receiver
		}
	}()

	// Goroutine 2: closes channel (too early?)
	go func() {
		time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)
		close(ch) // might close too early
	}()

	// Goroutine 3: Mutex misused
	go func() {
		time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)

		if time.Now().UnixNano()%2 == 0 {
			mu.Lock()
			defer mu.Unlock()
		} else {
			mu.Unlock() // unlock without lock → panic / potential bug
		}
	}()

	// Goroutine 4: potential leak
	go func() {
		select {
		case <-ch: // ok
		case <-time.After(200 * time.Millisecond): // timeout → potential leak
		}
	}()

	time.Sleep(1 * time.Second)
}
