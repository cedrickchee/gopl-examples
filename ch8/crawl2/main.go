// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
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

//!+semaphore
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

// In our crawler, items are URLs. The crawl function we’ll supply to
// breadthFirst prints the URL, extracts its links, and returns them so that
// they too are visited.
func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	return list
}

//!-semaphore

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
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- { // <-- fix program never terminates because the worklist is never closed.
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
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

$ go build gopl.io/ch8/crawl2
$ ./crawl2 http://gopl.io/
http://gopl.io/
https://golang.org/help/
https://golang.org/doc/
https://golang.org/blog/
...

The crawler is now highly concurrent and prints a storm of URLs, but it has two
problems. The first problem manifests itself as error messages in the log after
a few seconds of operation.

The initial error message is a surprising report of a DNS lookup failure for a
reliable domain. The subsequent error message reveals the cause: the program
created so many network connections at once that it exceeded the per-process
limit on the number of open files, causing operations such as DNS lookups and
calls to `net.Dial` to start failing.
*/

// Regarding the limit of 20 concurrent HTTP requests.
//
// This previous version quickly exhausts available file descriptors due to
// excessive concurrent calls to links.Extract.
//
// The previous program is _too_ parallel.
// Unbounded parallelism is rarely a good idea since there is always a limiting
// factor in the system.
// The solution is to limit the number of parallel uses of the resource to match
// the level of parallelism that is available.
// A simple way to do that in our example is to ensure that no more than n calls
// to `links.Extract` are active at once, where n is comfortably less than the
// file descriptor limit--20, say.
//
// We can limit parallelism using a buffered channel of capacity n to model a
// concurrency primitive called a _counting semaphore_. Conceptually, each of
// the n vacant slots in the channel buffer represents a token entitling the
// holder to proceed. Sending a value into the channel acquires a token, and
// receiving a value from the channel releases a token, creating a new vacant
// slot. This ensures that at most n sends can occur without an intervening
// receive.
//
// This version rewrite the `crawl` function so that the call to `links.Extract`
// is bracketed by operations to acquire and release a token, thus ensuring that
// at most 20 calls to it are active at one time.
// It’s good practice to keep the semaphore operations as close as possible to
// the I/O operation they regulate.

// The second problem is that the program never terminates, even when it has
// discovered all the links reachable from the initial URLs. (Of course, you’re
// unlikely to notice this problem unless you choose the initial URLs carefully
// or implement the depth-limiting feature of Exercise 8.6.) For the program to
// terminate, we need to break out of the main loop when the worklist is empty
// _and_ no crawl goroutines are active.
//
// In this version, the counter `n` keeps track of the number of sends to the
// worklist that are yet to occur. Each time we know that an item needs to be
// sent to the worklist, we increment `n`, once before we send the initial
// command-line arguments, and again each time we start a crawler goroutine. The
// main loop terminates when `n` falls to zero, since there is no more work to
// be done.

// Now the concurrent crawler runs about 20 times faster than the breadth-first
// crawler from Section 5.6, without errors, and terminates correctly if it
// should complete its task.
