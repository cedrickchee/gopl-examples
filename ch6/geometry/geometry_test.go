package geometry

import (
	"fmt"
	"testing"
)

// Run:
// $ go test -v -run Distance gopl.io/ch6/geometry
func TestTwoPointDistance(t *testing.T) {
	p := Point{1, 2}
	q := Point{4, 6}
	expected := 5.0

	// Test function call
	fdistance := Distance(p, q)
	if fdistance != expected {
		t.Fatalf(`expected '%f' but got '%f'`, expected, fdistance)
	}

	// Test method call
	mdistance := p.Distance(q)
	if mdistance != expected {
		t.Fatalf(`expected '%f' but got '%f'`, expected, mdistance)
	}
}

func TestPathDistance(t *testing.T) {
	perim := Path{
		{1, 1},
		{5, 1},
		{5, 4},
		{1, 1},
	}
	expected := 12.0
	distance := perim.Distance()
	if distance != expected {
		t.Fatalf(`expected '%f' but got '%f'`, expected, distance)
	}
}

func ExampleDistance() {
	p := Point{1, 2}
	q := Point{4, 6}
	fmt.Printf("%.1f\n", Distance(p, q))
	// Output: 5.0
}

func ExamplePoint_Distance() {
	p := Point{1, 2}
	q := Point{4, 6}
	fmt.Printf("%.1f\n", p.Distance(q))
	// Output: 5.0
}

func ExamplePath_Distance() {
	perim := Path{
		{1, 1},
		{5, 1},
		{5, 4},
		{1, 1},
	}
	fmt.Printf("%.1f\n", perim.Distance())
	// Output: 12.0
}
