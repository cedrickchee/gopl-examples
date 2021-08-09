// Various code samples for chapter 9.
package main

func main() {
	// =========================================================================
	raceCondition()
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
}
