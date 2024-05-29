package zraster

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/stretchr/testify/assert"
)

func TestZRasterRunTriangle(t *testing.T) {
	// Three lines, with cyclic overlapping relationship.
	// red over blue, blue over green, green over red.
	paths := []*SpacePath{
		{
			Segments: []*SpaceVertex{{
				Pos:   graphix.NewVec3(1, 5, -10),
				Color: color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff},
			}},
			End:       graphix.NewVec3(-3, -5, 10),
			LineWidth: 3,
		},
		{
			Segments: []*SpaceVertex{{
				Pos:   graphix.NewVec3(-1, 5, 10),
				Color: color.RGBA{R: 0, G: 0xff, B: 0, A: 0xff},
			}},
			End:       graphix.NewVec3(3, -5, -10),
			LineWidth: 5,
		},
		{
			Segments: []*SpaceVertex{{
				Pos:   graphix.NewVec3(-5, -3, 5),
				Color: color.RGBA{R: 0, G: 0, B: 0xff, A: 0xff},
			}},
			End:       graphix.NewVec3(5, -3, -5),
			LineWidth: 7,
		},
	}
	img := Run(Options{
		Camera: graphix.NewCamera(
			graphix.NewViewTransform(graphix.NewVec3(0, 0, 8), graphix.NewVec3(0, 0, -1), graphix.NewVec3(0, 1, 0)),
			graphix.NewOrthographic(),
			graphix.NewScreen(800, 800, -6, -6, 6, 6),
		),
		Paths:   paths,
		Workers: 1,
	})

	f, err := os.CreateTemp(os.TempDir(), "triangle*.png")
	assert.NoError(t, err)
	defer f.Close()
	assert.NoError(t, png.Encode(f, img))
	fmt.Fprintln(os.Stdout, f.Name())
}

// This test shows that the alpha of the internal vertices of a path is not double counted.
// The generated image should not have brighter spots at the internal vertices of the spiral path.
func TestZRasterRunSpiral(t *testing.T) {
	var segments []*SpaceVertex
	dl := 2.0
	dtheta := 5 * math.Pi / 180
	color := graphix.RandColor()
	// Set an non-opaque transparency.
	color.A = 0xff
	l := 10.0
	theta := 0.0
	for i := 0; i < 40; i++ {
		l += dl
		theta += dtheta
		color.A -= 5
		segments = append(segments, &SpaceVertex{
			Pos:   graphix.NewVec3(l*math.Cos(theta), l*math.Sin(theta), 0),
			Color: color,
		})
	}
	l += dl
	theta += dtheta
	paths := []*SpacePath{{Segments: segments, End: graphix.NewVec3(l*math.Cos(theta), l*math.Sin(theta), 0), LineWidth: 15}}

	img := Run(Options{
		Camera: graphix.NewCamera(
			graphix.NewViewTransform(graphix.NewVec3(0, 0, 8), graphix.NewVec3(0, 0, -1), graphix.NewVec3(0, 1, 0)),
			graphix.NewOrthographic(),
			graphix.NewScreen(800, 800, -100, -100, 100, 100),
		),
		Paths:   paths,
		Workers: 1,
	})
	f, err := os.CreateTemp(os.TempDir(), "spiral*.png")
	assert.NoError(t, err)
	defer f.Close()
	assert.NoError(t, png.Encode(f, img))
	fmt.Fprintln(os.Stdout, f.Name())
}
