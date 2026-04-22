# Leak: L07 - Block on select

The analyzer detected a Block on a select.
A Block on a select is a situation, where a select is still blocking at the end of the program.


## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestFuzzingPaths
- File: /workspaces/seminararbeit-go-concurrency/examples/FuzzingPaths_test.go
- Trace: advocateTrace_1

## Bug Elements
The elements involved in the found leak are located at the following positions:

###  Select:
-> /workspaces/seminararbeit-go-concurrency/examples/FuzzingPaths_test.go:18
```go
7 ...
8 
9 func TestFuzzingPaths(t *testing.T) {
10 	ch := make(chan int)
11 	var mu sync.Mutex
12 
13 	// Goroutine 1: arbeitet mit Channel
14 	go func() {
15 		time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)
16 
17 		select {
18 		case ch <- 1:           // <-------
19 			// normaler send
20 		default:
21 			// möglicherweise kein Receiver → alternative Pfad
22 		}
23 	}()
24 
25 	// Goroutine 2: schließt Channel (evtl. zu früh)
26 	go func() {
27 		time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)
28 		close(ch) // kann A01 triggern (wenn parallel send)
29 
30 ...
```


## Replay
**Replaying was not run**.

