// Comma prints its argument numbers with a comma at each power of 1000.
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
	n := len(s)
	// If its length is less than or equal to 3, no comma is necessary.
	if n <= 3 {
		return s
	}
	// calls itself recursively with a substring consisting of all but the last
	// 3 characters, and appends a comma and the last 3 characters to the
	// result of the recursive call.
	return comma(s[:n-3]) + "," + s[n-3:]
}
