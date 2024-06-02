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
	points []*point
}

// AddSymmetry adds a symmetry transform for a trajectory.
func (traj *Trajectory) AddSymmetry(transform graphix.Transform, color color.NRGBA64) {
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
	for w := 0; w < opts.Workers; w++ {
		go func(wk int) {
			traceWorker(&opts, trajs, wk)
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
			Camera:    opts.CameraOrbit.GetCamera(f),
			Paths:     paths,
			Workers:   opts.Workers,
			NearZClip: .1,
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

// See https://www.numerical.recipes/webnotes/nr3web20.pdf.
const (
	b1    = 5.42937341165687622380535766363e-2
	b6    = 4.45031289275240888144113950566e0
	b7    = 1.89151789931450038304281599044e0
	b8    = -5.8012039600105847814672114227e0
	b9    = 3.1116436695781989440891606237e-1
	b10   = -1.52160949662516078556178806805e-1
	b11   = 2.01365400804030348374776537501e-1
	b12   = 4.47106157277725905176885569043e-2
	bhh1  = 0.244094488188976377952755905512e+00
	bhh2  = 0.733846688281611857341361741547e+00
	bhh3  = 0.220588235294117647058823529412e-01
	er1   = 0.1312004499419488073250102996e-01
	er6   = -0.1225156446376204440720569753e+01
	er7   = -0.4957589496572501915214079952e+00
	er8   = 0.1664377182454986536961530415e+01
	er9   = -0.3503288487499736816886487290e+00
	er10  = 0.3341791187130174790297318841e+00
	er11  = 0.8192320648511571246570742613e-01
	er12  = -0.2235530786388629525884427845e-01
	a21   = 5.26001519587677318785587544488e-2
	a31   = 1.97250569845378994544595329183e-2
	a32   = 5.91751709536136983633785987549e-2
	a41   = 2.95875854768068491816892993775e-2
	a43   = 8.87627564304205475450678981324e-2
	a51   = 2.41365134159266685502369798665e-1
	a53   = -8.84549479328286085344864962717e-1
	a54   = 9.24834003261792003115737966543e-1
	a61   = 3.7037037037037037037037037037e-2
	a64   = 1.70828608729473871279604482173e-1
	a65   = 1.25467687566822425016691814123e-1
	a71   = 3.7109375e-2
	a74   = 1.70252211019544039314978060272e-1
	a75   = 6.02165389804559606850219397283e-2
	a76   = -1.7578125e-2
	a81   = 3.70920001185047927108779319836e-2
	a84   = 1.70383925712239993810214054705e-1
	a85   = 1.07262030446373284651809199168e-1
	a86   = -1.53194377486244017527936158236e-2
	a87   = 8.27378916381402288758473766002e-3
	a91   = 6.24110958716075717114429577812e-1
	a94   = -3.36089262944694129406857109825e0
	a95   = -8.68219346841726006818189891453e-1
	a96   = 2.75920996994467083049415600797e1
	a97   = 2.01540675504778934086186788979e1
	a98   = -4.34898841810699588477366255144e1
	a101  = 4.77662536438264365890433908527e-1
	a104  = -2.48811461997166764192642586468e0
	a105  = -5.90290826836842996371446475743e-1
	a106  = 2.12300514481811942347288949897e1
	a107  = 1.52792336328824235832596922938e1
	a108  = -3.32882109689848629194453265587e1
	a109  = -2.03312017085086261358222928593e-2
	a111  = -9.3714243008598732571704021658e-1
	a114  = 5.18637242884406370830023853209e0
	a115  = 1.09143734899672957818500254654e0
	a116  = -8.14978701074692612513997267357e0
	a117  = -1.85200656599969598641566180701e1
	a118  = 2.27394870993505042818970056734e1
	a119  = 2.49360555267965238987089396762e0
	a1110 = -3.0467644718982195003823669022e0
	a121  = 2.27331014751653820792359768449e0
	a124  = -1.05344954667372501984066689879e1
	a125  = -2.00087205822486249909675718444e0
	a126  = -1.79589318631187989172765950534e1
	a127  = 2.79488845294199600508499808837e1
	a128  = -2.85899827713502369474065508674e0
	a129  = -8.87285693353062954433549289258e0
	a1210 = 1.23605671757943030647266201528e1
	a1211 = 6.43392746015763530355970484046e-1
)

func traceWorker(opts *StreamLineOptions, trajs []*Trajectory, w int) {
	// Thread-local scratch area variables.
	var k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12, tmp1, tmp2 graphix.Vec3
	eval := func(xout, xerr, xerr2, x, tan *graphix.Vec3, h float64) {
		tmp1.Add(x, tmp2.Scale(tan, a21*h))
		opts.TangentAt(&k2, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a31*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k2, a32*h))
		opts.TangentAt(&k3, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a41*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k3, a43*h))
		opts.TangentAt(&k4, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a51*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k3, a53*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a54*h))
		opts.TangentAt(&k5, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a61*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a64*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k5, a65*h))
		opts.TangentAt(&k6, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a71*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a74*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k5, a75*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k6, a76*h))
		opts.TangentAt(&k7, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a81*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a84*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k5, a85*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k6, a86*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k7, a87*h))
		opts.TangentAt(&k8, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a91*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a94*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k5, a95*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k6, a96*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k7, a97*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k8, a98*h))
		opts.TangentAt(&k9, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a101*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a104*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k5, a105*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k6, a106*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k7, a107*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k8, a108*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k9, a109*h))
		opts.TangentAt(&k10, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a111*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a114*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k5, a115*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k6, a116*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k7, a117*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k8, a118*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k9, a119*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k10, a1110*h))
		opts.TangentAt(&k11, &tmp1)

		tmp1.Add(x, tmp2.Scale(tan, a121*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k4, a124*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k5, a125*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k6, a126*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k7, a127*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k8, a128*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k9, a129*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k10, a1210*h))
		tmp1.Add(&tmp1, tmp2.Scale(&k11, a1211*h))
		opts.TangentAt(&k12, &tmp1)

		xerr.Scale(tan, b1)
		xerr.Add(xerr, tmp1.Scale(&k6, b6))
		xerr.Add(xerr, tmp1.Scale(&k7, b7))
		xerr.Add(xerr, tmp1.Scale(&k8, b8))
		xerr.Add(xerr, tmp1.Scale(&k9, b9))
		xerr.Add(xerr, tmp1.Scale(&k10, b10))
		xerr.Add(xerr, tmp1.Scale(&k11, b11))
		xerr.Add(xerr, tmp1.Scale(&k12, b12))

		xout.Add(x, tmp1.Scale(xerr, h))

		xerr.Sub(xerr, tmp1.Scale(tan, bhh1))
		xerr.Sub(xerr, tmp1.Scale(&k9, bhh2))
		xerr.Sub(xerr, tmp1.Scale(&k12, bhh3))

		xerr2.Scale(tan, er1)
		xerr2.Add(xerr2, tmp1.Scale(&k6, er6))
		xerr2.Add(xerr2, tmp1.Scale(&k7, er7))
		xerr2.Add(xerr2, tmp1.Scale(&k8, er8))
		xerr2.Add(xerr2, tmp1.Scale(&k9, er9))
		xerr2.Add(xerr2, tmp1.Scale(&k10, er10))
		xerr2.Add(xerr2, tmp1.Scale(&k11, er11))
		xerr2.Add(xerr2, tmp1.Scale(&k12, er12))
	}

	evalError := func(xerr, xerr2 *graphix.Vec3, h, escale float64) float64 {
		tmp1.Scale(xerr, escale)
		tmp2.Scale(xerr2, escale)
		err2 := tmp1.Dot(&tmp1)
		err := tmp2.Dot(&tmp2)
		den := err + .01*err2
		if den <= 0.0 {
			den = 1.0
		}
		return math.Abs(h) * err * math.Sqrt(1/(3*den))
	}

	success := func(err float64, h *float64, rejected *bool) bool {
		const (
			alpha    = -1.0 / 8.0
			safe     = .9
			minscale = 1.0 / 3.0
			maxscale = 6.0
		)
		scale := 0.0
		if err <= 1.0 {
			if err == 0.0 {
				scale = maxscale
			} else {
				scale = safe * math.Pow(err, alpha)
				scale = min(max(scale, minscale), maxscale)
			}
			if *rejected {
				*h *= min(scale, 1.0)
			} else {
				*h *= scale
			}
			*rejected = false
			return true
		} else {
			scale = max(safe*math.Pow(err, alpha), minscale)
			*h *= scale
			*rejected = true
			return false
		}
	}

	var x1, x2, xerr, xerr2, tan graphix.Vec3
	for i := range trajs {
		if i%opts.Workers != w {
			continue
		}
		traj := trajs[i]
		x1.Copy(traj.Start)
		// Pointers of x/xout for this iteration.
		x, xout := &x1, &x2
		h := traj.InitStep
		rejected := false
		var last *point
		for {
			opts.TangentAt(&tan, x)
			if traj.AtEnd(x, &tan) {
				// Put as the last point of trajectory regardless of the min dist.
				traj.points = append(traj.points, &point{tan: tan.Norm(), pos: graphix.NewCopyVec3(x)})
				break
			}
			if last == nil || tmp1.Sub(x, last.pos).Norm() >= opts.MinDist {
				// Accummulate point onto trajectory if it is at least MinDist away from the previous point.
				pt := &point{tan: tan.Norm(), pos: graphix.NewCopyVec3(x)}
				traj.points = append(traj.points, pt)
				last = pt
			}

			// Run the adaptive Runge-Kutta (Dopr853) tracing.
			for {
				eval(xout, &xerr, &xerr2, x, &tan, h)
				err := evalError(&xerr, &xerr2, h, 1/traj.Epsilon)
				if success(err, &h, &rejected) {
					break
				}
			}
			// Swap the x/xout pointer for next iteration
			x, xout = xout, x
		}
	}
}
