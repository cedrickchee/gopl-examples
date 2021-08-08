// The du2 command computes the disk usage of the files in a directory.
package main

// The du2 variant uses select and a time.Ticker
// to print the totals periodically if -v is set.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")

// This program reports the disk usage of one or more directories specified on
// the command line, like the Unix `du` command. Most of its work is done by the
// `walkDir` function below, which enumerates the entries of the directory `dir`
// using the `dirents` helper function.
func main() {
	// The main function, shown below, uses two goroutines. The background
	// goroutine calls `walkDir` for each directory specified on the command
	// line and finally closes the `fileSizes` channel. The main goroutine
	// computes the sum of the file sizes it receives from the channel and
	// finally prints the total.

	// ...start background goroutine...

	// Determine the initial directories.
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
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
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
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
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

/*
The program now gives us a leisurely stream of updates:

$ go build gopl.io/ch8/du2
$ ./du2 $HOME
...
1248639 files  197.8 GB

$ ./du2
991 files  0.0 GB

$ ./du2 $HOME/Downloads .
38052 files  165.7 GB

$ ./du2 -v $HOME/Downloads
24869 files  38.2 GB
37061 files  165.6 GB

$ ./du2 -v $HOME
33738 files  0.4 GB
63552 files  0.8 GB
93921 files  1.0 GB
...
1219351 files  197.4 GB
1236915 files  197.7 GB
1248664 files  197.8 GB
*/
