// The du4 command computes the disk usage of the files in a directory.
package main

// The du4 variant includes cancellation:
// it terminates quickly when the user hits return.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "show verbose progress messages")

// A cancellation channel on which no values are ever sent, but whose closure
// indicates that it is time for the program to stop what it is doing.
var done = make(chan struct{})

// A utility function that checks or polls the cancellation state at the instant
// it is called.
func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// This program reports the disk usage of one or more directories specified on
// the command line, like the Unix `du` command. Most of its work is done by the
// `walkDir` function below, which enumerates the entries of the directory `dir`
// using the `dirents` helper function.
func main() {
	// ...determine roots...

	flag.Parse()

	// Determine the initial directories.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Cancel traversal when input is detected.
	go func() {
		// A goroutine that will read from the standard input, which is
		// typically connected to the terminal. As soon as any input is read
		// (for instance, the user presses the return key), this goroutine
		// broadcasts the cancellation by closing the `done` channel.
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
	}()

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var n sync.WaitGroup

	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically.
	//
	// The main goroutine now uses a ticker to generate events every 500ms, and
	// a select statement to wait for either a file size message, in which case
	// it updates the totals, or a tick event, in which case it prints the
	// current totals. If the `-v` flag is not specified, the `tick` channel
	// remains nil, and its case in the select is effectively disabled.
	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case <-done:
			// Drain fileSizes to allow existing goroutines to finish.
			// (see note 1 below for more explaination)
			for range fileSizes {
				// Do nothing.
			}
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes) // running totals
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals

	// Since the program no longer uses a `range` loop, the first `select` case
	// must explicitly test whether the `fileSizes` channel has been closed,
	// using the two-result form of receive operation. If the channel has been
	// closed, the program breaks out of the loop. The labeled `break` statement
	// breaks out of both the `select` and the `for` loop; an unlabeled `break`
	// would break out of only the `select`, causing the loop to begin the next
	// iteration.
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	// (see note 2 below for explanation)
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

//!+sema
// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	// (see note 3 below for explanation)
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token
	//!-sema

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}

/*
Now, when cancellation occurs, all the background goroutines quickly stop and
the main function returns. Of course, when main returns, a program exits, so it
can be hard to tell a main function that cleans up after itself from one that
does not.

$ go build gopl.io/ch8/du4
$ ./du4 $HOME
...
1248639 files  197.8 GB

$ ./du4
991 files  0.0 GB

$ ./du4 $HOME/Downloads .
38052 files  165.7 GB

$ ./du4 -v $HOME/Downloads
24869 files  38.2 GB
37061 files  165.6 GB

$ ./du4 -v $HOME
33738 files  0.4 GB
63552 files  0.8 GB
93921 files  1.0 GB
...
1219351 files  197.4 GB
1236915 files  197.7 GB
1248664 files  197.8 GB
*/

// Note 1:
// The function returns if this case is ever selected, but before it returns it
// must first drain the `fileSizes` channel, discarding all values until the
// channel is closed. It does this to ensure that any active calls to `walkDir`
// can run to completion without getting stuck sending to `fileSizes`.

// Note 2:
// The `walkDir` goroutine polls the cancellation status when it begins, and
// returns without doing anything if the status is set. This turns all
// goroutines created after cancellation into no-ops.

// Note 3:
// A little profiling of this program revealed that the bottleneck was the
// acquisition of a semaphore token in `dirents`. The `select` below makes this
// operation cancellable and reduces the typical cancellation latency of the
// program from hundreds of milliseconds to tens.
