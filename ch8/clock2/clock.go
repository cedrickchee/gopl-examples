// Clock is a TCP server that periodically writes the time.
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

		// The second client must wait until the first client is finished
		// because the server is sequential; it deals with only one client at a
		// time.
		// Just one small change is needed to make the server concurrent: adding
		// the `go` keyword to the call to `handleConn` causes each call to run
		// in its own goroutine.
		go handleConn(conn) // handle connections concurrently
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
Now, multiple clients can receive the time at once.
To connect to the server:

# Terminal 1
$ go build gopl.io/ch8/clock2
$ ./clock2 &
$ ./netcat1
16:40:52
16:40:53
16:40:54
16:40:55
16:40:56
16:40:57
16:40:58
16:40:59
16:41:00
16:41:01
16:41:02
16:41:03
^C

$ killall clock2

# Terminal 2
$ ./netcat1
16:40:55
16:40:56
16:40:57
^C

# Terminal 3
$ ./netcat1
16:40:59
16:41:00
16:41:01
^C
*/
