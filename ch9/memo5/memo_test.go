package memo_test

import (
	"testing"

	memo "gopl.io/ch9/memo5"
	"gopl.io/ch9/memotest"
)

var httpGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.Sequential(t, m)
}

func TestConcurrent(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
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

$ go test -v gopl.io/ch9/memo5
=== RUN   Test
https://golang.org, 353.024556ms, 9951 bytes
https://godoc.org, 1.542283403s, 13005 bytes
https://play.golang.org, 806.449366ms, 6364 bytes
http://gopl.io, 3.131659307s, 4154 bytes
https://golang.org, 3.335µs, 9951 bytes
https://godoc.org, 5.221µs, 13005 bytes
https://play.golang.org, 1.566µs, 6364 bytes
http://gopl.io, 1.447µs, 4154 bytes
--- PASS: Test (5.83s)
PASS
(... cut ...)
*/

/*
This test executes all calls to Get concurrently.

We run it with the race detector.

$ go test -run=TestConcurrent -race -v gopl.io/ch9/memo5
=== RUN   TestConcurrent
https://golang.org, 439.551826ms, 9951 bytes
https://golang.org, 440.126415ms, 9951 bytes
https://play.golang.org, 994.596884ms, 6364 bytes
https://play.golang.org, 993.955255ms, 6364 bytes
https://godoc.org, 1.245839121s, 13005 bytes
https://godoc.org, 1.245135702s, 13005 bytes
http://gopl.io, 2.066749361s, 4154 bytes
http://gopl.io, 2.065703332s, 4154 bytes
--- PASS: TestConcurrent (2.07s)
PASS
ok      gopl.io/ch9/memo5       2.091s
*/
