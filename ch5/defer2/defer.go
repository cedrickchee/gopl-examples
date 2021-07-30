// Defer2 demonstrates a deferred call to runtime.Stack during a panic.
package main

import (
	"fmt"
	"os"
	"runtime"
)

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

// For diagnostic purposes, the `runtime` package lets the programmer dump the
// stack using the same machinery. By deferring a call to `printStack` in main.
func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}

func main() {
	defer printStack()
	f(3)
}

/*
When run, the program prints the following to the standard output:

$ go run gopl.io/ch5/defer2
...
goroutine 1 [running]:
main.printStack()
        gopl-examples/ch5/defer2/defer.go:18 +0x5b
panic(0x4a8720, 0x546ae0)
        /usr/local/go/src/runtime/panic.go:965 +0x1b9
main.f(0x0)
        gopl-examples/ch5/defer2/defer.go:11 +0x1e5
main.f(0x1)
        gopl-examples/ch5/defer2/defer.go:13 +0x185
main.f(0x2)
        gopl-examples/ch5/defer2/defer.go:13 +0x185
main.f(0x3)
        gopl-examples/ch5/defer2/defer.go:13 +0x185
main.main()
        gopl-examples/ch5/defer2/defer.go:24 +0x4
...
*/
