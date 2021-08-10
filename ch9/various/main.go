// Various code samples for chapter 9.
package main

import (
	"fmt"
	"time"
)

func main() {
	// =========================================================================
	raceCondition()

	// =========================================================================
	memorySynchronization()

	// =========================================================================
	gomaxprocs()
}

func raceCondition() {
	// This program contains a particular kind of race condition called a _data
	// race_. A data race occurs whenever two goroutines access the same
	// variable concurrently and at least one of the accesses is a write.
	//
	// Things get even messier if the data race involves a variable of a type
	// that is larger than a single machine word, such as an interface, a
	// string, or a slice. This code updates `x` concurrently to two slices of
	// different lengths:
	var x []int
	go func() { x = make([]int, 10) }()
	go func() { x = make([]int, 1000000) }()
	// x[999999] = 1 // NOTE: undefined behavior; memory corruption possible!
	// panic: runtime error: index out of range [999999] with length 0
	fmt.Println("x =", x)
}

func memorySynchronization() {
	// There are two reasons we need a mutex.
	// The second (and more subtle) reason is that synchronization is about more
	// than just the order of execution of multiple goroutines; synchronization
	// also affects memory.
	//
	// Synchronization primitives like channel communications and mutex
	// operations cause the processor to flush out and commit all its
	// accumulated writes so that the effects of goroutine execution up to that
	// point are guaranteed to be visible to goroutines running on other
	// processors.

	// Consider the possible outputs of the following snippet of code:
	var x, y int
	go func() {
		x = 1                   // A1
		fmt.Print("y:", y, " ") // A2
	}()
	go func() {
		y = 1                   // B1
		fmt.Print("x:", x, " ") // B2
	}()

	time.Sleep(100 * time.Millisecond)

	// Since these two goroutines are concurrent and access shared variables
	// without mutual exclusion, there is a data race, so we should not be
	// surprised that the program is not deterministic. We might expect it to
	// print any one of these four results, which correspond to intuitive
	// interleavings of the labeled statements of the program:
	//
	// y:0 x:1
	// x:0 y:1
	// x:1 y:1
	// y:1 x:1
	//
	// but depending on the compiler, CPU, and many other factors, they can
	// happen too. What possible interleaving of the four statements could
	// explain them?
	//
	// Within a single goroutine, the effects of each statement are guaranteed
	// to occur in the order of execution; goroutines are _sequentially
	// consistent_. But in the absence of explicit synchronization using a
	// channel or mutex, there is no guarantee that events are seen in the same
	// order by all goroutines. Although goroutine `A` must observe the effect
	// of the write `x = 1` before it reads the value of `y`, it does not
	// necessarily observe the write to `y` done by goroutine `B`, so `A` may
	// print a stale value of `y`.
}

func gomaxprocs() {
	// The Go scheduler uses a parameter called `GOMAXPROCS` to determine how
	// many OS threads may be actively executing Go code simultaneously.
	// Its default value is the number of CPUs on the machine.
	//
	// You can explicitly control this parameter using the `GOMAXPROCS`
	// environment variable or the `runtime.GOMAXPROCS` function.
	// We can see the effect of `GOMAXPROCS` on this little program, which
	// prints an endless stream of zeros and ones:
	for {
		go fmt.Print(0)
		fmt.Print(1)
	}

	// Run:
	//
	// $ GOMAXPROCS=1 go run gopl.io/ch9/various
	// 111111111111111111110000000000000000000011111...
	//
	// $ GOMAXPROCS=2 go run gopl.io/ch9/various
	// 010101010101010101011001100101011010010100110...

	// In the first run, at most one goroutine was executed at a time.
	// Initially, it was the main goroutine, which prints ones. After a period
	// of time, the Go scheduler put it to sleep and woke up the goroutine that
	// prints zeros, giving it a turn to run on the OS thread. In the second
	// run, there were two OS threads available, so both goroutines ran
	// simultaneously, printing digits at about the same rate. We must stress
	// that many factors are involved in goroutine scheduling, and the runtime
	// is constantly evolving, so your results may differ from the ones above.
}
