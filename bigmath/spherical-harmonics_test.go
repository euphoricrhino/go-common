package bigmath

import (
	"math"
	"testing"
)

func TestSphericalHarmonics(t *testing.T) {
	prec := uint(100)
	errBnd := 1e-15
	thetaf64, phif64 := 0.5, 0.8
	ctf64, stf64 := math.Cos(thetaf64), math.Sin(thetaf64)
	sh := NewSphericalHarmonics(
		NewFloat(thetaf64, prec),
		NewFloat(phif64, prec),
		5,
		[]int{0, 1, 2, 3, 4, 5},
	)

	eval := func(m, n, d int, c ...float64) (float64, float64) {
		absn := n
		if n < 0 {
			absn = -n
		}
		f := math.Sqrt(float64(absn) / float64(d) / math.Pi)
		if n < 0 {
			f = -f
		}
		r := evalPolynomial(ctf64, c)
		absm := m
		if m < 0 {
			absm = -m
		}
		r *= f * math.Pow(stf64, float64(absm))
		return r * math.Cos(float64(m)*phif64), r * math.Sin(float64(m)*phif64)
	}

	// (l, m) = (0, 0)
	re, im := sh.Get(0, 0)
	r, i := eval(0, 1, 4, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (1, 0)
	re, im = sh.Get(1, 0)
	r, i = eval(0, 3, 4, 0, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (2, 0)
	re, im = sh.Get(2, 0)
	r, i = eval(0, 5, 16, -1, 0, 3)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (3, 0)
	re, im = sh.Get(3, 0)
	r, i = eval(0, 7, 16, 0, -3, 0, 5)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (4, 0)
	re, im = sh.Get(4, 0)
	r, i = eval(0, 9, 256, 3, 0, -30, 0, 35)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (5, 0)
	re, im = sh.Get(5, 0)
	r, i = eval(0, 11, 256, 0, 15, 0, -70, 0, 63)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (1, 1)
	re, im = sh.Get(1, 1)
	r, i = eval(1, -3, 8, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (2, 1)
	re, im = sh.Get(2, 1)
	r, i = eval(1, -15, 8, 0, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (3, 1)
	re, im = sh.Get(3, 1)
	r, i = eval(1, -21, 64, -1, 0, 5)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (4, 1)
	re, im = sh.Get(4, 1)
	r, i = eval(1, -45, 64, 0, -3, 0, 7)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (5, 1)
	re, im = sh.Get(5, 1)
	r, i = eval(1, -165, 512, 1, 0, -14, 0, 21)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (2, 2)
	re, im = sh.Get(2, 2)
	r, i = eval(2, 15, 32, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (3, 2)
	re, im = sh.Get(3, 2)
	r, i = eval(2, 105, 32, 0, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (4, 2)
	re, im = sh.Get(4, 2)
	r, i = eval(2, 45, 128, -1, 0, 7)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (5, 2)
	re, im = sh.Get(5, 2)
	r, i = eval(2, 1155, 128, 0, -1, 0, 3)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (3, 3)
	re, im = sh.Get(3, 3)
	r, i = eval(3, -35, 64, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (4, 3)
	re, im = sh.Get(4, 3)
	r, i = eval(3, -315, 64, 0, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (5, 3)
	re, im = sh.Get(5, 3)
	r, i = eval(3, -385, 1024, -1, 0, 9)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (4, 4)
	re, im = sh.Get(4, 4)
	r, i = eval(4, 315, 512, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (5, 4)
	re, im = sh.Get(5, 4)
	r, i = eval(4, 3465, 512, 0, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)

	// (l, m) = (5, 5)
	re, im = sh.Get(5, 5)
	r, i = eval(5, -693, 1024, 1)
	assertFloatEqualF64(t, r, re, errBnd)
	assertFloatEqualF64(t, i, im, errBnd)
}
