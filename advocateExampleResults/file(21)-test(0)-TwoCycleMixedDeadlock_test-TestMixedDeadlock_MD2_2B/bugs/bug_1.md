# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD2_2B
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_23

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Channel: Receive
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:108
```go
97 ...
98 
99 		m.Lock()
100 		c <- 1 // buffered send inside CS
101 		m.Unlock()
102 	}
103 
104 	receiver := func() {
105 		m.Lock()
106 		time.Sleep(10 * time.Millisecond)
107 		m.Unlock()
108 		<-c // receive after PCS           // <-------
109 	}
110 
111 	run2(sender, receiver)
112 }
113 
114 // ------------------------------------------------------------
115 // MD2-3: Sender with PCS, Receiver inside CS
116 // ------------------------------------------------------------
117 
118 // MD2-3U: Unbuffered Variant
119 
120 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:100
```go
89 ...
90 
91 
92 // MD2-2B: Buffered Variant
93 func TestMixedDeadlock_MD2_2B(t *testing.T) {
94 	var m sync.Mutex
95 	c := make(chan int, 1)
96 
97 	sender := func() {
98 		time.Sleep(50 * time.Millisecond) // sleep to let receiver complete PCS (not necessary, always non-blocking)
99 		m.Lock()
100 		c <- 1 // buffered send inside CS           // <-------
101 		m.Unlock()
102 	}
103 
104 	receiver := func() {
105 		m.Lock()
106 		time.Sleep(10 * time.Millisecond)
107 		m.Unlock()
108 		<-c // receive after PCS
109 	}
110 
111 
112 ...
```


###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:105
```go
94 ...
95 
96 
97 	sender := func() {
98 		time.Sleep(50 * time.Millisecond) // sleep to let receiver complete PCS (not necessary, always non-blocking)
99 		m.Lock()
100 		c <- 1 // buffered send inside CS
101 		m.Unlock()
102 	}
103 
104 	receiver := func() {
105 		m.Lock()           // <-------
106 		time.Sleep(10 * time.Millisecond)
107 		m.Unlock()
108 		<-c // receive after PCS
109 	}
110 
111 	run2(sender, receiver)
112 }
113 
114 // ------------------------------------------------------------
115 // MD2-3: Sender with PCS, Receiver inside CS
116 
117 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:99
```go
88 ...
89 
90 }
91 
92 // MD2-2B: Buffered Variant
93 func TestMixedDeadlock_MD2_2B(t *testing.T) {
94 	var m sync.Mutex
95 	c := make(chan int, 1)
96 
97 	sender := func() {
98 		time.Sleep(50 * time.Millisecond) // sleep to let receiver complete PCS (not necessary, always non-blocking)
99 		m.Lock()           // <-------
100 		c <- 1 // buffered send inside CS
101 		m.Unlock()
102 	}
103 
104 	receiver := func() {
105 		m.Lock()
106 		time.Sleep(10 * time.Millisecond)
107 		m.Unlock()
108 		<-c // receive after PCS
109 	}
110 
111 ...
```


## Replay
**Replaying was not run**.

