package visualizer

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"sync"
	"unsafe"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/graphix/zraster"
)

// Represents a sampled point along the traced trajectory to be rendered.
type renderPoint struct {
	// The length of the tangent vector at this point.
	tan float64
	pos *graphix.Vec3
}

// Defines a symmetry transform for a traced trajectory.
type symmetry struct {
	transform graphix.Transform
	color     color.NRGBA64
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
	AtEnd func(x, tan *graphix.Vec3, f int) bool
	// The sampled points along the trajectory for rendering.
	points []*renderPoint
}

// TrajectoryVisualAttributes defines visual attributes for rendering trajectories.
type TrajectoryVisualAttributes struct {
	LineWidth float64
	syms      []*symmetry
}

// NewTrajectoryVisualAttributes creates a new TrajectoryVisualAttributes with the specified line width and color.
func NewTrajectoryVisualAttributes(
	lineWidth float64,
	color color.NRGBA64,
) *TrajectoryVisualAttributes {
	return &TrajectoryVisualAttributes{
		LineWidth: lineWidth,
		syms:      []*symmetry{{transform: graphix.IdentityTransform(), color: color}},
	}
}

// AddSymmetry adds additional symmetry with corresponding color for the symmetry-transformed trajectory.
func (tva *TrajectoryVisualAttributes) AddSymmetry(
	transform graphix.Transform,
	color color.NRGBA64,
) *TrajectoryVisualAttributes {
	tva.syms = append(tva.syms, &symmetry{
		transform: transform,
		color:     color,
	})
	return tva
}

func (vt *Trajectory) spacePaths(
	settings *VisualizeSettings,
	vta *TrajectoryVisualAttributes,
	globalMinTan, globalMaxTan float64,
) []*zraster.SpacePath {
	var paths []*zraster.SpacePath

	if len(vt.points) == 0 {
		return paths
	}

	for _, sym := range vta.syms {
		path := &zraster.SpacePath{
			End: sym.transform.Apply(
				graphix.BlankVec3(),
				vt.points[len(vt.points)-1].pos,
			),
			LineWidth: vta.LineWidth,
		}
		for i := 0; i < len(vt.points)-1; i++ {
			// Take the average tangent between the two endpoints, then calculate the fading factor.
			tan := (vt.points[i].tan + vt.points[i+1].tan) / 2
			fading := settings.MinFading + (settings.MaxFading-settings.MinFading)*math.Pow(
				(tan-globalMinTan)/(globalMaxTan-globalMinTan),
				settings.FadingGamma,
			)

			path.Segments = append(path.Segments, &zraster.SpaceVertex{
				Pos: sym.transform.Apply(graphix.BlankVec3(), vt.points[i].pos),
				Color: color.NRGBA64{
					R: sym.color.R,
					G: sym.color.G,
					B: sym.color.B,
					A: uint16(float64(sym.color.A) * fading),
				},
			})
		}
		paths = append(paths, path)
	}
	return paths
}

type trajectoryStats struct {
	minTan float64
	maxTan float64
}

// TrajectoryFrame represents all trajectories traced for a single frame.
type TrajectoryFrame struct {
	Trajectories []*Trajectory
	stats        *trajectoryStats
	SwapFile     string
}

type fileHeader struct {
	minTan    float64
	maxTan    float64
	trajCount int64
}

type rawPoint struct {
	tan float64
	pos graphix.Vec3
}

func (tf *TrajectoryFrame) save(swapFile string, workers int) {
	file, err := os.Create(swapFile)
	if err != nil {
		panic(fmt.Sprintf("failed to create trajectory frame swap file %v: %v", swapFile, err))
	}
	defer func() { _ = file.Close() }()

	trajOffsets := make([]uintptr, len(tf.Trajectories))
	totalSize := unsafe.Sizeof(fileHeader{})

	for i, traj := range tf.Trajectories {
		trajOffsets[i] = totalSize
		totalSize += unsafe.Sizeof(int64(0)) // trajectory point count
		totalSize += uintptr(len(traj.points)) * unsafe.Sizeof(rawPoint{})
	}

	buffer := make([]byte, totalSize)

	header := (*fileHeader)(unsafe.Pointer(&buffer[0]))
	header.minTan = tf.stats.minTan
	header.maxTan = tf.stats.maxTan
	header.trajCount = int64(len(tf.Trajectories))

	var wg sync.WaitGroup
	wg.Add(workers)

	// Concurrent write into different parts of the buffer without locks.
	for w := range workers {
		go func(wk int) {
			defer wg.Done()

			for i := range tf.Trajectories {
				if i%workers != wk {
					continue
				}
				traj := tf.Trajectories[i]
				offset := trajOffsets[i]

				// Write trajectory point count
				*(*int64)(unsafe.Pointer(&buffer[offset])) = int64(len(traj.points))
				offset += unsafe.Sizeof(int64(0))

				if len(traj.points) > 0 {
					pointsSlice := unsafe.Slice(
						(*rawPoint)(unsafe.Pointer(&buffer[offset])),
						len(traj.points),
					)
					for j, pt := range traj.points {
						pointsSlice[j].tan = pt.tan
						pointsSlice[j].pos = *pt.pos
					}
				}

				// Clear points to save memory.
				traj.points = nil
			}
		}(w)
	}

	wg.Wait()

	if _, err := file.Write(buffer); err != nil {
		panic(fmt.Sprintf("failed to write to trajectory frame swap file %v: %v", swapFile, err))
	}

	tf.SwapFile = swapFile
}

