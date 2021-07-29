// Exercise 5.5: Implement `countWordsAndImages`. (See Exercise 4.9 for
// word-splitting.)
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Get URL from CLI args
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: ./exercise5-5 URL")
		os.Exit(1)
	}
	url := os.Args[1]

	words, images, err := CountWordsAndImages(url)
	if err != nil {
		log.Fatal(err)
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
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	u := make([]*html.Node, 0) // unvisited
	u = append(u, n)
	for len(u) > 0 {
		n = u[len(u)-1]
		u = u[:len(u)-1]
		switch n.Type {
		case html.TextNode:
			words += wordCount(n.Data)
		case html.ElementNode:
			if n.Data == "img" {
				images++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			u = append(u, c)
		}
	}
	return
}

func wordCount(s string) int {
	// Call `input.Split(bufio.ScanWords)` before the first call to `Scan` to
	// break the input into words instead of lines.
	n := 0
	inp := bufio.NewScanner(strings.NewReader(s))
	inp.Split(bufio.ScanWords)
	for inp.Scan() {
		n++
	}
	return n
}
