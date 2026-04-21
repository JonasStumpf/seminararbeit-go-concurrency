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

