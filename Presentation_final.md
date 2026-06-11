
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

Running the analysis gives a bug report, see [here](examples_presentation/analyse-SendClosedChannel/bugs/bug_1.md), where the bug is explained and the code parts involved in the bug are marked.  
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
