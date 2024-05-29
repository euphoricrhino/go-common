package zraster

import (
	"image/color"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/golang/freetype/raster"
)

// zBuffer stores for each pixel a list of zColors.
type zBuffer [][]*zColor

// strokeRecorder implements raster.Painter so every rasterized stroke will be recorded
// with the z-distance value of the pixels, which will later be sorted and rendered in order.
type strokeRecorder struct {
	width   int
	height  int
	zbuf    zBuffer
	p1      *graphix.Projection
	p2      *graphix.Projection
	dd      float64
	strokeR uint32
	strokeG uint32
	strokeB uint32
	strokeA uint32

	// These maps stores all the pixels touched by the previous stroke and the current stroke.
	front   int
	touched [2]map[int]struct{}
}

func newStrokeRecorder(width, height int) *strokeRecorder {
	rec := &strokeRecorder{
		width:  width,
		height: height,
		zbuf:   make(zBuffer, width*height),
	}
	rec.touched[0] = make(map[int]struct{})
	rec.touched[1] = make(map[int]struct{})
	return rec
}

func (rec *strokeRecorder) resetForPath() {
	clear(rec.touched[0])
	clear(rec.touched[1])
}

// Prepares the recorder for the rasterization of the next line segment stroke by
// storing the endpoints' position and stroke color.
func (rec *strokeRecorder) prepareForRasterization(p1, p2 *graphix.Projection, color color.Color) {
	// Prepare the recorder for this stroke.
	rec.p1 = p1
	rec.p2 = p2
	dx, dy := p1[0]-p2[0], p1[1]-p2[1]
	rec.dd = dx*dx + dy*dy
	rec.strokeR, rec.strokeG, rec.strokeB, rec.strokeA = color.RGBA()
	// Swap the front/back maps.
	rec.front = 1 - rec.front
	// Clear the touched map for the new stroke.
	clear(rec.touched[rec.front])
}

// Update the z-buffer with the pixel touched by the rasterizer.
func (rec *strokeRecorder) updateZBuf(x, y int, r, g, b, a uint32) {
	// Computes the z depth of the pixel by linearly interpolating between the two end points.
	// Due to rasterizing, (x,y) may not be on the line connecting the two end points.
	// The projection point from (x,y) to this line is used to interpolate z depth.
	z := rec.p1[2]
	// Degenerate case.
	if rec.dd == 0 {
		if z < rec.p2[2] {
			z = rec.p2[2]
		}
	} else {
		dx1, dy1 := float64(x)-rec.p1[0], float64(y)-rec.p1[1]
		dd1 := dx1*dx1 + dy1*dy1
		dx2, dy2 := float64(x)-rec.p2[0], float64(y)-rec.p2[1]
		dd2 := dx2*dx2 + dy2*dy2
		t := ((dd1-dd2)/rec.dd + 1) / 2
		z = rec.p1[2] + t*(rec.p2[2]-rec.p1[2])
	}

	i := y*rec.width + x
	if _, found := rec.touched[1-rec.front][i]; found {
		lastIdx := len(rec.zbuf[i]) - 1
		// The last stroke of the same path touched the same pixel, we will not record both to avoid making the
		// shared vertex brighter than other part of the path. We simply keep the one with greater opacity.
		if a > rec.zbuf[i][lastIdx].a {
			rec.zbuf[i][lastIdx].r = r
			rec.zbuf[i][lastIdx].g = g
			rec.zbuf[i][lastIdx].b = b
			rec.zbuf[i][lastIdx].a = a
			rec.zbuf[i][lastIdx].z = z
		}
	} else {
		rec.zbuf[i] = append(rec.zbuf[i], &zColor{
			r: r,
			g: g,
			b: b,
			a: a,
			z: z,
		})
	}
	rec.touched[rec.front][i] = struct{}{}
}

// Paint make strokeRecorder implement raster.Painter so we get the call for each rasterized span.
func (rec *strokeRecorder) Paint(ss []raster.Span, done bool) {
	for _, s := range ss {
		if s.Y < 0 {
			continue
		}
		if s.Y >= rec.height {
			return
		}
		if s.X0 < 0 {
			s.X0 = 0
		}
		if s.X1 > rec.width {
			s.X1 = rec.width
		}
		if s.X0 >= s.X1 {
			continue
		}
		r, g, b, a := rec.strokeR*s.Alpha, rec.strokeG*s.Alpha, rec.strokeB*s.Alpha, rec.strokeA*s.Alpha
		for x := s.X0; x < s.X1; x++ {
			rec.updateZBuf(x, s.Y, r, g, b, a)
		}
	}
}
