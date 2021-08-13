// Various code samples for chapter 12.
package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"

	"gopl.io/ch12/display"
	"gopl.io/ch7/eval"
)

func main() {
	// =========================================================================
	whyReflection()

	// =========================================================================
	reflectTypeAndValue()

	// =========================================================================
	recursiveValuePrinter()

	// =========================================================================
	displayExample()
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

func recursiveValuePrinter() {
	// Next we’ll take a look at how to improve the display of composite types.
	// Rather than try to copy `fmt.Sprint` exactly, we’ll build a debugging
	// utility function called `Display` that, given an arbitrarily complex
	// value `x`, prints the complete structure of that value, labeling each
	// element with the path by which it was found. Let’s start with an example.
	e, _ := eval.Parse("sqrt(A / pi)")
	display.Display("e", e)
	// In the call above, the argument to `Display` is a syntax tree from the
	// expression evaluator.
	// The output of `Display` is shown below:

	/*
		Display e (eval.call):
		e.fn = "sqrt"
		e.args[0].type = eval.binary
		e.args[0].value.op = 47
		e.args[0].value.x.type = eval.Var
		e.args[0].value.x.value = "A"
		e.args[0].value.y.type = eval.Var
		e.args[0].value.y.value = "pi"
	*/
}

func displayExample() {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	display.Display("strangelove", strangelove)
	// Output:
	// Display strangelove (main.Movie):
	// strangelove.Title = "Dr. Strangelove"
	// strangelove.Subtitle = "How I Learned to Stop Worrying and Love the Bomb"
	// strangelove.Year = 1964
	// strangelove.Color = false
	// strangelove.Actor["Pres. Merkin Muffley"] = "Peter Sellers"
	// strangelove.Actor["Gen. Buck Turgidson"] = "George C. Scott"
	// strangelove.Actor["Brig. Gen. Jack D. Ripper"] = "Sterling Hayden"
	// strangelove.Actor["Maj. T.J. \"King\" Kong"] = "Slim Pickens"
	// strangelove.Actor["Dr. Strangelove"] = "Peter Sellers"
	// strangelove.Actor["Grp. Capt. Lionel Mandrake"] = "Peter Sellers"
	// strangelove.Oscars[0] = "Best Actor (Nomin.)"
	// strangelove.Oscars[1] = "Best Adapted Screenplay (Nomin.)"
	// strangelove.Oscars[2] = "Best Director (Nomin.)"
	// strangelove.Oscars[3] = "Best Picture (Nomin.)"
	// strangelove.Sequel = nil

	// We can use `Display` to display the internals of library types, such as
	// `*os.File`:
	display.Display("os.Stderr", os.Stderr)
	// Output:
	// Display os.Stderr (*os.File):
	// (*(*os.Stderr).file).pfd.fdmu.state = 0
	// (*(*os.Stderr).file).pfd.fdmu.rsema = 0
	// (*(*os.Stderr).file).pfd.fdmu.wsema = 0
	// (*(*os.Stderr).file).pfd.Sysfd = 2
	// (*(*os.Stderr).file).pfd.pd.runtimeCtx = 0
	// (*(*os.Stderr).file).pfd.iovecs = nil
	// (*(*os.Stderr).file).pfd.csema = 0
	// (*(*os.Stderr).file).pfd.isBlocking = 1
	// (*(*os.Stderr).file).pfd.IsStream = true
	// (*(*os.Stderr).file).pfd.ZeroReadIsEOF = true
	// (*(*os.Stderr).file).pfd.isFile = true
	// (*(*os.Stderr).file).name = "/dev/stderr"
	// (*(*os.Stderr).file).dirinfo = nil
	// (*(*os.Stderr).file).nonblock = false
	// (*(*os.Stderr).file).stdoutOrErr = true
	// (*(*os.Stderr).file).appendMode = false

	// Notes
	//
	// Notice that even unexported fields are visible to reflection. Beware
	// that the particular output of this example may vary across platforms and
	// may change over time as libraries evolve. (Those fields are private for a
	// reason!) We can even apply `Display` to a `reflect.Value` and watch it
	// traverse the internal representation of the type descriptor for
	// `*os.File`.
	//
	// The output of the call `Display("rV", reflect.ValueOf(os.Stderr))` is
	// shown below, though of course your mileage may vary:
	display.Display("rV", reflect.ValueOf(os.Stderr))
	// Output:
	// Display rV (reflect.Value):
	// (*rV.typ).size = 8
	// (*rV.typ).ptrdata = 8
	// (*rV.typ).hash = 871609668
	// (*rV.typ).tflag = 9
	// (*rV.typ).align = 8
	// (*rV.typ).fieldAlign = 8
	// (*rV.typ).kind = 54
	// (*rV.typ).equal = func(unsafe.Pointer, unsafe.Pointer) bool 0x403800
	// (*(*rV.typ).gcdata) = 1
	// (*rV.typ).str = 6195
	// (*rV.typ).ptrToThis = 0
	// rV.ptr = unsafe.Pointer value
	// rV.flag = 22

	// Observe the difference between these two examples:
	var i interface{} = 3
	display.Display("i", i)
	// Output:
	// Display i (int):
	// i = 3
	display.Display("&i", &i)
	// Output:
	// Display &i (*interface {}):
	// (*&i).type = int
	// (*&i).value = 3

	// As currently implemented, `Display` will never terminate if it encounters
	// a cycle in the object graph, such as this linked list that eats its own
	// tail:
	// a struct that points to itself
	type Cycle struct {
		Value int
		Tail  *Cycle
	}
	var c Cycle
	c = Cycle{42, &c}
	// display.Display("c", c)
	// Display prints this ever-growing expansion:
	//
	// Display c (display.Cycle):
	// c.Value = 42
	// (*c.Tail).Value = 42
	// (*(*c.Tail).Tail).Value = 42
	// (*(*(*c.Tail).Tail).Tail).Value = 42
	// ...ad infinitum...
}
