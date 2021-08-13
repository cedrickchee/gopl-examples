// Various code samples for chapter 12.
package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

func main() {
	// =========================================================================
	whyReflection()

	// =========================================================================
	reflectTypeAndValue()
}

func Sprint(x interface{}) string {
	type stringer interface {
		String() string
	}
	switch x := x.(type) {
	case stringer:
		return x.String()
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	// ...similar cases for int16, uint32, and so on...
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		// array, chan, func, map, pointer, slice, struct
		return "???"
	}
}
func whyReflection() {
	// Sometimes we need to write a function capable of dealing uniformly with
	// values of types that don’t satisfy a common interface, don’t have a known
	// representation, or don’t exist at the time we design the function--or
	// even all three.
	//
	// A familiar example is the formatting logic within `fmt.Fprintf`, which
	// can usefully print an arbitrary value of any type, even a user-defined
	// one. Let’s try to implement a function like it using what we know
	// already.
	fmt.Println(`Sprint("hello") =`, Sprint("hello"))
	fmt.Println("Sprint(123) =", Sprint(123))
	fmt.Println("Sprint(false) =", Sprint(false))
}

func reflectTypeAndValue() {
	// The `reflect.TypeOf` function accepts any `interface{}` and returns its
	// dynamic type as a `reflect.Type`
	t := reflect.TypeOf(3)  // a reflect.Type
	fmt.Println(t.String()) // "int"
	fmt.Println(t)          // "int"
	// Recall that an assignment from a concrete value to an interface type
	// performs an implicit interface conversion, which creates an interface
	// value consisting of two components: its dynamic type and its dynamic
	// value.

	// Because `reflect.TypeOf` returns an interface value’s dynamic type, it
	// always returns a concrete type. So, for example, the code below prints
	// "`*os.File`", not "`io.Writer`". Later, we will see that `reflect.Type`
	// is capable of representing interface types too.
	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // "*os.File"

	// Notice that `reflect.Type` satisfies `fmt.Stringer`. Because printing the
	// dynamic type of an interface value is useful for debugging and logging,
	// `fmt.Printf` provides a shorthand, `%T`, that uses `reflect.TypeOf`
	// internally:
	fmt.Printf("%T\n", 3) // "int"

	// A `reflect.Value` can hold a value of any type. The `reflect.ValueOf`
	// function accepts any `interface{}` and returns a `reflect.Value`
	// containing the interface’s dynamic value. As with `reflect.TypeOf`, the
	// results of `reflect.ValueOf` are always concrete, but a `reflect.Value`
	// can hold interface values too.
	v := reflect.ValueOf(3) // a reflect.Value
	fmt.Println(v)          // "3"
	fmt.Printf("%v\n", v)   // "3"
	fmt.Println(v.String()) // NOTE: "<int Value>"
	// Like `reflect.Type`, `reflect.Value` also satisfies `fmt.Stringer`, but
	// unless the `Value` holds a string, the result of the `String` method
	// reveals only the type. Instead, use the `fmt` package’s `%v` verb, which
	// treats `reflect.Values` specially.

	// Calling the `Type` method on a `Value` returns its type as a
	// `reflect.Type`:
	{
		t := v.Type()           // a reflect.Type
		fmt.Println(t.String()) // "int"
	}

	// The inverse operation to `reflect.ValueOf` is the
	// `reflect.Value.Interface` method. It returns an `interface{}` holding the
	// same concrete value as the `reflect.Value`:
	{
		v := reflect.ValueOf(3) // a reflect.Value
		x := v.Interface()      // an interface{}
		i := x.(int)            // an int
		fmt.Printf("%d\n", i)   // "3"
	}
}
