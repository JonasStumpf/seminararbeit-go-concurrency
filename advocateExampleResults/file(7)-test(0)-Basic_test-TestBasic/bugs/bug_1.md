# Bug: P05 - Possible Cyclic Deadlock

The analysis detected a possible cyclic deadlock.
If this deadlock contains or influences the run of the main routine, this can result in the program getting stuck. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestBasic
- File: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/Basic_test.go
- Trace: advocateTrace_7

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Causing deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/Basic_test.go:19
```go
8 ...
9 
10 
11 	go func() {
12 		x.Lock()
13 		y.Lock()
14 		y.Unlock()
15 		x.Unlock()
16 	}()
17 
18 	y.Lock()
19 	x.Lock() // this SHOULD produce a deadlock           // <-------
20 	x.Unlock()
21 	y.Unlock()
22 }
23 
```


###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/Basic_test.go:13
```go
2 ...
3 
4 	"sync"
5 	"testing"
6 )
7 
8 func TestBasic(t *testing.T) {
9 	var x, y sync.Mutex
10 
11 	go func() {
12 		x.Lock()
13 		y.Lock()           // <-------
14 		y.Unlock()
15 		x.Unlock()
16 	}()
17 
18 	y.Lock()
19 	x.Lock() // this SHOULD produce a deadlock
20 	x.Unlock()
21 	y.Unlock()
22 }
23 
```


## Replay
The bug is a potential bug.
The analyzer has tried to rewrite the trace in such a way that the bug will be triggered when replaying the trace.

**Replaying confirmed the bug**.

It exited with the following code: 41

The replay reached the expected point and found stuck mutexes.The replay was therefore able to confirm that a deadlock can actually occur.

