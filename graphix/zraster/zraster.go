package zraster

import (
	"image"
	"image/color"
	"image/draw"
	"sort"
	"sync"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

// SpaceVertex represents the starting vertex of a line segment of a SpacePath.
// The color applies to the line segment.
type SpaceVertex struct {
	Pos   *graphix.Vec3
	Color color.Color
}

// SpacePath represents a 3D path.
type SpacePath struct {
	// All vertices up to the second-to-last vertex.
	Segments  []*SpaceVertex
	End       *graphix.Vec3
	LineWidth float64
}

// Options defines the options for zraster.Run().
type Options struct {
	Camera *graphix.Camera
	// All the 3D paths to render.
	Paths     []*SpacePath
	NearZClip float64
	// Concurrency.
	Workers int
}

// Run implements a specialized rasterizer for 3D paths.
// It renders the paths into an image while respecting their z-order.
func Run(opts Options) draw.Image {
	chs := make([]chan zBuffer, opts.Workers)
	for i := range chs {
		chs[i] = make(chan zBuffer)
	}

	for w := 0; w < opts.Workers; w++ {
		go zworker(w, opts.Workers, &opts, chs[w])
	}

	// Wait for all workers to finish updating their zbuffers.
	zbufs := make([]zBuffer, len(chs))
	for i, ch := range chs {
		zbufs[i] = <-ch
	}

	img := image.NewRGBA(image.Rect(0, 0, opts.Camera.Screen().Width(), opts.Camera.Screen().Height()))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0xff}}, image.Point{}, draw.Src)

	var wg sync.WaitGroup
	wg.Add(opts.Workers)
	for w := 0; w < opts.Workers; w++ {
		go func(wk int) {
			for x := 0; x < opts.Camera.Screen().Width(); x++ {
				for y := 0; y < opts.Camera.Screen().Height(); y++ {
					i := y*opts.Camera.Screen().Width() + x
					if i%opts.Workers != wk {
						continue
					}
					// Merge and sort zbuffers from all concurrent shards at pixel i.
					l := 0
					for _, zbuf := range zbufs {
						l += len(zbuf[i])
					}
					sorted := make([]*zColor, 0, l)
					for _, zbuf := range zbufs {
						sorted = append(sorted, zbuf[i]...)
					}
					sort.Sort(sortByZ(sorted))

					idx := y*img.Stride + x*4

					// Paint the pixels from far to near, multiplying the alpha at each layer in order.
					for _, zc := range sorted {
						const m = 1<<16 - 1
						dr := uint32(img.Pix[idx+0])
						dg := uint32(img.Pix[idx+1])
						db := uint32(img.Pix[idx+2])
						da := uint32(img.Pix[idx+3])
						a := (m - (zc.a / m)) * 0x101
						img.Pix[idx+0] = uint8((dr*a + zc.r) / m >> 8)
						img.Pix[idx+1] = uint8((dg*a + zc.g) / m >> 8)
						img.Pix[idx+2] = uint8((db*a + zc.b) / m >> 8)
						img.Pix[idx+3] = uint8((da*a + zc.a) / m >> 8)
					}
				}
			}
			wg.Done()
		}(w)
	}
	wg.Wait()

	return img
}

// One of the concurrent workers to work on a shard of the whole paths set, generating its own subset of z-buffers for each pixel.
// All z-buffers of the same pixel will be merged and sorted subsequently.
func zworker(w, workers int, opts *Options, ch chan<- zBuffer) {
	width, height := opts.Camera.Screen().Width(), opts.Camera.Screen().Height()
	rasterizer := raster.NewRasterizer(width, height)
	rasterizer.UseNonZeroWinding = true
	rec := newStrokeRecorder(width, height)
	// Thread-local scratch area variables.
	var v1, v2 graphix.Vec3
	var p1, p2 graphix.Projection
	var fp1, fp2 fixed.Point26_6

	for i, path := range opts.Paths {
		// Work only on worker's own shard.
		if i%opts.Workers != w {
			continue
		}

		// Degenerate path.
		if len(path.Segments) == 0 {
			continue
		}

		rec.resetForPath()
		// Strokes a 3D line segment from pos1 to pos2.
		stroke := func(pos1, pos2 *graphix.Vec3, color color.Color) {
			// View-transform to canonical camera coordinates.
			opts.Camera.ViewTransform().Apply(&v1, pos1)
			opts.Camera.ViewTransform().Apply(&v2, pos2)
			// Do the projection.
			opts.Camera.Projector().Project(&p1, &v1)
			opts.Camera.Projector().Project(&p2, &v2)
			// Discard the line if both ends are behind the z-clip plane.
			if p1[2] < opts.NearZClip && p2[2] < opts.NearZClip {
				return
			}
			// Clip the near end at the z-clip plane.
			if p1[2] < opts.NearZClip {
				zclip(&p1, &p2, opts.NearZClip)
			} else if p2[2] < opts.NearZClip {
				zclip(&p2, &p1, opts.NearZClip)
			}
			// Scale to screen dimensions.
			opts.Camera.Screen().Map(&p1, &p1)
			opts.Camera.Screen().Map(&p2, &p2)
			toFixedPoint(&fp1, &p1)
			toFixedPoint(&fp2, &p2)
			// Stroke the rasterizer path.
			var rasterPath raster.Path
			rasterPath.Start(fp1)
			rasterPath.Add1(fp2)

			rasterizer.Clear()
			rasterizer.AddStroke(rasterPath, toFixed(path.LineWidth), nil, nil)
			rec.prepareForRasterization(&p1, &p2, color)
			rasterizer.Rasterize(rec)
		}

		i := 0
		for ; i < len(path.Segments)-1; i++ {
			stroke(path.Segments[i].Pos, path.Segments[i+1].Pos, path.Segments[i].Color)
		}
		stroke(path.Segments[i].Pos, path.End, path.Segments[i].Color)
	}

	ch <- rec.zbuf
}

func zclip(near, far *graphix.Projection, z float64) {
	t := (z - far[2]) / (near[2] - far[2])
	near[0] = far[0] + t*(near[0]-far[0])
	near[1] = far[1] + t*(near[1]-far[1])
	near[2] = z
}

func toFixedPoint(fp *fixed.Point26_6, p *graphix.Projection) {
	fp.X = toFixed(p[0])
	fp.Y = toFixed(p[1])
}

func toFixed(f float64) fixed.Int26_6 {
	return fixed.Int26_6(f * 64)
}
