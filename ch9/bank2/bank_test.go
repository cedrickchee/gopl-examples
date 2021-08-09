package bank_test

import (
	"sync"
	"testing"

	bank "gopl.io/ch9/bank2"
)

func TestBank(t *testing.T) {
	// Deposit [1..1000] concurrently.
	var n sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		n.Add(1)
		go func(amount int) {
			bank.Deposit(amount)
			n.Done()
		}(i)
	}
	n.Wait()

	if got, want := bank.Balance(), (1000+1)*1000/2; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}

/*
Run test:
$ go test -v gopl.io/ch9/bank2
=== RUN   TestBank
--- PASS: TestBank (0.00s)
PASS
ok      gopl.io/ch9/bank2       0.004s
*/
