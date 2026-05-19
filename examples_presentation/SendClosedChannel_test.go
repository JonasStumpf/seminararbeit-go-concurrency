package main

import (
	"fmt"
	"testing"
	"time"
)

func TestSendOnClosedChannel(t *testing.T) {
	ch := make(chan int)

	go func() {
		time.Sleep(20 * time.Millisecond)
		ch <- 42
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(ch)
	}()

	v := <-ch
	fmt.Println("received:", v) // Outputs: received: 42

	time.Sleep(100 * time.Millisecond)
}
