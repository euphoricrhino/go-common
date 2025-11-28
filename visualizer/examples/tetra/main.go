package main

import (
	"flag"
	"image/draw"
	"math"
	"runtime"

	"github.com/euphoricrhino/go-common/graphix"
	"github.com/euphoricrhino/go-common/visualizer"
)

var outDir = flag.String("out-dir", "", "output file directory")

func main() {
	flag.Parse()

	a := 0.3
	// Location of the positive and negative charges.
	positives := []*graphix.Vec3{
		{a, a, a},
		{-a, a, -a},
	}
	negatives := []*graphix.Vec3{
		{a, -a, -a},
		{-a, -a, a},
	}

	tangentAt := func(tan, x *graphix.Vec3, f int) {
		tan[0], tan[1], tan[2] = 0, 0, 0
		tmp := graphix.BlankVec3()
		accum := func(sgn float64, charge *graphix.Vec3) {
			tmp.Sub(x, charge)
			d3 := math.Pow(tmp.Dot(tmp), 1.5)
			tan.Add(tan, tmp.Scale(tmp, sgn/d3))
		}
		for _, charge := range positives {
			accum(1.0, charge)
		}
		for _, charge := range negatives {
			accum(-1.0, charge)
		}
		tan.Scale(tan, 1.0/2000.0)
	}

	sr := 0.02
	atEnd := func(x, tan *graphix.Vec3, f int) bool {
		// Stop if the field is too weak.
		if tan.Dot(tan) < 1e-12 {
			// fmt.Printf("ended at x=%v,tan=%v, small tan\n", x, tan)
			return true
		}
		tmp := graphix.BlankVec3()
		// Stop if we are close the negative charges.
		for _, charge := range negatives {
			tmp.Sub(x, charge)
			if tmp.Dot(tmp) < sr*sr {
				// fmt.Printf("ended at x=%v,tan=%v, too close to negative\n", x, tan)
				return true
			}
		}
		return false
	}

	thetaDeg := []float64{30.0, 60.0, 90.0, 120.0, 150.0}
	phiDeg := []float64{0.0, 60.0, 120.0, 180.0, 240.0, 300.0}
	generateTraj := func(charge, localz, localx *graphix.Vec3) []*visualizer.Trajectory {
		lz := graphix.BlankVec3().Normalize(localz)
		xonz := graphix.BlankVec3().Scale(lz, localx.Dot(lz))
		lx := graphix.BlankVec3().Sub(localx, xonz)
		lx.Normalize(lx)
		ly := graphix.BlankVec3().Cross(lz, lx)

		lx.Scale(lx, sr)
		ly.Scale(ly, sr)
		lz.Scale(lz, sr)
		// North and south poles.
		ret := []*visualizer.Trajectory{
			{
				Start:    graphix.BlankVec3().Add(charge, lz),
				InitStep: .005,
				Epsilon:  1e-14,
				AtEnd:    atEnd,
			}, {
				Start:    graphix.BlankVec3().Sub(charge, lz),
				InitStep: .005,
				Epsilon:  1e-14,
				AtEnd:    atEnd,
			},
		}
		for _, theta := range thetaDeg {
			thetaRad := theta * math.Pi / 180.0
			for _, phi := range phiDeg {
				phiRad := phi * math.Pi / 180.0
				disp := graphix.BlankVec3().Scale(lx, math.Sin(thetaRad)*math.Cos(phiRad))
				disp.Add(disp, graphix.BlankVec3().Scale(ly, math.Sin(thetaRad)*math.Sin(phiRad)))
				disp.Add(disp, graphix.BlankVec3().Scale(lz, math.Cos(thetaRad)))
				ret = append(ret, &visualizer.Trajectory{
					Start:    graphix.BlankVec3().Add(charge, disp),
					InitStep: .005,
					Epsilon:  1e-14,
					AtEnd:    atEnd,
				})
			}
		}
		return ret
	}

	ang := 27 * math.Pi / 180
	traceSettings := visualizer.TraceSettings{
		MinDist:   0.01,
		MaxDist:   0.05,
		TangentAt: tangentAt,
		Workers:   runtime.NumCPU(),
	}
	visualizeSettings := visualizer.VisualizeSettings{
		CameraOrbit: graphix.NewCircularCameraOrbit(
			graphix.NewVec3(0, -math.Sin(ang), math.Cos(ang)),
			graphix.NewVec3(6, 0, 0),
			graphix.NewVec3(-1, 0, 0),
			graphix.NewVec3(0, 0, 1),
			180,
			-7*math.Pi/180,
			graphix.NewPerspective(4),
			graphix.NewScreen(1280, 1280, -2, -2, 2, 2),
		),
		MinFading:      0,
		MaxFading:      1,
		FadingGamma:    .2,
		Workers:        runtime.NumCPU(),
		FrameMapper:    func(f int) int { return 0 },
		ImageCallbacks: []func(img draw.Image, f int){visualizer.SavePNG(*outDir)},
	}

	trajs := generateTraj(positives[0], graphix.NewVec3(1, 1, 1), graphix.NewVec3(0, -1, -1))
	trajs = append(
		trajs,
		generateTraj(positives[1], graphix.NewVec3(-1, 1, -1), graphix.NewVec3(1, -1, 0))...)
	tfs := []*visualizer.TrajectoryFrame{{Trajectories: trajs}}
	visualizer.TraceStreamlines(traceSettings, tfs)

	vtf := tfs[0].ToVisual(
		func(idx int) *visualizer.TrajectoryVisualAttributes {
			return visualizer.NewTrajectoryVisualAttributes(2.5, graphix.RandColor())
		},
	)
	visualizer.VisualizeStreamlines(visualizeSettings, []*visualizer.VisualTrajectoryFrame{vtf})
}
