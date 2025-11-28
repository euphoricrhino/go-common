package main

import (
	"flag"
	"fmt"
	"image/draw"
	"math"
	"path/filepath"
	"runtime"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/visualizer"
)

const (
	retardationEpsilon = 1e-14
	traceEpsilon       = 1e-14
)

var (
	outDir     = flag.String("out-dir", "", "output file directory")
	width      = flag.Int("width", 480, "image width")
	height     = flag.Int("height", 480, "image height")
	maxDist    = flag.Float64("max-dist", 20, "maximum distance from origin to trace")
	lineWidth  = flag.Float64("line-width", 1.5, "line width")
	fovRange   = flag.Int("fov-range", 15, "field of view range")
	configName = flag.String("config-name", "", "configuration to use")
	gamma      = flag.Float64("gamma", 0.2, "fading gamma")
	renderOnly = flag.Bool("render-only", false, "only render the streamlines without tracing them")
)

const (
	a                  = 1.0
	br                 = 0.1
	slowOmega          = 0.45
	fastOmega          = 0.9
	slowFramesPerCycle = 120
	fastFramesPerCycle = 60
	initStep           = 0.001
)

type configuration struct {
	framesPerCycle int
	charges        []*movingCharge
	twoD           bool
}

var configs = map[string]configuration{
	"1-harmonic-2d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, slowOmega, 0, 1, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"1-harmonic-3d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, slowOmega, 0, 1, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
	"1-harmonic-2d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, fastOmega, 0, 1, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"1-harmonic-3d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, fastOmega, 0, 1, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
	"1-circular-2d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, slowOmega, 0, 2, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"1-circular-3d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, slowOmega, 0, 2, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
	"1-circular-2d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, fastOmega, 0, 2, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"1-circular-3d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, fastOmega, 0, 2, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
	"2-harmonic-2d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, slowOmega, 0, 1, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newHarmonicMotion(a, slowOmega, math.Pi/2, 0, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"2-harmonic-3d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, slowOmega, 0, 1, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newHarmonicMotion(a, slowOmega, math.Pi/2, 0, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
	"2-harmonic-2d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, fastOmega, 0, 1, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newHarmonicMotion(a, fastOmega, math.Pi/2, 0, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"2-harmonic-3d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newHarmonicMotion(a, fastOmega, 0, 1, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newHarmonicMotion(a, fastOmega, math.Pi/2, 0, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
	"2-circular-2d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, slowOmega, 0, 2, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newCircularMotion(a, slowOmega, math.Pi, 2, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"2-circular-3d-slow": {
		framesPerCycle: slowFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, slowOmega, 0, 2, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newCircularMotion(a, slowOmega, math.Pi, 2, slowFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
	"2-circular-2d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, fastOmega, 0, 2, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newCircularMotion(a, fastOmega, math.Pi, 2, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: true,
	},
	"2-circular-3d-fast": {
		framesPerCycle: fastFramesPerCycle,
		charges: []*movingCharge{
			{
				charge:  1,
				motion:  newCircularMotion(a, fastOmega, 0, 2, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
			{
				charge:  1,
				motion:  newCircularMotion(a, fastOmega, math.Pi, 2, fastFramesPerCycle),
				epsilon: retardationEpsilon,
			},
		},
		twoD: false,
	},
}

func main() {
	flag.Parse()

	config := configs[*configName]
	var tfs []*visualizer.TrajectoryFrame

	thetaDivides := 18
	if config.twoD {
		thetaDivides = 2
	}
	phiDivides := 60

	for f := range config.framesPerCycle {
		var trs []*visualizer.Trajectory
		var pos graphix.Vec3
		for _, mc := range config.charges {
			mc.pos(mc.frameToTime(f), &pos)
			trs = append(trs, generateTrajs(
				&pos,
				thetaDivides,
				phiDivides,
				config.twoD,
			)...)
		}
		tfs = append(tfs, &visualizer.TrajectoryFrame{
			Trajectories: trs,
		})
	}
	if !*renderOnly {
		traceSettings := visualizer.TraceSettings{
			MinDist: traceMinDist(&config),
			MaxDist: traceMaxDist(&config),
			TangentAt: func(tan, p *graphix.Vec3, f int) {
				tan.Clear()
				var tmp graphix.Vec3
				for _, mc := range config.charges {
					mc.evalElectric(p, mc.frameToTime(f), &tmp)
					tan.Add(tan, &tmp)
				}
			},
			Workers: runtime.NumCPU(),
			SwapDir: *outDir,
		}
		visualizer.TraceStreamlines(traceSettings, tfs)
	} else {
		// Traced data are stored in the swap files.
		for f, tf := range tfs {
			tf.SwapFile = filepath.Join(*outDir, fmt.Sprintf("traj-frame-%04v.swap", f))
		}
	}
	ang := 27 * math.Pi / 180
	orb2D := graphix.NewStationaryCamera(
		graphix.NewOrtho2DCamera(
			graphix.NewScreen(
				*width,
				*height,
				-float64(*fovRange),
				-float64(*fovRange),
				float64(*fovRange),
				float64(*fovRange),
			),
		),
		config.framesPerCycle,
	)

	offset := 75 * math.Pi / 180
	rot := graphix.NewAxisAngleRotation(graphix.NewVec3(0, -math.Sin(ang), math.Cos(ang)), offset)
	vt := graphix.NewViewTransform(
		rot.Apply(graphix.BlankVec3(), graphix.NewVec3(*maxDist, 0, 0)),
		rot.Apply(graphix.BlankVec3(), graphix.NewVec3(-1, 0, 0)),
		rot.Apply(graphix.BlankVec3(), graphix.NewVec3(0, 0, 1)),
	)
	orb3D := graphix.NewStationaryCamera(
		graphix.NewCamera(
			vt,
			graphix.NewOrthographic(),
			graphix.NewScreen(
				*width,
				*height,
				-float64(*fovRange),
				-float64(*fovRange),
				float64(*fovRange),
				float64(*fovRange),
			),
		),
		config.framesPerCycle,
	)

	visualizeSettings := visualizer.VisualizeSettings{
		FadingGamma:    *gamma,
		MaxFading:      1,
		MinFading:      0,
		Workers:        runtime.NumCPU(),
		FrameMapper:    func(f int) int { return f % config.framesPerCycle },
		ImageCallbacks: []func(img draw.Image, f int){visualizer.SavePNG(*outDir)},
	}
	if config.twoD {
		visualizeSettings.CameraOrbit = orb2D
	} else {
		visualizeSettings.CameraOrbit = orb3D
	}
	colorStride := thetaDivides + 1
	if config.twoD {
		colorStride = phiDivides
	}
	colorCnt := colorStride * len(config.charges)
	colors := graphix.RandColors(colorCnt)
	trajsPerCharge := phiDivides
	if !config.twoD {
		trajsPerCharge = 2 + (thetaDivides-1)*phiDivides
	}
	var vtfs []*visualizer.VisualTrajectoryFrame
	for _, tf := range tfs {
		vtfs = append(vtfs, tf.ToVisual(func(idx int) *visualizer.TrajectoryVisualAttributes {
			section, offset := idx/trajsPerCharge, idx%trajsPerCharge
			if !config.twoD {
				switch offset {
				case 0:
					offset = 0
				case 1:
					offset = thetaDivides
				default:
					offset = 1 + (offset-2)/phiDivides
				}
			}
			return visualizer.NewTrajectoryVisualAttributes(
				*lineWidth,
				colors[section*colorStride+offset],
			)
		}))
	}

	visualizer.VisualizeStreamlines(visualizeSettings, vtfs)
}

func atEnd(x, tan *graphix.Vec3, f int) bool {
	return x.Norm() > *maxDist
}

func generateTrajs(
	pos *graphix.Vec3,
	thetaDivides, phiDivides int,
	twoD bool,
) []*visualizer.Trajectory {
	var trs []*visualizer.Trajectory
	if !twoD {
		trs = append(trs,
			&visualizer.Trajectory{
				Start:    graphix.BlankVec3().Add(pos, graphix.NewVec3(0, 0, br)),
				InitStep: initStep,
				Epsilon:  traceEpsilon,
				AtEnd:    atEnd,
			},
			&visualizer.Trajectory{
				Start:    graphix.BlankVec3().Add(pos, graphix.NewVec3(0, 0, -br)),
				InitStep: initStep,
				Epsilon:  traceEpsilon,
				AtEnd:    atEnd,
			},
		)
	}
	for i := 1; i < thetaDivides; i++ {
		theta := math.Pi / float64(thetaDivides) * float64(i)
		rho := br * math.Sin(theta)
		z := br * math.Cos(theta)
		for j := range phiDivides {
			phi := 2 * math.Pi / float64(phiDivides) * float64(j)
			start := graphix.BlankVec3().Add(pos, graphix.NewVec3(
				rho*math.Cos(phi),
				rho*math.Sin(phi),
				z,
			))
			traj := &visualizer.Trajectory{
				Start:    start,
				InitStep: initStep,
				Epsilon:  traceEpsilon,
				AtEnd:    atEnd,
			}
			trs = append(trs, traj)
		}
	}
	return trs
}

func traceMinDist(config *configuration) float64 {
	if config.twoD {
		return 0.01
	}
	return 0.2
}

func traceMaxDist(config *configuration) float64 {
	if config.twoD {
		return 0.05
	}
	return 0.5
}
