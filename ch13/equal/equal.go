// Package equal provides a deep equivalence relation for arbitrary values.
package equal

import (
	"reflect"
	"unsafe"
)

func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() {
		return false
	}

	// cycle check
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true // identical references
		}
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}

	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()

	case reflect.String:
		return x.String() == y.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return x.Int() == y.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return x.Uint() == y.Uint()

	case reflect.Float32, reflect.Float64:
		return x.Float() == y.Float()

	case reflect.Complex64, reflect.Complex128:
		return x.Complex() == y.Complex()

	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()

	case reflect.Ptr, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)

	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true

	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if !equal(x.Field(i), y.Field(i), seen) {
				return false
			}
		}
		return true

	case reflect.Map:
		if x.Len() != y.Len() {
			return false
		}
		for _, k := range x.MapKeys() {
			if !equal(x.MapIndex(k), y.MapIndex(k), seen) {
				return false
			}
		}
		return true
	}
	panic("unreachable")
}

// Equal reports whether x and y are deeply equal.
//
// Map keys are always compared with ==, not deeply.
// (This matters for keys containing pointers or interfaces.)
func Equal(x, y interface{}) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

// Notes
//
// The function `Equal` compares arbitrary values. Like `DeepEqual`, it compares
// slices and maps based on their elements, but unlike `DeepEqual`, it considers
// a nil slice (or map) equal to a non-nil empty one.
// The basic recursion over the arguments can be done with reflection, using a
// similar approach to the `Display` program we saw.
// As usual, we define an unexported function, `equal`, for the recursion. Don’t
// worry about the `seen` parameter just yet. For each pair of values `x` and
// `y` to be compared, equal checks that both (or neither) are valid and checks
// that they have the same type. The result of the function is defined as a set
// of switch cases that compare two values of the same type.

// Regarding the cycle check
//
// To ensure that the algorithm terminates even for cyclic data structures, it
// must record which pairs of variables it has already compared and avoid
// comparing them a second time. `Equal` allocates a set of `comparison`
// structs, each holding the address of two variables (represented as
// `unsafe.Pointer` values) and the type of the comparison. We need to record
// the type in addition to the addresses because different variables can have
// the same address. For example, if `x` and `y` are both arrays, `x` and `x[0]`
// have the same address, as do `y` and `y[0]` , and it is important to
// distinguish whether we have compared `x` and `y` or `x[0]` and `y[0]`.
//
// Once `equal` has established that its arguments have the same type, and
// before it executes the switch, it checks whether it is comparing two
// variables it has already seen and, if so, terminates the recursion.
