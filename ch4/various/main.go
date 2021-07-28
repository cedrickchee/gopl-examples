// Various code samples for chapter 4.
package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	// =========================================================================
	arrays()

	// =========================================================================
	compareArray()

	// =========================================================================
	var a [32]byte = [32]byte{31: byte('a')}
	fmt.Println(a)
	zero(&a)
	fmt.Println(a)

	// =========================================================================
	var b [32]byte = [32]byte{31: byte('b')}
	fmt.Println(b)
	anotherZero(&b)
	fmt.Println(b)

	// =========================================================================
	arrayAndSlices()

	// =========================================================================
	compareSlices()

	// =========================================================================
	nilSlice()

	// =========================================================================
	appendSlices()

	// =========================================================================
	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s, 2)) // "[5 6 8 9]"
	// remove2
	s2 := []int{1, 2, 3, 4, 5}
	fmt.Println(remove2(s2, 2)) // "[1 2 5 4]"

	// =========================================================================
	maps()

	// =========================================================================
	sortMap()

	// =========================================================================
	zeroValueMap()

	// =========================================================================
	mapLookup()

	// =========================================================================
	mapComparison()

	// =========================================================================
	mapWhoseKeysAreSlices()

	// =========================================================================
	structs()

	// =========================================================================
	emptyStruct()

	// =========================================================================
	structLiterals()

	// =========================================================================
	structValues()

	// =========================================================================
	comparingStructs()

	// =========================================================================
	structEmbeddingAnonFields()
}

func arrays() {
	// Array literal
	var q [3]int = [3]int{1, 2, 3}
	fmt.Println(q[2]) // "3"
	var r [3]int = [3]int{1, 2}
	fmt.Println(r[2]) // "0"

	// In an array literal, if an ellipsis "..."" appears in place of the
	// length, the array length is determined by the number of initializers.
	p := [...]int{1, 2, 3, 4}
	fmt.Printf("%T\n", p) // "[4]int"

	// literal syntax
	type Currency int
	const (
		USD Currency = iota
		EUR
		GBP
		RMB
	)
	symbol := [...]string{USD: "$", EUR: "€", GBP: "£", RMB: "¥"}
	fmt.Println(RMB, symbol[RMB]) // "3 ¥"

	// Defines an array with 100 elements, all zero except for the last,
	// which has value -1.
	d := [...]int{99: -1}
	fmt.Println(d)
}

func compareArray() {
	a := [2]int{1, 2}
	b := [...]int{1, 2}
	c := [2]int{1, 3}
	fmt.Println(a == b, a == c, b == c) // "true false false"
	/*
		d := [3]int{1, 2}
		fmt.Println(a == d) // compile error: cannot compare [2]int == [3]int
	*/
}

// Explicitly pass a pointer to an array so that any modifications the function
// makes to array elements will be visible to the caller.
func zero(ptr *[32]byte) {
	// Zeroes the contents of a [32]byte array.
	for i := range ptr {
		ptr[i] = 0
	}
}

func anotherZero(ptr *[32]byte) {
	// The array literal [32]byte{} yields an array of 32 bytes.
	// Each element of the array has the zero value for byte, which is zero.
	*ptr = [32]byte{}
}

func arrayAndSlices() {
	// Array
	months := [...]string{1: "January", 2: "February", 3: "March", 4: "April",
		5: "May", 6: "June", 7: "July", 8: "August", 9: "September",
		10: "October", 11: "November", 12: "December"}
	fmt.Printf("months = %v, type = %[1]T\n", months)
	// Slices
	Q2 := months[4:7]
	summer := months[6:9]
	fmt.Printf("Q2 value=%v, type=%[1]T\n", Q2)         // "Q2 value=[April May June], type=[]string"
	fmt.Printf("summer value=%v, type=%[1]T\n", summer) // "summer value=[June July August], type=[]string"

	// June is included in each and is the sole output of this (inefficient)
	// test for common elements.
	for _, s := range summer {
		for _, q := range Q2 {
			if s == q {
				fmt.Printf("%s appears in both\n", s)
			}
		}
	}

	// Slicing beyond cap(s) causes a panic, but slicing beyond len(s) extends
	// the slice, so the result may be longer than the original.
	fmt.Printf("Q2 length = %d, capacity = %d\n", len(Q2), cap(Q2))
	fmt.Printf("summer length = %d, capacity = %d\n", len(summer), cap(summer))
	// fmt.Println(summer[:20])    // panic: runtime error: slice bounds out of range [:20] with capacity 7
	endlessSummer := summer[:5]                   // extend a slice (within capacity)
	fmt.Println("endlessSummer =", endlessSummer) // "[June July August September October]"
}

