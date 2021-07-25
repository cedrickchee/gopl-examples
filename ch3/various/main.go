// Various code samples for chapter 3.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func main() {
	// =========================================================================
	integerOverflow()

	// =========================================================================
	bitwiseOperations()

	// =========================================================================
	convertValueFromOneTypeToAnother()

	// =========================================================================
	conversionFromIntToFloat()

	// =========================================================================
	printDecimalOctalHexNumbers()

	// =========================================================================
	printRuneLiterals()

	// =========================================================================
	floatingPointNumbers()

	// =========================================================================
	complexNumbers()
}

func integerOverflow() {
	var u uint8 = 255
	fmt.Println(u, u+1, u*u) // "255 0 1"

	var i int8 = 127
	fmt.Println(i, i+1, i*i) // "127 -128 1"
}

// How bitwise operations can be used to interpret a uint8 value as a
// compact and efficient set of 8 independent bits.
func bitwiseOperations() {
	var x uint8 = 1<<1 | 1<<5
	var y uint8 = 1<<1 | 1<<2

	// Print a number’s binary digits; 08 modifies %b (an adverb!)
	// to pad the result with zeros to exactly 8 digits.
	fmt.Printf("%08b\n", x)    // "00100010", the set {1, 5}
	fmt.Printf("%08b\n", y)    // "00000110", the set {1, 2}
	fmt.Printf("%08b\n", x&y)  // "00000010", the intersection {1}
	fmt.Printf("%08b\n", x|y)  // "00100110", the union {1, 2, 5}
	fmt.Printf("%08b\n", x^y)  // the symmetric difference {2, 5}
	fmt.Printf("%08b\n", x&^y) // &^ is a bit clear operator (AND NOT). output: "00100000", the difference {5}

	for i := uint8(0); i < 8; i++ {
		if x&(1<<i) != 0 { // membership test
			fmt.Println(i) // "1", "5"
		}
	}
	fmt.Printf("%08b\n", x<<1) // "01000100", the set {2, 6}
	fmt.Printf("%08b\n", x>>1) // "00010001", the set {0, 4}
}

func convertValueFromOneTypeToAnother() {
	var apples int32 = 1
	var oranges int16 = 2
	/*
		var kombucha int = apples + oranges // compile error
	*/
	fmt.Printf("apples T = %T, oranges T = %T\n", apples, oranges)

	// This type mismatch can be fixed in several ways, most directly by
	// converting everything to a common type.
	var kombucha int = int(apples) + int(oranges)
	fmt.Println("kombucha =", kombucha)
}

// A conversion that narrows a big integer into a smaller one, or a conversion
// from integer to floating-point or vice versa, may change the value or
// lose precision.
func conversionFromIntToFloat() {
	// Float to integer conversion discards any fractional part, truncating
	// toward zero.
	f := 3.141 // a float64
	i := int(f)
	fmt.Println(f, i) // "3.141 3"
	f = 1.99
	fmt.Println(int(f)) // "1"

	// You should avoid conversions in which the operand is out of range for
	// the target type, because the behavior depends on the implementation.
	x := 1e100  // a float64
	y := int(x) // result is implementation-dependent
	fmt.Println("y =", y)
}

// Printing numbers and control the radix and format with the %d, %o,
// and %x verbs.
func printDecimalOctalHexNumbers() {
	o := 0666
	fmt.Printf("%d %[1]o %#[1]o\n", o) // "438 666 0666"
	x := int64(0xdeadbeef)
	fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", x)
	// Output:
	// 3735928559 deadbeef 0xdeadbeef 0XDEADBEE

	// Note the use of two fmt tricks. Usually a Printf format string containing
	// multiple % verbs would require the same number of extra operands,
	// but the [1] "adverbs" after % tell Printf to use the first operand
	// over and over again.
}

// Rune literals are written as a character within single quotes.
func printRuneLiterals() {
	ascii := 'a'
	unicode := '⌘'
	newline := '\n'
	fmt.Printf("%d %[1]c %[1]q\n", ascii)   // "97 a 'a'"
	fmt.Printf("%d %[1]c %[1]q\n", unicode) // "8984 ⌘ '⌘'"
	fmt.Printf("%d %[1]q\n", newline)       // "10 '\n'"
}

func floatingPointNumbers() {
	// The smallest positive integer that cannot be exactly represented as a
	// float32 is not large.

	var f float32 = 16777216 // 1<<24
	fmt.Println(f == f+1)    // "true"!

	// Floating-point numbers can be written literally using decimals, like this
	const e = 2.71828 // (approximately)
	fmt.Println(e)

	// Very small or very large numbers are better written in scientific
	// notation, with the letter e or E preceding the decimal exponent.
	const Avogadro = 6.02214129e23
	const Planck = 6.62606957e-34
	fmt.Println("Avogadro =", Avogadro, "Planck =", Planck)

	// Printing floating-point values
	for x := 0; x < 8; x++ {
		fmt.Printf("x = %d e^x = %8.3f\n", x, math.Exp(float64(x)))
	}

	// Creating and detecting the special values defined by IEEE 754: the
	// positive and negative infinities, division by zero; and
	// NaN ("not a number").
	var z float64
	fmt.Println(z, -z, 1/z, -1/z, z/z) // "0 -0 +Inf -Inf NaN"
}

func complexNumbers() {
	var x complex128 = complex(1, 2) // 1+2i
	var y complex128 = complex(3, 4) // 3+4i
	fmt.Println(x * y)               // "(-5+10i)"
	fmt.Println(real(x * y))         // "-5"
	fmt.Println(imag(x * y))         // "10"

	fmt.Println(cmplx.Sqrt(-1)) // "(0+1i)"
}
