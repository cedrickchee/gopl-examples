// Title2 prints the title of an HTML document specified by a URL.
// It uses defer to simplify closing the response body stream.
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check Content-Type is HTML (e.g., "text/html; charset=utf-8").
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	// Print doc's title element.
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	forEachNode(doc, visitNode, nil)
	return nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if err := title(arg); err != nil {
			fmt.Fprintf(os.Stderr, "title: %v", err)
		}
	}
}

/*
Run:
$ go run gopl.io/ch5/title2 https://golang.org
The Go Programming Language

$ go run gopl.io/ch5/title2 https://golang.org/doc/effective_go.html
Effective Go - The Go Programming Language

$ go run gopl.io/ch5/title2 https://github.com
GitHub: Where the world builds software · GitHub

$ go run gopl.io/ch5/title2 https://golang.org/doc/gopher/frontpage.png
title: https://golang.org/doc/gopher/frontpage.png has type image/png, not text/html
*/
