# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD_2_2B_RW
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_34

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:359
```go
348 ...
349 
350 
351 	writer := func() {
352 		time.Sleep(50 * time.Millisecond) // let receiver finish PCS
353 		rw.Lock()
354 		c <- 1 // send inside CS
355 		rw.Unlock()
356 	}
357 
358 	reader := func() {
359 		rw.RLock()           // <-------
360 		time.Sleep(10 * time.Millisecond)
361 		rw.RUnlock() // PCS
362 		<-c          // receive after PCS
363 	}
364 
365 	run2(reader, writer)
366 }
367 
368 // READ/WRTIE MD2-3U
369 func TestMixedDeadlock_MD_2_3_RW(t *testing.T) {
370 
371 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:353
```go
342 ...
343 
344 }
345 
346 // READ/WRTIE MD2-2B
347 func TestMixedDeadlock_MD_2_2B_RW(t *testing.T) {
348 	var rw sync.RWMutex
349 	c := make(chan int, 1)
350 
351 	writer := func() {
352 		time.Sleep(50 * time.Millisecond) // let receiver finish PCS
353 		rw.Lock()           // <-------
354 		c <- 1 // send inside CS
355 		rw.Unlock()
356 	}
357 
358 	reader := func() {
359 		rw.RLock()
360 		time.Sleep(10 * time.Millisecond)
361 		rw.RUnlock() // PCS
362 		<-c          // receive after PCS
363 	}
364 
365 ...
```


###  Channel: Receive
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:362
```go
351 ...
352 
353 		rw.Lock()
354 		c <- 1 // send inside CS
355 		rw.Unlock()
356 	}
357 
358 	reader := func() {
359 		rw.RLock()
360 		time.Sleep(10 * time.Millisecond)
361 		rw.RUnlock() // PCS
362 		<-c          // receive after PCS           // <-------
363 	}
364 
365 	run2(reader, writer)
366 }
367 
368 // READ/WRTIE MD2-3U
369 func TestMixedDeadlock_MD_2_3_RW(t *testing.T) {
370 	var rw sync.RWMutex
371 	c := make(chan int)
372 
373 
374 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:354
```go
343 ...
344 
345 
346 // READ/WRTIE MD2-2B
347 func TestMixedDeadlock_MD_2_2B_RW(t *testing.T) {
348 	var rw sync.RWMutex
349 	c := make(chan int, 1)
350 
351 	writer := func() {
352 		time.Sleep(50 * time.Millisecond) // let receiver finish PCS
353 		rw.Lock()
354 		c <- 1 // send inside CS           // <-------
355 		rw.Unlock()
356 	}
357 
358 	reader := func() {
359 		rw.RLock()
360 		time.Sleep(10 * time.Millisecond)
361 		rw.RUnlock() // PCS
362 		<-c          // receive after PCS
363 	}
364 
365 
366 ...
```


## Replay
**Replaying was not run**.

