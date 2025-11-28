package visualizer

import (
	"fmt"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/graphix/zraster"
)

// VisualizeSettings defines the settings to visualize traced streamlines.
type VisualizeSettings struct {
	CameraOrbit graphix.CameraOrbit
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
func VisualizeStreamlines(settings VisualizeSettings, vtfs []*VisualTrajectoryFrame) {
	// Calculate max and min of tangent lengths.
	maxTan, minTan := math.Inf(-1), math.Inf(1)
	for _, vtf := range vtfs {
		if vtf.stats == nil {
			vtf.load(true, settings.Workers)
		}
		minTan = math.Min(minTan, vtf.stats.minTan)
		maxTan = math.Max(maxTan, vtf.stats.maxTan)
	}
	// Degenerate case - all tan's are exactly the same.
	if maxTan == minTan {
		maxTan, minTan = 1, 0
	}

	for j, vtf := range vtfs {
		var cameraFrames []int
		for f := range settings.CameraOrbit.Frames() {
			if settings.FrameMapper(f) == j {
				cameraFrames = append(cameraFrames, f)
			}
		}
		if len(cameraFrames) == 0 {
			// No cameraFrame needs this trajectory frame, skip.
			continue
		}
		paths := vtf.spacePaths(&settings, minTan, maxTan)
		for _, cameraFrame := range cameraFrames {
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
