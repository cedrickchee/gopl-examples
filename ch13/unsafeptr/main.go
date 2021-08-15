// Package unsafeptr demonstrates basic use of unsafe.Pointer.
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var x struct {
		a bool
		b int16
		c []int
	}

	// equivalent to pb := &x.b
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42

	fmt.Println(x.b) // "42"

	// Let's break it down:
	fmt.Println(&x)                          // "&{false 0 []}"
	fmt.Println(unsafe.Pointer(&x))          // "0xc0000c0000"
	fmt.Println(uintptr(unsafe.Pointer(&x))) // "824634507264"
	fmt.Println(unsafe.Offsetof(x.b))        // "2"
	// equivalent to unsafe.Pointer(&x.b). Prints: "0xc0000c0002"
	fmt.Println(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b))) // "0xc0000c0002"
}

// Notes
//
// In general, `unsafe.Pointer` conversions let us write arbitrary values to
// memory and thus subvert the type system.
//
// An `unsafe.Pointer` may also be converted to a `uintptr` that holds the
// pointer’s numeric value, letting us perform arithmetic on addresses.
// Recall that a `uintptr` is an unsigned integer wide enough to represent an
// address.
//
// Many `unsafe.Pointer` values are thus intermediaries for converting ordinary
// pointers to raw numeric addresses and back again.
// The example above takes the address of variable `x`, adds the offset of its
// `b` field, converts the resulting address to `*int16`, and through that
// pointer updates `x.b`.
//
// Although the syntax is cumbersome--perhaps no bad thing since these features
// should be used sparingly.
