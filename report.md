# AdvocateGO: Systematic Concurrency Testing for Go <!-- omit in toc -->
<p style="font-size:12px">Jonas Stumpf</p>

This seminar report is about the tool [AdvocateGO](https://github.com/ErikKassubek/ADVOCATE/tree/main) and how it can be used to find concurrency bugs in Go programs. The focus is on the `analysis` mode.

- [Concurrency in Go](#concurrency-in-go)
  - [Communication between goroutines](#communication-between-goroutines)
  - [Mutex](#mutex)
- [Advocate](#advocate)
  - [Record](#record)
  - [Replay](#replay)
  - [Analysis](#analysis)
    - [Send on Closed Channel](#send-on-closed-channel)
      - [Bug report markdown file](#bug-report-markdown-file)
    - [(Cyclic) Deadlock](#cyclic-deadlock)
    - [Example without a bug](#example-without-a-bug)
  - [Fuzzing](#fuzzing)
- [Conclusion](#conclusion)


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
_"AdvocateGo is an analysis tool for concurrent Go programs. It tries to detects concurrency bugs and gives diagnostic insight."_    - [AdvocateGO](https://github.com/ErikKassubek/ADVOCATE/tree/main)

All bugs it tries to detect can be found [here](https://github.com/ErikKassubek/ADVOCATE/tree/main#what-is-advocatego).

It has 4 Modes:
- Record: Records execution of a program and saves it as a trace file.
- Replay: Replays a trace file.
- Analysis: Record a program, checks if trace can be rewritten in a way that a concurrency bug arises, replay the rewritten trace and confirm the bug arises.
- Fuzzing: Apply fuzzing to increase coverage of the analysis

> **Note**: Line numbers appear to be wrong in traces and results, probably because advocate adds lines before execution when the toolchain (commandline) is used:
> - adds `advocate` in import
> - adds at start of test/main method:
> ```go
> // ======= Preamble Start =======
>  advocate.InitTracing()
>  defer advocate.FinishTracing()
>// ======= Preamble End =======
> ```
> This sums up to 5 extra lines. Subtracting these from the reported line numbers seems to give the correct line numbers.  
> Keep that in mind when looking at the results in this report. This might be fixed in the future.

## Record
Records operations of a run of a program/test and produces a trace. The trace can then be used to replay the run deterministically.  
All operations and their formats can be found [here](https://github.com/ErikKassubek/ADVOCATE/blob/main/doc/recording.md).

Here is an example trace for a simple deadlock scenario.
The code:
```go
func TestBasicDeadlock(t *testing.T) {
	var x, y sync.Mutex

	go func() {
		x.Lock()
		y.Lock()
		y.Unlock()
		x.Unlock()
	}()

	y.Lock()
	x.Lock() // this SHOULD produce a deadlock
	x.Unlock()
	y.Unlock()
}
```
Advocate creates different trace files foreach routine. The trace for the main routine looks like this:
```
G,2,2,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:16
M,4,8,1000000001,f,L,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:23
M,10,14,1000000002,f,L,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:24
M,16,20,1000000002,f,U,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:25
M,22,26,1000000001,f,U,t,/workspaces/seminararbeit-go-concurrency/examples/BasicDeadlock_test.go:26
E,54
```

Definitions of the operations in the trace are:
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

The trace for the GoRoutine looks like this:
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
The replay mode takes a trace and replays the execution according to the trace. It forces the same order of events as in the trace, so it can be used to reproduce bugs.

## Analysis
First a run of the program is recorded. Then the analysis checks if the events can occur in a different order that leads to a concurrency bug. If it finds such a reordering, it rewrites the trace and replays it to confirm that the bug arises.

Following are some examples that show how it works and what the results look like.

### Send on Closed Channel
The code for the example:  
[Go Playground](https://go.dev/play/p/CNs_RU8MuZI)
```go
func TestSendOnClosedChannel(t *testing.T) {
	ch := make(chan int)

	go func() {
		time.Sleep(20 * time.Millisecond) //potential delay/timeout, e.g. server request
		ch <- 42
	}()

	go func() {
		time.Sleep(50 * time.Millisecond) //potential delay/timeout, e.g. server request
		close(ch)
	}()

	v := <-ch
	fmt.Println("received:", v) // Outputs: received: 42

	time.Sleep(100 * time.Millisecond)
}
```
The output is: `received: 42`.  
But the artificial delays may be different in real scenarios.

The execution of the program and it's trace show a correct execution without a bug. The analysis now checks if the events can occur in a different order that leads to a concurrency bug. It does so via happens-before relations.

It checks for every send if there is a close that happens or could happen before the send. Here it sees that both events are in go routines and therefore it can be reordered, so the close happens before the send, which leads to the bug.

The results are logged in a machine-readable and human-readable format. The human-readable format looks like this (paths shortened):
```
==================== Summary ====================

-------------------- Critical -------------------

1 Possible send on closed channel:tp:
	send: .../SendClosedChannel_test.go:19@10
	close: .../SendClosedChannel_test.go:24@28

```
There you can see the actions and their positions in the code that can lead to the bug as well as their timestamps: `[type]: [file]:[line]@[tPre]`.  

Advocare also generates a bug report markdown file, showing the code parts involved and explaining the bug.
#### Bug report markdown file

<div style="border: 1px solid black; padding: 10px; margin: 10px;">

# Bug: P01 - Possible Send on Closed Channel <!-- omit in toc -->

The analyzer detected a possible send on a closed channel.
Although the send on a closed channel did not occur during the recording, it is possible that it will occur, based on the happens before relation.
Such a send on a closed channel leads to a panic.

## Test/Program <!-- omit in toc -->
The bug was found in the following test/program:

- Test/Prog: TestSendOnClosedChannel
- File: /workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go
- Trace: advocateTrace_1

## Bug Elements <!-- omit in toc -->
The elements involved in the found bug are located at the following positions:

###  Channel: Send <!-- omit in toc -->
-> /workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#14
```go
3 ...
4 
5 	"testing"
6 	"time"
7 )
8 
9 func TestSendOnClosedChannel(t *testing.T) {
10 	ch := make(chan int)
11 
12 	go func() {
13 		time.Sleep(20 * time.Millisecond)
14 		ch <- 42           // <-------
15 	}()
16 
17 	go func() {
18 		time.Sleep(50 * time.Millisecond)
19 		close(ch)
20 	}()
21 
22 	v := <-ch
23 	fmt.Println("received:", v)
24 
25 
26 ...
```
###  Channel: Close <!-- omit in toc -->
-> /workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#19
```go
8 ...
9 
10 	ch := make(chan int)
11 
12 	go func() {
13 		time.Sleep(20 * time.Millisecond)
14 		ch <- 42
15 	}()
16 
17 	go func() {
18 		time.Sleep(50 * time.Millisecond)
19 		close(ch)           // <-------
20 	}()
21 
22 	v := <-ch
23 	fmt.Println("received:", v)
24 
25 	time.Sleep(100 * time.Millisecond)
26 }
27 
```
## Replay <!-- omit in toc -->
**Replaying was not run**.

</div>

### (Cyclic) Deadlock
Advocate uses lock trees to detect circles (deadlocks) and uses Happens-Before relations to check if the operations are concurrent.

Advocate uses lock trees to detect circles (deadlocks) and uses Happens-Before relations to check if the operations are concurrent.

The code for the example:
```go
var lock1 sync.Mutex
var lock2 sync.Mutex

func TestDeadlock(t *testing.T) {

	go func() {
		lock1.Lock()
		lock2.Lock()

		fmt.Println("goroutine 1 working")

		lock2.Unlock()
		lock1.Unlock()

		fmt.Println("goroutine 1 done")
	}()

	go func() {
		lock2.Lock()
		lock1.Lock()

		fmt.Println("goroutine 2 working")

		lock1.Unlock()
		lock2.Unlock()

		fmt.Println("goroutine 2 done")
	}()

	time.Sleep(2 * time.Second)
}
```

The output can be empty (deadlock) or it can be:
```
goroutine 1 working
goroutine 1 done
goroutine 2 working
goroutine 2 done
```

Assuming the following trace:
```
T1			T2
lock1
lock2

unlock2
unlock1

			lock2
			lock1

			unlock1
			unlock2
```
The lock trees advocate creates are:
```
T1	T2
1	2
|	|
2	1
```
And the detected circle advocate finds:
```
T1	T2
1	2
| X |
2	1
```

Can this happen? Yes, because the lock operations are concurrent. If the first go routine aquires the first but not the second lock and then the secong go routine aquires their first lock, a deadlock occurs. Advocate detects this and reports a possible deadlock.

The human-readable logged summary looks like this (paths shortened):
```
==================== Summary ====================

-------------------- Critical -------------------

1 Possible cyclic deadlock:tp:
	stuck: .../Deadlock_test.go:22@54
	cycle: .../Deadlock_test.go:34@12;.../Deadlock_test.go:22@54

```


### Example without a bug
The logged results for an example without a bug look like this:
```
==================== Summary ====================

No bugs found
```


## Fuzzing
If part of the code isn't executed, a potential bug in that section can't be found. With fuzzing, a program/test is executed multiple times with different inputs and traces to increase the code coverage and the chance of finding all potential Bugs.

Advocate has multiple [Fuzzing-Modes](https://github.com/ErikKassubek/ADVOCATE/blob/main/doc/usage.md#mode-fuzzing).

> **Note**: Documentation says the mode is set with `-fuzzingMode [mode]` but the flag is actually `-mode [mode]`.


# Conclusion
The powerful record and replay modes allow you to record an execution and understand what is happening and then you may modify the trace and replay the program to check different executions.

With the analysis mode, Advocate can find potential concurrency bugs that may not arise in the recorded execution but could arise in a different execution.

Advocate detects potential concurrency bugs you might miss and gives you insights into the causes of the bugs.