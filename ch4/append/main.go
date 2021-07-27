// Append illustrates the behavior of the built-in append function.
package main

import "fmt"

func main() {
	x := []int{1, 2, 3}
	fmt.Println(x)
	x2 := appendInt(x, 4)
	fmt.Println(x)
	fmt.Println(x2)
}

func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen < cap(x) {
		// There is room to grow. Extend the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space. Allocate a new array.
		// Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x) // a built-in function; see text
	}
	z[len(x)] = y
	return z
}
