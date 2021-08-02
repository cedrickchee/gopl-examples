// Various code samples for chapter 7.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"gopl.io/ch6/intset"
)

func main() {
	// =========================================================================
	interfacesAsContracts()

	// =========================================================================
	interfaceSatisfaction()

	// =========================================================================
	// Displaying the methods of a type

	// Example Print Duration
	Print(time.Hour)
	// Example Print Replacer
	Print(new(strings.Replacer))

	// =========================================================================
	interfaceAndConcreteTypes()

	// =========================================================================
	emptyInterfaceType()

	// =========================================================================
	assertValueOfTypeSatisfyInterface()
}

func interfacesAsContracts() {

}

func interfaceSatisfaction() {
	// A shorthand, Go programmers often say that a concrete type "is a"
	// particular interface type, meaning that it satisfies the interface. For
	// example, a *bytes.Buffer is an io.Writer; an *os.File is an
	// io.ReadWriter.

	// The assignability rule for interfaces is very simple: an expression may
	// be assigned to an interface only if its type satisfies the interface.

	var w io.Writer
	w = os.Stdout         // OK: *os.File has Write method
	w = new(bytes.Buffer) // OK: *bytes.Buffer has Write method
	// w = time.Second       // compile error: time.Duration lacks Write method

	var rwc io.ReadWriteCloser
	rwc = os.Stdout // OK: *os.File has Read, Write, Close methods
	// rwc = new(bytes.Buffer) // compile error: *bytes.Buffer lacks Close method

	// This rule applies even when the right-hand side is itself an interface:
	w = rwc // OK: io.ReadWriteCloser has Write method
	// rwc = w // compile error: io.Writer lacks Close method

	fmt.Printf("w type = %T\n", w)

	// The String method of the IntSet type requires a pointer receiver, so we
	// cannot call that method on a non-addressable IntSet value.
	/*
		var _ = intset.IntSet{}.String() // compile error: String requires *IntSet receiver
	*/

	// but we can call it on an IntSet variable:
	var s intset.IntSet // T is intset.IntSet in this example.
	// It is legal to call a *T method on an argument of type T so long as the
	// argument is a variable; the compiler implicitly takes its address.
	var _ = s.String() // OK: s is a variable and &s has a String method

	// However, since only *IntSet has a String method, only *IntSet satisfies
	// the fmt.Stringer interface:
	var _ fmt.Stringer = &s // OK
	// var _ fmt.Stringer = s // compile error: IntSet lacks String method
	var _ fmt.Stringer = &intset.IntSet{}   // OK
	var _ fmt.Stringer = new(intset.IntSet) // OK

	// Again, IntSet example copied from ch6
	var x intset.IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"
}

// Print prints the method set of the value x.
func Print(x interface{}) {
	// Print uses reflect.Type to print the type of an arbitrary value and enumerate
	// its methods.
	v := reflect.ValueOf(x)
	t := v.Type()
	fmt.Printf("type %s\n", t)

	for i := 0; i < v.NumMethod(); i++ {
		methType := v.Method(i).Type()
		fmt.Printf("func (%s) %s%s\n", t, t.Method(i).Name,
			strings.TrimPrefix(methType.String(), "func"))
	}
}

func interfaceAndConcreteTypes() {
	// Only the methods revealed by the interface type may be called, even if
	// the concrete type has others:

	os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
	// os.Stdout.Close()                // OK: *os.File has Close method
	// Note: if you close stdout, whatever you write to it won't show up in your
	// terminal/CLI from this point of time.

	var w io.Writer
	w = os.Stdout
	w.Write([]byte("hey\n")) // OK: io.Writer has Write method
	/*
		w.Close()              // compile error: io.Writer lacks Close method
	*/

	// An interface with more methods, such as `io.ReadWriter`, tells us more
	// about the values it contains, and places greater demands on the types
	// that implement it, than does an interface with fewer methods such as
	// `io.Reader`.
}

func emptyInterfaceType() {
	// Because the empty interface type places no demands on the types that
	// satisfy it, we can assign any value to the empty interface.
	var any interface{}
	fmt.Printf("any type = %T\n", any) // <nil>
	any = true
	fmt.Printf("any type = %T\n", any) // bool
	any = 12.34
	fmt.Printf("any type = %T\n", any) // float64
	any = map[string]int{"one": 1}
	fmt.Printf("any type = %T\n", any) // map[string]int
	any = new(bytes.Buffer)
	fmt.Printf("any type = %T\n", any) // *bytes.Buffer
}

func assertValueOfTypeSatisfyInterface() {
	// The declaration below asserts at compile time that a value of type
	// *bytes.Buffer satisfies io.Writer:

	// *bytes.Buffer must satisfy io.Writer
	var w io.Writer = new(bytes.Buffer) // allocates memory and the value returned is a pointer to a newly allocated zero value of that type.
	fmt.Printf("w type = %T\n", w)      // w type = *bytes.Buffer (a pointer to a struct)

	// We needn’t allocate a new variable since any value of type `*bytes.Buffer` will do, even `nil`,
	// which we write as `(*bytes.Buffer)(nil)` using an explicit conversion. And since we never
	// intend to refer to `w`, we can replace it with the blank identifier. Together, these changes give us
	// this more frugal variant:
	var _ io.Writer = (*bytes.Buffer)(nil)
	// fmt.Printf("v type = %T, value = %[1]v\n", v) // v type = *bytes.Buffer, value = <nil>
}
