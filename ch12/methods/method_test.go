package methods_test

import (
	"strings"
	"time"

	"gopl.io/ch12/methods"
)

func ExamplePrintDuration() {
	methods.Print(time.Hour)
	// Output:
	// type time.Duration
	// func (time.Duration) Hours() float64
	// func (time.Duration) Microseconds() int64
	// func (time.Duration) Milliseconds() int64
	// func (time.Duration) Minutes() float64
	// func (time.Duration) Nanoseconds() int64
	// func (time.Duration) Round(time.Duration) time.Duration
	// func (time.Duration) Seconds() float64
	// func (time.Duration) String() string
	// func (time.Duration) Truncate(time.Duration) time.Duration
}

func ExamplePrintReplacer() {
	methods.Print(new(strings.Replacer))
	// Output:
	// type *strings.Replacer
	// func (*strings.Replacer) Replace(string) string
	// func (*strings.Replacer) WriteString(io.Writer, string) (int, error)
}

/*
Here are the methods belonging to two types, `time.Duration` and
`*strings.Replacer`:

$ go test -v gopl.io/ch12/methods

methods.Print(time.Hour)
// Output:
// type time.Duration
// func (time.Duration) Hours() float64
// func (time.Duration) Microseconds() int64
// func (time.Duration) Milliseconds() int64
// func (time.Duration) Minutes() float64
// func (time.Duration) Nanoseconds() int64
// func (time.Duration) Round(time.Duration) time.Duration
// func (time.Duration) Seconds() float64
// func (time.Duration) String() string
// func (time.Duration) Truncate(time.Duration) time.Duration

methods.Print(new(strings.Replacer))
// Output:
// type *strings.Replacer
// func (*strings.Replacer) Replace(string) string
// func (*strings.Replacer) WriteString(io.Writer, string) (int, error)
*/
