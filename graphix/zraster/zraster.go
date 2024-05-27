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

// SpaceLine represents a 3D line segment with its own color and line width.
type SpaceLine struct {
	Start     *graphix.Vec3
	End       *graphix.Vec3
	Color     color.Color
	LineWidth float64
}

// Options defines the options for zraster.Run().
type Options struct {
	Camera *graphix.Camera
	// The near-Z clipping - only objects further than this distance will be drawn.
	NearZClip float64
	// All the 3D line segments to render.
	Lines []*SpaceLine
	// Concurrency.
	Workers int
}

// Run implements a specialized rasterizer for 3D line segments.
// It renders the line segments into an image while respecting their z-order.
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

// One of the concurrent workers to work on a shard of the whole line set, generating its own subset of z-buffers for each pixel.
// All z-buffers of the same pixel will be merged and sorted subsequently.
func zworker(w, workers int, opts *Options, ch chan<- zBuffer) {
	width, height := opts.Camera.Screen().Width(), opts.Camera.Screen().Height()
	rasterizer := raster.NewRasterizer(width, height)
	rasterizer.UseNonZeroWinding = true
	rec := &strokeRecorder{
		width:  width,
		height: height,
		zbuf:   make(zBuffer, width*height),
	}
	// Thread-local scratch area variables.
	v1 := graphix.BlankVec3()
	v2 := graphix.BlankVec3()
	p1 := graphix.BlankProjection()
	p2 := graphix.BlankProjection()
	fp1 := &fixed.Point26_6{}
	fp2 := &fixed.Point26_6{}
	for i, line := range opts.Lines {
		// Work only on worker's own shard.
		if i%opts.Workers != w {
			continue
		}
		// View-transform to canonical camera coordinates.
		opts.Camera.ViewTransform().Apply(v1, line.Start)
		opts.Camera.ViewTransform().Apply(v2, line.End)
		// Do the projection.
		opts.Camera.Projector().Project(p1, v1)
		opts.Camera.Projector().Project(p2, v2)
		// Discard the line if both ends are behind the z-clip plane.
		if p1[2] < opts.NearZClip && p2[2] < opts.NearZClip {
			continue
		}
		// Clip the near end at the z-clip plane.
		if p1[2] < opts.NearZClip {
			zclip(p1, p2, opts.NearZClip)
		} else if p2[2] < opts.NearZClip {
			zclip(p2, p1, opts.NearZClip)
		}
		// Scale to screen dimensions.
		opts.Camera.Screen().Map(p1, p1)
		opts.Camera.Screen().Map(p2, p2)
		toFixedPoint(fp1, p1)
		toFixedPoint(fp2, p2)
		// Stroke the rasterizer path.
		var path raster.Path
		path.Start(*fp1)
		path.Add1(*fp2)
		rec.prepareForStroke(p1, p2, line.Color)
		rasterizer.Clear()
		rasterizer.AddStroke(path, toFixed(line.LineWidth), nil, nil)
		rasterizer.Rasterize(rec)
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
