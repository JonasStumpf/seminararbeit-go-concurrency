
# AdvocateGo

_"AdvocateGo is an analysis tool for concurrent Go programs. It tries to detects concurrency bugs and gives diagnostic insight."_   
[source](https://github.com/ErikKassubek/ADVOCATE/tree/main)

It has 4 Modes:
- Record: Records execution of a program and saves it as a trace file.
- Replay: Replays a trace file.
- Analysis: Record a program, checks if trace can be rewritten in a way that a concurrency bug arises, replay the rewritten trace and confirm the bug arises.
- Fuzzing: Apply fuzzing to increase coverage of the analysis


# Analysis



## Send on Closed Channel Example
The code for the example:
```go
func TestSendOnClosedChannel(t *testing.T) {
	ch := make(chan int)

	go func() {
		time.Sleep(20 * time.Millisecond) //delay
		ch <- 42
	}()

	go func() {
		time.Sleep(50 * time.Millisecond) //delay
		close(ch)
	}()

	v := <-ch
	fmt.Println("received:", v) // Outputs: received: 42

	time.Sleep(100 * time.Millisecond)
}
```
The output is: `received: 42`.  
But the artificial delays may be different on real delays.

Running the analysis gives a bug report, see [here](examples_presentation/advocateResult/analyse-SendClosedChannel/bugs/bug_1.md), where the bug is explained and the code parts involved in the bug are marked.  
It also gives the results logged:
```
==================== Summary ====================

-------------------- Critical -------------------

1 Possible send on closed channel:tp:
	send: /workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go:19@10
	close: /workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go:24@28

```
There you can see the actions and their positions in the code that lead to the bug.

Looking at the generated traces, we can see that in order to find the bug, the analysis had to rewrite the traces in a way that the send on the channel happens after the channel is closed, which is not the case in the original trace. Combined recorded trace:
```
N,2,1000000001,NC,0,/workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#15
G,4,2,/workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#17
G,6,3,/workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#22
C,8,13,1000000001,R,f,1,0,0,/workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#27
E,32
// sends on the channel
C,10,12,1000000001,S,f,1,0,0,/workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#19
E,14
// closes the channel
C,28,28,1000000001,C,f,0,0,0,/workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go#24
E,30
```


## Deadlock Example
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

The output can be empty or it can be:
```
goroutine 1 working
goroutine 1 done
goroutine 2 working
goroutine 2 done
```

