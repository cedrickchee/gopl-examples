package memo_test

import (
	"testing"

	memo "gopl.io/ch9/memo2"
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

$ go test -v gopl.io/ch9/memo2
=== RUN   Test
https://golang.org, 458.99344ms, 9951 bytes
https://godoc.org, 1.986294098s, 13005 bytes
https://play.golang.org, 921.596661ms, 6364 bytes
http://gopl.io, 3.917147764s, 4154 bytes
https://golang.org, 644ns, 9951 bytes
https://godoc.org, 238ns, 13005 bytes
https://play.golang.org, 194ns, 6364 bytes
http://gopl.io, 227ns, 4154 bytes
--- PASS: Test (7.28s)
PASS
(... cut ...)
*/

/*
This test executes all calls to Get concurrently.

We run it with the race detector.
The test runs much slower.

$ go test -run=TestConcurrent -race -v gopl.io/ch9/memo2
=== RUN   TestConcurrent
https://golang.org, 1.041694592s, 9951 bytes
https://godoc.org, 2.190742852s, 13005 bytes
https://play.golang.org, 3.234881659s, 6364 bytes
http://gopl.io, 8.539341916s, 4154 bytes
https://godoc.org, 8.538987952s, 13005 bytes
https://golang.org, 8.539252883s, 9951 bytes
https://play.golang.org, 8.538840134s, 6364 bytes
http://gopl.io, 8.538693967s, 4154 bytes
--- PASS: TestConcurrent (8.54s)
PASS
ok      gopl.io/ch9/memo2       8.566s
*/

// Now the race detector is silent, even when running the tests concurrently.
// Unfortunately this change to `Memo` reverses our earlier performance gains.
// By holding the lock for the duration of each call to `f`, `Get` serializes
// all the I/O operations we intended to parallelize. What we need is a
// _non-blocking_ cache, one that does not serialize calls to the function it
// memoizes.
