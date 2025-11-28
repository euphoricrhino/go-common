package visualizer

import (
	"math"

	"github.com/euphoricrhino/go-common/graphix"
)

// See https://www.numerical.recipes/webnotes/nr3web20.pdf.
const (
	b1  = 5.42937341165687622380535766363e-2
	b6  = 4.45031289275240888144113950566e0
	b7  = 1.89151789931450038304281599044e0
	b8  = -5.8012039600105847814672114227e0
	b9  = 3.1116436695781989440891606237e-1
	b10 = -1.52160949662516078556178806805e-1
	b11 = 2.01365400804030348374776537501e-1
	b12 = 4.47106157277725905176885569043e-2

	bhh1 = 0.244094488188976377952755905512e+00
	bhh2 = 0.733846688281611857341361741547e+00
	bhh3 = 0.220588235294117647058823529412e-01

	er1  = 0.1312004499419488073250102996e-01
	er6  = -0.1225156446376204440720569753e+01
	er7  = -0.4957589496572501915214079952e+00
	er8  = 0.1664377182454986536961530415e+01
	er9  = -0.3503288487499736816886487290e+00
	er10 = 0.3341791187130174790297318841e+00
	er11 = 0.8192320648511571246570742613e-01
	er12 = -0.2235530786388629525884427845e-01

	a21 = 5.26001519587677318785587544488e-2
	a31 = 1.97250569845378994544595329183e-2
	a32 = 5.91751709536136983633785987549e-2
	a41 = 2.95875854768068491816892993775e-2
	a43 = 8.87627564304205475450678981324e-2
	a51 = 2.41365134159266685502369798665e-1
	a53 = -8.84549479328286085344864962717e-1
	a54 = 9.24834003261792003115737966543e-1
	a61 = 3.7037037037037037037037037037e-2
	a64 = 1.70828608729473871279604482173e-1
	a65 = 1.25467687566822425016691814123e-1
	a71 = 3.7109375e-2
	a74 = 1.70252211019544039314978060272e-1
	a75 = 6.02165389804559606850219397283e-2
	a76 = -1.7578125e-2

	a81  = 3.70920001185047927108779319836e-2
	a84  = 1.70383925712239993810214054705e-1
	a85  = 1.07262030446373284651809199168e-1
	a86  = -1.53194377486244017527936158236e-2
	a87  = 8.27378916381402288758473766002e-3
	a91  = 6.24110958716075717114429577812e-1
	a94  = -3.36089262944694129406857109825e0
	a95  = -8.68219346841726006818189891453e-1
	a96  = 2.75920996994467083049415600797e1
	a97  = 2.01540675504778934086186788979e1
	a98  = -4.34898841810699588477366255144e1
	a101 = 4.77662536438264365890433908527e-1
	a104 = -2.48811461997166764192642586468e0
	a105 = -5.90290826836842996371446475743e-1
	a106 = 2.12300514481811942347288949897e1
	a107 = 1.52792336328824235832596922938e1
	a108 = -3.32882109689848629194453265587e1
	a109 = -2.03312017085086261358222928593e-2

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

	a141  = 5.61675022830479523392909219681e-2
	a147  = 2.53500210216624811088794765333e-1
	a148  = -2.46239037470802489917441475441e-1
	a149  = -1.24191423263816360469010140626e-1
	a1410 = 1.5329179827876569731206322685e-1
	a1411 = 8.20105229563468988491666602057e-3
	a1412 = 7.56789766054569976138603589584e-3
	a1413 = -8.298e-3

	a151  = 3.18346481635021405060768473261e-2
	a156  = 2.83009096723667755288322961402e-2
	a157  = 5.35419883074385676223797384372e-2
	a158  = -5.49237485713909884646569340306e-2
	a1511 = -1.08347328697249322858509316994e-4
	a1512 = 3.82571090835658412954920192323e-4
	a1513 = -3.40465008687404560802977114492e-4
	a1514 = 1.41312443674632500278074618366e-1
	a161  = -4.28896301583791923408573538692e-1
	a166  = -4.69762141536116384314449447206e0
	a167  = 7.68342119606259904184240953878e0
	a168  = 4.06898981839711007970213554331e0
	a169  = 3.56727187455281109270669543021e-1
	a1613 = -1.39902416515901462129418009734e-3
	a1614 = 2.9475147891527723389556272149e0
	a1615 = -9.15095847217987001081870187138e0

	d41  = -0.84289382761090128651353491142e+01
	d46  = 0.56671495351937776962531783590e+00
	d47  = -0.30689499459498916912797304727e+01
	d48  = 0.23846676565120698287728149680e+01
	d49  = 0.21170345824450282767155149946e+01
	d410 = -0.87139158377797299206789907490e+00
	d411 = 0.22404374302607882758541771650e+01
	d412 = 0.63157877876946881815570249290e+00
	d413 = -0.88990336451333310820698117400e-01
	d414 = 0.18148505520854727256656404962e+02
	d415 = -0.91946323924783554000451984436e+01
	d416 = -0.44360363875948939664310572000e+01

	d51  = 0.10427508642579134603413151009e+02
	d56  = 0.24228349177525818288430175319e+03
	d57  = 0.16520045171727028198505394887e+03
	d58  = -0.37454675472269020279518312152e+03
	d59  = -0.22113666853125306036270938578e+02
	d510 = 0.77334326684722638389603898808e+01
	d511 = -0.30674084731089398182061213626e+02
	d512 = -0.93321305264302278729567221706e+01
	d513 = 0.15697238121770843886131091075e+02
	d514 = -0.31139403219565177677282850411e+02
	d515 = -0.93529243588444783865713862664e+01
	d516 = 0.35816841486394083752465898540e+02

	d61  = 0.19985053242002433820987653617e+02
	d66  = -0.38703730874935176555105901742e+03
	d67  = -0.18917813819516756882830838328e+03
	d68  = 0.52780815920542364900561016686e+03
	d69  = -0.11573902539959630126141871134e+02
	d610 = 0.68812326946963000169666922661e+01
	d611 = -0.10006050966910838403183860980e+01
	d612 = 0.77771377980534432092869265740e+00
	d613 = -0.27782057523535084065932004339e+01
	d614 = -0.60196695231264120758267380846e+02
	d615 = 0.84320405506677161018159903784e+02
	d616 = 0.11992291136182789328035130030e+02

	d71  = -0.25693933462703749003312586129e+02
	d76  = -0.15418974869023643374053993627e+03
	d77  = -0.23152937917604549567536039109e+03
	d78  = 0.35763911791061412378285349910e+03
	d79  = 0.93405324183624310003907691704e+02
	d710 = -0.37458323136451633156875139351e+02
	d711 = 0.10409964950896230045147246184e+03
	d712 = 0.29840293426660503123344363579e+02
	d713 = -0.43533456590011143754432175058e+02
	d714 = 0.96324553959188282948394950600e+02
	d715 = -0.39177261675615439165231486172e+02
	d716 = -0.14972683625798562581422125276e+03
)

