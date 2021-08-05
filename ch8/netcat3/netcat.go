// Netcat is a simple read/write client for TCP servers.
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
	conn, err := net.Dial("tcp", "localhost:8000")
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
	conn.Close()
	<-done // wait for background goroutine to finish
}

// mustCopy function is a utility used in several examples in this section.
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

/*
In the session below, the client’s input is left-aligned and the server’s
responses are indented. The client shouts at the echo server three times:

$ go build gopl.io/ch8/reverb2
$ ./reverb2 &
$ go build gopl.io/ch8/netcat3
$ ./netcat3
Is there anybody there?
    IS THERE ANYBODY THERE?
Yooo-hooo!
    Is there anybody there?
YOOO-HOOO!
    is there anybody there?
    Yooo-hooo!
    yooo-hooo!
^D
$ killall reverb2
*/

// When the user closes the standard input stream, `mustCopy` returns and the
// main goroutine calls `conn.Close()`, closing both halves of the network
// connection. Closing the write half of the connection causes the server to see
// an end-of-file condition. Closing the read half causes the background
// goroutine’s call to `io.Copy` to return a "read from closed connection"
// error, which is why we’ve removed the error logging; Exercise 8.3 suggests a
// better solution. (Notice that the `go` statement calls a literal function, a
// common construction.)
//
// Before it returns, the background goroutine logs a message, then sends a
// value on the `done` channel. The main goroutine waits until it has received
// this value before returning. As a result, the program always logs the "done"
// message before exiting.
