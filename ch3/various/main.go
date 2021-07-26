// Various code samples for chapter 3.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"time"
	"unicode/utf8"
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

	// =========================================================================
	booleans()

	// =========================================================================
	strings()

	// =========================================================================
	hp := HasPrefix("hello, world", "abc")
	fmt.Println("HasPrefix =", hp)
	hp = HasPrefix("世界", "世")
	fmt.Println("HasPrefix =", hp)

	// =========================================================================
	hs := HasSuffix("hello, world", "world")
	fmt.Println("HasSuffix =", hs)
	hs = HasSuffix("世界", "界")
	fmt.Println("HasSuffix =", hs)

	// =========================================================================
	ct := Contains("hello, world", "ello")
	fmt.Println("Contains =", ct)
	ct = Contains("hello, world", "abc")
	fmt.Println("Contains =", ct)

	// =========================================================================
	unicode()

	// =========================================================================
	// External examples: https://www.practical-go-lessons.com/chap-7-hexadecimal-octal-ascii-utf8-unicode-runes#runes
	runes()

	// =========================================================================
	runeConversion()

	// =========================================================================
	convertStringToByteSlice()

	// =========================================================================
	convertIntToString()

	// =========================================================================
	constants()
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

func booleans() {
	// There is no implicit conversion from a boolean value to a numeric value
	// like 0 or 1, or vice versa.
	i := 0
	b := true
	if b {
		i = 1
	}
	fmt.Println("boolean i:", i)

	// A conversion function if this operation were needed often.
	btoi := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}
	fmt.Println("btoi(true):", btoi(true))
}

func strings() {
	// The built-in **`len` function returns the number of bytes (not runes)
	// in a string**, and the index operation `s[i]` retrieves the _i_-th
	// byte of string `s`, where `0 ≤ i < len(s)`.
	s := "hello, world"
	fmt.Println(len(s))     // "12"
	fmt.Println(s[0], s[7]) // "104 119" ('h' and 'w')

	/*
		c := s[len(s)] // panic: index out of range
	*/

	// The substring operation `s[i:j]` yields a new string consisting of the
	// bytes of the original string starting at index `i` and continuing up to,
	// but not including, the byte at index `j`.
	// The result contains `j-i` bytes.
	fmt.Println(s[0:5]) //  "hello"
	fmt.Println(s[:5])  // "hello"
	fmt.Println(s[7:])  // "world"
	fmt.Println(s[:])   // "hello, world"

	// The + operator makes a new string by concatenating two strings:
	fmt.Println("goodbye" + s[5:]) // "goodbye, world"

	// String literals
	fmt.Println("Hello, 世界")

	// The string "hello, world" and two substrings.
	// The string and two of its substrings sharing the same underlying
	// byte array.
	str := "hello, world"
	hello := str[:5]
	world := str[7:]
	str += ". How are you?"
	fmt.Println(hello, world)
	fmt.Println(str)

	// Raw string literal
	const GoUsage = `Go is a tool for managing Go source code.
	Usage:
		go command [arguments]`
	fmt.Println(GoUsage)

	// Unicode and UTF-8
	//
	// Unicode escapes in Go string literals allow us to specify them by their
	// numeric code point value.
	// There are two forms, `\uhhhh` for a 16-bit value and `\Uhhhhhhhh` for
	// a 32-bit value, where each h is a hexadecimal digit; the need for the
	// 32-bit form arises very infrequently. Each denotes the UTF-8 encoding
	// of the specified code point.
	// Thus, for example, the following string literals all represent the
	// same six-byte string.
	// The three escape sequences provide alternative notations for the
	// first string, but the values they denote are identical.
	fmt.Println("世界")
	fmt.Println("\xe4\xb8\x96\xe7\x95\x8c") // UTF-8 bytestring
	fmt.Println("\u4e16\u754c")             // Unicode string (16-bit)
	fmt.Println("\U00004e16\U0000754c")     // Unicode string (32-bit)
	// Unicode escapes may also be used in rune literals.
	// These three literals are equivalent:
	fmt.Println('世' == '\u4e16')
	fmt.Println('世' == '\U00004e16')
	fmt.Println('\U00004e16' == '\u4e16')

	// A rune whose value is less than 256 may be written with a single
	// hexadecimal escape, but for higher values, a \u or \U escape
	// must be used.
	fmt.Println('\x41') // "65" for 'A'

	// Consequently, '\xe4\xb8\x96' is not a legal rune literal, even though
	// those three bytes are a valid UTF-8 encoding of a single code point.
	// fmt.Println('\xe4\xb8\x96') // compiler error "illegal rune literal"
}

