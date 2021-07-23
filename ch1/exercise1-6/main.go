// Exercise 1.6: Modify the Lissajous program to produce images in multiple
// colors by adding more values to palette and then displaying them by changing
// the third argument of SetColorIndex in some interesting way.

package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

// Palette is emerald green on black
var palette = []color.Color{color.Black, color.White,
	color.RGBA{0xff, 0, 0, 0xff},
	color.RGBA{63, 195, 128, 1},
	color.RGBA{0, 0, 0xff, 0xff},
}

const (
	blackIndex = 0 // first color in palette
	whiteIndex = 1 // next color in palette
	redIndex   = 2 // third color in palette
	greenIndex = 3 // fourth color in palette
	blueIndex  = 4 // fifth color in palette
)

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		var colorIndex uint8 = whiteIndex
		if i%3 == 0 {
			colorIndex = uint8(rand.Uint32()%5 + 1)
		}
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), colorIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
