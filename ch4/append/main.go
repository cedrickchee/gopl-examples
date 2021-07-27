// Append illustrates the behavior of the built-in append function.
package main

import "fmt"

func main() {
	// append
	x := []int{1, 2, 3}
	fmt.Println(x)
	x2 := appendInt(x, 4)
	fmt.Println(x)
	fmt.Println(x2)
	x2 = appendInt(x2, 5, 6)
	fmt.Println(x2)
	x2 = appendInt(x2, []int{7, 8}...)
	fmt.Println(x2)

	// growth
	var p, q []int
	for i := 0; i < 10; i++ {
		q = appendInt(p, i)
		fmt.Printf("%d cap=%d\t%v\n", i, cap(q), q)
		p = q
	}
	/*
		Output:
		0 cap=1   [0]
		1 cap=2   [0 1]
		2 cap=4   [0 1 2]
		3 cap=6   [0 1 2 3]
		4 cap=6   [0 1 2 3 4]
		5 cap=10  [0 1 2 3 4 5]
		6 cap=10  [0 1 2 3 4 5 6]
		7 cap=10  [0 1 2 3 4 5 6 7]
		8 cap=10  [0 1 2 3 4 5 6 7 8]
		9 cap=18  [0 1 2 3 4 5 6 7 8 9]
	*/
}

func appendInt(x []int, y ...int) []int {
	// The ellipsis ""..."" in the declaration of appendInt makes the
	// function variadic: it accepts any number of final arguments.

	var z []int
	zlen := len(x) + len(y)
	// expand z to at least zlen
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

	copy(z[len(x):], y)
	return z
}
