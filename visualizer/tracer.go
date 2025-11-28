package visualizer

import (
	"fmt"
	"math"
	"path/filepath"
	"sync"

	"github.com/euphoricrhino/go-common/graphix"
)

// TraceSettings defines the settings for tracing streamlines.
type TraceSettings struct {
	// Adjacent points for rendering will have a distance between MinDist and MaxDist.
	MinDist float64
	MaxDist float64

	// User-provided tangent function at position x and frame f (i.e., discrete time t). Result should be written to tan.
	TangentAt func(tan, x *graphix.Vec3, f int)

	// Concurrency.
	Workers int

	// Directory for swap files.
	SwapDir string
}

// TraceStreamlines runs the streamline tracing given the settings and trajectory frames.
// tfs[f] contains all the trajectories for discrete time/frame f.
func TraceStreamlines(settings TraceSettings, tfs []*TrajectoryFrame) {
	for f, tf := range tfs {
		var wg sync.WaitGroup
		wg.Add(settings.Workers)
		for w := range settings.Workers {
			go func(wk int) {
				newTraceWorker(&settings, f).run(wk, tf.Trajectories)
				wg.Done()
			}(w)
		}
		wg.Wait()
		pointsCount := 0
		tf.stats = &trajectoryStats{
			minTan: math.Inf(1),
			maxTan: math.Inf(-1),
		}
		for _, traj := range tf.Trajectories {
			pointsCount += len(traj.points)
			for _, pt := range traj.points {
				tf.stats.minTan = math.Min(tf.stats.minTan, pt.tan)
				tf.stats.maxTan = math.Max(tf.stats.maxTan, pt.tan)
			}
		}
		suffix := ""
		// Save to swap file if needed.
		if settings.SwapDir != "" {
			swapFile := filepath.Join(settings.SwapDir, fmt.Sprintf("traj-frame-%04v.swap", f))
			tf.save(swapFile, settings.Workers)
			suffix = fmt.Sprintf(", saved to swap file %v", swapFile)
		}
		fmt.Printf(
			"completed tracing all trajectories for frame %v, total points: %v%v\n",
			f,
			pointsCount,
			suffix,
		)
	}
}
