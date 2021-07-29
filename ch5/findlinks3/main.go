// Findlinks3 crawls the web, starting with the URLs on the command line.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

// In our crawler, items are URLs. The crawl function we’ll supply to
// breadthFirst prints the URL, extracts its links, and returns them so that
// they too are visited.
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

/*
Let’s crawl the web starting from https://golang.org. Here are some of the
resulting links:

$ go run gopl.io/ch5/findlinks3 https://golang.org
https://golang.org
https://support.eji.org/give/153413/#!/donation/checkout
https://golang.org/
https://golang.org/doc/
https://golang.org/pkg/
https://golang.org/project/
https://golang.org/help/
https://golang.org/blog/
https://play.golang.org/
https://golang.org/dl/
https://tour.golang.org/
https://blog.golang.org/
https://golang.org/doc/copyright.html
https://golang.org/doc/tos.html
http://www.google.com/intl/en/policies/privacy/
http://golang.org/issues/new?title=x/website:
https://google.com
https://golang.org/doc/install
https://golang.org/doc/tutorial/getting-started.html
...

The process ends when all reachable web pages have been crawled or the memory
of the computer is exhausted.
*/
