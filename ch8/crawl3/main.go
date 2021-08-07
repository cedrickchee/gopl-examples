// Crawl3 crawls web links starting with the command-line arguments.
//
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
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
// func breadthFirst(f func(item string) []string, worklist []string) {
// 	seen := make(map[string]bool)
// 	for len(worklist) > 0 {
// 		items := worklist
// 		for _, item := range items {
// 			if !seen[item] {
// 				seen[item] = true
// 				worklist = append(worklist, f(item)...)
// 			}
// 		}
// 	}
// }

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

// A simple web crawler program that explored the link graph of the web in
// breadth-first order.
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	// OLD VERSION
	// breadthFirst(crawl, os.Args[1:])

	// NEW VERSION: concurrent web crawler
	// The main function resembles breadthFirst. As before, a worklist records
	// the queue of items that need processing, each item being a list of URLs
	// to crawl, but this time, instead of representing the queue using a slice,
	// we use a channel. Each call to `crawl` occurs in its own goroutine and
	// sends the links it discovers back to the worklist.
	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() { worklist <- os.Args[1:] }()

	// Alternative solution to the problem of excessive concurrency.
	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
				// go func(link string) {
				// 	worklist <- crawl(link)
				// }(link)
			}
		}
	}

	// Notice that the initial send of the command-line arguments to the
	// worklist must run in its own goroutine to avoid _deadlock_, a stuck
	// situation in which both the main goroutine and a crawler goroutine
	// attempt to send to each other while neither is receiving.
	// An alternative solution would be to use a buffered channel.
}

/*
Let’s crawl the web starting from http://gopl.io/. Here are some of the
resulting links:

$ go build gopl.io/ch8/crawl3
$ ./crawl3 http://gopl.io/
http://gopl.io/
https://golang.org/help/
https://golang.org/doc/
https://golang.org/blog/
...
2015/07/15 18:22:12 Get ...: dial tcp: lookup blog.golang.org: no such host
2015/07/15 18:22:12 Get ...: dial tcp 23.21.222.120:443: socket:
too many open files

The crawler is now highly concurrent and prints a storm of URLs, but it has two
problems. The first problem manifests itself as error messages in the log after
a few seconds of operation.

The initial error message is a surprising report of a DNS lookup failure for a
reliable domain. The subsequent error message reveals the cause: the program
created so many network connections at once that it exceeded the per-process
limit on the number of open files, causing operations such as DNS lookups and
calls to `net.Dial` to start failing.
*/

// The crawler goroutines are all fed by the same channel, `unseenLinks`. The
// main goroutine is responsible for de-duplicating items it receives from the
// worklist, and then sending each unseen one over the `unseenLinks` channel to
// a crawler goroutine.
//
// The `seen` map is _confined_ within the main goroutine; that is, it can be
// accessed only by that goroutine. Like other forms of information hiding,
// confinement helps us reason about the correctness of a program.
// In all cases, information hiding helps to limit unintended interactions
// between parts of the program.
//
// Links found by `crawl` are sent to the worklist from a dedicated goroutine to
// avoid deadlock.
//
// To save space, we have not addressed the problem of termination in this
// example.
