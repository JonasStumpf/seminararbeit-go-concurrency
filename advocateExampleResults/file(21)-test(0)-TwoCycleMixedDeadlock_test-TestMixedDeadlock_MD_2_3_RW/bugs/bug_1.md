# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD_2_3_RW
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_35

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:381
```go
370 ...
371 
372 
373 	reader := func() {
374 		time.Sleep(50 * time.Millisecond) // let sender finish PCS
375 		rw.RLock()
376 		<-c // receive inside CS
377 		rw.RUnlock()
378 	}
379 
380 	writer := func() {
381 		rw.Lock()           // <-------
382 		time.Sleep(10 * time.Millisecond)
383 		rw.Unlock() // PCS
384 		c <- 1      // send after PCS
385 	}
386 
387 	run2(reader, writer)
388 }
389 
390 // READ/WRTIE MD2-3B
391 func TestMixedDeadlock_MD_2_3B_RW(t *testing.T) {
392 
393 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:375
```go
364 ...
365 
366 }
367 
368 // READ/WRTIE MD2-3U
369 func TestMixedDeadlock_MD_2_3_RW(t *testing.T) {
370 	var rw sync.RWMutex
371 	c := make(chan int)
372 
373 	reader := func() {
374 		time.Sleep(50 * time.Millisecond) // let sender finish PCS
375 		rw.RLock()           // <-------
376 		<-c // receive inside CS
377 		rw.RUnlock()
378 	}
379 
380 	writer := func() {
381 		rw.Lock()
382 		time.Sleep(10 * time.Millisecond)
383 		rw.Unlock() // PCS
384 		c <- 1      // send after PCS
385 	}
386 
387 ...
```


###  Channel: Send
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:384
```go
373 ...
374 
375 		rw.RLock()
376 		<-c // receive inside CS
377 		rw.RUnlock()
378 	}
379 
380 	writer := func() {
381 		rw.Lock()
382 		time.Sleep(10 * time.Millisecond)
383 		rw.Unlock() // PCS
384 		c <- 1      // send after PCS           // <-------
385 	}
386 
387 	run2(reader, writer)
388 }
389 
390 // READ/WRTIE MD2-3B
391 func TestMixedDeadlock_MD_2_3B_RW(t *testing.T) {
392 	var rw sync.RWMutex
393 	c := make(chan int, 1)
394 
395 
396 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:376
```go
365 ...
366 
367 
368 // READ/WRTIE MD2-3U
369 func TestMixedDeadlock_MD_2_3_RW(t *testing.T) {
370 	var rw sync.RWMutex
371 	c := make(chan int)
372 
373 	reader := func() {
374 		time.Sleep(50 * time.Millisecond) // let sender finish PCS
375 		rw.RLock()
376 		<-c // receive inside CS           // <-------
377 		rw.RUnlock()
378 	}
379 
380 	writer := func() {
381 		rw.Lock()
382 		time.Sleep(10 * time.Millisecond)
383 		rw.Unlock() // PCS
384 		c <- 1      // send after PCS
385 	}
386 
387 
388 ...
```


## Replay
**Replaying was not run**.

