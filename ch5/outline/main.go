// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
		os.Exit(1)
	}
	outline(nil, doc)
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}

/*
Run:
$ go build gopl.io/ch5/outline
$ ./fetch https://golang.org | ./outline

Output:
[html]
[html head]
[html head meta]
[html head meta]
[html head meta]
[html head meta]
[html head title]
[html head link]
[html head link]
[html head link]
[html head script]
[html head script]
[html head script]
[html head script]
[html head script]
[html head script]
[html body]
[html body header]
[html body header div]
[html body header div a]
[html body header nav]
[html body header nav a]
[html body header nav a img]
[html body header nav button]
...
*/
