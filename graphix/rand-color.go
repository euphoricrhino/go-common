package graphix

import (
	"image/color"
	"math"
	"math/rand"
)

// The following formula is from wikipedia: https://en.wikipedia.org/wiki/HSL_and_HSV#HSV_to_RGB_alternative
// With S=V=1.
func hueAsRGB(h float64) color.NRGBA64 {
	rr := math.Mod(5.0+h*6.0, 6.0)
	gg := math.Mod(3.0+h*6.0, 6.0)
	bb := math.Mod(1.0+h*6.0, 6.0)

	r := 1.0 - math.Max(min(rr, 4.0-rr, 1.0), 0.0)
	g := 1.0 - math.Max(min(gg, 4.0-gg, 1.0), 0.0)
	b := 1.0 - math.Max(min(bb, 4.0-bb, 1.0), 0.0)

	m := float64(math.MaxUint16)
	return color.NRGBA64{R: uint16(m * r), G: uint16(m * g), B: uint16(m * b), A: math.MaxUint16}
}

// Pre-generated full-saturation colors in RGB.
var hueRGB []color.NRGBA64

func init() {
	for h := 0.0; h < 1; h += 1.0 / 1024.0 {
		hueRGB = append(hueRGB, hueAsRGB(h))
	}
}

// RandColor returns a full-saturation color of random hue.
func RandColor() color.NRGBA64 {
	return hueRGB[rand.Intn(len(hueRGB))]
}