// Unlike arrays, slices are not comparable, so we cannot use == to test whether
// two slices contain the same elements.
// We must do the comparison ourselves.
func compareSlices() {
	// comparing two slices of string
	equal := func(x, y []string) bool {
		// Deep equality test.
		if len(x) != len(y) {
			return false
		}
		for i := range x {
			if x[i] != y[i] {
				return false
			}
		}
		return true
	}

	a := []string{"apple", "orange", "lemon"}
	b := []string{"apple", "orange", "lemon"}
	c := []string{"orange", "lemon", "mango"}
	fmt.Println(equal(a, b)) // "true"
	fmt.Println(equal(a, c)) // "false"
}

func nilSlice() {
	// Slices comparison

	// In Go, the safest choice is to disallow slice comparisons altogether.
	// The only legal slice comparison is against nil.

	// The zero value of a slice type is nil.
	// A nil slice has no underlying array. The nil slice has length and
	// capacity zero, but there are also non-nil slices of length and
	// capacity zero, such as `[]int{}` or `make([]int, 3)[3:]`.
	var s []int
	if s == nil {
		fmt.Println("s is nil slice")
		fmt.Printf("s len = %d, cap = %d\n", len(s), cap(s))
	}
	p := []int{}
	if p != nil {
		fmt.Println("p is non-nil slice")
		fmt.Printf("p len = %d, cap = %d\n", len(p), cap(p))
	}
	q := make([]int, 3)[3:]
	if q != nil {
		fmt.Println("q is non-nil slice")
		fmt.Printf("q len = %d, cap = %d\n", len(q), cap(q))
	}

	// As with any type that can have `nil` values, the nil value of a
	// particular slice type can be written using a conversion expression
	// such as `[]int(nil)`.
	var r []int
	fmt.Printf("r = %v, len = %d, is nil = %t\n", r, len(r), (r == nil)) // "r = [], len = 0, is nil = true"
	r = nil
	fmt.Printf("r = %v, len = %d, is nil = %t\n", r, len(r), (r == nil)) // "r = [], len = 0, is nil = true"
	r = []int(nil)
	fmt.Printf("r = %v, len = %d, is nil = %t\n", r, len(r), (r == nil)) // "r = [], len = 0, is nil = true"
	r = []int{}
	fmt.Printf("r = %v, len = %d, is nil = %t\n", r, len(r), (r == nil)) // "r = [], len = 0, is nil = false"

	// So, if you need to test whether a slice is empty, use len(s) == 0,
	// not s == nil .
}

func appendSlices() {
	var runes []rune
	for _, r := range "Hello, 世界" {
		runes = append(runes, r)
	}
	fmt.Printf("%q\n", runes) // "['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']"

	// This specific problem is more conveniently solved by using the built-in
	// conversion.
	fmt.Printf("%q\n", []rune("Hello, 世界"))

	// Add more than one new element, or even a whole slice of them.
	var x []int
	x = append(x, 1)
	x = append(x, 2, 3)
	x = append(x, x...) // append the slice x
	fmt.Println(x)      // "[1 2 3 1 2 3]"
}

// To remove an element from the middle of a slice, preserving the order of the
// remaining elements, use copy to slide the higher-numbered elements down by
// one to fill the gap.
func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

