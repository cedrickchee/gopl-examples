// Crawl1 crawls web links starting with the command-line arguments.
//
// We’ll make this version concurrent so that independent calls to crawl can
// exploit the I/O parallelism available in the web.
//
// This version quickly exhausts available file descriptors
// due to excessive concurrent calls to links.Extract.
//
// Also, it never terminates because the worklist is never closed.
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
	worklist := make(chan []string)

	// Start with the command-line arguments.
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- crawl(link)
				}(link)
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

$ go build gopl.io/ch8/crawl1
$ ./crawl1 http://gopl.io/
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
