package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	declarations()

	// =========================================================================
	multipleReturnValues()

	// =========================================================================
	functionValues()

	// =========================================================================
	zeroValueFunctionType()

	// =========================================================================
	functionValuesBehavior()
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

func multipleReturnValues() {
	// Bare return
	words, images, err := CountWordsAndImages("https://golang.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "countWordsAndImages failed: %v", err)
	}
	fmt.Printf("%d words and %d images\n", words, images)
}

// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	return rand.Int() % 100, rand.Int() % 100
}

func functionValues() {
	// Functions are first-class values in Go.
	fmt.Println("square(2) =", square(2))
	fmt.Println("negative(5) =", negative(5))
	fmt.Println("product(3,5) =", product(3, 5))

	f := square
	fmt.Println("f(3) =", f(3)) // "9"
	f = negative
	fmt.Println("f(3) =", f(3)) // "-3"
	fmt.Printf("%T\n", f)       // "func(int) int"
	/*
		f = product                 // compile error: cannot use product (value of type func(m int, n int) int) as func(n int) int value in assignment
	*/
}
func square(n int) int     { return n * n }
func negative(n int) int   { return -n }
func product(m, n int) int { return m * n }

func zeroValueFunctionType() {
	// The zero value of a function type is nil. Calling a nil function value
	// causes a panic.
	var f func(int) int
	/*
		f(3) // panic: runtime error: invalid memory address or nil pointer dereference
	*/

	// Function values may be compared with `nil`, but they are not comparable,
	// so they may not be compared against each other or used as keys in a map.
	if f != nil {
		f(3)
	}
}

func functionValuesBehavior() {
	fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
	fmt.Println(strings.Map(add1, "VMS"))      // "WNT"
	fmt.Println(strings.Map(add1, "Admix"))    // "Benjy"
}
func add1(r rune) rune { return r + 1 }