// Thanks to the nice properties of UTF-8, many string operations don’t
// require decoding. We can test whether one string contains another as
// a prefix.
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// Test whether one string contains another as a substring.
func Contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if HasPrefix(s[i:], substr) {
			return true
		}
	}
	return false
}

// UTF-8 is a variable-length encoding of Unicode code points as bytes.
// Dealing with the individual Unicode characters.
func unicode() {
	// The string contains 13 bytes, but interpreted as UTF-8, it encodes
	// only nine code points or runes.
	s := "Hello, 世界"
	fmt.Println("bytes =", len(s))                    // "13"
	fmt.Println("runes =", utf8.RuneCountInString(s)) // "9"

	// To process those characters, we need a UTF-8 decoder.
	fmt.Printf("byte index\trune\tvalue(decimal)\tsize\n")
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%d\t\t%c\t%d\t\t%d\n", i, r, r, size)
		i += size
	}

	// A range loop decodes a UTF-8-encoded string.
	fmt.Printf("\nbyte index\trune\tvalue(decimal)\n")
	for i, r := range s {
		fmt.Printf("%d\t\t%c\t%d\n", i, r, r)
	}

	// We could use a simple range loop to count the number of runes in a string.
	n := 0
	for range s {
		n++
	}
	fmt.Println("runes count =", n)

	// rune literals
	// print the hexadecimal representation of a number
	// '世' equivalent to '\u4e16' (a 16-bit value)
	// 'A' is mapped to the code point (as hex) "0041" (U+0041) (65 is decimal for 0x41)
	// 'A' equivalent to '\u0041' (a 16-bit value)
	// Unicode code point search, e.g.,
	// https://unicode.scarfboy.com/?s=0041
	// https://unicode.scarfboy.com/?s=U%2b4E16
	fmt.Printf("%q: hex = %x\n", '世', 19990)
}

// Behind the scene, a string is a collection of bytes.
// We can iterate over the bytes of a string with a for loop.
func runes() {
	s := "我爱 Golang" // "I love Golang"
	for _, r := range s {
		// r is of type rune
		fmt.Printf("Unicode code point: %U - character '%c' - binary %b - hex %X - decimal %d\n", r, r, r, r, r)
	}
}

// A `[]rune` conversion applied to a UTF-8-encoded string returns the sequence
// of Unicode code points that the string encodes.
func runeConversion() {
	// "program" in Japanese katakana
	s := "プログラム"
	fmt.Println(len(s)) // "15"
	// (The verb % x inserts a space between each pair of hex digits.)
	fmt.Printf("% x\n", s) // "e3 83 97 e3 83 ad e3 82 b0 e3 83 a9 e3 83 a0"
	// Returns the sequence of Unicode code points that the string encodes.
	r := []rune(s)
	fmt.Printf("%x\n", r) // "[30d7 30ed 30b0 30e9 30e0]"
	// If a slice of runes is converted to a string, it produces the
	// concatenation of the UTF-8 encodings of each rune.
	fmt.Println(string(r)) // "プログラム"
	// Converting an integer value to a string interprets the integer as a
	// rune value, and yields the UTF-8 representation of that rune.
	fmt.Println(string(65))     // "A", not "65"
	fmt.Println(rune('A'))      // "65", not "A"
	fmt.Println(string(0x4eac)) // "京"
	// If the rune is invalid, the replacement character is substituted.
	fmt.Println(string(1234567)) // "�"
}

// Strings can be converted to byte slices and back again.
func convertStringToByteSlice() {
	s := "abc"
	b := []byte(s)
	fmt.Println("byte slice =", b)
	s2 := string(b)
	fmt.Println("string =", s2)
}

// Convert between numeric values and their string representations.
func convertIntToString() {
	// Convert an integer to a string
	x := 123
	y := fmt.Sprintf("%d", x)
	// integer to ASCII
	fmt.Println(y, strconv.Itoa(x)) // "123 123"
	// FormatInt and FormatUint can be used to format numbers in a
	// different base.
	fmt.Println(strconv.FormatInt(int64(x), 2)) // "1111011"

	// To parse a string representing an integer, use the strconv functions
	// Atoi or ParseInt, or ParseUint for unsigned integers.
	m, _ := strconv.Atoi("123") // m is an int
	fmt.Printf("m is type of %T\n", m)
	n, _ := strconv.ParseInt("123", 10, 64) // base 10, up to 64 bits
	fmt.Printf("n is type of %T\n", n)
}

