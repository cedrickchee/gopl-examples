// Package bank provides a concurrency-safe bank with one account.
package bank

var (
	sema    = make(chan struct{}, 1) // a binary semaphore guarding balance
	balance int
)

func Deposit(amount int) {
	sema <- struct{}{} // acquire token
	balance = balance + amount
	<-sema // release token
}

func Balance() int {
	sema <- struct{}{} // acquire token
	b := balance
	<-sema // release token
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
// In Section 8.6, we used a buffered channel as a _counting semaphore_ to
// ensure that no more than 20 goroutines made simultaneous HTTP requests. With
// the same idea, we can use a channel of capacity 1 to ensure that at most one
// goroutine accesses a shared variable at a time. A semaphore that counts only
// to 1 is called a _binary semaphore_.
