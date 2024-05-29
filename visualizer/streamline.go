package visualizer

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/graphix/zraster"
)

// Represents a point along the traced trajectory.
type point struct {
	// The length of the tangent vector at this point.
	tan float64
	pos *graphix.Vec3
}

// Symmetry defines a symmetry transform for a traced trajectory.
type Symmetry struct {
	Transform graphix.Transform
	Color     color.Color
}

// Trajectory represents a series of points traced from start.
type Trajectory struct {
	Start *graphix.Vec3
	// User-provided callback that returns whether the tracing of streamline should
	// terminate at point x whose tangent is tan.
	AtEnd func(x, tan *graphix.Vec3) bool
	Color color.Color
	// Symmetry transforms applicable to this trajectory.
	// This saves us from tracing these multiple times while applying the symmetry to the original
	// one will achieve the same result.
	syms   []*Symmetry
	points []*point
}

// AddSymmetry adds a symmetry transform for a trajectory.
func (traj *Trajectory) AddSymmetry(transform graphix.Transform, color color.Color) {
	traj.syms = append(traj.syms, &Symmetry{
		Transform: transform,
		Color:     color,
	})
}

// StreamLineOptions defines the options to run VisualizeStreamLines().
type StreamLineOptions struct {
	// Generated images will be saved to this directory.
	OutDir      string
	CameraOrbit graphix.CameraOrbit
	// Step size for the Runge Kutta-4 tracing.
	Step float64
	// Points within this distance will not be recorded in the trajectory (but will still participate in the Runge Kutta calculation).
	MinDist float64
	// User-provided tangent function at x. Result should be written to tan.
	TangentAt func(tan, x *graphix.Vec3)
	LineWidth float64
	// Fading factor for max/min tangent values, intermediate tangent values will be linearly interpolated and then gamma corrected by FadingGamma.
	MinFading   float64
	MaxFading   float64
	FadingGamma float64
	// Concurrency
	Workers int
}

// VisualizeStreamLines runs the stream line tracing and rendering given the options and trajectory settings. Upon completion
// trajs internal data structure would have been modified.
func VisualizeStreamLines(opts StreamLineOptions, trajs []*Trajectory) {
	var wg sync.WaitGroup
	wg.Add(opts.Workers)
	// See multi variable Runge Kutta-4 at https://www.myphysicslab.com/explain/runge-kutta-en.html
	h := opts.Step
	h2 := opts.Step / 2
	h6 := opts.Step / 6
	h3 := opts.Step / 3
	for w := 0; w < opts.Workers; w++ {
		go func(wk int) {
			defer wg.Done()
			// Thread-local scratch area variables.
			var x, a, xb, b, xc, c, xd, d, tmp graphix.Vec3
			for i := range trajs {
				if i%opts.Workers != wk {
					continue
				}
				traj := trajs[i]
				x.Copy(traj.Start)
				var last *point
				for {
					opts.TangentAt(&a, &x)
					if traj.AtEnd(&x, &a) {
						// Put as the last point of trajectory regardless of the min dist.
						traj.points = append(traj.points, &point{tan: a.Norm(), pos: graphix.NewCopyVec3(&x)})
						break
					}
					if last == nil || tmp.Sub(&x, last.pos).Norm() >= opts.MinDist {
						// Accummulate point onto trajectory if it is at least step away from the previous point.
						pt := &point{tan: a.Norm(), pos: graphix.NewCopyVec3(&x)}
						traj.points = append(traj.points, pt)
						last = pt
					}
					xb.Add(&x, tmp.Scale(&a, h2))
					opts.TangentAt(&b, &xb)
					xc.Add(&x, tmp.Scale(&b, h2))
					opts.TangentAt(&c, &xc)
					xd.Add(&x, tmp.Scale(&c, h))
					opts.TangentAt(&d, &xd)
					x.Add(&x, tmp.Scale(&a, h6))
					x.Add(&x, tmp.Scale(&b, h3))
					x.Add(&x, tmp.Scale(&c, h3))
					x.Add(&x, tmp.Scale(&d, h6))
				}
			}
		}(w)
	}
	wg.Wait()

	fmt.Fprintln(os.Stdout, "completed tracing all trajectories.")

	seenPoints := false
	// Calculate max and min of tangent lengths.
	maxTan, minTan := math.Inf(-1), math.Inf(1)
	for _, traj := range trajs {
		for _, pt := range traj.points {
			maxTan = math.Max(maxTan, pt.tan)
			minTan = math.Min(minTan, pt.tan)
			seenPoints = true
		}
	}
	// Degenerate case.
	if !seenPoints {
		maxTan, minTan = 1.0, 0.0
	}

	// Create line segments for rendering.
	var lines []*zraster.SpaceLine
	for _, traj := range trajs {
		r, g, b, _ := traj.Color.RGBA()
		for i := 0; i < len(traj.points)-1; i++ {
			// Take the average tangent between the two endpoints, then calculate the fading factor.
			tan := (traj.points[i].tan + traj.points[i+1].tan) / 2
			fading := opts.MinFading + (opts.MaxFading-opts.MinFading)*math.Pow((tan-minTan)/(maxTan-minTan), opts.FadingGamma)
			// Premultiply for the uint32 color fading below: 255/65535=1/257.
			fading /= 257

			// Instead of applying fading to the alpha channel, we apply them to the RGB channels and set alpha to 1.
			// This is because zraster will stroke individual line segment of the path, causing double-rendering of the connecting
			// points. If alpha is not 1, the connecting points will appear brighter than the rest of the path.
			lines = append(lines, &zraster.SpaceLine{
				Start:     traj.points[i].pos,
				End:       traj.points[i+1].pos,
				Color:     color.NRGBA{R: uint8(float64(r) * fading), G: uint8(float64(g) * fading), B: uint8(float64(b) * fading), A: 0xff},
				LineWidth: opts.LineWidth,
			})
			// Apply the symmetry transforms.
			for _, sym := range traj.syms {
				sr, sg, sb, _ := sym.Color.RGBA()
				lines = append(lines, &zraster.SpaceLine{
					Start:     sym.Transform.Apply(graphix.BlankVec3(), traj.points[i].pos),
					End:       sym.Transform.Apply(graphix.BlankVec3(), traj.points[i+1].pos),
					Color:     color.NRGBA{R: uint8(float64(sr) * fading), G: uint8(float64(sg) * fading), B: uint8(float64(sb) * fading), A: 0xff},
					LineWidth: opts.LineWidth,
				})
			}
		}
	}

	for f := 0; f < opts.CameraOrbit.NumPositions(); f++ {
		img := zraster.Run(zraster.Options{
			Camera:  opts.CameraOrbit.GetCamera(f),
			Lines:   lines,
			Workers: opts.Workers,
		})
		fn := filepath.Join(opts.OutDir, fmt.Sprintf("frame-%04v.png", f))
		file, err := os.Create(fn)
		if err != nil {
			panic(fmt.Sprintf("failed to create output file '%v': %v", fn, err))
		}
		if err := png.Encode(file, img); err != nil {
			panic(fmt.Sprintf("failed to encode to PNG: %v", err))
		}
		file.Close()
		fmt.Fprintf(os.Stdout, "generated %v.\n", fn)
	}
}