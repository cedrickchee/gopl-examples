package cake_test

import (
	"flag"
	"os"
	"testing"
	"time"

	"gopl.io/ch8/cake"
)

var defaults = cake.Shop{
	Cakes:        20,
	BakeTime:     10 * time.Millisecond,
	NumIcers:     1,
	IceTime:      10 * time.Millisecond,
	InscribeTime: 10 * time.Millisecond,
}

// Reference:
// https://stackoverflow.com/questions/38484942/golang-testing-with-init-func

// func init() {
// 	fmt.Println("initializing...")
// 	// testing.Init()

// 	// defaults.Verbose = testing.Verbose()
// }

// Fix: panic: testing: Verbose called before Init
// Reference:
// https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	defaults.Verbose = testing.Verbose()
	os.Exit(m.Run())
}

func Benchmark(b *testing.B) {
	// Baseline: one baker, one icer, one inscriber.
	// Each step takes exactly 10ms.  No buffers.
	cakeshop := defaults
	// Each benchmark must execute the code under test b.N times.
	cakeshop.Work(b.N) // 224 ms
}

func BenchmarkBuffers(b *testing.B) {
	// Adding buffers has no effect.
	cakeshop := defaults
	cakeshop.BakeBuf = 10
	cakeshop.IceBuf = 10
	cakeshop.Work(b.N) // 224 ms
}

func BenchmarkVariable(b *testing.B) {
	// Adding variability to rate of each step
	// increases total time due to channel delays.
	cakeshop := defaults
	cakeshop.BakeStdDev = cakeshop.BakeTime / 4
	cakeshop.IceStdDev = cakeshop.IceTime / 4
	cakeshop.InscribeStdDev = cakeshop.InscribeTime / 4
	cakeshop.Work(b.N) // 259 ms
}

func BenchmarkVariableBuffers(b *testing.B) {
	// Adding channel buffers reduces
	// delays resulting from variability.
	cakeshop := defaults
	cakeshop.BakeStdDev = cakeshop.BakeTime / 4
	cakeshop.IceStdDev = cakeshop.IceTime / 4
	cakeshop.InscribeStdDev = cakeshop.InscribeTime / 4
	cakeshop.BakeBuf = 10
	cakeshop.IceBuf = 10
	cakeshop.Work(b.N) // 244 ms
}

func BenchmarkSlowIcing(b *testing.B) {
	// Making the middle stage slower
	// adds directly to the critical path.
	cakeshop := defaults
	cakeshop.IceTime = 50 * time.Millisecond
	cakeshop.Work(b.N) // 1.032 s
}

func BenchmarkSlowIcingManyIcers(b *testing.B) {
	// Adding more icing cooks reduces the cost of icing
	// to its sequential component, following Amdahl's Law.
	cakeshop := defaults
	cakeshop.IceTime = 50 * time.Millisecond
	cakeshop.NumIcers = 5
	cakeshop.Work(b.N) // 288ms
}

/*
Benchmarks:
$ go test -bench=. gopl.io/ch8/cake
goos: linux
goarch: amd64
pkg: gopl.io/ch8/cake
cpu: Intel(R) Core(TM) i7-1065G7 CPU @ 1.30GHz
Benchmark-8                                    5         226120896 ns/op
BenchmarkBuffers-8                             5         226709955 ns/op
BenchmarkVariable-8                            4         270192422 ns/op
BenchmarkVariableBuffers-8                     4         250861534 ns/op
BenchmarkSlowIcing-8                           1        1031919583 ns/op
BenchmarkSlowIcingManyIcers-8                  4         271209112 ns/op
PASS
ok      gopl.io/ch8/cake        11.143s
*/
