// Exercise 1.8: Modify fetch to add the prefix http:// to each argument URL if
// it is missing. You might want to use strings.HasPrefix.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		fullUrl := url
		if !strings.HasPrefix(url, "http://") {
			fullUrl = "http://" + url
		}
		resp, err := http.Get(fullUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch123: %v\n", err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch123: reading %s: %v\n", fullUrl, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}
