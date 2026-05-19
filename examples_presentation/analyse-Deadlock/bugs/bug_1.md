# Bug: P05 - Possible Cyclic Deadlock

The analysis detected a possible cyclic deadlock.
If this deadlock contains or influences the run of the main routine, this can result in the program getting stuck. Otherwise it can lead to an unnecessary use of resources.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestDeadlock
- File: /workspaces/seminararbeit-go-concurrency/examples_presentation/Deadlock_test.go
- Trace: advocateTrace_1

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Mutex: Causing deadlock
-> /workspaces/seminararbeit-go-concurrency/examples_presentation/Deadlock_test.go#17
```go
6 ...
7 
8 )
9 
10 var lock1 sync.Mutex
11 var lock2 sync.Mutex
12 
13 func TestDeadlock(t *testing.T) {
14 
15 	go func() {
16 		lock1.Lock()
17 		lock2.Lock()           // <-------
18 
19 		fmt.Println("goroutine 1 working")
20 
21 		lock2.Unlock()
22 		lock1.Unlock()
23 
24 		fmt.Println("goroutine 1 done")
25 	}()
26 
27 	go func() {
28 
29 ...
```


###  Mutex: Part of deadlock
-> /workspaces/seminararbeit-go-concurrency/examples_presentation/Deadlock_test.go#29
```go
18 ...
19 
20 
21 		lock2.Unlock()
22 		lock1.Unlock()
23 
24 		fmt.Println("goroutine 1 done")
25 	}()
26 
27 	go func() {
28 		lock2.Lock()
29 		lock1.Lock()           // <-------
30 
31 		fmt.Println("goroutine 2 working")
32 
33 		lock1.Unlock()
34 		lock2.Unlock()
35 
36 		fmt.Println("goroutine 2 done")
37 	}()
38 
39 	time.Sleep(2 * time.Second)
40 
41 ...
```


## Replay
**Replaying was not run**.

