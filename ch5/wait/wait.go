// The wait program waits for an HTTP server to start responding.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetPrefix("WAIT: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

// WaitForServer attempts to contact the server of a URL.
// It tries for one minute using exponential back-off.
// It reports an error if all attempts fail.
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // success
		}
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: wait URL\n")
		os.Exit(1)
	}
	url := os.Args[1]
	if err := WaitForServer(url); err != nil {
		log.Fatalf("Site is down: %v\n", err)
	}
}

/*
Run:
$go run gopl.io/ch5/wait https://githubabc.com

Output:
2021/07/29 15:07:42 server not responding (Head "https://githubabc.com": dial tcp: lookup githubabc.com on 10.6.252.1:53: no such host); retrying...
2021/07/29 15:07:42 exp back-off 1s
2021/07/29 15:07:43 server not responding (Head "https://githubabc.com": dial tcp: lookup githubabc.com on 10.6.252.1:53: no such host); retrying...
2021/07/29 15:07:43 exp back-off 2s
2021/07/29 15:07:45 server not responding (Head "https://githubabc.com": dial tcp: lookup githubabc.com on 10.6.252.1:53: no such host); retrying...
2021/07/29 15:07:45 exp back-off 4s
2021/07/29 15:07:49 server not responding (Head "https://githubabc.com": dial tcp: lookup githubabc.com on 10.6.252.1:53: no such host); retrying...
2021/07/29 15:07:49 exp back-off 8s
2021/07/29 15:07:57 server not responding (Head "https://githubabc.com": dial tcp: lookup githubabc.com on 10.6.252.1:53: no such host); retrying...
2021/07/29 15:07:57 exp back-off 16s
2021/07/29 15:08:13 server not responding (Head "https://githubabc.com": dial tcp: lookup githubabc.com on 10.6.252.1:53: no such host); retrying...
2021/07/29 15:08:13 exp back-off 32s
Site is down: server https://githubabc.com failed to respond after 1m0s
exit status 1
*/
