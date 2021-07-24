// Various code samples for chapter 2.
package main

import "fmt"

func main() {
	pointers()

	// =========================================================================
	zeroptr()

	// =========================================================================
	var p = f()
	fmt.Println(p)  // 0xc000014100
	fmt.Println(*p) // 1
	// Each call of f returns a distinct value
	fmt.Println(f() == f()) // "false"

	// =========================================================================
	v := 1
	incr(&v)              // side effect: v is now 2
	fmt.Println(incr(&v)) // "3" (and v is 3)

	// =========================================================================
	theNewFunc()

	// =========================================================================
	newDistinct()

	// =========================================================================
	variableLifetime()

	// =========================================================================
	tupleAssignment()
}

// Pointers
func pointers() {
	x := 1
	p := &x         // p, of type *int, points to x
	fmt.Println(*p) // "1"
	*p = 2          // equivalent to x = 2
	fmt.Println(x)  // "2"
}

// The zero value for a point er of any type is nil.
func zeroptr() {
	var x, y int
	fmt.Println(&x == &x, &x == &y, &x == nil) // "true false false"
}

// It is perfectly safe for a function to return the address of a local variable.
func f() *int {
	v := 1
	return &v
}

func incr(p *int) int {
	*p++ // increments what p points to; does not change p
	return *p
}

func theNewFunc() {
	p := new(int)   // p, of type *int, points to an unnamed int variable
	fmt.Println(*p) // "0"
	*p = 2          // sets the unnamed int to 2
	fmt.Println(*p) // "2"
}

func newDistinct() {
	p := new(int)
	q := new(int)
	fmt.Println("new return distinct var:", p == q) // "false"
}

var global *int

func variableLifetime() {
	// Here, x must be heap-allocated because it is still reachable from the
	// variable global after variableLifetime has returned, despite being declared
	// as a local variable; we say x escapes from variableLifetime
	var x int
	x = 1
	global = &x
}

// Another form of assignment, known as tuple assignment, allows several
// variables to be assigned at once.
func tupleAssignment() {
	fmt.Println("gcd:", gcd(32, 4))
	fmt.Println("8-th Fibonacci number:", fib(8))
}

// gcd compute the greatest common divisor (GCD) of two integers
func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// fib compute the n-th Fibonacci number iteratively
// 0, 1, 1, 2, 3, 5, 8, 13, 21, ...
func fib(n int) int {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		x, y = y, x+y
	}
	return x
}
