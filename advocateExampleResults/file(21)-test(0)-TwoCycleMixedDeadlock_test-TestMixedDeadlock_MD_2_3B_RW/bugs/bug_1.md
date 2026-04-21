# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD_2_3B_RW
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_36

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:403
```go
392 ...
393 
394 
395 	reader := func() {
396 		time.Sleep(50 * time.Millisecond) // let sender finish PCS
397 		rw.RLock()
398 		<-c // receive inside CS
399 		rw.RUnlock()
400 	}
401 
402 	writer := func() {
403 		rw.Lock()           // <-------
404 		time.Sleep(10 * time.Millisecond)
405 		rw.Unlock() // PCS
406 		c <- 1      // send after PCS
407 	}
408 
409 	run2(reader, writer)
410 }
411 
412 // ------------------------------------------------------------
413 // FALSE POSITIVE TESTS
414 
415 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:397
```go
386 ...
387 
388 }
389 
390 // READ/WRTIE MD2-3B
391 func TestMixedDeadlock_MD_2_3B_RW(t *testing.T) {
392 	var rw sync.RWMutex
393 	c := make(chan int, 1)
394 
395 	reader := func() {
396 		time.Sleep(50 * time.Millisecond) // let sender finish PCS
397 		rw.RLock()           // <-------
398 		<-c // receive inside CS
399 		rw.RUnlock()
400 	}
401 
402 	writer := func() {
403 		rw.Lock()
404 		time.Sleep(10 * time.Millisecond)
405 		rw.Unlock() // PCS
406 		c <- 1      // send after PCS
407 	}
408 
409 ...
```


###  Channel: Send
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:406
```go
395 ...
396 
397 		rw.RLock()
398 		<-c // receive inside CS
399 		rw.RUnlock()
400 	}
401 
402 	writer := func() {
403 		rw.Lock()
404 		time.Sleep(10 * time.Millisecond)
405 		rw.Unlock() // PCS
406 		c <- 1      // send after PCS           // <-------
407 	}
408 
409 	run2(reader, writer)
410 }
411 
412 // ------------------------------------------------------------
413 // FALSE POSITIVE TESTS
414 // ------------------------------------------------------------
415 
416 func TestMixedDeadlock_No_MD_BeforeCS(t *testing.T) {
417 
418 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:398
```go
387 ...
388 
389 
390 // READ/WRTIE MD2-3B
391 func TestMixedDeadlock_MD_2_3B_RW(t *testing.T) {
392 	var rw sync.RWMutex
393 	c := make(chan int, 1)
394 
395 	reader := func() {
396 		time.Sleep(50 * time.Millisecond) // let sender finish PCS
397 		rw.RLock()
398 		<-c // receive inside CS           // <-------
399 		rw.RUnlock()
400 	}
401 
402 	writer := func() {
403 		rw.Lock()
404 		time.Sleep(10 * time.Millisecond)
405 		rw.Unlock() // PCS
406 		c <- 1      // send after PCS
407 	}
408 
409 
410 ...
```


## Replay
**Replaying was not run**.

