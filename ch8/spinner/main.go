// Spinner displays an animation while computing the 45th Fibonacci number.
package main

import (
	"fmt"
	"time"
)

// The main goroutine computes the 45th Fibonacci number.
//
// Since it uses the terribly inefficient recursive algorithm, it runs for an
// appreciable time, during which we’d like to provide the user with a visual
// indication that the program is still running, by displaying an animated
// textual "spinner".
func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// Each number in the Fibonacci sequence is the sum of the two numbers that
// precede it.
// An example of fib sequence: 0, 1, 1, 2, 3, 5, 8, 13, 21, and so on.
// fib(2) = fib(0) + fib(1) = 0 + 1 = 1
// fib(3) = fib(1) + fib(2) = 1 + 1 = 2
// fib(4) = fib(2) + fib(3) = 1 + 2 = 3
func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

/*
After several seconds of animation, the fib(45) call returns and the main
function prints its result:

$ go run gopl.io/ch8/spinner
Fibonacci(45) = 1134903170
*/

// The `main` function then returns. When this happens, all goroutines are
// abruptly terminated and the program exits. Other than by returning from
// `main` or exiting the program, there is no programmatic way for one goroutine
// to stop another, but as we will see later, there are ways to communicate with
// a goroutine to request that it stop itself.
//
// Notice how the program is expressed as the composition of two autonomous
// activities, spinning and Fibonacci computation. Each is written as a separate
// function but both make progress concurrently.