// Remove an element from the middle of a slice, don’t need to preserve the
// order. We can just move the last element into the gap.
func remove2(slice []int, i int) []int {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func maps() {
	ages := make(map[string]int) // mapping from strings to ints
	fmt.Println(ages)            // "map[]"

	ages2 := map[string]int{
		"alice":   31,
		"charlie": 34,
	}
	fmt.Println(ages2) // "map[alice:31 charlie:34]"

	// empty map
	em := map[string]int{}
	fmt.Printf("em value = %v, len = %d\n", em, len(em))

	// Access map elements
	em["orange"] = 2
	em["apple"] = 3
	fmt.Println("orange count =", em["orange"]) // "2"

	// Remove map element
	delete(em, "orange")    // remove element em["orange"]
	fmt.Println("em =", em) // "em = map[apple:3]"

	// A map element is not a variable, and we cannot take its address.
	/*
		_ = &em["apple"] // compile error: cannot take address of map element
	*/

	// Enumerate all the key/value pairs in the map.
	fruits := map[string]int{"strawberry": 2, "lime": 5, "kiwi": 1}
	for fruit, count := range fruits {
		fmt.Printf("%s\t%d\n", fruit, count)
	}
}

func sortMap() {
	ages := map[string]int{
		"bob":   30,
		"alice": 28,
		"john":  32,
	}
	// var names []string
	// Replace the previous line with a more efficient slice. This allocate an
	// array of the required size up front.
	names := make([]string, 0, len(ages))
	for name := range ages {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}
}

// The zero value for a map type is nil, that is,
// a reference to no hash table at all.
func zeroValueMap() {
	var ages map[string]int
	fmt.Println("ages is zero value map?", ages == nil) // "true"
	fmt.Println("ages length is zero?", len(ages) == 0) // "true"

	// Most operations on maps are safe to perform on a nil map reference.
	// But storing to a nil map causes a panic.
	/*
		ages["carol"] = 21 // panic: assignment to entry in nil map
	*/

	// Accessing a map element by subscripting always yields a value.
	fmt.Println(ages["carol"])
}

func mapLookup() {
	// Know whether the element was really there or not.
	// For example, if the element type is numeric, you might have to
	// distinguish between a nonexistent element and an element that happens
	// to have the value zero, using a test like this.
	ages := map[string]int{
		"bob":   30,
		"alice": 28,
		"john":  32,
	}
	age, ok := ages["arya"]
	if !ok {
		fmt.Println("arya is not a key in this map; age =", age)
	}
	// You’ll often see these two statements combined, like this.
	if age, ok := ages["julia"]; !ok {
		fmt.Println("julia is not a key in this map; age =", age)
	}
}

func mapComparison() {
	// Maps cannot be compared to each other; the only legal comparison
	// is with nil.

	// To test whether two maps contain the same keys and the same associated
	// values, we must write a loop.
	equal := func(x, y map[string]int) bool {
		if len(x) != len(y) {
			return false
		}
		for k, xv := range x {
			if yv, ok := y[k]; !ok || yv != xv {
				return false
			}
		}
		return true
	}

	m1 := map[string]int{"A": 0}
	m2 := map[string]int{"B": 42}
	m3 := map[string]int{"A": 42}
	m4 := map[string]int{"A": 0, "B": 42}
	m5 := map[string]int{"A": 0, "B": 42}
	c := equal(m1, m2)
	fmt.Printf("m1 and m2 equal = %t\n", c)
	c = equal(m1, m3)
	fmt.Printf("m1 and m3 equal = %t\n", c)
	c = equal(m4, m3)
	fmt.Printf("m4 and m3 equal = %t\n", c)
	c = equal(m4, m5)
	fmt.Printf("m4 and m5 equal = %t\n", c)
}

func mapWhoseKeysAreSlices() {
	var m = make(map[string]int)

	// Helper function k that maps each key to a string, with the property
	// that k(x) == k(y) if and only if we consider x and y equivalent.
	k := func(list []string) string {
		return fmt.Sprintf("%q", list)
	}
	Add := func(list []string) { m[k(list)]++ }
	Count := func(list []string) int { return m[k(list)] }

	kFruits := []string{"fruit", "apple", "orange"}
	kNames := []string{"name", "alice", "john"}
	Add(kFruits)
	Add(kFruits)
	Add(kNames)
	fmt.Println(m)
	nFruits := Count(kFruits)
	fmt.Println("fruits count =", nFruits)
}

func structs() {
	// Declarations
	type Employee struct {
		ID        int
		Name      string
		Address   string
		DoB       time.Time
		Position  string
		Salary    int
		ManagerID int
	}

	var dilbert Employee
	fmt.Printf("%v\n", dilbert)
}

func emptyStruct() {
	// struct type with no fields is called the empty struct, written struct{}.
	// It has size zero.
	seen := make(map[string]struct{}) // set of strings
	// ...
	s := "aKey"
	if _, ok := seen[s]; !ok {
		seen[s] = struct{}{}
		// ...first time seeing s...
	}
}

func structLiterals() {
	type Point struct{ X, Y int }
	p := Point{1, 2}
	fmt.Println("Point struct p =", p)
}

func structValues() {
	type Point struct{ X, Y int }

	// Because structs are so commonly dealt with through pointers, it’s
	// possible to use this shorthand notation to create and initialize a struct
	// variable and obtain its address.
	pp := &Point{1, 2}
	fmt.Println(pp) // "&{1 2}"
	// It is exactly equivalent to
	p := new(Point)
	*p = Point{1, 2}
	fmt.Println(p) // "&{1 2}"
}

func comparingStructs() {
	// If all the fields of a struct are comparable, the struct itself is
	// comparable.

	// Comparable struct types, like other comparable types, may be used as the
	// key type of a map.
	type address struct {
		hostname string
		port     int
	}
	hits := make(map[address]int)
	hits[address{"golang.org", 443}]++
	fmt.Println("hits =", hits)
}

// Go's unusual struct embedding mechanism lets us use one named struct type as
// an _anonymous_ field of another struct type, providing a convenient syntactic
// shortcut so that a simple dot expression like `x.f` can stand for a chain of
// fields like `x.d.e.f`.
func structEmbeddingAnonFields() {
	// Consider a 2-D drawing program that provides a library of shapes, such as
	// rectangles, ellipses, stars, and wheels. Here are two of the types it
	// might define.

	// type Circle struct {
	// 	X, Y, Radius int
	// }

	// type Wheel struct {
	// 	X, Y, Radius, Spokes int
	// }
	// A Circle has fields for the X and Y coordinates of its center, and a
	// Radius. A Wheel has all the features of a Circle, plus Spokes, the number
	// of inscribed radial spokes.

	// var w Wheel
	// w.X = 8
	// w.Y = 8
	// w.Radius = 5
	// w.Spokes = 20

	// As the set of shapes grows, we’re bound to notice similarities and
	// repetition among them, so it may be convenient to factor out their common
	// parts.
	// type Point struct {
	// 	X, Y int
	// }

	// type Circle struct {
	// 	Center Point
	// 	Radius int
	// }

	// type Wheel struct {
	// 	Circle Circle
	// 	Spokes int
	// }

	// The application may be clearer for it, but this change makes accessing
	// the fields of a Wheel more verbose:
	// var w Wheel
	// w.Circle.Center.X = 8
	// w.Circle.Center.Y = 8
	// w.Circle.Radius = 5
	// w.Spokes = 20

	// Go lets us declare a field with a type but no name; such fields are
	// called _anonymous fields_. The type of the field must be a named type or
	// a pointer to a named type. Below, `Circle` and `Wheel` have one anonymous
	// field each. We say that a `Point` is embedded within `Circle`, and a
	// `Circle` is embedded within `Wheel`.
	type Point struct {
		X, Y int
	}

	type Circle struct {
		Point
		Radius int
	}

	type Wheel struct {
		Circle
		Spokes int
	}

	// Thanks to embedding, we can refer to the names at the leaves of the
	// implicit tree without giving the intervening names:
	var w Wheel
	w.X = 8      // equivalent to w.Circle.Point.X = 8
	w.Y = 8      // equivalent to w.Circle.Point.Y = 8
	w.Radius = 5 // equivalent to w.Circle.Radius = 5
	w.Spokes = 20

	// Unfortunately, there’s no corresponding shorthand for the struct literal
	// syntax, so neither of these will compile:
	/*
		w = Wheel{8, 8, 5, 20}                       // compile error: unknown fields
		w = Wheel{X: 8, Y: 8, Radius: 5, Spokes: 20} // compile error: unknown fields
	*/
}
