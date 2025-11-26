package zraster

import (
	"bytes"
	"image/color"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/stretchr/testify/assert"
)

func zrasterTestHelper(t *testing.T, paths []*SpacePath, benchmarkFile string) {
	img := Run(Settings{
		Camera: graphix.NewCamera(
			graphix.NewViewTransform(
				graphix.NewVec3(0, 0, 8),
				graphix.NewVec3(0, 0, -1),
				graphix.NewVec3(0, 1, 0),
			),
			graphix.NewOrthographic(),
			graphix.NewScreen(800, 800, -6, -6, 6, 6),
		),
		Paths:   paths,
		Workers: 1,
	})

	benchmarkBytes, err := os.ReadFile(benchmarkFile)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	assert.NoError(t, png.Encode(buf, img))
	assert.Equal(t, benchmarkBytes, buf.Bytes())
}

// Three lines, with cyclic overlapping relationship.
// red over blue, blue over green, green over red.
func TestZRasterRunTriangle(t *testing.T) {
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

	zrasterTestHelper(t, paths, "testdata/triangle.png")
}

// This test shows that the alpha of the internal vertices of a path is not double counted.
// The generated image should not have brighter spots at the internal vertices of the spiral path.
func TestZRasterRunSpiral(t *testing.T) {
	var segments []*SpaceVertex
	dl := .1
	dtheta := 5 * math.Pi / 180
	color := color.NRGBA{R: 0, G: 0, B: 0xff, A: 0xff}
	l := .5
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
	paths := []*SpacePath{
		{
			Segments:  segments,
			End:       graphix.NewVec3(l*math.Cos(theta), l*math.Sin(theta), 0),
			LineWidth: 15,
		},
	}

	zrasterTestHelper(t, paths, "testdata/spiral.png")
}

// This test shows a path with three segments: AB-BC-CD, where the projection of AB and CD intersect.
// Their projected intersection should have a brighter color due to the alpha blending (at different z),
// but B and C should not because they belong to adjacent strokes.
func TestZRasterRunSelfIntersectingPath(t *testing.T) {
	color := color.NRGBA{R: 0, G: 0xff, B: 0, A: 0x80}
	paths := []*SpacePath{{
		Segments: []*SpaceVertex{{
			Pos:   graphix.NewVec3(-5, -5, -5),
			Color: color,
		}, {
			Pos:   graphix.NewVec3(5, 5, 5),
			Color: color,
		}, {
			Pos:   graphix.NewVec3(-5, 5, -5),
			Color: color,
		}},
		End:       graphix.NewVec3(5, -5, -5),
		LineWidth: 15,
	}}
	zrasterTestHelper(t, paths, "testdata/self-intersecting.png")
}

// Line segments will be clipped by z-clip plane.
func TestZRasterRunZClip(t *testing.T) {
	color := color.NRGBA{R: 0, G: 0, B: 0xff, A: 0xff}
	paths := []*SpacePath{{
		Segments: []*SpaceVertex{{
			Pos:   graphix.NewVec3(4, 2, 6),
			Color: color,
		}, {
			Pos:   graphix.NewVec3(-4, 2, 6),
			Color: color,
		}, {
			Pos:   graphix.NewVec3(0, -2, 10),
			Color: color,
		}},
		End:       graphix.NewVec3(4, 2, 6),
		LineWidth: 3,
	}, {
		Segments: []*SpaceVertex{{
			Pos:   graphix.NewVec3(5, -3, 10),
			Color: color,
		}, {
			Pos:   graphix.NewVec3(5, 3, 6),
			Color: color,
		}, {
			Pos:   graphix.NewVec3(-5, 3, 6),
			Color: color,
		}, {
			Pos:   graphix.NewVec3(-5, -3, 10),
			Color: color,
		}},
		End:       graphix.NewVec3(5, -3, 10),
		LineWidth: 3,
	}}
	zrasterTestHelper(t, paths, "testdata/zclip.png")
}
