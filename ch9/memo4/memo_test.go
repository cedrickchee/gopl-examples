package memo_test

import (
	"testing"

	memo "gopl.io/ch9/memo4"
	"gopl.io/ch9/memotest"
)

var httpGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := memo.New(httpGetBody)
	memotest.Sequential(t, m)
}

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

$ go test -v gopl.io/ch9/memo4
=== RUN   Test
https://golang.org, 405.120488ms, 9951 bytes
https://godoc.org, 1.399513367s, 13005 bytes
https://play.golang.org, 838.223279ms, 6364 bytes
http://gopl.io, 2.958558268s, 4154 bytes
https://golang.org, 1.272µs, 9951 bytes
https://godoc.org, 685ns, 13005 bytes
https://play.golang.org, 522ns, 6364 bytes
http://gopl.io, 255ns, 4154 bytes
--- PASS: Test (5.60s)
PASS
(... cut ...)
*/

/*
This test executes all calls to Get concurrently.

We run it with the race detector.

$ go test -run=TestConcurrent -race -v gopl.io/ch9/memo4
=== RUN   TestConcurrent
https://golang.org, 542.468726ms, 9951 bytes
https://golang.org, 542.956336ms, 9951 bytes
https://play.golang.org, 1.120891609s, 6364 bytes
https://play.golang.org, 1.119739071s, 6364 bytes
https://godoc.org, 1.355244237s, 13005 bytes
https://godoc.org, 1.354138542s, 13005 bytes
http://gopl.io, 2.918568965s, 4154 bytes
http://gopl.io, 2.917629606s, 4154 bytes
--- PASS: TestConcurrent (2.92s)
PASS
ok      gopl.io/ch9/memo4       2.945s
*/

// Our concurrent, duplicate-suppressing, non-blocking cache is complete.
