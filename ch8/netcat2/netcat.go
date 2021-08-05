// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

// While the main goroutine reads the standard input and sends it to the server,
// a second goroutine reads and prints the server’s response. When the main
// goroutine encounters the end of the input, for example, after the user types
// Control-D (`^D`) at the terminal (or the equivalent Control-Z on Microsoft
// Windows), the program stops, even if the other goroutine still has work to
// do.
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// Reads data from the connection and writes it to the standard
	// output until an end-of-file condition or an error occurs.
	go mustCopy(os.Stdout, conn)
	// Reads data from the standard input and writes it to the connection.
	mustCopy(conn, os.Stdin)
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

$ go build gopl.io/ch8/reverb1
$ ./reverb1 &
$ go build gopl.io/ch8/netcat2
$ ./netcat2
Hello?
    HELLO?
    Hello?
    hello?
Is there anybody there?
    IS THERE ANYBODY THERE?
Yooo-hooo!
    Is there anybody there?
    is there anybody there?
    YOOO-HOOO!
    Yooo-hooo!
    yooo-hooo!
^D
$ killall reverb1
*/
