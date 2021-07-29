package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	declarations()

	// =========================================================================
	multipleReturnValues()

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
