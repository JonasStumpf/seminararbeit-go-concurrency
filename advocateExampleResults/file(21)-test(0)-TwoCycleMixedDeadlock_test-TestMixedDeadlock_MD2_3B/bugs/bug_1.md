# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD2_3B
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_25

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:146
```go
135 ...
136 
137 	run2(sender, receiver)
138 }
139 
140 // MD2-3B: Buffered Variant
141 func TestMixedDeadlock_MD2_3B(t *testing.T) {
142 	var m sync.Mutex
143 	c := make(chan int, 1)
144 
145 	sender := func() {
146 		m.Lock()           // <-------
147 		time.Sleep(10 * time.Millisecond)
148 		m.Unlock()
149 		c <- 1 // send after PCS
150 	}
151 
152 	receiver := func() {
153 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
154 		m.Lock()
155 		<-c // receive inside CS
156 		m.Unlock()
157 
158 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:154
```go
143 ...
144 
145 	sender := func() {
146 		m.Lock()
147 		time.Sleep(10 * time.Millisecond)
148 		m.Unlock()
149 		c <- 1 // send after PCS
150 	}
151 
152 	receiver := func() {
153 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
154 		m.Lock()           // <-------
155 		<-c // receive inside CS
156 		m.Unlock()
157 	}
158 
159 	run2(sender, receiver)
160 }
161 
162 // ------------------------------------------------------------
163 // CLOSE TESTS: MDX-Y-CLOSE VARIANTS
164 // ------------------------------------------------------------
165 
166 ...
```


###  Channel: Send
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:149
```go
138 ...
139 
140 // MD2-3B: Buffered Variant
141 func TestMixedDeadlock_MD2_3B(t *testing.T) {
142 	var m sync.Mutex
143 	c := make(chan int, 1)
144 
145 	sender := func() {
146 		m.Lock()
147 		time.Sleep(10 * time.Millisecond)
148 		m.Unlock()
149 		c <- 1 // send after PCS           // <-------
150 	}
151 
152 	receiver := func() {
153 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
154 		m.Lock()
155 		<-c // receive inside CS
156 		m.Unlock()
157 	}
158 
159 	run2(sender, receiver)
160 
161 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:155
```go
144 ...
145 
146 		m.Lock()
147 		time.Sleep(10 * time.Millisecond)
148 		m.Unlock()
149 		c <- 1 // send after PCS
150 	}
151 
152 	receiver := func() {
153 		time.Sleep(50 * time.Millisecond) // sleep to let sender complete PCS
154 		m.Lock()
155 		<-c // receive inside CS           // <-------
156 		m.Unlock()
157 	}
158 
159 	run2(sender, receiver)
160 }
161 
162 // ------------------------------------------------------------
163 // CLOSE TESTS: MDX-Y-CLOSE VARIANTS
164 // ------------------------------------------------------------
165 
166 
167 ...
```


## Replay
**Replaying was not run**.

