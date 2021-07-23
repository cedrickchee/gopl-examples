// Dup1 prints the text of each line that appears more than
// once in the standard input, preceded by its count.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// Ref this SO thread for breaking the loop:
		// https://stackoverflow.com/a/49715256/206570
		text := input.Text()
		if len(text) != 0 {
			counts[text]++
		} else {
			break
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	for line, n := range counts {
		if n > 0 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
