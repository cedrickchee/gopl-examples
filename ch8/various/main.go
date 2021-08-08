// Various code samples for chapter 8.
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// =========================================================================
	goroutines()

	// =========================================================================
	channels()

	// =========================================================================
	bufferedChannels()

	// =========================================================================
	multiplexingWithSelect()

	// =========================================================================
	tickerPattern()

	// =========================================================================
	pollingChannel()
}

func f(name string) {
	fmt.Printf("%s\nspinning...\n", name)
	time.Sleep(1 * time.Second)
}
func goroutines() {
	f("func call")    // call f(); wait for it to return
	go f("goroutine") // create a new goroutine that calls f(); don't wait
}

// If goroutines are the activities of a concurrent Go program, _channels_ are
// the connections between them. A channel is a communication mechanism that
// lets one goroutine send values to another goroutine. Each channel is a
// conduit for values of a particular type, called the channel’s _element_ type.
// The type of a channel whose elements have type `int` is written `chan int`.
func channels() {
	// To create a channel, we use the built-in `make` function:
	ch := make(chan int)
	fmt.Printf("ch type = %T, value = %[1]v\n", ch)

	// Two channels of the same type may be compared using `==`.

	// A channel has two principal operations, _send_ and _receive_,
	// collectively known as _communications_. A send statement transmits a
	// value from one goroutine, through the channel, to another goroutine
	// executing a corresponding receive expression. Both operations are written
	// using the `<-` operator. In a send statement, the `<-` separates the
	// channel and value operands. In a receive expression, `<-` precedes the
	// channel operand. A receive expression whose result is not used is a valid
	// statement.
	// x := 42
	// ch <- x  // a send statement
	// x = <-ch // a receive expression in an assignment statement
	// <-ch // a receive statement; result is discarded

	// Channels support a third operation, _close_, which sets a flag
	// indicating that no more values will ever be sent on this channel;
	// subsequent attempts to send will panic.
	// To close a channel, we call the built-in `close` function:
	// close(ch)

	// A channel created with a simple call to `make` is called an _unbuffered_
	// channel, but `make` accepts an optional second argument, an integer
	// called the channel’s _capacity_. If the capacity is non-zero, `make`
	// creates a _buffered_ channel.
	ch = make(chan int)    // unbuffered channel
	ch = make(chan int, 0) // unbuffered channel
	ch = make(chan int, 3) // buffered channel with capacity 3
}

func bufferedChannels() {
	// A buffered channel has a queue of elements. The queue’s maximum size is
	// determined when it is created, by the capacity argument to `make`.
	ch := make(chan string, 3)

	// A send operation on a buffered channel inserts an element at the back of
	// the queue, and a receive operation removes an element from the front. If
	// the channel is full, the send operation blocks its goroutine until space
	// is made available by another goroutine’s receive.

	// We can send up to three values on this channel without the goroutine
	// blocking:
	ch <- "A"
	ch <- "B"
	ch <- "C"
	// At this point, the channel is full, and a fourth send statement would
	// block.

	// If we receive one value,
	fmt.Println(<-ch) // "A"

	// the channel is neither full nor empty, so either a send operation or a
	// receive operation could proceed without blocking. In this way, the
	// channel’s buffer decouples the sending and receiving goroutines.

	// In the unlikely event that a program needs to know the channel’s buffer
	// capacity, it can be obtained by calling the built-in `cap` function:
	fmt.Println(cap(ch)) // "3"

	// When applied to a channel, the built-in `len` function returns the number
	// of elements currently buffered. Since in a concurrent program this
	// information is likely to be stale as soon as it is retrieved, its value
	// is limited, but it could conceivably be useful during fault diagnosis or
	// performance optimization.
	fmt.Println(len(ch)) // "2"

	// After two more receive operations the channel is empty again, and a
	// fourth would block:
	fmt.Println(<-ch) // "B"
	fmt.Println(<-ch) // "C"

	// In this example, the send and receive operations were all performed by
	// the same goroutine, but in real programs they are usually executed by
	// different goroutines.

	// The example below shows an application of a buffered channel.
	quickestMirror := mirroredQuery()
	fmt.Println("Quickest Mirror =", quickestMirror)
}

// mirroredQuery makes parallel requests to three _mirrors_, that is, equivalent
// but geographically distributed servers. It tells you the quickest mirror.
func mirroredQuery() string {
	// It sends their responses over a buffered channel, then receives and
	// returns only the first response, which is the quickest one to arrive.
	responses := make(chan string, 3)
	go func() { responses <- request("asia.gopl.io") }()
	go func() { responses <- request("europe.gopl.io") }()
	go func() { responses <- request("americas.gopl.io") }()
	return <-responses // return the quickest response

	// Had we used an unbuffered channel, the two slower goroutines would have
	// gotten stuck trying to send their responses on a channel from which no
	// goroutine will ever receive. This situation, called a goroutine leak,
	// would be a bug.
}

// Simulate request to hostname.
func request(hostname string) (response string) {
	time.Sleep(1 * time.Second)
	response = hostname
	return
}

func multiplexingWithSelect() {
	// A `select` waits until a communication for some case is ready to proceed.
	// It then performs that communication and executes the case’s associated
	// statements; the other communications do not happen. A `select` with no
	// cases, `select{}`, waits forever.

	// The example below is more subtle. The channel `ch`, whose buffer size is
	// 1, is alternately empty then full, so only one of the cases can proceed,
	// either the send when `i` is even, or the receive when `i` is odd. It
	// always prints `0 2 4 6 8`.
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x) // "0" "2" "4" "6" "8"
		case ch <- i:
		}
	}

	// If multiple cases are ready, `select` picks one at random, which ensures
	// that every channel has an equal chance of being selected. Increasing the
	// buffer size of the previous example makes its output nondeterministic,
	// because when the buffer is neither full nor empty, the select statement
	// figuratively tosses a coin.
}

func tickerPattern() {
	// A solution for the ticker goroutine leak.

	// The `Tick` function is convenient, but it’s appropriate only when the
	// ticks will be needed throughout the lifetime of the application.
	// Otherwise, we should use this pattern:
	ticker := time.NewTicker(1 * time.Second)
	<-ticker.C // receive from the ticker's channel

	// Stop the ticker to release associated resources.
	ticker.Stop() // cause the ticker's goroutine to terminate

	// My example:
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	fmt.Println("Commencing countdown.  Press return to abort.")
	cnttick := time.NewTicker(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		// The select statement below causes each iteration of the loop to wait
		// up to 1 second for an abort, but no longer.
		select {
		case <-cnttick.C:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			cnttick.Stop()
			return
		}
	}
	launch := func() { fmt.Println("Lift off!") }
	launch()
}

func pollingChannel() {
	// Sometimes we want to try to send or receive on a channel but avoid
	// blocking if the channel is not ready—a _non-blocking_ communication. A
	// select statement can do that too. A `select` may have a `default`, which
	// specifies what to do when none of the other communications can proceed
	// immediately.

	// The select statement below receives a value from the `abort` channel if
	// there is one to receive; otherwise it does nothing. This is a
	// non-blocking receive operation; doing it repeatedly is called _polling_ a
	// channel.
	abort := make(chan struct{})
	select {
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	default:
		// do nothing
	}
}
