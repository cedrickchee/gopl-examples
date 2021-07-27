// Nonempty is an example of an in-place slice algorithm.
package main

import "fmt"

func main() {
	data := []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty(data)) // "["one" "three"]"
	fmt.Printf("%q\n", data)           // "["one" "three" "three"]"

	// nonempty implemented using `append`.
	fruits := []string{"apple", "", "orange"}
	fmt.Printf("%q\n", nonempty2(fruits)) // "["apple" "orange"]"
	fmt.Printf("%q\n", fruits)            // "["apple" "orange" "orange"]"
}

// nonempty returns a slice holding only the non-empty strings.
// The underlying array is modified during the call.
//
// The subtle part is that the input slice and the output slice share the same
// underlying array. This avoids the need to allocate another array, though
// of course the contents of data are partly overwritten.
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

// The nonempty function can also be written using append
func nonempty2(strings []string) []string {
	out := strings[:0] // zero-length slice of original
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
