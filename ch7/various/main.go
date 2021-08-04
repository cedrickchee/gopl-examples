// Various code samples for chapter 7.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"syscall"
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

	// =========================================================================
	interfaceValues()

	// =========================================================================
	compareInterfaceValues()

	// =========================================================================
	caveatInterfaceContainingNilPointer()

	// =========================================================================
	sortingWithSortInterface()

	// =========================================================================
	testingWhetherASequenceIsSorted()

	// =========================================================================
	errorInterface()
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

func interfaceValues() {
	var w io.Writer
	fmt.Printf("interface value w is nil = %t\n", w == nil)
	// w.Write([]byte("hello")) // panic: runtime error: invalid memory address or nil pointer dereference

	var x io.Writer
	// The statement assigns a value of type `*os.File` to `x`:
	// This assignment involves an implicit conversion from a concrete type to
	// an interface type.
	//
	// The interface value’s dynamic type is set to the type descriptor for the
	// pointer type `*os.File`, and its dynamic value holds a copy of
	// `os.Stdout`, which is a pointer to the `os.File` variable representing
	// the standard output of the process.
	x = os.Stdout
	// Calling the `Write` method on an interface value containing an `*os.File`
	// pointer causes the `(*os.File).Write` method to be called.
	x.Write([]byte("hello")) // "hello"
	// The effect is as if we had made this call directly:
	os.Stdout.Write([]byte("hello")) // "hello"

	var z io.Writer
	z = new(bytes.Buffer)
	// The dynamic type is now `*bytes.Buffer` and the dynamic value is a
	// pointer to the newly allocated buffer.
	z.Write([]byte("hello")) // writes "hello" to the bytes.Buffer
	// This time, the type descriptor is `*bytes.Buffer`, so the
	// `(*bytes.Buffer).Write` method is called, with the address of the buffer
	// as the value of the receiver parameter. The call appends "hello" to the
	// buffer.

	fmt.Printf("\ninterface value z is not nil = %t\n", z != nil)
	// Assigns nil to the interface value.
	z = nil
	// This resets both its components (dynamic type and dynamic value) to
	// `nil`, restoring `w` to the same state as when it was declared.
	fmt.Printf("interface value z is nil = %t\n", z == nil)

	// An interface value can hold arbitrarily large dynamic values.
	// Create an interface value from `time.Time` type.
	var q interface{} = time.Now() // an interface value holding a time.Time struct (a large type)
	fmt.Printf("interface value q is not nil = %t\n", q != nil)
}

func compareInterfaceValues() {
	// Interface values may be compared using `==` and `!=`. Two interface
	// values are equal if both are `nil`, or if their dynamic types are
	// identical and their dynamic values are equal.

	// However, if two interface values are compared and have the same dynamic
	// type, but that type is not comparable (a slice, for instance), then the
	// comparison fails with a panic:
	var x interface{} = []int{1, 2, 3}
	fmt.Printf("interface value x is not nil = %t\n", x != nil)
	/*
		fmt.Println(x == x) // panic: runtime error: comparing uncomparable type []int
	*/

	// When handling errors, or during debugging, it is often helpful to report
	// the dynamic type of an interface value.
	var w io.Writer
	fmt.Printf("%T\n", w) // "<nil>"

	w = os.Stdout
	fmt.Printf("%T\n", w) // "*os.File"

	w = new(bytes.Buffer)
	fmt.Printf("%T\n", w) // "*bytes.Buffer"
}

// With `debug` set to `true`, the `collectOutput` function collects the output
// of the function `f` in a `bytes.Buffer`.
const debug = true

func collectOutput() {
	var buf *bytes.Buffer
	if debug {
		buf = new(bytes.Buffer) //  enable collection of output
	}
	// fmt.Println("buf value =", buf) // "<nil>" if debug set to false, empty of debug set to true.
	f(buf) // NOTE: subtly incorrect! Solution: change the type of buf to io.Writer (var buf io.Writer)
	if debug {
		// ...use buf...
	}
}

// If out is non-nil, output will be written to it.
func f(out io.Writer) {
	// ...do something...
	if out != nil {
		out.Write([]byte("done!\n"))
	}
	// We might expect that changing debug to false would disable the collection
	// of the output, but in fact it causes the program to panic during the
	// out.Write call. -- panic: invalid memory address or nil pointer
	// dereference
}
func caveatInterfaceContainingNilPointer() {
	// An interface containing a nil pointer is non-nil.
	//
	// A nil interface value, which contains no value at all, is not the same as
	// an interface value containing a pointer that happens to be nil. This
	// subtle distinction creates a trap into which every Go programmer has
	// stumbled.
	collectOutput()
}

// An in-place sort algorithm needs three things—the length of the sequence, a
// means of comparing two elements, and a way to swap two elements—so they are
// the three methods of `sort.Interface`.
//
// To sort any sequence, we need to define a type that implements these three
// methods, then apply `sort.Sort` to an instance of that type. As perhaps the
// simplest example, consider sorting a slice of strings. The new type
// `StringSlice` and its `Len`, `Less`, and `Swap` methods are shown below.
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sorting with sort.Interface
func sortingWithSortInterface() {
	// Now we can sort a slice of strings, names, by converting the slice to a
	// `StringSlice` like this:
	names := []string{"John", "Bob", "Alice", "Tracy"}
	// The conversion yields a slice value with the same length, capacity, and
	// underlying array as `names` but with a type that has the three methods
	// required for sorting.
	sort.Sort(StringSlice(names))
	fmt.Printf("%#v\n", names) // "[]string{"Alice", "Bob", "John", "Tracy"}"

	// Sorting a slice of strings is so common that the `sort` package provides
	// the `StringSlice` type, as well as a function called `Strings` so that
	// the call above can be simplified to `sort.Strings(names)`.
	s := []string{"Go", "Bravo", "Gopher", "Alpha", "Grin", "Delta"}
	sort.Strings(s)
	fmt.Println(s)
}

func testingWhetherASequenceIsSorted() {
	values := []int{3, 1, 4, 1}
	fmt.Println(sort.IntsAreSorted(values)) // "false"
	sort.Ints(values)
	fmt.Println(values)                     // "[1 1 3 4]"
	fmt.Println(sort.IntsAreSorted(values)) // "true"
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	fmt.Println(values)                     // "[4 3 1 1]"
	fmt.Println(sort.IntsAreSorted(values)) // "false"
}

func errorInterface() {
	// Every call to `New` allocates a distinct `error` instance that is equal
	// to no other. We would not want a distinguished error such as `io.EOF` to
	// compare equal to one that merely happened to have the same message.
	fmt.Println("error compare =", errors.New("EOF") == errors.New("EOF")) // "false"

	// Although `*errorString` may be the simplest type of `error`, it is far
	// from the only one. For example, the `syscall` package provides Go’s
	// low-level system call API. On many platforms, it defines a numeric type
	// `Errno` that satisfies `error`, and on Unix platforms, `Errno`’s `Error`
	// method does a lookup in a table of strings.
	//
	// The following statement creates an interface value holding the `Errno`
	// value 2, signifying the POSIX `ENOENT` condition:
	var err error = syscall.Errno(2)
	fmt.Println(err.Error()) // "no such file or directory"
	fmt.Println(err)         // "no such file or directory"
}
