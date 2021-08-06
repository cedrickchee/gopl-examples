// Pipeline2 demonstrates a finite 3-stage pipeline.
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
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	// Squarer
	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	for x := range squares {
		fmt.Println(x)
	}
}

/*
Run:
$ go run gopl.io/ch8/pipeline2
0
1
4
9
16
25
36
49
64
81
100
121
(cut)
9025
9216
9409
9604
9801
*/
