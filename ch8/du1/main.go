// The du1 command computes the disk usage of the files in a directory.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

	// Print the results.
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
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

// The du1 variant uses two goroutines and
// prints the total after every file is found.

/*
This program pauses for a long while before printing its result:

$ go build gopl.io/ch8/du1
$ ./du1 $HOME
du1: open /dev/scratch/experiment/my-shell-scripts/container-root/root: permission denied
du1: open /dev/scratch/experiment/my-shell-scripts/container-root/tmp/fish.root: permission denied
...
1248496 files  197.8 GB

$ ./du1
983 files  0.0 GB

$ ./du1 $HOME/Downloads .
38044 files  165.7 GB
*/
