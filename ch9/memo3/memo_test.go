package memo_test

import (
	"testing"

	memo "gopl.io/ch9/memo3"
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

$ go test -v gopl.io/ch9/memo3
=== RUN   Test
https://golang.org, 363.970824ms, 9951 bytes
https://godoc.org, 1.558894391s, 13005 bytes
https://play.golang.org, 881.142262ms, 6364 bytes
http://gopl.io, 2.883468994s, 4154 bytes
https://golang.org, 707ns, 9951 bytes
https://godoc.org, 324ns, 13005 bytes
https://play.golang.org, 259ns, 6364 bytes
http://gopl.io, 229ns, 4154 bytes
--- PASS: Test (5.69s)
PASS
(... cut ...)
*/

/*
This test executes all calls to Get concurrently.

We run it with the race detector.

$ go test -run=TestConcurrent -race -v gopl.io/ch9/memo3
=== RUN   TestConcurrent
https://golang.org, 545.080001ms, 9951 bytes
https://golang.org, 553.03443ms, 9951 bytes
https://play.golang.org, 1.101588445s, 6364 bytes
https://play.golang.org, 1.126848279s, 6364 bytes
https://godoc.org, 1.358854701s, 13005 bytes
http://gopl.io, 1.498545902s, 4154 bytes
http://gopl.io, 1.506229733s, 4154 bytes
https://godoc.org, 2.070201618s, 13005 bytes
--- PASS: TestConcurrent (2.07s)
PASS
ok      gopl.io/ch9/memo3       2.101s
*/

// The performance improves again, but now we notice that some URLs are being fetched twice.
// This happens when two or more goroutines call `Get` for the same URL at about the same time.
// Both consult the cache, find no value there, and then call the slow function `f`. Then both of
// them update the map with the result they obtained. One of the results is overwritten by the
// other.
