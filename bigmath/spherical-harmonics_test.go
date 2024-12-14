package bigmath

import (
	"math"
	"math/big"
	"testing"
)

func TestSphericalHarmonics(t *testing.T) {
	prec := uint(100)
	errBnd := 1e-15
	thetaf64, phif64 := 0.5, 0.8
	ctf64, stf64 := math.Cos(thetaf64), math.Sin(thetaf64)
	theta := NewFloat(thetaf64, prec)
	phi := NewFloat(phif64, prec)
	sh := NewSphericalHarmonics(
		5,
		[]int{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5},
	)

	verify := func(m, n, d int, re, im *big.Float, c ...float64) {
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
		mphi := float64(m) * phif64
		assertFloatEqualF64(t, r*math.Cos(mphi), re, errBnd)
		assertFloatEqualF64(t, r*math.Sin(mphi), im, errBnd)
	}

	// (l, m) = (0, 0)
	re, im := sh.Get(0, 0, theta, phi)
	verify(0, 1, 4, re, im, 1)

	// (l, m) = (1, 0)
	re, im = sh.Get(1, 0, theta, phi)
	verify(0, 3, 4, re, im, 0, 1)

	// (l, m) = (2, 0)
	re, im = sh.Get(2, 0, theta, phi)
	verify(0, 5, 16, re, im, -1, 0, 3)

	// (l, m) = (3, 0)
	re, im = sh.Get(3, 0, theta, phi)
	verify(0, 7, 16, re, im, 0, -3, 0, 5)

	// (l, m) = (4, 0)
	re, im = sh.Get(4, 0, theta, phi)
	verify(0, 9, 256, re, im, 3, 0, -30, 0, 35)

	// (l, m) = (5, 0)
	re, im = sh.Get(5, 0, theta, phi)
	verify(0, 11, 256, re, im, 0, 15, 0, -70, 0, 63)

	// (l, m) = (1, 1)
	re, im = sh.Get(1, 1, theta, phi)
	verify(1, -3, 8, re, im, 1)

	// (l, m) = (1, -1)
	re, im = sh.Get(1, -1, theta, phi)
	verify(-1, 3, 8, re, im, 1)

	// (l, m) = (2, 1)
	re, im = sh.Get(2, 1, theta, phi)
	verify(1, -15, 8, re, im, 0, 1)

	// (l, m) = (2, -1)
	re, im = sh.Get(2, -1, theta, phi)
	verify(-1, 15, 8, re, im, 0, 1)

	// (l, m) = (3, 1)
	re, im = sh.Get(3, 1, theta, phi)
	verify(1, -21, 64, re, im, -1, 0, 5)

	// (l, m) = (3, -1)
	re, im = sh.Get(3, -1, theta, phi)
	verify(-1, 21, 64, re, im, -1, 0, 5)

	// (l, m) = (4, 1)
	re, im = sh.Get(4, 1, theta, phi)
	verify(1, -45, 64, re, im, 0, -3, 0, 7)

	// (l, m) = (4, -1)
	re, im = sh.Get(4, -1, theta, phi)
	verify(-1, 45, 64, re, im, 0, -3, 0, 7)

	// (l, m) = (5, 1)
	re, im = sh.Get(5, 1, theta, phi)
	verify(1, -165, 512, re, im, 1, 0, -14, 0, 21)

	// (l, m) = (5, -1)
	re, im = sh.Get(5, -1, theta, phi)
	verify(-1, 165, 512, re, im, 1, 0, -14, 0, 21)

	// (l, m) = (2, 2)
	re, im = sh.Get(2, 2, theta, phi)
	verify(2, 15, 32, re, im, 1)

	// (l, m) = (2, -2)
	re, im = sh.Get(2, -2, theta, phi)
	verify(-2, 15, 32, re, im, 1)

	// (l, m) = (3, 2)
	re, im = sh.Get(3, 2, theta, phi)
	verify(2, 105, 32, re, im, 0, 1)

	// (l, m) = (3, -2)
	re, im = sh.Get(3, -2, theta, phi)
	verify(-2, 105, 32, re, im, 0, 1)

	// (l, m) = (4, 2)
	re, im = sh.Get(4, 2, theta, phi)
	verify(2, 45, 128, re, im, -1, 0, 7)

	// (l, m) = (4, -2)
	re, im = sh.Get(4, -2, theta, phi)
	verify(-2, 45, 128, re, im, -1, 0, 7)

	// (l, m) = (5, 2)
	re, im = sh.Get(5, 2, theta, phi)
	verify(2, 1155, 128, re, im, 0, -1, 0, 3)

	// (l, m) = (5, -2)
	re, im = sh.Get(5, -2, theta, phi)
	verify(-2, 1155, 128, re, im, 0, -1, 0, 3)

	// (l, m) = (3, 3)
	re, im = sh.Get(3, 3, theta, phi)
	verify(3, -35, 64, re, im, 1)

	// (l, m) = (3, -3)
	re, im = sh.Get(3, -3, theta, phi)
	verify(-3, 35, 64, re, im, 1)

	// (l, m) = (4, 3)
	re, im = sh.Get(4, 3, theta, phi)
	verify(3, -315, 64, re, im, 0, 1)

	// (l, m) = (4, -3)
	re, im = sh.Get(4, -3, theta, phi)
	verify(-3, 315, 64, re, im, 0, 1)

	// (l, m) = (5, 3)
	re, im = sh.Get(5, 3, theta, phi)
	verify(3, -385, 1024, re, im, -1, 0, 9)

	// (l, m) = (5, -3)
	re, im = sh.Get(5, -3, theta, phi)
	verify(-3, 385, 1024, re, im, -1, 0, 9)

	// (l, m) = (4, 4)
	re, im = sh.Get(4, 4, theta, phi)
	verify(4, 315, 512, re, im, 1)

	// (l, m) = (4, -4)
	re, im = sh.Get(4, -4, theta, phi)
	verify(-4, 315, 512, re, im, 1)

	// (l, m) = (5, 4)
	re, im = sh.Get(5, 4, theta, phi)
	verify(4, 3465, 512, re, im, 0, 1)

	// (l, m) = (5, -4)
	re, im = sh.Get(5, -4, theta, phi)
	verify(-4, 3465, 512, re, im, 0, 1)

	// (l, m) = (5, 5)
	re, im = sh.Get(5, 5, theta, phi)
	verify(5, -693, 1024, re, im, 1)

	// (l, m) = (5, -5)
	re, im = sh.Get(5, -5, theta, phi)
	verify(-5, 693, 1024, re, im, 1)
}
