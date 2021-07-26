// Exercise 3.10: Write a non-recursive version of comma, using bytes.Buffer
// instead of string concatenation.
//
// Example:
// 	$ go build gopl.io/ch3/comma
//	$ ./comma 1 12 123 1234 1234567890
// 	1
// 	12
// 	123
// 	1,234
// 	1,234,567,890
//
package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("%s\n", comma(os.Args[i]))
	}
}

// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	d := 3 // number of digit before/after comma
	n := len(s)
	// Base case
	if n <= 3 {
		return s
	}

	var buf bytes.Buffer
	// Example (has remainder): 1,234,567,890 (remainder is the number 1.
	// 2 commas/3 groups after the first number)
	// Example (no remainder): 123,456,789 (2 commas/3 groups)
	remainder, group := n%3, n/3
	fmt.Println("group =", group)
	fmt.Println("remainder =", remainder)

	if remainder != 0 {
		buf.WriteString(s[:remainder])
		buf.WriteByte(',')
	}

	// Test: non-loop
	// buf.WriteString(s[1:4])
	// buf.WriteByte(',')
	// buf.WriteString(s[4:7])
	// buf.WriteByte(',')
	// buf.WriteString(s[7:10])
	// 1,234,567,890
	// 123,456,789
	for i := 0; i < group; i++ {
		start := i*d + remainder
		end := start + d // or i*d+d+remainder
		buf.WriteString(s[start:end])
		if i != group-1 {
			buf.WriteByte(',')
		}
	}
	return buf.String()
}
