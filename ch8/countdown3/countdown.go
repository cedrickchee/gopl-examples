// Countdown implements the countdown for a rocket launch.
package main

// NOTE: the ticker goroutine never terminates if the launch is aborted.
// This is a "goroutine leak".

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// ...create abort channel...

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown.  Press return to abort.")
	tick := time.Tick(1 * time.Second) // like a metronome
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		// The select statement below causes each iteration of the loop to wait
		// up to 1 second for an abort, but no longer.
		select {
		case <-tick:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	launch()
}

func launch() {
	fmt.Println("Lift off!")
}

/*
Run:
$ go run gopl.io/ch8/countdown3
Commencing countdown.  Press return to abort.

Launch aborted!

$ go run gopl.io/ch8/countdown3
Commencing countdown.  Press return to abort.
10
9
8
7
6
5
4
3
2
1
Lift off!
*/

// About the goroutine leak
//
// The `time.Tick` function behaves as if it creates a goroutine that calls
// `time.Sleep` in a loop, sending an event each time it wakes up. When the
// countdown function above returns, it stops receiving events from tick, but
// the ticker goroutine is still there, trying in vain to send on a channel from
// which no goroutine is receiving--a _goroutine leak_
