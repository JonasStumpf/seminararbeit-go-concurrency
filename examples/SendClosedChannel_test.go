package main

import (
	"testing"
	"time"
)

func TestSendOnClosedChannel(t *testing.T) {
	ch := make(chan int)

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(ch)
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- 42 // should send on closed channel
	}()

	time.Sleep(1 * time.Second)
}
