# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD_2_3CloseU
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_29

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:237
```go
226 ...
227 
228 	run2(receiver, closer)
229 }
230 
231 // MD-CloseU: Unbuffered Variant (Mirror of MD-2-3U)
232 func TestMixedDeadlock_MD_2_3CloseU(t *testing.T) {
233 	var m sync.Mutex
234 	c := make(chan int) // unbuffered
235 
236 	closer := func() {
237 		m.Lock()           // <-------
238 		time.Sleep(10 * time.Millisecond)
239 		m.Unlock()
240 		close(c) // close after PCS
241 	}
242 
243 	receiver := func() {
244 		time.Sleep(50 * time.Millisecond)
245 		m.Lock()
246 		<-c // receive inside CS (blocked until close)
247 		m.Unlock()
248 
249 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:245
```go
234 ...
235 
236 	closer := func() {
237 		m.Lock()
238 		time.Sleep(10 * time.Millisecond)
239 		m.Unlock()
240 		close(c) // close after PCS
241 	}
242 
243 	receiver := func() {
244 		time.Sleep(50 * time.Millisecond)
245 		m.Lock()           // <-------
246 		<-c // receive inside CS (blocked until close)
247 		m.Unlock()
248 	}
249 
250 	run2(receiver, closer)
251 }
252 
253 // MD-CloseB: Buffered Variant (Mirror of MD-2-3B)
254 func TestMixedDeadlock_MD_2_3CloseB(t *testing.T) {
255 	var m sync.Mutex
256 
257 ...
```


###  Channel: Close
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:240
```go
229 ...
230 
231 // MD-CloseU: Unbuffered Variant (Mirror of MD-2-3U)
232 func TestMixedDeadlock_MD_2_3CloseU(t *testing.T) {
233 	var m sync.Mutex
234 	c := make(chan int) // unbuffered
235 
236 	closer := func() {
237 		m.Lock()
238 		time.Sleep(10 * time.Millisecond)
239 		m.Unlock()
240 		close(c) // close after PCS           // <-------
241 	}
242 
243 	receiver := func() {
244 		time.Sleep(50 * time.Millisecond)
245 		m.Lock()
246 		<-c // receive inside CS (blocked until close)
247 		m.Unlock()
248 	}
249 
250 	run2(receiver, closer)
251 
252 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:246
```go
235 ...
236 
237 		m.Lock()
238 		time.Sleep(10 * time.Millisecond)
239 		m.Unlock()
240 		close(c) // close after PCS
241 	}
242 
243 	receiver := func() {
244 		time.Sleep(50 * time.Millisecond)
245 		m.Lock()
246 		<-c // receive inside CS (blocked until close)           // <-------
247 		m.Unlock()
248 	}
249 
250 	run2(receiver, closer)
251 }
252 
253 // MD-CloseB: Buffered Variant (Mirror of MD-2-3B)
254 func TestMixedDeadlock_MD_2_3CloseB(t *testing.T) {
255 	var m sync.Mutex
256 	c := make(chan int, 1)
257 
258 ...
```


## Replay
**Replaying was not run**.