// Represents a stateful worker goroutine to work on a subset of trajectories using Runge-Kutta Dopr853 variant.
type traceWorker struct {
	settings   *TraceSettings
	f          int
	invEpsilon float64

	// Stateful variables.
	x1, x2, xint, xerr, xerr2, tan1, tan2, tanint graphix.Vec3
	rejected                                      bool
	// Current step and step for next iteration.
	h, hnext float64
	// General purpose temp Vec3s.
	tmp1, tmp2 graphix.Vec3
	// For eval.
	k2, k3, k4, k5, k6, k7, k8, k9, k10, k11, k12 graphix.Vec3
	// For prepareDense.
	rcont1, rcont2, rcont3, rcont4, rcont5, rcont6, rcont7, rcont8, k14, k15, k16 graphix.Vec3
}

func newTraceWorker(settings *TraceSettings, f int) *traceWorker {
	return &traceWorker{settings: settings, f: f}
}

func (tw *traceWorker) eval(xout, x, tan *graphix.Vec3) {
	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a21*tw.h))
	tw.settings.TangentAt(&tw.k2, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a31*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k2, a32*tw.h))
	tw.settings.TangentAt(&tw.k3, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a41*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k3, a43*tw.h))
	tw.settings.TangentAt(&tw.k4, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a51*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k3, a53*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a54*tw.h))
	tw.settings.TangentAt(&tw.k5, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a61*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a64*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k5, a65*tw.h))
	tw.settings.TangentAt(&tw.k6, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a71*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a74*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k5, a75*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a76*tw.h))
	tw.settings.TangentAt(&tw.k7, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a81*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a84*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k5, a85*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a86*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a87*tw.h))
	tw.settings.TangentAt(&tw.k8, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a91*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a94*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k5, a95*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a96*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a97*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k8, a98*tw.h))
	tw.settings.TangentAt(&tw.k9, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a101*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a104*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k5, a105*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a106*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a107*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k8, a108*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k9, a109*tw.h))
	tw.settings.TangentAt(&tw.k10, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a111*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a114*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k5, a115*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a116*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a117*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k8, a118*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k9, a119*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k10, a1110*tw.h))
	tw.settings.TangentAt(&tw.k11, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a121*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k4, a124*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k5, a125*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a126*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a127*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k8, a128*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k9, a129*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k10, a1210*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k11, a1211*tw.h))
	tw.settings.TangentAt(&tw.k12, &tw.tmp1, tw.f)

	tw.xerr.Scale(tan, b1)
	tw.xerr.Add(&tw.xerr, tw.tmp1.Scale(&tw.k6, b6))
	tw.xerr.Add(&tw.xerr, tw.tmp1.Scale(&tw.k7, b7))
	tw.xerr.Add(&tw.xerr, tw.tmp1.Scale(&tw.k8, b8))
	tw.xerr.Add(&tw.xerr, tw.tmp1.Scale(&tw.k9, b9))
	tw.xerr.Add(&tw.xerr, tw.tmp1.Scale(&tw.k10, b10))
	tw.xerr.Add(&tw.xerr, tw.tmp1.Scale(&tw.k11, b11))
	tw.xerr.Add(&tw.xerr, tw.tmp1.Scale(&tw.k12, b12))

	xout.Add(x, tw.tmp1.Scale(&tw.xerr, tw.h))

	tw.xerr.Sub(&tw.xerr, tw.tmp1.Scale(tan, bhh1))
	tw.xerr.Sub(&tw.xerr, tw.tmp1.Scale(&tw.k9, bhh2))
	tw.xerr.Sub(&tw.xerr, tw.tmp1.Scale(&tw.k12, bhh3))

	tw.xerr2.Scale(tan, er1)
	tw.xerr2.Add(&tw.xerr2, tw.tmp1.Scale(&tw.k6, er6))
	tw.xerr2.Add(&tw.xerr2, tw.tmp1.Scale(&tw.k7, er7))
	tw.xerr2.Add(&tw.xerr2, tw.tmp1.Scale(&tw.k8, er8))
	tw.xerr2.Add(&tw.xerr2, tw.tmp1.Scale(&tw.k9, er9))
	tw.xerr2.Add(&tw.xerr2, tw.tmp1.Scale(&tw.k10, er10))
	tw.xerr2.Add(&tw.xerr2, tw.tmp1.Scale(&tw.k11, er11))
	tw.xerr2.Add(&tw.xerr2, tw.tmp1.Scale(&tw.k12, er12))
}

