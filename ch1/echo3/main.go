// Echo3 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// If the amount of data involved is large, this could be costly.
	// A simpler and more efficient solution would be to use the Join function
	// from the strings package.
	fmt.Println(strings.Join(os.Args[1:], " "))
}
