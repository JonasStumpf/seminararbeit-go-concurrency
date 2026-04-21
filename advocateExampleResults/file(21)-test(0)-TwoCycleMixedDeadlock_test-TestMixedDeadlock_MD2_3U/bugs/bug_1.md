# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD2_3U
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_24

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:124
```go
113 ...
114 
115 // MD2-3: Sender with PCS, Receiver inside CS
116 // ------------------------------------------------------------
117 
118 // MD2-3U: Unbuffered Variant
119 func TestMixedDeadlock_MD2_3U(t *testing.T) {
120 	var m sync.Mutex
121 	c := make(chan int)
122 
123 	sender := func() {
124 		m.Lock()           // <-------
125 		time.Sleep(10 * time.Millisecond)
126 		m.Unlock()
127 		c <- 1 // send after PCS
128 	}
129 
130 	receiver := func() {
131 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
132 		m.Lock()
133 		<-c // receive inside CS
134 		m.Unlock()
135 
136 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:132
```go
121 ...
122 
123 	sender := func() {
124 		m.Lock()
125 		time.Sleep(10 * time.Millisecond)
126 		m.Unlock()
127 		c <- 1 // send after PCS
128 	}
129 
130 	receiver := func() {
131 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
132 		m.Lock()           // <-------
133 		<-c // receive inside CS
134 		m.Unlock()
135 	}
136 
137 	run2(sender, receiver)
138 }
139 
140 // MD2-3B: Buffered Variant
141 func TestMixedDeadlock_MD2_3B(t *testing.T) {
142 	var m sync.Mutex
143 
144 ...
```


###  Channel: Send
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:127
```go
116 ...
117 
118 // MD2-3U: Unbuffered Variant
119 func TestMixedDeadlock_MD2_3U(t *testing.T) {
120 	var m sync.Mutex
121 	c := make(chan int)
122 
123 	sender := func() {
124 		m.Lock()
125 		time.Sleep(10 * time.Millisecond)
126 		m.Unlock()
127 		c <- 1 // send after PCS           // <-------
128 	}
129 
130 	receiver := func() {
131 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
132 		m.Lock()
133 		<-c // receive inside CS
134 		m.Unlock()
135 	}
136 
137 	run2(sender, receiver)
138 
139 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:133
```go
122 ...
123 
124 		m.Lock()
125 		time.Sleep(10 * time.Millisecond)
126 		m.Unlock()
127 		c <- 1 // send after PCS
128 	}
129 
130 	receiver := func() {
131 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
132 		m.Lock()
133 		<-c // receive inside CS           // <-------
134 		m.Unlock()
135 	}
136 
137 	run2(sender, receiver)
138 }
139 
140 // MD2-3B: Buffered Variant
141 func TestMixedDeadlock_MD2_3B(t *testing.T) {
142 	var m sync.Mutex
143 	c := make(chan int, 1)
144 
145 ...
```


## Replay
**Replaying was not run**.