func constants() {
	// Constant declaration.
	const noDelay time.Duration = 0
	const timeout = 5 * time.Minute
	fmt.Printf("%T %[1]v\n", noDelay)     // "time.Duration 0s"
	fmt.Printf("%T %[1]v\n", timeout)     // "time.Duration 5m0s"
	fmt.Printf("%T %[1]v\n", time.Minute) // "time.Duration 1m0s"

	// When a sequence of constants is declared as a group, the right-hand side
	// expression may be omitted for all but the first of the group, implying
	// that the previous expression and its type should be used again.
	const (
		a = 1
		b
		c = 2
		d
	)
	fmt.Println(a, b, c, d) // "1 1 2 2"

	// The constant generator iota
	// Types of this kind are often called enumerations, or enums for short.
	type Weekday int

	const (
		Sunday Weekday = iota
		Monday
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
	)
	fmt.Println("Tuesday =", Tuesday) // "2"

	// We can use iota in more complex expressions too, as in this example from
	// the net package where each of the lowest 5 bits of an unsigned integer
	// is given a distinct name and boolean interpretation.
	type Flags uint

	const (
		FlagUp           Flags = 1 << iota // is up
		FlagBroadcast                      // supports broadcast access capability
		FlagLoopback                       // is a loopback interface
		FlagPointToPoint                   // belongs to a point-to-point link
		FlagMulticast                      // supports multicast access capability
	)
	// 1 << iota, which evaluates to successive powers of two,
	// each corresponding to a single bit.
	fmt.Println("FlagBroadcast =", FlagBroadcast) // "2"
	fmt.Println("FlagLoopback =", FlagLoopback)   // "4"

	// As a more complex example of iota,
	// this declaration names the powers of 1024.
	const (
		_   = 1 << (10 * iota)
		KiB // 1024
		MiB // 1048576
		GiB // 1073741824
		TiB // 1099511627776 (exceeds 1 << 32)
		PiB // 1125899906842624
		EiB // 1152921504606846976
		ZiB // 1180591620717411303424 (exceeds 1 << 64)
		YiB // 1208925819614629174706176

	)
	fmt.Printf("1 MiB = %d Byte\n", MiB)
	fmt.Printf("1 TiB = %d Byte\n", TiB)

	// Untyped constants
	// The values ZiB and YiB in the example above are too big to store in any
	// integer variable, but they are legitimate constants that may be used in
	// expressions like this one.
	fmt.Println(YiB / ZiB) // "1024"

	// As another example, the floating-point constant math.Pi may be used
	// wherever any floating-point or complex value is needed.
	var p float32 = math.Pi
	var q float64 = math.Pi
	var r complex128 = math.Pi
	fmt.Println("p, q, r = ", p, q, r)
	const Pi64 float64 = math.Pi
	var s float32 = float32(Pi64)
	var t float64 = Pi64
	var u complex128 = complex128(Pi64)
	fmt.Println("s, t, u = ", s, t, u)

	// Recall that `/` may represent integer or floating-point division
	// depending on its operands. Consequently, the choice of literal may
	// affect the result of a constant division expression.
	var f float64 = 212
	div1 := (f - 32) * 5 / 9
	fmt.Printf("result1 = %v %[1]T\n", div1) // "100"; (f - 32) * 5 is a float64
	div2 := 5 / 9 * (f - 32)
	fmt.Printf("result2 = %v %[1]T\n", div2) // "0";  5/9 is an untyped integer, 0
	div3 := 5.0 / 9.0 * (f - 32)
	fmt.Printf("result3 = %v %[1]T\n", div3) // "100"; 5.0/9.0 is an untyped float

	// Whether implicit or explicit, converting a constant from one type to
	// another requires that the target type can represent the original value.
	// Rounding is allowed for real and complex floating-point numbers.
	const (
		deadbeef = 0xdeadbeef        // untyped int with value 3735928559
		i        = uint32(deadbeef)  // uint32 with value 3735928559
		j        = float32(deadbeef) // float32 with value 3735928576 (rounded up)
		k        = float64(deadbeef) // float64 with value 3735928559 (exact)
		/*
			l        = int32(deadbeef)   // compile error: constant overflows int32
			m        = float64(1e309)    // compile error: constant overflows float64
			n        = uint(-1)          // compile error: constant underflows uint
		*/
	)
	fmt.Println("deadbeef, i, j, k = ", deadbeef, i, j, k)
}
