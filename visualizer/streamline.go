package visualizer

import (
	"fmt"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/graphix/zraster"
)

// Represents a sampled point along the traced trajectory to be rendered.
type renderPoint struct {
	// The length of the tangent vector at this point.
	tan float64
	pos *graphix.Vec3
}

// Symmetry defines a symmetry transform for a traced trajectory.
type Symmetry struct {
	Transform graphix.Transform
	Color     color.NRGBA64
}

// Trajectory represents a series of points traced from a start point.
type Trajectory struct {
	Start *graphix.Vec3
	// Initial step size for the adaptive Runge-Kutta tracing.
	// This value can be negative (reverse tracing).
	InitStep float64
	// Error bound for the adaptive Runge-Kutta tracing.
	Epsilon float64
	// User-provided callback that returns whether the tracing of streamline should
	// terminate at point x whose tangent is tan.
	AtEnd func(x, tan *graphix.Vec3) bool
	Color color.NRGBA64
	// Symmetry transforms applicable to this trajectory.
	// This saves us from tracing these multiple times while applying the symmetry to the original
	// one will achieve the same result.
	syms []*Symmetry
	// The sampled points along the trajectory for rendering.
	points []*renderPoint
}

// AddSymmetry adds a symmetry transform for a trajectory.
func (traj *Trajectory) AddSymmetry(transform graphix.Transform, color color.NRGBA64) {
	traj.syms = append(traj.syms, &Symmetry{
		Transform: transform,
		Color:     color,
	})
}

// StreamlineOptions defines the options to run VisualizeStreamLines().
type StreamlineOptions struct {
	CameraOrbit graphix.CameraOrbit
	// Adjacent points for rendering will have a distance between MinDist and MaxDist.
	MinDist float64
	MaxDist float64
	// User-provided tangent function at x. Result should be written to tan.
	TangentAt func(tan, x *graphix.Vec3)
	LineWidth float64
	// Fading factor for max/min tangent values, intermediate tangent values will be linearly interpolated and then gamma corrected by FadingGamma.
	MinFading   float64
	MaxFading   float64
	FadingGamma float64
	// Concurrency
	Workers int
	// User-provided callback functions for each generated image, together with the frame index.
	ImageCallbacks []func(img draw.Image, f int)
}

// VisualizeStreamlines runs the stream line tracing and rendering given the options and trajectory settings. Upon completion
// trajs internal data structure would have been modified.
func VisualizeStreamlines(opts StreamlineOptions, trajs []*Trajectory) {
	var wg sync.WaitGroup
	wg.Add(opts.Workers)
	for w := 0; w < opts.Workers; w++ {
		go func(wk int) {
			newTraceWorker(&opts).run(wk, trajs)
			wg.Done()
		}(w)
	}
	wg.Wait()

	fmt.Fprintln(os.Stdout, "completed tracing all trajectories.")

	// Calculate max and min of tangent lengths.
	maxTan, minTan := math.Inf(-1), math.Inf(1)
	for _, traj := range trajs {
		for _, pt := range traj.points {
			maxTan = math.Max(maxTan, pt.tan)
			minTan = math.Min(minTan, pt.tan)
		}
	}
	// Degenerate case - all tan's are exactly the same.
	if maxTan == minTan {
		maxTan, minTan = 1, 0
	}

	// Create zraster.SpacePaths for rendering.
	var paths []*zraster.SpacePath
	id := graphix.IdentityTransform()
	for _, traj := range trajs {
		if len(traj.points) == 0 {
			continue
		}
		syms := append([]*Symmetry{{Transform: id, Color: traj.Color}}, traj.syms...)
		for _, sym := range syms {
			path := &zraster.SpacePath{
				End:       sym.Transform.Apply(graphix.BlankVec3(), traj.points[len(traj.points)-1].pos),
				LineWidth: opts.LineWidth,
			}
			for i := 0; i < len(traj.points)-1; i++ {
				// Take the average tangent between the two endpoints, then calculate the fading factor.
				tan := (traj.points[i].tan + traj.points[i+1].tan) / 2
				fading := opts.MinFading + (opts.MaxFading-opts.MinFading)*math.Pow((tan-minTan)/(maxTan-minTan), opts.FadingGamma)

				path.Segments = append(path.Segments, &zraster.SpaceVertex{
					Pos: sym.Transform.Apply(graphix.BlankVec3(), traj.points[i].pos),
					Color: color.NRGBA64{
						R: sym.Color.R,
						G: sym.Color.G,
						B: sym.Color.B,
						A: uint16(float64(sym.Color.A) * fading),
					},
				})
			}
			paths = append(paths, path)
		}
	}

	for f := 0; f < opts.CameraOrbit.NumPositions(); f++ {
		img := zraster.Run(zraster.Options{
			Camera:  opts.CameraOrbit.GetCamera(f),
			Paths:   paths,
			Workers: opts.Workers,
		})
		for _, cb := range opts.ImageCallbacks {
			cb(img, f)
		}
	}
}

// An image callback that saves the image as png file in the given directory.
func SavePNG(outDir string) func(draw.Image, int) {
	return func(img draw.Image, f int) {
		fn := filepath.Join(outDir, fmt.Sprintf("frame-%04v.png", f))
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
