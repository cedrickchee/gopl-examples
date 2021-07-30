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

	// =========================================================================
	variadicFunctions()

	// =========================================================================
	deferFunctionCalls()

	// =========================================================================
	dangerousDeferFuncExec()

	// =========================================================================
	fixDangerousDeferFuncExec()
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

func f(vals ...int) {}
func g(vals []int)  {}
func variadicFunctions() {
	fmt.Printf("%T\n", f) // "func(...int)"
	fmt.Printf("%T\n", g) // "func([]int)"
}

func double(x int) (result int) {
	// By naming its result variable and adding a defer statement, we can make
	// the function print its arguments and results each time it is called.
	defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
	return x + x
}
func triple(x int) (result int) {
	// A deferred anonymous function can even change the values that the
	// enclosing function returns to its caller:
	defer func() { result += x }()
	return double(x)
}
func deferFunctionCalls() {
	_ = double(4)
	// Output:
	// "double(4) = 8"

	fmt.Println(triple(4)) // "12"
}

func readFiles(filenames []string) error {
	// Because deferred functions aren’t executed until the very end of a
	// function’s execution, a defer statement in a loop deserves extra
	// scrutiny. The code below could run out of file descriptors since no file
	// will be closed until all files have been processed:

	fmt.Println("dangerously reading files...")

	for _, filename := range filenames {
		fmt.Printf("reading filename: %s\n", filename)

		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close() // NOTE: risky; could run out of file descriptors

		defer func() { fmt.Println("file close") }()
		// ...process f...
	}

	defer func() { fmt.Println("done reading files") }()

	return nil
}
func dangerousDeferFuncExec() {
	filenames := []string{"bin/out.gif", "bin/test.txt", "bin/test2.txt"}
	err := readFiles(filenames)
	if err != nil {
		fmt.Printf("dangerousDeferFuncExec: %v", err)
	}
	/*
		Output (notice the sequence of the log messages):

		dangerously reading files...
		reading filename: bin/out.gif
		reading filename: bin/test.txt
		reading filename: bin/test2.txt
		done reading files
		file close
		file close
		file close
	*/
}

func fixDangerousDeferFuncExec() {
	// One solution is to move the loop body, including the `defer` statement,
	// into another function that is called on each iteration.

	filenames := []string{"bin/out.gif", "bin/test.txt", "bin/test2.txt"}
	err := safeReadFiles(filenames)
	if err != nil {
		fmt.Printf("fixDangerousDeferFuncExec: %v", err)
	}
	/*
		Output (notice the sequence of the log messages):

		safely reading files...
		reading filename: bin/out.gif
		file close
		reading filename: bin/test.txt
		file close
		reading filename: bin/test2.txt
		file close
		done reading files
	*/
}
func safeReadFiles(filenames []string) error {
	fmt.Println("safely reading files...")

	for _, filename := range filenames {
		fmt.Printf("reading filename: %s\n", filename)

		if err := doFile(filename); err != nil {
			return err
		}
	}

	defer func() { fmt.Println("done reading files") }()

	return nil
}
func doFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	defer func() { fmt.Println("file close") }()
	// ...process f...

	return nil
}