func (tf *TrajectoryFrame) load(statsOnly bool, workers int) {
	file, err := os.Open(tf.SwapFile)
	if err != nil {
		panic(fmt.Sprintf("failed to open trajectory frame swap file %v: %v", tf.SwapFile, err))
	}
	defer func() { _ = file.Close() }()

	// For stats only, just read the header.
	if statsOnly {
		headerSize := unsafe.Sizeof(fileHeader{})
		buffer := make([]byte, headerSize)
		if _, err := file.Read(buffer); err != nil {
			panic(
				fmt.Sprintf(
					"failed to read from trajectory frame swap file %v: %v",
					tf.SwapFile,
					err,
				),
			)
		}

		tf.loadStatsFromHeader((*fileHeader)(unsafe.Pointer(&buffer[0])))
		return
	}

	stat, err := file.Stat()
	if err != nil {
		panic(fmt.Sprintf("failed to stat trajectory frame swap file %v: %v", tf.SwapFile, err))
	}
	fileSize := stat.Size()

	buffer := make([]byte, fileSize)
	if _, err := file.Read(buffer); err != nil {
		panic(
			fmt.Sprintf(
				"failed to read from trajectory frame swap file %v: %v",
				tf.SwapFile,
				err,
			),
		)
	}
	header := (*fileHeader)(unsafe.Pointer(&buffer[0]))
	tf.loadStatsFromHeader(header)
	if header.trajCount != int64(len(tf.Trajectories)) {
		panic(fmt.Sprintf(
			"trajectory count mismatch when loading from swap file %v: expected %v, got %v",
			tf.SwapFile,
			len(tf.Trajectories),
			header.trajCount,
		))
	}

	trajOffsets := make([]uintptr, header.trajCount)
	offset := unsafe.Sizeof(fileHeader{})

	for i := range header.trajCount {
		trajOffsets[i] = offset
		pointCount := *(*int64)(unsafe.Pointer(&buffer[offset]))
		offset += unsafe.Sizeof(int64(0))
		offset += uintptr(pointCount) * unsafe.Sizeof(rawPoint{})
	}

	// Concurrent read of trajectories from buffer without locks.
	var wg sync.WaitGroup
	wg.Add(workers)

	for w := range workers {
		go func(wk int) {
			defer wg.Done()

			for i := range int(header.trajCount) {
				if i%workers != wk {
					continue
				}
				offset := trajOffsets[i]

				// Read trajectory point count
				pointCount := *(*int64)(unsafe.Pointer(&buffer[offset]))
				offset += unsafe.Sizeof(int64(0))

				traj := tf.Trajectories[i]
				traj.points = make([]*renderPoint, pointCount)

				if pointCount > 0 {
					// Get points directly from buffer
					pointsSlice := unsafe.Slice(
						(*rawPoint)(unsafe.Pointer(&buffer[offset])),
						pointCount,
					)
					for j := range pointCount {
						traj.points[j] = &renderPoint{
							tan: pointsSlice[j].tan,
							pos: graphix.NewCopyVec3(&pointsSlice[j].pos),
						}
					}
				}
			}
		}(w)
	}

	wg.Wait()
}

func (tf *TrajectoryFrame) loadStatsFromHeader(fh *fileHeader) {
	if tf.stats == nil {
		tf.stats = &trajectoryStats{}
	}
	tf.stats.minTan = fh.minTan
	tf.stats.maxTan = fh.maxTan
}

func (tf *TrajectoryFrame) ToVisual(
	tov func(int) *TrajectoryVisualAttributes,
) *VisualTrajectoryFrame {
	vta := make([]*TrajectoryVisualAttributes, len(tf.Trajectories))
	for i := range tf.Trajectories {
		vta[i] = tov(i)
	}
	return &VisualTrajectoryFrame{
		TrajectoryFrame: tf,
		vta:             vta,
	}
}

// VisualTrajectoryFrame represents a trajectory frame with visual attributes for visualization.
type VisualTrajectoryFrame struct {
	*TrajectoryFrame
	// One to one mapping to TrajectoryFrame.Trajectories.
	vta []*TrajectoryVisualAttributes
}

func (vtf *VisualTrajectoryFrame) spacePaths(
	settings *VisualizeSettings,
	globalMinTan, globalMaxTan float64,
) []*zraster.SpacePath {
	if vtf.SwapFile != "" {
		vtf.load(false, settings.Workers)
	}
	defer func() {
		if vtf.SwapFile != "" {
			// Clear points to save memory.
			for _, traj := range vtf.Trajectories {
				traj.points = nil
			}
		}
	}()
	var paths []*zraster.SpacePath
	for i, traj := range vtf.Trajectories {
		paths = append(paths, traj.spacePaths(
			settings,
			vtf.vta[i],
			globalMinTan,
			globalMaxTan,
		)...)
	}
	return paths
}
