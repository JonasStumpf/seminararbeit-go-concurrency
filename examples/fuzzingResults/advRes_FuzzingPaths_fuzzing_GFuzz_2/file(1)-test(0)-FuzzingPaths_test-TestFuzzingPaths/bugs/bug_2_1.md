# Bug: A06 - Actual Unlock of Not Locked Mutex

During the execution, a not locked mutex was unlocked.
The occurrence of this lead to a panic.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestFuzzingPaths
- File: /workspaces/seminararbeit-go-concurrency/examples/FuzzingPaths_test.go
- Trace: advocateTrace_3

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Lock
-> /workspaces/seminararbeit-go-concurrency/examples/FuzzingPaths_test.go:37
```go
26 ...
27 
28 
29 	// Goroutine 3: Mutex misused
30 	go func() {
31 		time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)
32 
33 		if time.Now().UnixNano()%2 == 0 {
34 			mu.Lock()
35 			defer mu.Unlock()
36 		} else {
37 			mu.Unlock() // unlock without lock → panic / potential bug           // <-------
38 		}
39 	}()
40 
41 	// Goroutine 4: potential leak
42 	go func() {
43 		select {
44 		case <-ch: // ok
45 		case <-time.After(200 * time.Millisecond): // timeout → potential leak
46 		}
47 	}()
48 
49 ...
```


## Replay
**Replaying was not run**.

