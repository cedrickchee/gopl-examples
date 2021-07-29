// The toposort program prints the nodes of a Directed Acyclic Graph (DAG) in
// topological order.
//
// As a somewhat academic example of anonymous functions, consider the problem
// of computing a sequence of computer science courses that satisfies the
// prerequisite requirements of each one. The prerequisites are given in the
// `prereqs` table below, which is a mapping from each course to the list of
// courses that must be completed before it.
package main

import (
	"fmt"
	"sort"
)

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

// This kind of problem is known as topological sorting. Conceptually, the
// prerequisite information forms a directed graph with a node for each course
// and edges from each course to the courses that it depends on. The graph is
// acyclic: there is no path from a course that leads back to itself. We can
// compute a valid sequence using depth-first search through the graph with the
// code below.
func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

/*
The output of the `toposort` program is shown below. It is deterministic, an
often-desirable property that doesn’t always come for free. Here, the values of
the prereqs map are slices, not more maps, so their iteration order is
deterministic, and we sor ted the keys of prereqs before making the initial
calls to visitAll.

$ go run gopl.io/ch5/toposort

1:      intro to programming
2:      discrete math
3:      data structures
4:      algorithms
5:      linear algebra
6:      calculus
7:      formal languages
8:      computer organization
9:      compilers
10:     databases
11:     operating systems
12:     networks
13:     programming languages
*/
