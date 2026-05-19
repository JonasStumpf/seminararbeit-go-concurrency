package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var lock1 sync.Mutex
var lock2 sync.Mutex

func TestDeadlock(t *testing.T) {

	go func() {
		lock1.Lock()
		lock2.Lock()

		fmt.Println("goroutine 1 working")

		lock2.Unlock()
		lock1.Unlock()

		fmt.Println("goroutine 1 done")
	}()

	go func() {
		lock2.Lock()
		lock1.Lock()

		fmt.Println("goroutine 2 working")

		lock1.Unlock()
		lock2.Unlock()

		fmt.Println("goroutine 2 done")
	}()

	time.Sleep(2 * time.Second)
}
