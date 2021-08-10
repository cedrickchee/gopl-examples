// The jpeg command reads a PNG image from the standard input
// and writes it as a JPEG image to the standard output.
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // register PNG decoder
	"io"
	"os"
)

func main() {
	if err := toJPEG(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

/*
Run:

# With blank import
If we feed the output of `gopl.io/ch3/mandelbrot` to the converter program, it
detects the PNG input format and writes a JPEG version.

$ go build gopl.io/ch3/mandelbrot
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
*/

/*
Run:

# Without blank import
Notice the blank import of `image/png`. Without that line, the program compiles
and links as usual but can no longer recognize or decode input in PNG format:

$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format
*/

// Notes
//
// Here’s how it works. The standard library provides decoders for GIF, PNG,
// and JPEG, and users may provide others, but to keep executables small,
// decoders are not included in an application unless explicitly requested.
//
// The effect is that an application need only blank-import the package for the
// format it needs to make the `image.Decode` function able to decode it.
