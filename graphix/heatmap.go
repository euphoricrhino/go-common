package graphix

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
)

const max16f = float64(math.MaxUint16)

// LoadHeatmap loads the heatmap from the given PNG file, uses its first row of pixels as color spectrum, and returns the gamma-corrected color spectrum.
func LoadHeatmap(file string, gamma float64) ([]color.Color, error) {
	// Load heatmap file.
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open heatmap file: %v", err)
	}
	defer f.Close()
	hm, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode heatmap as PNG: %v", err)
	}
	rect := hm.Bounds()
	width := rect.Max.X - rect.Min.X
	heatmap := make([]color.Color, width)
	for i := 0; i < width; i++ {
		r, g, b, _ := hm.At(i+rect.Min.X, rect.Min.Y).RGBA()
		r16 := uint16(math.Pow(float64(r)/max16f, gamma) * max16f)
		g16 := uint16(math.Pow(float64(g)/max16f, gamma) * max16f)
		b16 := uint16(math.Pow(float64(b)/max16f, gamma) * max16f)
		heatmap[i] = color.RGBA64{R: r16, G: g16, B: b16, A: math.MaxUint16}
	}
	return heatmap, nil
}
