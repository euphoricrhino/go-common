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

// TraceSettings defines the settings for tracing streamlines.
type TraceSettings struct {
	// Adjacent points for rendering will have a distance between MinDist and MaxDist.
	MinDist float64
	MaxDist float64

	// User-provided tangent function at position x and frame f (i.e., discrete time t). Result should be written to tan.
	TangentAt func(tan, x *graphix.Vec3, f int)

	// Concurrency.
	Workers int
}

// TraceStreamlines runs the streamline tracing given the settings and trajectories.
// trajs[f] contains all the trajectories for discrete time/frame f.
func TraceStreamlines(settings TraceSettings, trajs [][]*Trajectory) {
	for f, trs := range trajs {
		var wg sync.WaitGroup
		wg.Add(settings.Workers)
		for w := range settings.Workers {
			go func(wk int) {
				newTraceWorker(&settings, f).run(wk, trs)
				wg.Done()
			}(w)
		}
		wg.Wait()
		pointsCount := 0
		for _, traj := range trs {
			pointsCount += len(traj.points)
		}
		fmt.Printf(
			"completed tracing all trajectories for frame %v, total points: %v\n",
			f,
			pointsCount,
		)
	}
}

// VisualizeSettings defines the settings to visualize traced streamlines.
type VisualizeSettings struct {
	CameraOrbit graphix.CameraOrbit
	LineWidth   float64
	// Fading factor for max/min tangent values, intermediate tangent values will be linearly interpolated and then gamma corrected by FadingGamma.
	MinFading   float64
	MaxFading   float64
	FadingGamma float64
	// Concurrency.
	Workers int
	// Map from camera frame index to trajectory frame index.
	FrameMapper func(f int) int
	// User-provided callback functions for each generated image, together with the camera frame index.
	ImageCallbacks []func(img draw.Image, f int)
}

// VisualizeStreamlines visualizes the traced streamlines.
func VisualizeStreamlines(settings VisualizeSettings, tracedTrajs [][]*Trajectory) {
	// Calculate max and min of tangent lengths.
	maxTan, minTan := math.Inf(-1), math.Inf(1)
	for _, trs := range tracedTrajs {
		for _, traj := range trs {
			for _, pt := range traj.points {
				maxTan = math.Max(maxTan, pt.tan)
				minTan = math.Min(minTan, pt.tan)
			}
		}
	}
	// Degenerate case - all tan's are exactly the same.
	if maxTan == minTan {
		maxTan, minTan = 1, 0
	}

	for cameraFrame := range settings.CameraOrbit.Frames() {
		trajFrame := settings.FrameMapper(cameraFrame)
		// Create zraster.SpacePaths for rendering.
		var paths []*zraster.SpacePath
		id := graphix.IdentityTransform()
		for _, traj := range tracedTrajs[trajFrame] {
			if len(traj.points) == 0 {
				continue
			}
			syms := append([]*Symmetry{{Transform: id, Color: traj.Color}}, traj.syms...)
			for _, sym := range syms {
				path := &zraster.SpacePath{
					End: sym.Transform.Apply(
						graphix.BlankVec3(),
						traj.points[len(traj.points)-1].pos,
					),
					LineWidth: settings.LineWidth,
				}
				for i := 0; i < len(traj.points)-1; i++ {
					// Take the average tangent between the two endpoints, then calculate the fading factor.
					tan := (traj.points[i].tan + traj.points[i+1].tan) / 2
					fading := settings.MinFading + (settings.MaxFading-settings.MinFading)*math.Pow(
						(tan-minTan)/(maxTan-minTan),
						settings.FadingGamma,
					)

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
		img := zraster.Run(zraster.Settings{
			Camera:  settings.CameraOrbit.GetCamera(cameraFrame),
			Paths:   paths,
			Workers: settings.Workers,
		})
		for _, cb := range settings.ImageCallbacks {
			cb(img, cameraFrame)
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
		_ = file.Close()
		fmt.Printf("generated %v\n", fn)
	}
}
