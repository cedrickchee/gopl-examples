// Exercise 8.3: In `netcat3`, the interface value `conn` has the concrete type
// `*net.TCPConn`, which represents a TCP connection. A TCP connection consists
// of two halves that may be closed independently using its `CloseRead` and
// `CloseWrite` methods. Modify the main goroutine of `netcat3` to close only
// the write half of the connection so that the program will continue to print
// the final echoes from the `reverb1` server even after the standard input has
// been closed. (Doing this for the `reverb2` server is harder ; see Exercise
// 8.4.)
package main

import (
	"io"
	"log"
	"net"
	"os"
)

// The program copies input to the server (connection) in its main goroutine, so
// the program terminates as soon as the input stream closes, even if the
// background goroutine is still working.
// To make the program wait for the background goroutine to complete before
// exiting, we use a channel to synchronize the two goroutines.
func main() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		// Reads data from the connection and writes it to the standard output
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	// Reads data from the standard input and writes it to the connection.
	mustCopy(conn, os.Stdin)
	conn.CloseWrite()
	<-done // wait for background goroutine to finish
}

// mustCopy function is a utility used in several examples in this section.
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}

/*
In the session below, the client’s input is left-aligned and the server’s
responses are indented. The client shouts at the echo server three times:

$ go build gopl.io/ch8/reverb1
$ ./reverb1 &
$ go build gopl.io/ch8/exercise8-3
$ ./exercise8-3
Is there anybody there?
    IS THERE ANYBODY THERE?
Yooo-hooo!
    Is there anybody there?
YOOO-HOOO!
    is there anybody there?
^D
    Yooo-hooo!
    yooo-hooo!
$ killall reverb1
*/