func (tw *traceWorker) evalError() float64 {
	tw.tmp1.Scale(&tw.xerr, tw.invEpsilon)
	tw.tmp2.Scale(&tw.xerr2, tw.invEpsilon)
	err2 := tw.tmp1.Dot(&tw.tmp1)
	err := tw.tmp2.Dot(&tw.tmp2)
	den := err + .01*err2
	if den <= 0.0 {
		den = 1.0
	}
	return math.Abs(tw.h) * err * math.Sqrt(1/(3*den))
}

func (tw *traceWorker) success(err float64) bool {
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
		if tw.rejected {
			tw.hnext = tw.h * min(scale, 1)
		} else {
			tw.hnext = tw.h * scale
		}
		tw.rejected = false
		return true
	} else {
		scale = max(safe*math.Pow(err, alpha), minscale)
		tw.h *= scale
		tw.rejected = true
		return false
	}
}

func (tw *traceWorker) prepareDense(x, xout, tan, tanout *graphix.Vec3) {
	tw.rcont1.Copy(x)

	tw.rcont2.Sub(xout, x)

	tw.rcont3.Sub(tw.tmp1.Scale(tan, tw.h), &tw.rcont2)

	tw.rcont4.Sub(&tw.rcont2, tw.tmp1.Scale(tanout, tw.h))
	tw.rcont4.Sub(&tw.rcont4, &tw.rcont3)

	tw.rcont5.Scale(tan, d41)
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k6, d46))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k7, d47))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k8, d48))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k9, d49))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k10, d410))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k11, d411))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k12, d412))

	tw.rcont6.Scale(tan, d51)
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k6, d56))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k7, d57))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k8, d58))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k9, d59))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k10, d510))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k11, d511))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k12, d512))

	tw.rcont7.Scale(tan, d61)
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k6, d66))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k7, d67))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k8, d68))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k9, d69))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k10, d610))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k11, d611))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k12, d612))

	tw.rcont8.Scale(tan, d71)
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k6, d76))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k7, d77))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k8, d78))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k9, d79))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k10, d710))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k11, d711))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k12, d712))

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a141*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a147*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k8, a148*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k9, a149*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k10, a1410*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k11, a1411*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k12, a1412*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(tanout, a1413*tw.h))
	tw.settings.TangentAt(&tw.k14, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a151*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a156*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a157*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k8, a158*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k11, a1511*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k12, a1512*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(tanout, a1513*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k14, a1514*tw.h))
	tw.settings.TangentAt(&tw.k15, &tw.tmp1, tw.f)

	tw.tmp1.Add(x, tw.tmp2.Scale(tan, a161*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k6, a166*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k7, a167*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k8, a168*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k9, a169*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(tanout, a1613*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k14, a1614*tw.h))
	tw.tmp1.Add(&tw.tmp1, tw.tmp2.Scale(&tw.k15, a1615*tw.h))
	tw.settings.TangentAt(&tw.k16, &tw.tmp1, tw.f)

	tw.rcont5.Scale(&tw.rcont5, tw.h)
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(tanout, d413*tw.h))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k14, d414*tw.h))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k15, d415*tw.h))
	tw.rcont5.Add(&tw.rcont5, tw.tmp1.Scale(&tw.k16, d416*tw.h))

	tw.rcont6.Scale(&tw.rcont6, tw.h)
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(tanout, d513*tw.h))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k14, d514*tw.h))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k15, d515*tw.h))
	tw.rcont6.Add(&tw.rcont6, tw.tmp1.Scale(&tw.k16, d516*tw.h))

	tw.rcont7.Scale(&tw.rcont7, tw.h)
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(tanout, d613*tw.h))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k14, d614*tw.h))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k15, d615*tw.h))
	tw.rcont7.Add(&tw.rcont7, tw.tmp1.Scale(&tw.k16, d616*tw.h))

	tw.rcont8.Scale(&tw.rcont8, tw.h)
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(tanout, d713*tw.h))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k14, d714*tw.h))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k15, d715*tw.h))
	tw.rcont8.Add(&tw.rcont8, tw.tmp1.Scale(&tw.k16, d716*tw.h))
}

