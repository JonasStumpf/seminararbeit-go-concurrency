# seminararbeit-go-concurrency
Seminararbeit Go Concurrency

[Advocate](https://github.com/ErikKassubek/ADVOCATE/tree/main)


Running advocate in this container needs full paths.  
e.g. running an analysis on Basic_lockdep_test from the advocate deadlocks example:
```bash
/workspaces/seminararbeit-go-concurrency/ADVOCATE/advocate/advocate analysis -path /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/ -exec TestBasicLockdep
```
e.g. recording of BasicDeadlock_test:
```bash
/workspaces/seminararbeit-go-concurrency/ADVOCATE/advocate/advocate record -path /workspaces/seminararbeit-go-concurrency/examples -exec TestBasicDeadlock
```
e.g. fuzzing of BasicDeadlock_test:
```bash
/workspaces/seminararbeit-go-concurrency/ADVOCATE/advocate/advocate fuzzing -path /workspaces/seminararbeit-go-concurrency/examples -exec TestBasicDeadlock -mode Flow
```

# Concurrency in Go
To run concurrent execution in Go, add the `go` keyword before a function call. This will execute the function in a new goroutine, allowing it to run concurrently with the rest of the program.
```go
func main() {
    go func() {
        fmt.Println("Hello from a goroutine!")
    }()
    fmt.Println("Hello from the main function!")
}
```
## Communication between goroutines
Using channels, you can send and receive values. Sends and Receives are blocking by default, waiting until the other side is ready, allowing for synchronization between goroutines. Channels can be closed to signal that no more values will be sent. By adding a second parameter to the receive operation, you can check if the channel is closed.
```go
func main() {
    ch := make(chan string)
    
    go func() {
        ch <- "Hello from a goroutine!"
        close(ch)
    }()
    
    msg, ok := <-ch
	fmt.Println(ok, msg) // true Hello from a goroutine! 
	
	msg, ok = <-ch
	fmt.Println(ok, msg) // false
}
```
## Mutex
Mutexes exist in Go, e.g. if no communication is needed.
```go
var mu sync.Mutex
func main() {
    mu.Lock()
    // critical section
    mu.Unlock()
}
```

# Advocate
[Advocate](https://github.com/ErikKassubek/ADVOCATE/tree/main) tries to identify potential concurrency issues.  
All bugs it tries to detect can be found [here](https://github.com/ErikKassubek/ADVOCATE/tree/main#what-is-advocatego).


> **Note**: Line numbers appear to be wrong in traces and results, probably because advocate adds lines before execution when the toolchain (commandline) is used:
> - adds "advocate" in import
> - adds at start of test/main method:
> ```go
> // ======= Preamble Start =======
>  advocate.InitTracing()
>  defer advocate.FinishTracing()
>// ======= Preamble End =======
> ```
> This sums up to 5 extra lines. Subtracting these from the reported line numbers seems to give the correct line numbers.

## Record
Records operations of a run of a program/test and produces a trace. The trace can then be used to replay the run deterministically.  
All operations and their formats can be found [here](https://github.com/ErikKassubek/ADVOCATE/blob/main/doc/recording.md).

Here is an Example for `example/BasicDeadlock_test`, the traces can be found in [examples/advRes_BasicDeadlock/advocateTrace](examples/advRes_BasicDeadlock/advocateTrace/).  
Example trace_1.log:
```
G,2,2,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:16
M,4,8,1000000001,f,L,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:23
M,10,14,1000000002,f,L,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:24
M,16,20,1000000002,f,U,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:25
M,22,26,1000000001,f,U,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:26
E,54
```

Definitions of the operations in the trace:
- `G`: GoRoutine (creation), Format:
  > G,[tPost],[id],[pos]
- `M`: Mutex, Format:
  > M,[tPre],[tPost],[id],[rw],[opM],[suc],[pos]
- `E`: End of routine, Format:
  > E,[t]

| Field | Description |
|-------|-------------|
| t | time, value of global counter at the end of the routine |
| tPre | time, value of global counter before the event |
| tPost | time, value of global counter after/at the event |
| id | id of element |
| pos | position in program (code) of element |
| rw | t/f, true if mutex is rw mutex |
| opM | operation of element |
| suc | t/f, true if lock was acquired |

So the above trace can be interpreted as:
1. GoRoutine created (line 11)
2. Mutex 1: lock acquired (line 18)
3. Mutex 2: lock acquired (line 19)
4. Mutex 2: unlock (line 20)
5. Mutex 1: unlock (line 21)
6. End of routine

The second trace:
```
M,28,32,1000000002,f,L,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:17
M,34,38,1000000001,f,L,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:18
M,40,44,1000000001,f,U,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:19
M,46,50,1000000002,f,U,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:20
E,52
```
1. Mutex 2: lock acquired (line 12)
2. Mutex 1: lock acquired (line 13)
3. Mutex 1: unlock (line 14)
4. Mutex 2: unlock (line 15)
5. End of routine

If you look at the timestamps and combine the traces, you see that the execution of the GoRoutine (second trace) happens after the main routine (first trace) has acquired and released both locks. Here is the code commented with the order of execution from the traces:
```go
8   func TestBasicDeadlock(t *testing.T) {
9   	var x, y sync.Mutex
10
11  	go func() { // 1. routine created
12  		x.Lock() // 6. mutex 2 lock acquired
13  		y.Lock() // 7. mutex 1 lock acquired
14  		y.Unlock() // 8. mutex 1 unlock
15  		x.Unlock() // 9. mutex 2 unlock
16  	}()
17
18      y.Lock() // 2. mutex 1 lock acquired
19      x.Lock() // 3. mutex 2 lock acquired
20      x.Unlock() // 4.  mutex 2 unlock
21      y.Unlock() // 5. mutex 1 unlock
22  }
```

## Replay
Runs a program/test according to a trace. This allows for deterministic execution.

## Analysis
Advocate runs the program/test and records the trace. It analyzes the trace and tries to find bugs and then rewrites the traces to produce those bugs. Then it replays the rewritten traces to check if the bugs exist.

Following are some examples of the results of running advocate on some example programs.

Format in result info is: `path:line@timestamp`  
timestamp is the value of the vector clock at the time of the event.

### Advocate deadlock example Basic_test results (also copied to examples/BasicDeadlock_test):
Results in results_readable.log:
```
==================== Summary ====================

-------------------- Critical -------------------

1 Possible cyclic deadlock:tp:
	stuck: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/Basic_test.go:24@10
	cycle: /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/Basic_test.go:18@38;/workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/Basic_test.go:24@10

```
Stuck 24 -> 19
Block cycle: line 24 <-> line 18 -> line 19 <-> line 13
```go
8   func TestBasic(t *testing.T) {
9       var x, y sync.Mutex
10
11      go func() {
12          x.Lock()
13          y.Lock() // <--
14          y.Unlock()
15          x.Unlock()
16      }()
17
18      y.Lock()
19      x.Lock() // <--
20      x.Unlock()
21      y.Unlock()
22  }
```
Sometimes a bug report markdown is created, see for this [Example](advocateExampleResults/file%287%29-test%280%29-Basic_test-TestBasic/bugs/bug_1.md).


### Advocate example without a Bug, Guard_test:
```
==================== Summary ====================

No bugs found
```

### Example, SendClosedChannel_test:
```
==================== Summary ====================

-------------------- Critical -------------------

1 Actual Send on Closed Channel:tp:
	send: /workspaces/seminararbeit-go-concurrency/examples/SendClosedChannel_test.go:23@0


2 Actual Send on Closed Channel:tp:
	send: /workspaces/seminararbeit-go-concurrency/examples/SendClosedChannel_test.go:23@12
	close: /workspaces/seminararbeit-go-concurrency/examples/SendClosedChannel_test.go:18@8

3 Block on routine or unknown element:tp:
	elem: 


```

It shows where the channel is closed and where the send on the closed channel happens.
(Block might be detected because there is no receive to the send operation).
Send: Line 23 -> line 18
Closed: Line 18 -> line 13

```go
8   func TestSendOnClosedChannel(t *testing.T) {
9	    ch := make(chan int)
10
11	    go func() {
12		    time.Sleep(100 * time.Millisecond)
13		    close(ch) // <--
14	    }()
15
16	    go func() {
17		    time.Sleep(200 * time.Millisecond)
18		    ch <- 42 // <--     // should send on closed channel
19    	}()
20
21	    time.Sleep(1 * time.Second)
22  }
```


## Fuzzing
If part of the code isn't executed, a potential bug in that section can't be found. With fuzzing, a program/test is executed multiple times to increase the chance of executing all code and finding all potential bugs. Advocate tries to influence execution to increase the chance of executing new program paths.

Advocate has multiple [Fuzzing-Modes](https://github.com/ErikKassubek/ADVOCATE/blob/main/doc/usage.md#mode-fuzzing).

> **Note**: Documentation says the mode is set with `-fuzzingMode [mode]` but the flag is actually `-mode [mode]`.

The example `example/MutexFuzzing_test` (copied and timeout changed from `advocate/doc/examples/fuzzing/Flow/mutex`) shows no bugs with `analysis` but with `fuzzing` (mode: Flow) it finds the bugs.  
See the results for fuzzing in [examples/fuzzingResults/advRes_MutexFuzzing_fuzzing](examples/fuzzingResults/advRes_MutexFuzzing_fuzzing/) and for analysis in [examples/fuzzingResults/advRes_MutexFuzzing_analysis](examples/fuzzingResults/advRes_MutexFuzzing_analysis/).  
You can see the multiple runs (2 in this case) in `total_results_readable.log`, at first it found no bugs and in the second run it found the bugs.

Example `example/FuzzingPaths_test`:  
Testing the fuzzing modes Flow, GoPie and GFuzz. They don't always find all bugs, Flow even found no bugs on a few tries (same for analysis).

