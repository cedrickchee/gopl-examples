// Package bank provides a concurrency-safe bank with one account.
package bank

import "sync"

var (
	mu      sync.Mutex
	balance int
)

func Deposit(amount int) {
	mu.Lock()
	balance = balance + amount
	mu.Unlock()
}

func Balance() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}

// Notes
//
// How do we avoid data races in our programs?
//
// The third way to avoid a data race is to allow many goroutines to access the
// variable, but only one at a time. This approach is known as _mutual
// exclusion_.
//
// This pattern of _mutual exclusion_ is so useful that it is supported
// directly by the `Mutex` type from the `sync` package. Its `Lock` method
// acquires the token (called a `lock`) and its `Unlock` method releases it.
