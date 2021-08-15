// Various code samples for chapter 12.
package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	// =========================================================================
	unsafeSizeOfAlignOfOffsetof()

	// =========================================================================
	unsafePointer()

	// =========================================================================
	deepEquivalence()
}

func unsafeSizeOfAlignOfOffsetof() {
	// The `unsafe.Sizeof` function reports the size in bytes of the
	// representation of its operand, which may be an expression of any type;
	// the expression is not evaluated. A call to `Sizeof` is a constant
	// expression of type `uintptr`, so the result may be used as the dimension
	// of an array type, or to compute other constants.
	fmt.Println(unsafe.Sizeof(float64(0))) // "8" Equivalent to 64/8=8

	// Figure 13.1 shows a struct variable `x` and its memory layout on typical
	// 32- and 64-bit Go implementations. The gray regions are holes.
	var x struct {
		a bool
		b int16
		c []int
	}
	// The table below shows the results of applying the three `unsafe`
	// functions to `x` itself and to each of its three fields:

	fmt.Println("Sizeof(x) =", unsafe.Sizeof(x))
	fmt.Println("Sizeof(x.a) =", unsafe.Sizeof(x.a))
	fmt.Println("Sizeof(x.b) =", unsafe.Sizeof(x.b))
	fmt.Println("Sizeof(x.c) =", unsafe.Sizeof(x.c))
	fmt.Println("Alignof(x) =", unsafe.Alignof(x))
	fmt.Println("Alignof(x.a) =", unsafe.Alignof(x.a))
	fmt.Println("Alignof(x.b) =", unsafe.Alignof(x.b))
	fmt.Println("Alignof(x.c) =", unsafe.Alignof(x.c))
	fmt.Println("Offsetof(x.a) =", unsafe.Offsetof(x.a))
	fmt.Println("Offsetof(x.b) =", unsafe.Offsetof(x.b))
	fmt.Println("Offsetof(x.c) =", unsafe.Offsetof(x.c))
	// Prints

	// Typical 32-bit platform:
	// Sizeof(x)   = 16		Alignof(x)   = 4
	// Sizeof(x.a) = 1		Alignof(x.a) = 1	Offsetof(x.a) = 0
	// Sizeof(x.b) = 2		Alignof(x.b) = 2	Offsetof(x.b) = 2
	// Sizeof(x.c) = 12		Alignof(x.c) = 4	Offsetof(x.c) = 4

	// Typical 64-bit platform:
	// Sizeof(x)   = 32		Alignof(x)   = 8
	// Sizeof(x.a) = 1		Alignof(x.a) = 1	Offsetof(x.a) = 0
	// Sizeof(x.b) = 2		Alignof(x.b) = 2	Offsetof(x.b) = 2
	// Sizeof(x.c) = 24		Alignof(x.c) = 8	Offsetof(x.c) = 8

	// Despite their names, these functions are not in fact `unsafe`, and they
	// may be helpful for understanding the layout of raw memory in a program
	// when optimizing for space.
}

func Float64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }
func unsafePointer() {
	// The `unsafe.Pointer` type is a special kind of pointer that can hold the
	// address of any variable. Like ordinary pointers, `unsafe.Pointer`s are
	// comparable and may be compared with `nil`, which is the zero value of the
	// type.
	//
	// An ordinary `*T` pointer may be converted to an `unsafe.Pointer`, and an
	// `unsafe.Pointer` may be converted back to an ordinary pointer, not
	// necessarily of the same type `*T`. By converting a `*float64` pointer to
	// a `*uint64`, for instance, we can inspect the bit pattern of a
	// floating-point variable:
	fmt.Printf("%#016x\n", Float64bits(1.0)) // "0x3ff0000000000000"

	// Demonstrates basic use of unsafe.Pointer.
	var x struct {
		a bool
		b int16
		c []int
	}

	{
		// equivalent to pb := &x.b
		pb := (*int16)(unsafe.Pointer(
			uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
		*pb = 42
		fmt.Println(x.b) // "42"
	}
	{
		// Do not be tempted to introduce temporary variables of type `uintptr` to
		// break the lines. This code is incorrect:
		tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
		pb := (*int16)(unsafe.Pointer(tmp)) // linter warning: possible misuse of unsafe.Pointer
		*pb = 9
		fmt.Println(x.b) // "9"

		// Notes
		//
		// The reason is **very subtle**. **Some garbage collectors move
		// variables around in memory to reduce fragmentation or bookkeeping**.
		// Garbage collectors of this kind are known as _moving GCs_. When a
		// variable is moved, all pointers that hold the address of the old
		// location must be updated to point to the new one. From the
		// perspective of the garbage collector, an `unsafe.Pointer` is a
		// pointer and thus its value must change as the variable moves, but a
		// `uintptr` is just a number so its value must not change. **The
		// incorrect code above hides a pointer from the garbage collector in
		// the non-pointer variable `tmp`**. By the time the second statement
		// executes, the variable `x` could have moved and the number in `tmp`
		// would no longer be the address `&x.b`. The third statement clobbers
		// an arbitrary memory location with the value 42.

	}
}

func testSplit() {
	got := strings.Split("a:b:c", ":")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		fmt.Println("not equal")
	}
}
func deepEquivalence() {
	// The `DeepEqual` function from the `reflect` package reports whether two
	// values are "deeply" equal.
	testSplit()

	// Although `DeepEqual` is convenient, its distinctions can seem arbitrary.
	// For example, it doesn’t consider a nil map equal to a non-nil empty map,
	// nor a nil slice equal to a non-nil empty one:

	// A non-nil empty slice and a nil slice (for example, []byte{} and
	// []byte(nil)) are not deeply equal.
	var a, b []string = nil, []string{}
	fmt.Println(reflect.DeepEqual(a, b)) // "false"

	// Map values are deeply equal when all of the following are true: they are
	// both nil or both non-nil, they have the same length, and either they are
	// the same map object or their corresponding keys (matched using Go
	// equality) map to deeply equal values.
	var c, d map[string]int = nil, make(map[string]int)
	fmt.Println(reflect.DeepEqual(c, d)) // "false"
}
