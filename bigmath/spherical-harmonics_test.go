package bigmath

import (
	"math"
	"math/big"
	"testing"
)

func TestSphericalHarmonics(t *testing.T) {
	prec := uint(100)
	cpr := newBigFloatComparator(t, prec)
	thetaf64, phif64 := 0.5, 0.8
	ct := NewFloat(math.Cos(thetaf64), prec)
	st := BlankFloat(prec).Mul(ct, ct)
	st.Sub(NewFloat(1, prec), st)
	st.Sqrt(st)
	sh := NewSphericalHarmonics(
		NewFloat(thetaf64, prec),
		NewFloat(phif64, prec),
		5,
		[]int{0, 1, 2, 3, 4, 5},
	)

	pi := NewFloat(math.Pi, prec)
	eval := func(m, n, d int, c ...int) (*big.Float, *big.Float) {
		r := evalPolynomial(fromInts(c, prec), ct, prec)
		absm := m
		if m < 0 {
			absm = -m
		}
		r.Mul(r, PowerN(st, absm))
		absn := n
		if n < 0 {
			absn = -n
		}
		f := NewFloatFromRat(absn, d, prec)
		f.Quo(f, pi)
		r.Mul(r, f.Sqrt(f))
		if n < 0 {
			r.Neg(r)
		}
		cp := NewFloat(math.Cos(float64(m)*phif64), prec)
		sp := NewFloat(math.Sin(float64(m)*phif64), prec)
		return cp.Mul(cp, r), sp.Mul(sp, r)
	}

	// (l, m) = (0, 0)
	re, im := sh.Get(0, 0)
	r, i := eval(0, 1, 4, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (1, 0)
	re, im = sh.Get(1, 0)
	r, i = eval(0, 3, 4, 0, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (2, 0)
	re, im = sh.Get(2, 0)
	r, i = eval(0, 5, 16, -1, 0, 3)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (3, 0)
	re, im = sh.Get(3, 0)
	r, i = eval(0, 7, 16, 0, -3, 0, 5)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (4, 0)
	re, im = sh.Get(4, 0)
	r, i = eval(0, 9, 256, 3, 0, -30, 0, 35)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (5, 0)
	re, im = sh.Get(5, 0)
	r, i = eval(0, 11, 256, 0, 15, 0, -70, 0, 63)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (1, 1)
	re, im = sh.Get(1, 1)
	r, i = eval(1, -3, 8, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (2, 1)
	re, im = sh.Get(2, 1)
	r, i = eval(1, -15, 8, 0, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (3, 1)
	re, im = sh.Get(3, 1)
	r, i = eval(1, -21, 64, -1, 0, 5)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (4, 1)
	re, im = sh.Get(4, 1)
	r, i = eval(1, -45, 64, 0, -3, 0, 7)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (5, 1)
	re, im = sh.Get(5, 1)
	r, i = eval(1, -165, 512, 1, 0, -14, 0, 21)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (2, 2)
	re, im = sh.Get(2, 2)
	r, i = eval(2, 15, 32, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (3, 2)
	re, im = sh.Get(3, 2)
	r, i = eval(2, 105, 32, 0, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (4, 2)
	re, im = sh.Get(4, 2)
	r, i = eval(2, 45, 128, -1, 0, 7)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (5, 2)
	re, im = sh.Get(5, 2)
	r, i = eval(2, 1155, 128, 0, -1, 0, 3)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (3, 3)
	re, im = sh.Get(3, 3)
	r, i = eval(3, -35, 64, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (4, 3)
	re, im = sh.Get(4, 3)
	r, i = eval(3, -315, 64, 0, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (5, 3)
	re, im = sh.Get(5, 3)
	r, i = eval(3, -385, 1024, -1, 0, 9)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (4, 4)
	re, im = sh.Get(4, 4)
	r, i = eval(4, 315, 512, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (5, 4)
	re, im = sh.Get(5, 4)
	r, i = eval(4, 3465, 512, 0, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)

	// (l, m) = (5, 5)
	re, im = sh.Get(5, 5)
	r, i = eval(5, -693, 1024, 1)
	cpr.assertFloatEqual(r, re)
	cpr.assertFloatEqual(i, im)
}
