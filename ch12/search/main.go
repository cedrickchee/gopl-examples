// Search is a demo of the params.Unpack function.
package main

import (
	"fmt"
	"log"
	"net/http"

	"gopl.io/ch12/params"
)

// search implements the /search URL endpoint.
// The search function below is an HTTP handler. It defines a variable called
// `data` of an anonymous struct type whose fields correspond to the HTTP
// request parameters. The struct’s field tags specify the parameter names,
// which are often short and cryptic since space is precious in a URL. The
// `Unpack` function populates the struct from the request so that the
// parameters can be accessed conveniently and with an appropriate type.
func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10 // set default
	// Unpack populates the struct from the request so that the parameters can
	// be accessed conveniently and with an appropriate type.
	if err := params.Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// ...rest of handler...
	fmt.Fprintf(resp, "Search: %+v\n", data)
}

func main() {
	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

/*
Output:

$ go build gopl.io/ch12/search
$ ./search &
$ ./fetch 'http://localhost:12345/search'
Search: {Labels:[] MaxResults:10 Exact:false}
$ ./fetch 'http://localhost:12345/search?l=golang&l=programming'
Search: {Labels:[golang programming] MaxResults:10 Exact:false}
$ ./fetch 'http://localhost:12345/search?l=golang&l=programming&max=100'
Search: {Labels:[golang programming] MaxResults:100 Exact:false}
$ ./fetch 'http://localhost:12345/search?x=true&l=golang&l=programming'
Search: {Labels:[golang programming] MaxResults:10 Exact:true}
$ ./fetch 'http://localhost:12345/search?q=hello&x=123'
x: strconv.ParseBool: parsing "123": invalid syntax
$ ./fetch 'http://localhost:12345/search?q=hello&max=lots'
max: strconv.ParseInt: parsing "lots": invalid syntax
*/
