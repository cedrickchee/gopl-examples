// Exercise 8.6: Add depth-limiting to the concurrent crawler. That is, if the
// user sets `-depth=3`, then only URLs reachable by at most three links will be
// fetched.

// Code is based on Crawl2.

// This program is a depth-limited web crawler.
// It crawls web links starting with the command-line arguments.
//
// This version uses:
// - a buffered channel as a counting semaphore to limit the number of
//   concurrent calls to links.Extract. (the `tokens` chan as a
//   semaphore to limit concurrent requests)
// - a WaitGroup to determine when the work is done
// - a mutex around the `seen` map to avoid concurrent reads and writes.
package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"gopl.io/ch5/links"
)

//!+semaphore
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

var maxDepth int
var seen = make(map[string]bool)
var seenLock = sync.Mutex{}

// In our crawler, items are URLs. The crawl function prints the URL, extracts
// its links, and returns them so that they too are visited.
func crawl(url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(depth, url)
	if depth >= maxDepth {
		return
	}
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	//!-semaphore

	for _, link := range list {
		seenLock.Lock()
		if seen[link] {
			seenLock.Unlock()
			continue
		}
		seen[link] = true
		seenLock.Unlock()
		wg.Add(1)
		go crawl(link, depth+1, wg)
	}
}

// A simple web crawler program that explored the link graph of the web in
// breadth-first order.
func main() {
	// Parse command-line arguments
	flag.IntVar(&maxDepth, "depth", 3, "max crawl depth")
	flag.Parse()

	wg := &sync.WaitGroup{}
	for _, link := range flag.Args() {
		wg.Add(1)
		go crawl(link, 0, wg)
	}
	wg.Wait()

	// Older Crawl2 version:
	// The main function resembles breadthFirst. As before, a worklist records
	// the queue of items that need processing, each item being a list of URLs
	// to crawl, but this time, instead of representing the queue using a slice,
	// we use a channel. Each call to `crawl` occurs in its own goroutine and
	// sends the links it discovers back to the worklist.
	// worklist := make(chan []string)
	// var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	// n++
	// go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	// seen := make(map[string]bool)
	// for ; n > 0; n-- { // <-- fix program never terminates because the worklist is never closed.
	// 	list := <-worklist
	// 	for _, link := range list {
	// 		if !seen[link] {
	// 			seen[link] = true
	// 			n++
	// 			go func(link string) {
	// 				worklist <- crawl(link)
	// 			}(link)
	// 		}
	// 	}
	// }
}

/*
Let’s crawl the web starting from http://gopl.io/. Here are some of the
resulting links:

$ go build gopl.io/ch8/exercise8-6
$ ./exercise8-6 -depth=1 http://gopl.io/
0 http://gopl.io/
1 http://www.amazon.com/dp/020161586X?tracking_id=disfordig-20
1 http://golang.org/lib/godoc/analysis/help.html
1 http://www.informit.com/store/go-programming-language-9780134190440
1 http://www.amazon.com/dp/0134190440
...

$ ./exercise8-6 -depth=1 https://golang.org
0 https://golang.org
1 https://golang.org/pkg/
1 https://support.eji.org/give/153413/#!/donation/checkout
1 https://google.com
1 https://golang.org/
1 https://golang.org/help/
1 https://golang.org/project/
1 https://golang.org/doc/
1 https://golang.org/doc/copyright.html
1 https://golang.org/doc/tos.html
1 https://play.golang.org/
1 https://tour.golang.org/
1 https://blog.golang.org/
1 http://golang.org/issues/new?title=x/website:
1 https://golang.org/dl/
1 http://www.google.com/intl/en/policies/privacy/
1 https://golang.org/blog/

$ ./exercise8-6 -depth=1 https://golang.org http://gopl.io/
0 http://gopl.io/
0 https://golang.org
1 https://golang.org/pkg/
1 https://tour.golang.org/
1 https://blog.golang.org/
1 https://google.com
1 https://golang.org/doc/copyright.html
1 https://golang.org/project/
1 https://support.eji.org/give/153413/#!/donation/checkout
1 http://www.google.com/intl/en/policies/privacy/
1 http://golang.org/issues/new?title=x/website:
1 https://golang.org/
1 https://golang.org/dl/
1 https://golang.org/blog/
1 https://golang.org/doc/tos.html
1 https://golang.org/doc/
1 https://golang.org/help/
1 https://play.golang.org/
1 http://www.amazon.com/dp/020161586X?tracking_id=disfordig-20
1 http://www.informit.com/store/go-programming-language-9780134190440
1 http://www.amazon.com/dp/0134190440
1 http://www.barnesandnoble.com/w/1121601944
1 http://www.gopl.io/ch1.pdf
...

The crawler is now highly concurrent and prints a storm of URLs.
*/

// Now the concurrent crawler runs about 20 times faster than the breadth-first
// crawler from Section 5.6, without errors, and terminates correctly if it
// should complete its task.
