# seminararbeit-go-concurrency
Seminararbeit Go Concurrency




Running advocate in this container needs full paths.  
e.g. running an analysis on Basic_lockdep_test from the advocate deadlocks example:
```bash
/workspaces/seminararbeit-go-concurrency/ADVOCATE/advocate/advocate analysis -path /workspaces/seminararbeit-go-concurrency/ADVOCATE/doc/examples/deadlocks/ -exec TestBasicLockdep
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


## Results from `analysis` run
Format is:
path:line@timestamp
timestamp is the value of the vector clock at the time of the event.

> **Note**: Line numbers appear to be wrong, probably because advocate adds lines before execution when the toolchain (commandline) is used:
> - adds "advocate" in import
> - adds at start of test/main method:
> ```go
> // ======= Preamble Start =======
>  advocate.InitTracing()
>  defer advocate.FinishTracing()
>// ======= Preamble End =======
> ```
> This sums up to 5 extra lines. Subtracting these from the reported line numbers seems to give the correct line numbers.


Advocate deadlock example Basic_test results:
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
19      x.Lock() // <--     // this SHOULD produce a deadlock
20      x.Unlock()
21      y.Unlock()
22  }
```
Sometimes a bug report markdown is created, see for this [Example](advocateExampleResults/file%287%29-test%280%29-Basic_test-TestBasic/bugs/bug_1.md).


Advocate example without a Bug, Guard_test:
```
==================== Summary ====================

No bugs found
```

Example, SendClosedChannel_test:
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

