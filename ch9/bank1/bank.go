// Package bank provides a concurrency-safe bank with one account.
package bank

// The bank example rewritten with the `balance` variable confined to a monitor
// goroutine called `teller`.

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

// Notes
//
// How do we avoid data races in our programs?
// The second way to avoid a data race is to avoid accessing the variable from
// multiple goroutines.
//
// The `balance` variable is _confined_ to a single goroutine, `teller`.
//
// Since other goroutines cannot access the variable directly, they must use a
// channel to send the confining goroutine a request to query or update the
// variable.
//
// A goroutine that brokers access to a confined variable using channel requests
// is called a _monitor goroutine_ for that variable. In this example, the
// `teller` goroutine monitors access to the `balance` variable.
