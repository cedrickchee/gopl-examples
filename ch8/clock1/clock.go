// Clock1 is a TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// The `Listen` function creates a `net.Listener`, an object that listens
	// for incoming connections on a network port, in this case TCP port
	// `localhost:8000`. The listener’s `Accept` method blocks until an incoming
	// connection request is made, then returns a `net.Conn` object representing
	// the connection.
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		handleConn(conn) // handle one connection at a time
	}
}

// handleConn function handles one complete client connection. In a loop, it
// writes the current time, `time.Now()`, to the client. Since `net.Conn`
// satisfies the `io.Writer` interface, we can write directly to it. The loop
// ends when the write fails, most likely because the client has disconnected,
// at which point `handleConn` closes its side of the connection using a
// deferred call to `Close` and goes back to waiting for another connection
// request.
func handleConn(c net.Conn) {
	// net.Conn is an interface type.
	// It is a generic stream-oriented network connection.
	// Multiple goroutines may invoke methods on a Conn simultaneously.

	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

/*
To connect to the server, we’ll need a client program such as nc ("netcat"), a
standard utility program for manipulating network connections:

$ go build gopl.io/ch8/clock1
$ ./clock1 &
15:58:23
15:58:24
15:58:25
^C
*/
