// Various code samples for chapter 4.
package main

import "fmt"

func main() {
	// =========================================================================
	arrays()

	// =========================================================================
	compareArray()
}

func arrays() {
	// Array literal
	var q [3]int = [3]int{1, 2, 3}
	fmt.Println(q[2]) // "3"
	var r [3]int = [3]int{1, 2}
	fmt.Println(r[2]) // "0"

	// In an array literal, if an ellipsis "..."" appears in place of the
	// length, the array length is determined by the number of initializers.
	p := [...]int{1, 2, 3, 4}
	fmt.Printf("%T\n", p) // "[4]int"

	// literal syntax
	type Currency int
	const (
		USD Currency = iota
		EUR
		GBP
		RMB
	)
	symbol := [...]string{USD: "$", EUR: "€", GBP: "£", RMB: "¥"}
	fmt.Println(RMB, symbol[RMB]) // "3 ¥"

	// Defines an array with 100 elements, all zero except for the last,
	// which has value -1.
	d := [...]int{99: -1}
	fmt.Println(d)
}

func compareArray() {
	a := [2]int{1, 2}
	b := [...]int{1, 2}
	c := [2]int{1, 3}
	fmt.Println(a == b, a == c, b == c) // "true false false"
	/*
		d := [3]int{1, 2}
		fmt.Println(a == d) // compile error: cannot compare [2]int == [3]int
	*/
}
