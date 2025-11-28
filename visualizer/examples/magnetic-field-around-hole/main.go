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
	const a = .35
	const a2 = a * a
	const h0 = 1
	const hh = h0 * 2 * a / math.Pi

	// See https://github.com/euphoricrhino/jackson-em-notes/blob/main/pdf/ch-5/pp203-circular-hole-conducting-plane-magnetic.pdf
	// for the formulas.
	tangentAt := func(tan, p *graphix.Vec3, f int) {
		x, y, z := p[0], p[1], p[2]
		rho := math.Sqrt(x*x + y*y)
		sphi, cphi := y/rho, x/rho
		neg := false
		if z < 0 {
			z = -z
			neg = true
		}
		lambda := (z*z + rho*rho - a2) / a2
		r := math.Sqrt(lambda*lambda + 4*z*z/a2)
		v1 := math.Sqrt((r - lambda) / 2)
		v2 := math.Sqrt((r + lambda) / 2)
		c1 := z / (8 * rho) / v1
		c2 := -rho / (8 * a) / (1 + 1/(v2*v2)) / (v2 * v2 * v2)
		c3 := -a / (8 * rho) / v2
		dlambdadz := 2 * z / a2
		dlambdadrho := 2 * rho / a2
		drdz := lambda/r*dlambdadz + 4*z/a2/r
		drdrho := lambda / r * dlambdadrho
		hz := v1/(2*rho) + c1*(drdz-dlambdadz)
		hz += c2 * (drdz + dlambdadz)
		hz += c3 * (drdz + dlambdadz)
		hz *= -sphi

		hrho := -z/(2*rho*rho)*v1 + c1*(drdrho-dlambdadrho)
		hrho += 1/(2*a)*math.Atan(1/v2) + c2*(drdrho+dlambdadrho)
		hrho += a/(2*rho*rho)*v2 + c2*(drdrho+dlambdadrho)
		hrho *= -sphi

		hphi := z / (2 * rho * rho) * v1
		hphi += 1 / (2 * a) * math.Atan(1/v2)
		hphi -= a / (2 * rho * rho) * v2
		hphi *= -cphi

		hx := hrho*cphi - hphi*sphi
		hy := hrho*sphi + hphi*cphi
		if !neg {
			tan[0], tan[1], tan[2] = hh*hx, h0+hh*hy, hh*hz
		} else {
			tan[0], tan[1], tan[2] = -hh*hx, -hh*hy, hh*hz
		}
	}

	atEnd := func(x, tan *graphix.Vec3, f int) bool {
		return x[1] >= .99
	}

	zs := []float64{.02, .05, .1, .2}
	zcolors := graphix.RandColors(len(zs))

	straceSettings := visualizer.TraceSettings{
		MinDist:   .001,
		MaxDist:   .005,
		TangentAt: tangentAt,
		Workers:   runtime.NumCPU(),
	}

	ang := 27 * math.Pi / 180
	visualizeSettings := visualizer.VisualizeSettings{
		CameraOrbit: graphix.NewCircularCameraOrbit(
			graphix.NewVec3(0, -math.Sin(ang), math.Cos(ang)),
			graphix.NewVec3(6, 0, 0),
			graphix.NewVec3(-1, 0, 0),
			graphix.NewVec3(0, 0, 1),
			180,
			0,
			graphix.NewPerspective(6),
			graphix.NewScreen(1280, 1280, -1.5, -1.5, 1.5, 1.5),
		),
		FadingGamma:    .3,
		MaxFading:      1,
		MinFading:      0,
		Workers:        runtime.NumCPU(),
		FrameMapper:    func(f int) int { return 0 },
		ImageCallbacks: []func(img draw.Image, f int){visualizer.SavePNG(*outDir)},
	}

	const samples = 50
	gap := 1.0 / samples
	var trajs []*visualizer.Trajectory
	for _, z := range zs {
		for j := 0; j <= samples; j++ {
			xstart := gap * float64(j)
			traj := &visualizer.Trajectory{
				Start:    graphix.NewVec3(xstart, -.99, z),
				InitStep: 0.001,
				Epsilon:  1e-14,
				AtEnd:    atEnd,
			}
			trajs = append(trajs, traj)
		}
	}

	tfs := []*visualizer.TrajectoryFrame{{Trajectories: trajs}}
	// Trace the trajectories.
	visualizer.TraceStreamlines(straceSettings, tfs)

	// Configure visual attributes and visualize the trajectories.
	vtf := tfs[0].ToVisual(
		func(idx int) *visualizer.TrajectoryVisualAttributes {
			tva := visualizer.NewTrajectoryVisualAttributes(1.5, zcolors[idx/(samples+1)])
			if idx%(samples+1) != 0 {
				tva.AddSymmetry(
					graphix.TransformFunc(func(v, u *graphix.Vec3) *graphix.Vec3 {
						v[0], v[1], v[2] = -u[0], u[1], u[2]
						return v
					}), zcolors[idx/(samples+1)],
				)
			}
			return tva
		},
	)

	visualizer.VisualizeStreamlines(visualizeSettings, []*visualizer.VisualTrajectoryFrame{vtf})
}
