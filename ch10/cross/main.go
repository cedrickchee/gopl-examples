// The cross command prints the values of GOOS and GOARCH for this target.
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOOS, runtime.GOARCH)
}

// Notes
//
// It is straightforward to cross-compile a Go program, that is, to build an
// executable intended for a different operating system or CPU. Just set the
// `GOOS` or `GOARCH` variables during the build. The cross program prints the
// operating system and architecture for which it was built.

/*
The following commands produce 64-bit and 32-bit executables respectively:

$ go build gopl.io/ch10/cross
$ ./cross
linux amd64

$ GOARCH=386 go build gopl.io/ch10/cross
$ ./cross
linux 386

$ GOOS=darwin go build gopl.io/ch10/cross
$ ./cross
darwin amd64
*/
