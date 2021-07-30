// Fetch saves the contents of a URL into a local file.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// The example below is an improved fetch program (§1.5) that writes the HTTP
// response to a local file instead of to the standard output. It derives the
// file name from the last component of the URL path, which it obtains using the
// `path.Base` function.

// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" || local == "." {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return local, n, err
}

func main() {
	for _, url := range os.Args[1:] {
		local, n, err := fetch(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s => %s (%d bytes)\n", url, local, n)
	}
}

/*
Run:
$ go build gopl.io/ch5/fetch

$ ./fetch https://golang.org
https://golang.org => index.html (9951 bytes)

$ ./fetch https://golang.org/doc
https://golang.org/doc/ => doc (18722 bytes)
*/
