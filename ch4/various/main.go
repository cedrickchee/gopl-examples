// Various code samples for chapter 4.
package main

import "fmt"

func main() {
	// =========================================================================
	arrays()

	// =========================================================================
	compareArray()

	// =========================================================================
	var a [32]byte = [32]byte{31: byte('a')}
	fmt.Println(a)
	zero(&a)
	fmt.Println(a)

	// =========================================================================
	var b [32]byte = [32]byte{31: byte('b')}
	fmt.Println(b)
	anotherZero(&b)
	fmt.Println(b)

	// =========================================================================
	arrayAndSlices()
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

// Explicitly pass a pointer to an array so that any modifications the function
// makes to array elements will be visible to the caller.
func zero(ptr *[32]byte) {
	// Zeroes the contents of a [32]byte array.
	for i := range ptr {
		ptr[i] = 0
	}
}

func anotherZero(ptr *[32]byte) {
	// The array literal [32]byte{} yields an array of 32 bytes.
	// Each element of the array has the zero value for byte, which is zero.
	*ptr = [32]byte{}
}

func arrayAndSlices() {
	// Array
	months := [...]string{1: "January", 2: "February", 3: "March", 4: "April",
		5: "May", 6: "June", 7: "July", 8: "August", 9: "September",
		10: "October", 11: "November", 12: "December"}
	fmt.Printf("months = %v, type = %[1]T\n", months)
	// Slices
	Q2 := months[4:7]
	summer := months[6:9]
	fmt.Printf("Q2 value=%v, type=%[1]T\n", Q2)         // "Q2 value=[April May June], type=[]string"
	fmt.Printf("summer value=%v, type=%[1]T\n", summer) // "summer value=[June July August], type=[]string"

	// June is included in each and is the sole output of this (inefficient)
	// test for common elements.
	for _, s := range summer {
		for _, q := range Q2 {
			if s == q {
				fmt.Printf("%s appears in both\n", s)
			}
		}
	}

	// Slicing beyond cap(s) causes a panic, but slicing beyond len(s) extends
	// the slice, so the result may be longer than the original.
	fmt.Printf("Q2 length = %d, capacity = %d\n", len(Q2), cap(Q2))
	fmt.Printf("summer length = %d, capacity = %d\n", len(summer), cap(summer))
	// fmt.Println(summer[:20])    // panic: runtime error: slice bounds out of range [:20] with capacity 7
	endlessSummer := summer[:5]                   // extend a slice (within capacity)
	fmt.Println("endlessSummer =", endlessSummer) // "[June July August September October]"
}
