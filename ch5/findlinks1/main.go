// Findlinks1 prints the links in an HTML document read from standard input.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

//  visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

/*
Run:
$ go mod tidy
go: finding module for package golang.org/x/net/html
go: downloading golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985
go: found golang.org/x/net/html in golang.org/x/net
v0.0.0-20210726213435-c6fcb2dbf985

$ go build gopl.io/ch1/fetch
$ go build gopl.io/ch5/findlinks1

$ ./fetch https://golang.org | ./findlinks1
https://support.eji.org/give/153413/#!/donation/checkout
/
/doc/
/pkg/
/project/
/help/
/blog/
https://play.golang.org/
/dl/
https://tour.golang.org/
https://blog.golang.org/
/doc/copyright.html
/doc/tos.html
http://www.google.com/intl/en/policies/privacy/
http://golang.org/issues/new?title=x/website:
https://google.com
*/
