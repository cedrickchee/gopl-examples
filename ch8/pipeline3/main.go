// Pipeline3 demonstrates a finite 3-stage pipeline
// with range, close, and unidirectional channel types.
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

	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}

// Counter
func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

// Squarer
func squarer(out chan<- int, in <-chan int) {
	for x := range in {
		out <- x * x
	}
	close(out)
}

// Printer (run in main goroutine)
func printer(squares <-chan int) {
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
