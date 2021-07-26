package byteconstant_test

import (
	"testing"

	byteconstant "gopl.io/ch3/exercise3-13"
)

// Run:
// $ go test -v -run ByteConstant gopl.io/ch3/exercise3-13
func TestByteConstant(t *testing.T) {
	t.Log("1 MB =", byteconstant.MB)
	got := byteconstant.MB
	want := 1_000_000
	if got != want {
		t.Fatalf("want %v, got %v", want, got)
	}
}
