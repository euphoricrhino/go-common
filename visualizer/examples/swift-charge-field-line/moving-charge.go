package main

import (
	"math"

	"github.com/euphoricrhino/go-common/graphix"
)

type movingCharge struct {
	charge float64
	motion
	epsilon float64
}

type retardation struct {
	n  graphix.Vec3
	tt float64
	R  float64
}

func (ret *retardation) eval(z func(float64, *graphix.Vec3), x *graphix.Vec3, t float64) {
	z(t, &ret.n)
	ret.n.Sub(x, &ret.n)
	ret.R = ret.n.Norm()
}

// Use binary search to find the retarded quantity.
func (mc *movingCharge) solveRetardation(x *graphix.Vec3, t float64, ret *retardation) {
	ret.eval(mc.pos, x, t)
	dt := ret.R
	tr, tl := t, t-dt
	for {
		ret.eval(mc.pos, x, tl)
		if dt > ret.R {
			break
		}
		dt *= 2
		tl = t - dt
	}

	// Now with t-tl > R(tl), t-tr < R(tr), we do a binary search between tl, tr until toleration is satisfied (or binary search gets stuck due to limited float64 precision).
	for {
		tm := (tl + tr) / 2
		ret.eval(mc.pos, x, tm)
		if math.Abs((t-tm)-ret.R) < mc.epsilon || tm == tl || tm == tr {
			ret.tt = tm
			// Normalize n before returning.
			ret.n.Scale(&ret.n, 1/ret.R)
			return
		}
		if t-tm < ret.R {
			tr = tm
		} else {
			tl = tm
		}
	}
}

// Evaluate Liénard–Wiechert electric field at position x and time t.
func (mc *movingCharge) evalElectric(x *graphix.Vec3, t float64, E *graphix.Vec3) {
	var vec1, vec2 graphix.Vec3
	var ret retardation
	mc.solveRetardation(x, t, &ret)
	// vec1=beta, vec2=betaDot
	mc.velAcc(ret.tt, &vec1, &vec2)

	// 1/gamma^2
	gammaFactor := 1.0 - vec1.Dot(&vec1)

	kappa := 1.0 - vec1.Dot(&ret.n)
	kappa3 := kappa * kappa * kappa

	aFactor := 1.0 / (kappa3 * ret.R)

	vFactor := aFactor * gammaFactor / ret.R

	// vec1=n-beta
	vec1.Sub(&ret.n, &vec1)

	// Velocity term.
	E.Scale(&vec1, vFactor)

	// vec2=(n-beta) x betaDot
	vec2.Cross(&vec1, &vec2)
	// vec2=n x [(n-beta) x betaDot]
	vec2.Cross(&ret.n, &vec2)
	vec2.Scale(&vec2, aFactor)

	// Velocity term + acceleration term.
	E.Add(E, &vec2)
	E.Scale(E, mc.charge)
}
