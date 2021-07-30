// Defer1 demonstrates a deferred call being invoked during a panic.
package main

import "fmt"

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

func main() {
	f(3)
}

/*
When run, the program prints the following to the standard output:

$ go run gopl.io/ch5/defer1
f(3)
f(2)
f(1)
defer 1
defer 2
defer 3
*/

//
// A panic occurs during the call to f(0), causing the three deferred calls to
// fmt.Printf to run. Then the runtime terminates the program, printing the
// panic message and a stack dump to the standard error stream (simplified for
// clarity):
//
/*
panic: runtime error: integer divide by zero
goroutine 1 [running]:
main.f(0x0)
        gopl-examples/ch5/defer1/defer.go:7 +0x1e5
main.f(0x1)
        gopl-examples/ch5/defer1/defer.go:9 +0x185
main.f(0x2)
        gopl-examples/ch5/defer1/defer.go:9 +0x185
main.f(0x3)
        gopl-examples/ch5/defer1/defer.go:9 +0x185
main.main()
        gopl-examples/ch5/defer1/defer.go:13 +0x2a
*/
