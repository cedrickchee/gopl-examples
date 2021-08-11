package main

import (
	"math/rand"
	"testing"
	"time"

	word "gopl.io/ch11/word2"
)

func TestRandomPalindromes(t *testing.T) {
	// Initialize a pseudo-random number generator.
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !word.IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}

// Notes
//
// How to test the main package functions in golang?
// https://stackoverflow.com/questions/31352239/how-to-test-the-main-package-functions-in-golang

/*
Test:
$ go test -v gopl.io/ch11/various
=== RUN   TestRandomPalindromes
    main_test.go:14: Random seed: 1628686320342718423
--- PASS: TestRandomPalindromes (0.00s)
PASS
ok   gopl.io/ch11/various       (cached)
*/
