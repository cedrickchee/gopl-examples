// Exercise 2.1: Add types, constants, and functions to tempconv for processing
// temperatures in the Kelvin scale, where zero Kelvin is −273.15°C and a
// difference of 1K has the same magnitude as 1°C.
package tempconv

import "testing"

// $ go test -v -run StringFormat gopl.io/ch2/exercise2-1
func TestStringFormat(t *testing.T) {
	k := Kelvin(9)
	t.Log(k.String()) // "9K"
	t.Logf("%v\n", k) // "9K"; no need to call String explicitly
}

func TestConversion(t *testing.T) {
	got := KToC(1)
	want := Celcius(-272.15)
	if got != want {
		t.Fatalf("want %g, got %g", want, got)
	}
	t.Log(got) // "-272.15°C"
}
