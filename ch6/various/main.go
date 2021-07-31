// Various code samples for chapter 6.
package main

import "fmt"

func main() {
	// =========================================================================
	pointerReceiver()
}

type Point struct{ X, Y float64 }

func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}
func pointerReceiver() {
	r := &Point{1, 2}
	r.ScaleBy(2)
	fmt.Printf("%+v", *r) // "{2, 4}"
}
