package zraster

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/stretchr/testify/assert"
)

func TestZRasterizerRun(t *testing.T) {
	// Three lines, with cyclic overlapping relationship.
	// red over blue, blue over green, green over red.
	lines := []*SpaceLine{
		{
			Start:     graphix.NewVec3(1, 5, -10),
			End:       graphix.NewVec3(-3, -5, 10),
			Color:     color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff},
			LineWidth: 3,
		},
		{
			Start:     graphix.NewVec3(-1, 5, 10),
			End:       graphix.NewVec3(3, -5, -10),
			Color:     color.RGBA{R: 0, G: 0xff, B: 0, A: 0xff},
			LineWidth: 5,
		},
		{
			Start:     graphix.NewVec3(-5, -3, 5),
			End:       graphix.NewVec3(5, -3, -5),
			Color:     color.RGBA{R: 0, G: 0, B: 0xff, A: 0xff},
			LineWidth: 7,
		},
	}
	img := Run(Options{
		Camera: graphix.NewCamera(
			graphix.NewViewTransform(graphix.NewVec3(0, 0, 8), graphix.NewVec3(0, 0, -1), graphix.NewVec3(0, 1, 0)),
			graphix.NewOrthographic(),
			graphix.NewScreen(800, 800, -6, -6, 6, 6),
		),
		Lines:   lines,
		Workers: 4,
	})

	f, err := os.CreateTemp(os.TempDir(), "test-zrasterizer*.png")
	assert.NoError(t, err)
	defer f.Close()
	assert.NoError(t, png.Encode(f, img))
	fmt.Fprintln(os.Stdout, f.Name())
}
