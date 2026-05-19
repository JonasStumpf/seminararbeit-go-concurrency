# Bug: P01 - Possible Send on Closed Channel

The analyzer detected a possible send on a closed channel.
Although the send on a closed channel did not occur during the recording, it is possible that it will occur, based on the happens before relation.
Such a send on a closed channel leads to a panic.

## Test/Program
The bug was found in the following test/program:

- Test/Prog: TestSendOnClosedChannel
- File: /workspaces/seminararbeit-go-concurrency/examples_presentation/SendClosedChannel_test.go
- Trace: advocateTrace_1

## Bug Elements
The elements involved in the found bug are located at the following positions:

###  Channel: Send
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


###  Channel: Close
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


## Replay
**Replaying was not run**.

