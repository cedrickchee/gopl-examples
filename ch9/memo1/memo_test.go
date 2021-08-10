package memo_test

import (
	"testing"

	memo "gopl.io/ch9/memo1"
	"gopl.io/ch9/memotest"
)

var httpGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := memo.New(httpGetBody)
	memotest.Sequential(t, m)
}

// NOTE: not concurrency-safe!  Test fails.
func TestConcurrent(t *testing.T) {
	m := memo.New(httpGetBody)
	memotest.Concurrent(t, m)
}

// Notes
//
// We use the testing package to systematically investigate the effect of
// memoization. From the test output below, we see that the URL stream contains
// duplicates, and that although the first call to `(*Memo).Get` for each URL
// takes hundreds of milliseconds, the second request returns the same amount of
// data in under a millisecond.

/*
This test executes all calls to Get sequentially:

$ go test -v gopl.io/ch9/memo1
=== RUN   Test
https://golang.org, 359.646118ms, 9951 bytes
https://godoc.org, 1.457182462s, 13005 bytes
https://play.golang.org, 1.100813355s, 6364 bytes
http://gopl.io, 2.682721174s, 4154 bytes
https://golang.org, 605ns, 9951 bytes
https://godoc.org, 268ns, 13005 bytes
https://play.golang.org, 194ns, 6364 bytes
http://gopl.io, 212ns, 4154 bytes
--- PASS: Test (5.60s)
PASS
ok      gopl.io/ch9/memo1       5.605s
*/

/*
This test executes all calls to Get concurrently.

The test runs much faster, but unfortunately it is unlikely to work correctly
all the time. We may notice unexpected cache misses, or cache hits that return
incorrect values, or even crashes.

Worse, it is likely to work correctly _some_ of the time, so we may not even
notice that it has a problem. But if we run it with the `-race` flag, the race
detector often prints a report such as this one:

$ go test -run=TestConcurrent -race -v gopl.io/ch9/memo1
=== RUN   TestConcurrent
https://play.golang.org, 1.113718178s, 6364 bytes
==================
WARNING: DATA RACE
Write at 0x00c00011ede0 by goroutine 23:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:202 +0x0
  gopl.io/ch9/memo1.(*Memo).Get()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memo1/memo.go:30 +0x212
  gopl.io/ch9/memotest.Concurrent.func1()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memotest/memotest.go:70 +0xe1

Previous write at 0x00c00011ede0 by goroutine 11:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:202 +0x0
  gopl.io/ch9/memo1.(*Memo).Get()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memo1/memo.go:30 +0x212
  gopl.io/ch9/memotest.Concurrent.func1()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memotest/memotest.go:70 +0xe1

Goroutine 23 (running) created at:
  gopl.io/ch9/memotest.Concurrent()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memotest/memotest.go:67 +0x10c
  gopl.io/ch9/memo1_test.TestConcurrent()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memo1/memo_test.go:20 +0xdd
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1194 +0x202

Goroutine 11 (finished) created at:
  gopl.io/ch9/memotest.Concurrent()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memotest/memotest.go:67 +0x10c
  gopl.io/ch9/memo1_test.TestConcurrent()
      /home/neo/dev/work/repo/github/gopl-examples/ch9/memo1/memo_test.go:20 +0xdd
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1194 +0x202
==================

(... cut ...)

https://play.golang.org, 1.114864038s, 6364 bytes
https://godoc.org, 1.189437307s, 13005 bytes
https://godoc.org, 1.730792453s, 13005 bytes
https://golang.org, 2.457317472s, 9951 bytes
https://golang.org, 2.46091782s, 9951 bytes
http://gopl.io, 3.265350269s, 4154 bytes
http://gopl.io, 3.267201029s, 4154 bytes
    testing.go:1093: race detected during execution of test
--- FAIL: TestConcurrent (3.27s)
=== CONT
    testing.go:1093: race detected during execution of test
FAIL
FAIL    gopl.io/ch9/memo1       3.283s
FAIL
*/

// The reference to `memo.go:30` tells us that two goroutines have updated the
// `cache` map without any intervening synchronization. `Get` is not
// concurrency-safe: it has a data race.
