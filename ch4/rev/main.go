// Rev reverses a slice.
package main

import "fmt"

func main() {
	a := [...]int{0, 1, 2, 3, 4, 5}
	// reverse the whole array a.
	reverse(a[:])
	fmt.Println(a) // ""[5 4 3 2 1 0]"

	s := []int{0, 1, 2, 3, 4, 5}
	// Rotate s left by two positions.
	reverse(s[:2])
	fmt.Println(s) // "[1 0 2 3 4 5]"
	reverse(s[2:])
	fmt.Println(s) // "[1 0 5 4 3 2]"
	reverse(s)
	fmt.Println(s) // "[2 3 4 5 0 1]"

	// (To rotate to the right, make the third call first.)
	s2 := []int{0, 1, 2, 3, 4, 5}
	reverse(s2)
	reverse(s2[:2])
	reverse(s2[2:])
	fmt.Println(s2) // "[4 5 0 1 2 3]"
}

// reverse reverses a slice of ints in place.
func reverse(s []int) {
	// start: 5 4 3 2 1 0
	// step: i=0, j=5
	// 		  0 4 3 2 1 5
	// step: i=1, j=4
	// 		  0 1 3 2 4 5
	// step: i=2, j=3
	// end:   0 1 2 3 4 5
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
