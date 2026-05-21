
# AdvocateGo

_"AdvocateGo is an analysis tool for concurrent Go programs. It tries to detects concurrency bugs and gives diagnostic insight."_   
[source](https://github.com/ErikKassubek/ADVOCATE/tree/main)

It has 4 Modes:
- Record: Records execution of a program and saves it as a trace file.
- Replay: Replays a trace file.
- Analysis: Record a program, checks if trace can be rewritten in a way that a concurrency bug arises, replay the rewritten trace and confirm the bug arises.
- Fuzzing: Apply fuzzing to increase coverage of the analysis


# Analysis
First a run of the program is recorded. Then the analysis checks if the events can occur in a different order that leads to a concurrency bug. If it finds such a reordering, it rewrites the trace and replays it to confirm that the bug arises.


## Send on Closed Channel Example
How does it work?
It checks for a close if there is an earlier send that could happen concurrently (send could occur after the close).

The code for the example:
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

Running the analysis gives a bug report, see [here](examples_presentation/advocateResult/analyse-SendClosedChannel/bugs/bug_1.md), where the bug is explained and the code parts involved in the bug are marked.  
It also gives the results in a shorter machine-readable and human-readable format (paths shortened):
```
==================== Summary ====================

-------------------- Critical -------------------

1 Possible send on closed channel:tp:
	send: .../SendClosedChannel_test.go:19@10
	close: .../SendClosedChannel_test.go:24@28

```
There you can see the actions and their positions in the code that lead to the bug as well as their timestamps: `[type]: [file]:[line]@[tPre]`.  
It shows that the send could possibly happen after the channel is closed, leading to a bug.

Looking at the recorded traces, we can see the following execution order:
```go
func TestSendOnClosedChannel(t *testing.T) {
	ch := make(chan int) 					// 1. channel created

	go func() { 							// 2. go routine created
		time.Sleep(20 * time.Millisecond)
		ch <- 42 							// 4. send on channel
	}()

	go func() { 							// 3. go routine created
		time.Sleep(50 * time.Millisecond)
		close(ch) 							// 6. close channel
	}()

	v := <-ch 								// 5. receive from channel
	fmt.Println("received:", v)

	time.Sleep(100 * time.Millisecond)
}
```
Advocate now checks if the send can happen after the close, which is possible.


## (Cyclic) Deadlock Example
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
The lock trees are:
```
T1	T2
1	2
|	|
2	1
```
And the detected circle:
```
T1	T2
1	2
| X |
2	1
```

Can this happen? Yes, because the lock operations are concurrent. Advocate detects this and reports a possible deadlock.

See the bug report [here](examples_presentation/analyse-Deadlock/bugs/bug_1.md).
And the human-readable summary (paths shortened):
```
==================== Summary ====================

-------------------- Critical -------------------

1 Possible cyclic deadlock:tp:
	stuck: .../Deadlock_test.go:22@54
	cycle: .../Deadlock_test.go:34@12;.../Deadlock_test.go:22@54

```
When looking at the code we see that we can run into a deadlock if both routines acquire their first lock:
```go
var lock1 sync.Mutex
var lock2 sync.Mutex

func TestDeadlock(t *testing.T) {

	go func() {
		lock1.Lock()
		lock2.Lock() // stuck - cycle

		fmt.Println("goroutine 1 working")

		lock2.Unlock()
		lock1.Unlock()

		fmt.Println("goroutine 1 done")
	}()

	go func() {
		lock2.Lock()
		lock1.Lock() // cycle

		fmt.Println("goroutine 2 working")

		lock1.Unlock()
		lock2.Unlock()

		fmt.Println("goroutine 2 done")
	}()

	time.Sleep(2 * time.Second)
}
```