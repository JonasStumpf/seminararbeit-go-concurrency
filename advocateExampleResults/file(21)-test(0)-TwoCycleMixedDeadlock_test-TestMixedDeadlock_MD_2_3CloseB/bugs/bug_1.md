# Bug: P06 - Possible Mixed Deadlock

The analysis detected a Possible Mixed Deadlock.
A mixed deadlock is a situation, where two routines are blocked on each other, because they are waiting to send or receive on a channel, while holding locks that the other routine needs to proceed.
This can lead to the program getting stuck, if one of the routines is the main routine. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestMixedDeadlock_MD_2_3CloseB
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go
- Trace: advocateTrace_30

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:259
```go
248 ...
249 
250 	run2(receiver, closer)
251 }
252 
253 // MD-CloseB: Buffered Variant (Mirror of MD-2-3B)
254 func TestMixedDeadlock_MD_2_3CloseB(t *testing.T) {
255 	var m sync.Mutex
256 	c := make(chan int, 1)
257 
258 	closer := func() {
259 		m.Lock()           // <-------
260 		time.Sleep(10 * time.Millisecond)
261 		m.Unlock()
262 		close(c)
263 	}
264 
265 	receiver := func() {
266 		time.Sleep(50 * time.Millisecond)
267 		m.Lock()
268 		<-c // receive inside CS (blocked until close)
269 		m.Unlock()
270 
271 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:267
```go
256 ...
257 
258 	closer := func() {
259 		m.Lock()
260 		time.Sleep(10 * time.Millisecond)
261 		m.Unlock()
262 		close(c)
263 	}
264 
265 	receiver := func() {
266 		time.Sleep(50 * time.Millisecond)
267 		m.Lock()           // <-------
268 		<-c // receive inside CS (blocked until close)
269 		m.Unlock()
270 	}
271 
272 	run2(receiver, closer)
273 }
274 
275 // ------------------------------------------------------------
276 // LOCKTYPE TESTS: READ/READ
277 // ------------------------------------------------------------
278 
279 ...
```


###  Channel: Close
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:262
```go
251 ...
252 
253 // MD-CloseB: Buffered Variant (Mirror of MD-2-3B)
254 func TestMixedDeadlock_MD_2_3CloseB(t *testing.T) {
255 	var m sync.Mutex
256 	c := make(chan int, 1)
257 
258 	closer := func() {
259 		m.Lock()
260 		time.Sleep(10 * time.Millisecond)
261 		m.Unlock()
262 		close(c)           // <-------
263 	}
264 
265 	receiver := func() {
266 		time.Sleep(50 * time.Millisecond)
267 		m.Lock()
268 		<-c // receive inside CS (blocked until close)
269 		m.Unlock()
270 	}
271 
272 	run2(receiver, closer)
273 
274 ...
```


-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/TwoCycleMixedDeadlock_test.go:268
```go
257 ...
258 
259 		m.Lock()
260 		time.Sleep(10 * time.Millisecond)
261 		m.Unlock()
262 		close(c)
263 	}
264 
265 	receiver := func() {
266 		time.Sleep(50 * time.Millisecond)
267 		m.Lock()
268 		<-c // receive inside CS (blocked until close)           // <-------
269 		m.Unlock()
270 	}
271 
272 	run2(receiver, closer)
273 }
274 
275 // ------------------------------------------------------------
276 // LOCKTYPE TESTS: READ/READ
277 // ------------------------------------------------------------
278 
279 
280 ...
```


## Replay
**Replaying was not run**.

