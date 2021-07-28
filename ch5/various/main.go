package main

import "fmt"

func main() {
	declarations()

	// =========================================================================
}

func declarations() {
	{
		add := func(x int, y int) int { return x + y }
		sub := func(x, y int) (z int) { z = x - y; return }
		first := func(x int, _ int) int { return x }
		zero := func(int, int) int { return 0 }

		fmt.Printf("%T\n", add)   // "func(int, int) int"
		fmt.Printf("%T\n", sub)   // "func(int, int) int"
		fmt.Printf("%T\n", first) // "func(int, int) int"
		fmt.Printf("%T\n", zero)  // "func(int, int) int"
	}

	{
		var fac func(int) int
		fac = func(n int) int {
			if n == 0 {
				return 1
			}
			return n * fac(n-1)
		}
		fmt.Println("factorial(7) =", fac(7))
	}
}
