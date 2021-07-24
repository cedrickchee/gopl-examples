package tempconv

import (
	"testing"
)

// $ go test -v -run Example_one gopl.io/ch2/tempconv0
func TestExample_one(t *testing.T) {
	got := BoilingC - FreezingC
	want := Celcius(100)
	if got != want {
		t.Fatalf("want %g, got %g", want, got)
	}
	t.Logf("%g\n", got) // "100" °C
	boilingF := CToF(BoilingC)
	t.Logf("%g\n", boilingF-CToF(FreezingC)) // "180" °F
	/*
		fmt.Printf("%g\n", boilingF-FreezingC)       // compile error: type mismatch
	*/
}

func TestExample_three(t *testing.T) {
	var c Celcius
	var f Fahrenheit
	t.Log(c == 0) // "true"
	t.Log(f >= 0) // "true"
	// t.Log(c == f) // "compile error: type mismatch"
	t.Log(c == Celcius(f)) // "true"!
}

func TestExample_two(t *testing.T) {
	c := FToC(212.0)
	t.Log(c.String()) // "100°C"
	t.Logf("%v\n", c) // "100°C"; no need to call String explicitly
	t.Logf("%s\n", c) // "100°C"
	t.Log(c)          // "100°C"
	t.Logf("%g\n", c) // "100"; does not call String
	t.Log(float64(c)) // "100"; does not call String
}
