// Pipeline1 demonstrates an infinite 3-stage pipeline.
package main

import "fmt"

// The first goroutine, _counter_, generates the integers 0, 1, 2, ..., and
// sends them over a channel to the second goroutine, _squarer_, which receives
// each value, squares it, and sends the result over another channel to the
// third goroutine, _printer_, which receives the squared values and prints
// them. For clarity of this example, we have intentionally chosen very simple
// functions, though of course they are too computationally trivial to warrant
// their own goroutines in a realistic program.
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; ; x++ {
			naturals <- x
		}
	}()

	// Squarer
	go func() {
		for {
			x := <-naturals
			squares <- x * x
		}
	}()

	// Printer (in main goroutine)
	for {
		fmt.Println(<-squares)
	}
}

/*
Run:
$ go run gopl.io/ch8/pipeline1
*/
