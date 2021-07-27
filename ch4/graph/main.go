// Graph shows how to use a map of maps to represent a directed graph.
package main

import "fmt"

// The key type of graph is string and the value type is
// representing a set of strings.
// So, a string map to a set of strings.
var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	edges := graph[from]
	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}

func main() {
	// edges:
	// a -- b
	// c -- d
	// a -- d
	// d -- a

	// graph:
	//   a -- b
	//    \
	// c -- d
	addEdge("a", "b")
	addEdge("c", "d")
	addEdge("a", "d")
	addEdge("d", "a")
	fmt.Println(hasEdge("a", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("a", "d"))
	fmt.Println(hasEdge("d", "a"))
	fmt.Println(hasEdge("x", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("x", "d"))
	fmt.Println(hasEdge("d", "x"))

	fmt.Println(hasEdge("c", "a"))
	fmt.Println(hasEdge("b", "d"))
}
