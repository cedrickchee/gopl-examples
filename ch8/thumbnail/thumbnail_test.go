// This file is just a place to put example code from the book.
// It does not actually run any code in gopl.io/ch8/thumbnail.

package thumbnail_test

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopl.io/ch8/thumbnail"
)

const folder = "/home/neo/dev/work/repo/github/gopl-examples/bin/"

var filenames = []string{folder + "pic-1.jpg", folder + "pic-2.jpg"}

// makeThumbnails makes thumbnails of the specified files.
func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

// NOTE: incorrect!
func makeThumbnails2(filenames []string) {
	// makeThumbnails2 returns before it has finished doing what it was supposed
	// to do. It starts all the goroutines, one per file name, but doesn’t wait
	// for them to finish.
	for _, f := range filenames {
		go thumbnail.ImageFile(f) // NOTE: ignoring errors
	}
}

// makeThumbnails3 makes thumbnails of the specified files in parallel.
func makeThumbnails3(filenames []string) {
	// We change the inner goroutine to report its completion to the outer
	// goroutine by sending an event on a shared channel. Since we know that
	// there are exactly `len(filenames)` inner goroutines, the outer goroutine
	// need only count that many events before it returns.

	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f) // NOTE: ignoring errors
			ch <- struct{}{}
		}(f)
	}

	// Wait for goroutines to complete.
	for range filenames {
		<-ch
	}
}

// makeThumbnails4 makes thumbnails for the specified files in parallel.
// It returns an error if any step failed.
func makeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := thumbnail.ImageFile(f)
			errors <- err
		}(f)
	}

	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: incorrect: goroutine leak!
		}
	}

	return nil

	// This function has a subtle bug. When it encounters the first non-nil
	// error, it returns the error to the caller, leaving no goroutine draining
	// the `errors` channel. Each remaining worker goroutine will block forever
	// when it tries to send a value on that channel, and will never terminate.
	// This situation, a goroutine leak, may cause the whole program to get
	// stuck or to run out of memory.
}

// makeThumbnails5 makes thumbnails for the specified files in parallel.
// It returns the generated file names in an arbitrary order,
// or an error if any step failed.
func makeThumbnails5(filenames []string) (thumbfiles []string, err error) {
	type item struct {
		thumbfile string
		err       error
	}
	ch := make(chan item, len(filenames))
	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

// makeThumbnails6 makes thumbnails for each file received from the channel.
// It returns the number of bytes occupied by the files it creates.
func makeThumbnails6(filenames []string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines
	for _, f := range filenames {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb) // OK to ignore error
			sizes <- info.Size()
		}(f)
	}

	// closer
	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}
	return total

	// Our final version of `makeThumbnails`, below, returns the total number of
	// bytes occupied by the new files. Unlike the previous versions, however,
	// it receives the file names not as a slice but over a channel of strings,
	// so we cannot predict the number of loop iterations.
	//
	// To know when the last goroutine has finished (which may not be the last
	// one to start), we need to increment a counter before each goroutine
	// starts and decrement it as each goroutine finishes. This demands a
	// special kind of counter, one that can be safely manipulated from multiple
	// goroutines and that provides a way to wait until it becomes zero. This
	// counter type is known as `sync.WaitGroup`, and the code above shows how
	// to use it.
}

func Example_makeThumbnails() {
	// makeThumbnails(filenames)

	// Output:
	//
}

func Example_makeThumbnails2() {
	// makeThumbnails2(filenames)

	// Output:
	//
}

func Example_makeThumbnails3() {
	// makeThumbnails3(filenames)

	// Output:
	//
}

func Example_makeThumbnails4() {
	// err := makeThumbnails4(filenames)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Output:
	//
}

func Example_makeThumbnails5() {
	// thumbfiles, err := makeThumbnails5(filenames)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(thumbfiles)

	// Output:
	//
}

func Example_makeThumbnails6() {
	n := makeThumbnails6(filenames)
	fmt.Println("num of bytes occupied by the files =", n)

	// Output:
	// num of bytes occupied by the files = 6151
}

/*
Run:
$ go test -v gopl.io/ch8/thumbnail
*/
