// Basename1 reads file names from stdin and prints the base name of each one.
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(basename("a/b/c.go")) // "c"
	fmt.Println(basename("c.d.go"))   // "c.d"
	fmt.Println(basename("abc"))      // "abc"
}

// basename removes directory components and a .suffix.
// e.g., a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c
func basename(s string) string {
	// A simpler version uses the `strings.LastIndex` library function.

	// Discard last '/' and everything before.
	slash := strings.LastIndex(s, "/")
	s = s[slash+1:]

	// Preserve everything before last '.'.
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}
