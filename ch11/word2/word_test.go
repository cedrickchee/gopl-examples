package word

import "testing"

func TestIsPalindrome(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"kayak", true},
		{"detartrated", true},
		{"A man, a plan, a canal: Panama", true},
		{"Evil I did dwell; lewd did I live.", true},
		{"Able was I ere I saw Elba", true},
		{"été", true},
		{"Et se resservir, ivresse reste.", true},
		{"palindrome", false}, // non-palindrome
		{"desserts", false},   // semi-palindrome
	}
	for _, test := range tests {
		if got := IsPalindrome(test.input); got != test.want {
			t.Errorf("IsPalindrome(%q) = %v", test.input, got)
		}
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	// A benchmark for `IsPalindrome` that calls it `N` times in a loop.
	for i := 0; i < b.N; i++ {
		IsPalindrome("A man, a plan, a canal: Panama")
	}
}

/*
We run the benchmark with the command below.
$ cd $GOPATH/src/gopl.io/ch11/word2
$ go test -bench=.
goos: linux
goarch: amd64
pkg: gopl.io/ch11/word2
cpu: Intel(R) Core(TM) i7-1065G7 CPU @ 1.30GHz
BenchmarkIsPalindrome-8          2261197               523.7 ns/op
PASS
ok   gopl.io/ch11/word2 1.725s
*/

/*
Performance Optimization 1:

Now that we have a benchmark and tests, it’s easy to try out ideas for making
the program faster. Perhaps the most obvious optimization is to make
`IsPalindrome`'s second loop stop checking at the midpoint, to avoid doing each
comparison twice.

But as is often the case, an obvious optimization doesn’t always yield the
expected benefit. This one delivered a mere 4% improvement in one experiment.

$ go test -bench=. gopl.io/ch11/word2
goos: linux
goarch: amd64
pkg: gopl.io/ch11/word2
cpu: Intel(R) Core(TM) i7-1065G7 CPU @ 1.30GHz
BenchmarkIsPalindrome-8          2543506               482.0 ns/op
PASS
ok   gopl.io/ch11/word2 1.710s

$ go test -bench=. -benchmem gopl.io/ch11/word2
goos: linux
goarch: amd64
pkg: gopl.io/ch11/word2
cpu: Intel(R) Core(TM) i7-1065G7 CPU @ 1.30GHz
BenchmarkIsPalindrome-8          2434116               487.5 ns/op           248 B/op          5 allocs/op
PASS
ok   gopl.io/ch11/word2 1.692s

and after performance optimization 2:

$ go test -bench=. -benchmem gopl.io/ch11/word2
goos: linux
goarch: amd64
pkg: gopl.io/ch11/word2
cpu: Intel(R) Core(TM) i7-1065G7 CPU @ 1.30GHz
BenchmarkIsPalindrome-8          4855659               254.8 ns/op           128 B/op          1 allocs/op
PASS
ok   gopl.io/ch11/word2 1.493s
*/

// Consolidating the allocations in a single call to `make` eliminated 75% of
// the allocations and halved the quantity of allocated memory.
