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
	theta := NewFloat(thetaf64, prec)
	phi := NewFloat(phif64, prec)
	sh := NewSphericalHarmonics(
		5,
		[]int{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5},
	)

	verify := func(m, n, d int, z *ModArg, c ...float64) {
		f := math.Sqrt(float64(abs(n)) / float64(d) / math.Pi)
		if n < 0 {
			f = -f
		}
		mod := evalPolynomial(ctf64, c)
		mod *= f * math.Pow(stf64, float64(abs(m)))
		assertModArgEqualF64(t, mod, float64(m)*phif64, z, errBnd)
	}

	// (l, m) = (0, 0)
	z := sh.Get(0, 0, theta, phi)
	verify(0, 1, 4, z, 1)

	// (l, m) = (1, 0)
	z = sh.Get(1, 0, theta, phi)
	verify(0, 3, 4, z, 0, 1)

	// (l, m) = (2, 0)
	z = sh.Get(2, 0, theta, phi)
	verify(0, 5, 16, z, -1, 0, 3)

	// (l, m) = (3, 0)
	z = sh.Get(3, 0, theta, phi)
	verify(0, 7, 16, z, 0, -3, 0, 5)

	// (l, m) = (4, 0)
	z = sh.Get(4, 0, theta, phi)
	verify(0, 9, 256, z, 3, 0, -30, 0, 35)

	// (l, m) = (5, 0)
	z = sh.Get(5, 0, theta, phi)
	verify(0, 11, 256, z, 0, 15, 0, -70, 0, 63)

	// (l, m) = (1, 1)
	z = sh.Get(1, 1, theta, phi)
	verify(1, -3, 8, z, 1)

	// (l, m) = (1, -1)
	z = sh.Get(1, -1, theta, phi)
	verify(-1, 3, 8, z, 1)

	// (l, m) = (2, 1)
	z = sh.Get(2, 1, theta, phi)
	verify(1, -15, 8, z, 0, 1)

	// (l, m) = (2, -1)
	z = sh.Get(2, -1, theta, phi)
	verify(-1, 15, 8, z, 0, 1)

	// (l, m) = (3, 1)
	z = sh.Get(3, 1, theta, phi)
	verify(1, -21, 64, z, -1, 0, 5)

	// (l, m) = (3, -1)
	z = sh.Get(3, -1, theta, phi)
	verify(-1, 21, 64, z, -1, 0, 5)

	// (l, m) = (4, 1)
	z = sh.Get(4, 1, theta, phi)
	verify(1, -45, 64, z, 0, -3, 0, 7)

	// (l, m) = (4, -1)
	z = sh.Get(4, -1, theta, phi)
	verify(-1, 45, 64, z, 0, -3, 0, 7)

	// (l, m) = (5, 1)
	z = sh.Get(5, 1, theta, phi)
	verify(1, -165, 512, z, 1, 0, -14, 0, 21)

	// (l, m) = (5, -1)
	z = sh.Get(5, -1, theta, phi)
	verify(-1, 165, 512, z, 1, 0, -14, 0, 21)

	// (l, m) = (2, 2)
	z = sh.Get(2, 2, theta, phi)
	verify(2, 15, 32, z, 1)

	// (l, m) = (2, -2)
	z = sh.Get(2, -2, theta, phi)
	verify(-2, 15, 32, z, 1)

	// (l, m) = (3, 2)
	z = sh.Get(3, 2, theta, phi)
	verify(2, 105, 32, z, 0, 1)

	// (l, m) = (3, -2)
	z = sh.Get(3, -2, theta, phi)
	verify(-2, 105, 32, z, 0, 1)

	// (l, m) = (4, 2)
	z = sh.Get(4, 2, theta, phi)
	verify(2, 45, 128, z, -1, 0, 7)

	// (l, m) = (4, -2)
	z = sh.Get(4, -2, theta, phi)
	verify(-2, 45, 128, z, -1, 0, 7)

	// (l, m) = (5, 2)
	z = sh.Get(5, 2, theta, phi)
	verify(2, 1155, 128, z, 0, -1, 0, 3)

	// (l, m) = (5, -2)
	z = sh.Get(5, -2, theta, phi)
	verify(-2, 1155, 128, z, 0, -1, 0, 3)

	// (l, m) = (3, 3)
	z = sh.Get(3, 3, theta, phi)
	verify(3, -35, 64, z, 1)

	// (l, m) = (3, -3)
	z = sh.Get(3, -3, theta, phi)
	verify(-3, 35, 64, z, 1)

	// (l, m) = (4, 3)
	z = sh.Get(4, 3, theta, phi)
	verify(3, -315, 64, z, 0, 1)

	// (l, m) = (4, -3)
	z = sh.Get(4, -3, theta, phi)
	verify(-3, 315, 64, z, 0, 1)

	// (l, m) = (5, 3)
	z = sh.Get(5, 3, theta, phi)
	verify(3, -385, 1024, z, -1, 0, 9)

	// (l, m) = (5, -3)
	z = sh.Get(5, -3, theta, phi)
	verify(-3, 385, 1024, z, -1, 0, 9)

	// (l, m) = (4, 4)
	z = sh.Get(4, 4, theta, phi)
	verify(4, 315, 512, z, 1)

	// (l, m) = (4, -4)
	z = sh.Get(4, -4, theta, phi)
	verify(-4, 315, 512, z, 1)

	// (l, m) = (5, 4)
	z = sh.Get(5, 4, theta, phi)
	verify(4, 3465, 512, z, 0, 1)

	// (l, m) = (5, -4)
	z = sh.Get(5, -4, theta, phi)
	verify(-4, 3465, 512, z, 0, 1)

	// (l, m) = (5, 5)
	z = sh.Get(5, 5, theta, phi)
	verify(5, -693, 1024, z, 1)

	// (l, m) = (5, -5)
	z = sh.Get(5, -5, theta, phi)
	verify(-5, 693, 1024, z, 1)
}
