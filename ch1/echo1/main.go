// Echo1 prints its command-line arguments.
package main

// Import two packages, which are given as a parenthesized list form.
import (
	"fmt"
	"os"
)

func main() {
	// Declares two variables s and sep, of type string.
	// It is not explicitly initialized. It is implicitly initialize d to the
	// zero value for its type, which is the empty string "" for strings.
	var s, sep string

	// The echo program could have printed its out put in a loop one piece at
	// a time, but this version instead builds up a string by repeatedly
	// appending new text to the end.
	for i := 1; i < len(os.Args); i++ {
		// The := symbol is part of a short variable declaration, a statement
		// that declares one or more variables and gives them appropriate types
		// based on the initializer values.

		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}
