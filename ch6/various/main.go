// Various code samples for chapter 6.
package main

import (
	"fmt"
	"sync"
)

func main() {
	// =========================================================================
	pointerReceiver()

	// =========================================================================
	pointerReceiverRule()

	// =========================================================================
	nilReceiverValue()

	// =========================================================================
	composingTypesByStructEmbedding()

	// =========================================================================
	methodExpressions()

	// =========================================================================
	encapsulation()
}

type Point struct{ X, Y float64 }

func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}
func pointerReceiver() {
	r := &Point{1, 2}
	r.ScaleBy(2)
	fmt.Printf("%+v\n", *r) // "{2, 4}"
}

type P *int

// func (P) f() { /* ... */ } // compile error: invalid receiver type
func pointerReceiverRule() {
	// The (*Point).ScaleBy method can be called by providing a *Point receiver,
	// like this:
	r := &Point{1, 2}
	r.ScaleBy(2)
	fmt.Printf("%+v\n", *r) // "{2, 4}"

	// or this:
	p := Point{1, 2}
	pptr := &p
	pptr.ScaleBy(2)
	fmt.Printf("%+v\n", *pptr) // "{2, 4}"

	// or this:
	q := Point{1, 2}
	// &q.ScaleBy(2) // compile error: (no value) used as value
	(&q).ScaleBy(2)
	fmt.Printf("%+v\n", q) // "{2, 4}"

	// Shorthand
	// The compiler will perform an implicit &x on the variable. This works only
	// for variables.
	x := Point{1, 2}
	x.ScaleBy(2)
	fmt.Printf("%+v\n", x) // "{2, 4}"
	// We cannot call a *Point method on a non-addressable Point receiver,
	// because there’s no way to obtain the address of a temporary value.
	// Point{1, 2}.ScaleBy(2) // compile error: can't take address of Point literal
}

// An IntList is a linked list of integers.
// A nil *IntList represents the empty list.
type IntList struct {
	Value int
	Tail  *IntList
}

// Sum returns the sum of the list elements.
func (list *IntList) Sum() int {
	if list == nil {
		return 0
	}
	return list.Value + list.Tail.Sum()
}

func nilReceiverValue() {
	// When you define a type whose methods allow nil as a receiver value, it’s
	// worth pointing this out explicitly in its documentation comment, as we
	// did above.
	list := IntList{
		Value: 1,
		Tail: &IntList{
			Value: 2,
			Tail: &IntList{
				Value: 3,
			},
		},
	}
	fmt.Println("sum linked list of integers =", list.Sum())

	list2 := IntList{2, &IntList{4, &IntList{6, nil}}}
	fmt.Println("sum linked list of integers =", list2.Sum())
}

func composingTypesByStructEmbedding() {
	// ### Trick ###
	// Thanks to embedding, it's possible and sometimes useful for unnamed
	// struct types to have methods too.
	//
	// To illustrate. This example shows part of a simple cache implemented
	// using two package-level variables, a mutex and the map that it guards:
	{
		var (
			mu      sync.Mutex // guards mapping
			mapping = make(map[string]string)
		)
		Lookup := func(key string) string {
			mu.Lock()
			v := mapping[key]
			mu.Unlock()
			return v
		}
		fmt.Println("cache lookup =", Lookup("hello"))
	}

	// The version below is functionally equivalent but groups together the two
	// related variables in a single package-level variable, `cache`:
	{
		var cache = struct {
			sync.Mutex
			mapping map[string]string
		}{
			mapping: make(map[string]string),
		}

		Lookup := func(key string) string {
			cache.Lock()
			v := cache.mapping[key]
			cache.Unlock()
			return v
		}
		fmt.Println("cache lookup =", Lookup("hello"))
	}
}

// type Point struct{ X, Y float64 }

func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

type Path []Point

func (path Path) TranslateBy(offset Point, add bool) {
	var op func(p, q Point) Point
	if add {
		op = Point.Add
	} else {
		op = Point.Sub
	}
	for i := range path {
		// Call either path[i].Add(offset) or path[i].Sub(offset).
		path[i] = op(path[i], offset)
	}
}
func methodExpressions() {
	paths := Path{Point{1, 2}, Point{4, 6}}
	fmt.Println(paths) // "[{1 2} {4 6}]"
	paths.TranslateBy(Point{1, 1}, true)
	fmt.Println(paths) // "[{2 3} {5 7}]"

	paths2 := Path{Point{1, 2}, Point{4, 6}}
	paths2.TranslateBy(Point{1, 1}, false)
	fmt.Println(paths2) // "[{0 1} {3 5}]"
}

// Encapsulation provides three benefits.
type Counter struct{ n int }

func (c *Counter) N() int     { return c.n } // getter
func (c *Counter) Increment() { c.n++ }      // setter
func (c *Counter) Reset()     { c.n = 0 }    // setter
func encapsulation() {
	// The third benefit of encapsulation, is that it prevents clients from
	// setting an object’s variables arbitrarily.
	c := &Counter{} // or c := new(Counter)
	fmt.Println("Initial Counter =", c.N())
	c.Increment()
	fmt.Println("Increment Counter =", c.N())
	c.Increment()
	fmt.Println("Increment Counter =", c.N())
	c.Reset()
	fmt.Println("Reset Counter =", c.N())
}
