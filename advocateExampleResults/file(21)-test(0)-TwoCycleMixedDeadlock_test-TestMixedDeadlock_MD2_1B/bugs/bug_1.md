# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD2_1B
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_21

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Channel: Send
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:52
```go
41 ...
42 
43 // ------------------------------------------------------------
44 
45 // MD2-1B: Buffered Variant
46 func TestMixedDeadlock_MD2_1B(t *testing.T) {
47 	var m sync.Mutex
48 	c := make(chan int, 1) // buffered
49 
50 	sender := func() {
51 		m.Lock()
52 		c <- 1 // send inside CS           // <-------
53 		m.Unlock()
54 	}
55 
56 	receiver := func() {
57 		time.Sleep(50 * time.Millisecond)
58 		m.Lock()
59 		<-c // receive inside CS
60 		m.Unlock()
61 	}
62 
63 
64 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:59
```go
48 ...
49 
50 	sender := func() {
51 		m.Lock()
52 		c <- 1 // send inside CS
53 		m.Unlock()
54 	}
55 
56 	receiver := func() {
57 		time.Sleep(50 * time.Millisecond)
58 		m.Lock()
59 		<-c // receive inside CS           // <-------
60 		m.Unlock()
61 	}
62 
63 	run2(sender, receiver)
64 }
65 
66 // ------------------------------------------------------------
67 // MD2-2: Sender inside CS, Receiver with PCS
68 // ------------------------------------------------------------
69 
70 
71 ...
```


###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:51
```go
40 ...
41 
42 // MD2-1: Both sender and receiver in CS
43 // ------------------------------------------------------------
44 
45 // MD2-1B: Buffered Variant
46 func TestMixedDeadlock_MD2_1B(t *testing.T) {
47 	var m sync.Mutex
48 	c := make(chan int, 1) // buffered
49 
50 	sender := func() {
51 		m.Lock()           // <-------
52 		c <- 1 // send inside CS
53 		m.Unlock()
54 	}
55 
56 	receiver := func() {
57 		time.Sleep(50 * time.Millisecond)
58 		m.Lock()
59 		<-c // receive inside CS
60 		m.Unlock()
61 	}
62 
63 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:58
```go
47 ...
48 
49 
50 	sender := func() {
51 		m.Lock()
52 		c <- 1 // send inside CS
53 		m.Unlock()
54 	}
55 
56 	receiver := func() {
57 		time.Sleep(50 * time.Millisecond)
58 		m.Lock()           // <-------
59 		<-c // receive inside CS
60 		m.Unlock()
61 	}
62 
63 	run2(sender, receiver)
64 }
65 
66 // ------------------------------------------------------------
67 // MD2-2: Sender inside CS, Receiver with PCS
68 // ------------------------------------------------------------
69 
70 ...
```


## Replay
**Replaying was not run**.

