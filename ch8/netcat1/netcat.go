// Netcat1 is a read-only TCP client.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

// This program reads data from the connection and writes it to the standard
// output until an end-of-file condition or an error occurs.
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

// mustCopy function is a utility used in several examples in this section.
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

/*
Let’s run two clients at the same time on different terminals:

# Left terminal
$ go build gopl.io/ch8/netcat1
$ ./netcat1
13:58:54
13:58:55
13:58:56
^C

$ killall clock1

# Right terminal
$ ./netcat1
13:58:57
13:58:58
13:58:59
^C
*/
