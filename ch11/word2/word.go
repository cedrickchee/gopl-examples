// Package word provides utilities for word games.
package word

import "unicode"

// IsPalindrome reports whether s reads the same forward and backward.
// Letter case is ignored, as are non-letters.
func IsPalindrome(s string) bool {
	// Performance optimization 2:
	// Pre-allocate a sufficiently large array for use by letters, rather than
	// expand it by successive calls to append.
	// var letters []rune // previously
	letters := make([]rune, 0, len(s)) // now
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	// Performance optimization 1:
	// Make the loop stop checking at the midpoint, to avoid doing each
	// comparison twice.
	n := len(letters) / 2
	for i := 0; i < n; i++ {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}