func (tw *traceWorker) denseOut(t float64) {
	t1 := 1 - t
	tw.xint.Scale(&tw.rcont8, t)
	tw.xint.Add(&tw.xint, &tw.rcont7)
	tw.xint.Scale(&tw.xint, t1)
	tw.xint.Add(&tw.xint, &tw.rcont6)
	tw.xint.Scale(&tw.xint, t)
	tw.xint.Add(&tw.xint, &tw.rcont5)
	tw.xint.Scale(&tw.xint, t1)
	tw.xint.Add(&tw.xint, &tw.rcont4)
	tw.xint.Scale(&tw.xint, t)
	tw.xint.Add(&tw.xint, &tw.rcont3)
	tw.xint.Scale(&tw.xint, t1)
	tw.xint.Add(&tw.xint, &tw.rcont2)
	tw.xint.Scale(&tw.xint, t)
	tw.xint.Add(&tw.xint, &tw.rcont1)
}

func (tw *traceWorker) run(w int, trajs []*Trajectory) {
	for i := range trajs {
		if i%tw.settings.Workers != w {
			continue
		}
		traj := trajs[i]

		// Reset states for the new trajectory.
		tw.x1.Copy(traj.Start)
		// Pointers of x/xout for this iteration.
		x, xout := &tw.x1, &tw.x2
		tan, tanout := &tw.tan1, &tw.tan2
		tw.invEpsilon = 1 / traj.Epsilon
		tw.h, tw.hnext = traj.InitStep, 0.0
		tw.rejected = false
		tw.settings.TangentAt(tan, x, tw.f)
		var tail *renderPoint
		for !traj.AtEnd(x, tan, tw.f) {
			if tail == nil || tw.tmp1.Sub(x, tail.pos).Norm() >= tw.settings.MinDist {
				tail = &renderPoint{tan: tan.Norm(), pos: graphix.NewCopyVec3(x)}
				traj.points = append(traj.points, tail)
			}
			// Run the adaptive Runge-Kutta (Dopr853) tracing.
			for {
				tw.eval(xout, x, tan)
				err := tw.evalError()
				if tw.success(err) {
					break
				}
			}

			if tw.tmp1.Sub(xout, tail.pos).Norm() > tw.settings.MaxDist {
				// xout is too far from the tail render point.
				// Prepare the dense interpolation coefficients.
				tw.prepareDense(x, xout, tan, tanout)

				binaryInterpolate := func() bool {
					// Precondition before calling binaryInterpolate:
					//
					// tail      x  tail+min           tail+max        xout
					//  |--------|-----|------------------|-------------|
					//
					// where x and tail may coincide.
					left, right, t := 0.0, 1.0, .5
					for {
						// Do a binary search to find xint such that MinDist <= |xint-tail| <= MaxDist.
						//
						// tail      x  tail+min  xint     tail+max        xout
						//  |--------|-----|-------⭐---------|-------------|
						//
						// Note xint must exist if the path is continuous and the precondition is satisfied.
						for {
							tw.denseOut(t)
							distFromTail := tw.tmp1.Sub(&tw.xint, tail.pos).Norm()
							if distFromTail < tw.settings.MinDist {
								left, t = t, (t+right)/2
								continue
							}
							if distFromTail > tw.settings.MaxDist {
								t, right = (t+left)/2, t
								continue
							}
							break
						}
						// At this point, we have found xint such that MinDist <= |xint-tail| <= MaxDist.
						// Corner case: is the tracing ending at xint?
						tw.settings.TangentAt(&tw.tanint, &tw.xint, tw.f)
						if traj.AtEnd(&tw.xint, &tw.tanint, tw.f) {
							return false
						}
						// Make xint the new tail and add it to the render points.
						tail = &renderPoint{
							tan: tw.tanint.Norm(),
							pos: graphix.NewCopyVec3(&tw.xint),
						}
						traj.points = append(traj.points, tail)
						// If |xout-tail| > MaxDist, repeat the above to insert more interpolated points closer to xout.
						if tw.tmp1.Sub(xout, &tw.xint).Norm() <= tw.settings.MinDist {
							return true
						}
						left, right, t = t, 1, (1+t)/2
					}
				}

				if !binaryInterpolate() {
					break
				}
			}

			tw.settings.TangentAt(tanout, xout, tw.f)
			// Swap the x/xout pointer for next iteration.
			x, xout = xout, x
			tan, tanout = tanout, tan
			tw.h = tw.hnext
		}
	}
}
