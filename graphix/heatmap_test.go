package graphix

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestImage(t *testing.T) (string, []color.Color) {
	f, err := os.CreateTemp(os.TempDir(), "test-heatmap")
	assert.NoError(t, err)
	defer f.Close()

	cnt := 200
	colors := make([]color.Color, 0, cnt)
	c := 0
	for i := 0; i < cnt; i++ {
		switch c % 3 {
		case 0:
			colors = append(colors, color.RGBA64{R: 0xffff, G: 0, B: 0, A: 0xffff})
		case 1:
			colors = append(colors, color.RGBA64{R: 0, G: 0xffff, B: 0, A: 0xffff})
		case 2:
			colors = append(colors, color.RGBA64{R: 0, G: 0, B: 0xffff, A: 0xffff})
		}
		c++
	}
	img := image.NewRGBA(image.Rect(0, 0, cnt, 300))
	for x := 0; x < cnt; x++ {
		img.Set(x, 0, colors[x])
	}
	assert.NoError(t, png.Encode(f, img))
	return f.Name(), colors
}

func TestLoadHeatmapErrors(t *testing.T) {
	_, err := LoadHeatmap("?@$>", 1.0)
	assert.ErrorContains(t, err, "failed to open heatmap file: ")

	f, err := os.CreateTemp(os.TempDir(), "empty")
	assert.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	_, err = LoadHeatmap(f.Name(), 1.0)
	assert.ErrorContains(t, err, "failed to decode heatmap as PNG")
}

func TestLoadHeatmap(t *testing.T) {
	fn, colors := createTestImage(t)
	defer os.Remove(fn)

	hm, err := LoadHeatmap(fn, 1.0)
	assert.NoError(t, err)
	assert.Equal(t, colors, hm)
}
