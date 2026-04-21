# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD_2_2CloseU
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_27

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Channel: Receive
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:203
```go
192 ...
193 
194 		m.Lock()
195 		close(c) // close in CS
196 		m.Unlock()
197 	}
198 
199 	receiver := func() {
200 		m.Lock()
201 		time.Sleep(10 * time.Millisecond)
202 		m.Unlock()
203 		<-c // recv with PCS           // <-------
204 	}
205 
206 	run2(receiver, closer)
207 }
208 
209 // MD-CloseB: Buffered Variant (Mirror of MD-2-2B)
210 func TestMixedDeadlock_MD_2_2CloseB(t *testing.T) {
211 	var m sync.Mutex
212 	c := make(chan int, 1) // unbuffered
213 
214 
215 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:195
```go
184 ...
185 
186 
187 // MD-CloseU: Unbuffered Variant (Mirror of MD-2-2U)
188 func TestMixedDeadlock_MD_2_2CloseU(t *testing.T) {
189 	var m sync.Mutex
190 	c := make(chan int) // unbuffered
191 
192 	closer := func() {
193 		time.Sleep(50 * time.Millisecond) // let receiver finish first
194 		m.Lock()
195 		close(c) // close in CS           // <-------
196 		m.Unlock()
197 	}
198 
199 	receiver := func() {
200 		m.Lock()
201 		time.Sleep(10 * time.Millisecond)
202 		m.Unlock()
203 		<-c // recv with PCS
204 	}
205 
206 
207 ...
```


###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:200
```go
189 ...
190 
191 
192 	closer := func() {
193 		time.Sleep(50 * time.Millisecond) // let receiver finish first
194 		m.Lock()
195 		close(c) // close in CS
196 		m.Unlock()
197 	}
198 
199 	receiver := func() {
200 		m.Lock()           // <-------
201 		time.Sleep(10 * time.Millisecond)
202 		m.Unlock()
203 		<-c // recv with PCS
204 	}
205 
206 	run2(receiver, closer)
207 }
208 
209 // MD-CloseB: Buffered Variant (Mirror of MD-2-2B)
210 func TestMixedDeadlock_MD_2_2CloseB(t *testing.T) {
211 
212 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:194
```go
183 ...
184 
185 }
186 
187 // MD-CloseU: Unbuffered Variant (Mirror of MD-2-2U)
188 func TestMixedDeadlock_MD_2_2CloseU(t *testing.T) {
189 	var m sync.Mutex
190 	c := make(chan int) // unbuffered
191 
192 	closer := func() {
193 		time.Sleep(50 * time.Millisecond) // let receiver finish first
194 		m.Lock()           // <-------
195 		close(c) // close in CS
196 		m.Unlock()
197 	}
198 
199 	receiver := func() {
200 		m.Lock()
201 		time.Sleep(10 * time.Millisecond)
202 		m.Unlock()
203 		<-c // recv with PCS
204 	}
205 
206 ...
```


## Replay
**Replaying was not run**.

