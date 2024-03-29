// Various code samples for chapter 9.
package main

import (
	"math/rand"
)

func main() {
	// =========================================================================
	// Randomized Testing
	// Create input values (for test cases) according to a pattern so that we
	// know what output to expect.
	// randomPalindrome()
}

// randomPalindrome generates words that are known to be palindromes by
// construction.
// randomPalindrome returns a palindrome whose length and contents
// are derived from the pseudo-random number generator rng.
func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) // random length up to 24
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000)) // random rune up to '\u0999'
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}
