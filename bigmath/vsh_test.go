package bigmath

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVSH(t *testing.T) {
	prec := uint(100)
	errBnd := 1e-15
	thetaf64, phif64 := 0.5, 0.8
	ct, st := math.Cos(thetaf64), math.Sin(thetaf64)
	theta := NewFloat(thetaf64, prec)
	phi := NewFloat(phif64, prec)
	vsh := NewVSH(
		3,
		[]int{-3, -2, -1, 0, 1, 2, 3},
	)

	verify := func(n, d int, z *ModArg, poly, expArg float64) {
		f := math.Sqrt(float64(abs(n)) / float64(d) / math.Pi)
		if n < 0 {
			f = -f
		}
		mod := poly * f
		assertModArgEqualF64(t, mod, expArg, z, errBnd)
	}

	// (l, m) = (0, 0)
	vshY := vsh.GetY(0, 0, theta, phi)
	verify(1, 4, vshY[0], 1, 0)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi := vsh.GetPsi(0, 0, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(0, 1, vshPsi[1], 0, 0)
	verify(0, 1, vshPsi[2], 0, math.Pi/2)

	vshPhi := vsh.GetPhi(0, 0, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(0, 1, vshPhi[1], 0, math.Pi/2)
	verify(0, 1, vshPhi[2], 0, 0)

	// (l, m) = (1, 0)
	vshY = vsh.GetY(1, 0, theta, phi)
	verify(3, 4, vshY[0], ct, 0)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(1, 0, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(-3, 4, vshPsi[1], st, 0)
	verify(0, 1, vshPsi[2], 0, math.Pi/2)

	vshPhi = vsh.GetPhi(1, 0, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(0, 1, vshPhi[1], 0, math.Pi/2)
	verify(-3, 4, vshPhi[2], st, 0)

	// (l, m) = (2, 0)
	vshY = vsh.GetY(2, 0, theta, phi)
	verify(5, 16, vshY[0], 3*ct*ct-1, 0)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(2, 0, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(-45, 4, vshPsi[1], st*ct, 0)
	verify(0, 1, vshPsi[2], 0, math.Pi/2)

	vshPhi = vsh.GetPhi(2, 0, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(0, 1, vshPhi[1], 0, math.Pi/2)
	verify(-45, 4, vshPhi[2], st*ct, 0)

	// (l, m) = (3, 0)
	vshY = vsh.GetY(3, 0, theta, phi)
	verify(7, 16, vshY[0], -3*ct+5*ct*ct*ct, 0)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(3, 0, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(7, 16, vshPsi[1], -st*(-3+15*ct*ct), 0)
	verify(0, 1, vshPsi[2], 0, math.Pi/2)

	vshPhi = vsh.GetPhi(3, 0, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(0, 1, vshPhi[1], 0, math.Pi/2)
	verify(7, 16, vshPhi[2], -st*(-3+15*ct*ct), 0)

	// (l, m) = (1, 1)
	vshY = vsh.GetY(1, 1, theta, phi)
	verify(-3, 8, vshY[0], st, phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(1, 1, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(-3, 8, vshPsi[1], ct, phif64)
	verify(-3, 8, vshPsi[2], 1, phif64+math.Pi/2)

	vshPhi = vsh.GetPhi(1, 1, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(3, 8, vshPhi[1], 1, phif64+math.Pi/2)
	verify(-3, 8, vshPhi[2], ct, phif64)

	// (l, m) = (1, -1)
	vshY = vsh.GetY(1, -1, theta, phi)
	verify(3, 8, vshY[0], st, -phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(1, -1, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(3, 8, vshPsi[1], ct, -phif64)
	verify(3, 8, vshPsi[2], 1, -(phif64 + math.Pi/2))

	vshPhi = vsh.GetPhi(1, -1, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-3, 8, vshPhi[1], 1, -(phif64 + math.Pi/2))
	verify(3, 8, vshPhi[2], ct, -phif64)

	// (l, m) = (2, 1)
	vshY = vsh.GetY(2, 1, theta, phi)
	verify(-15, 8, vshY[0], st*ct, phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(2, 1, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(-15, 8, vshPsi[1], 2*ct*ct-1, phif64)
	verify(-15, 8, vshPsi[2], ct, phif64+math.Pi/2)

	vshPhi = vsh.GetPhi(2, 1, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(15, 8, vshPhi[1], ct, phif64+math.Pi/2)
	verify(-15, 8, vshPhi[2], 2*ct*ct-1, phif64)

	// (l, m) = (2, -1)
	vshY = vsh.GetY(2, -1, theta, phi)
	verify(15, 8, vshY[0], st*ct, -phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(2, -1, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(15, 8, vshPsi[1], 2*ct*ct-1, -phif64)
	verify(15, 8, vshPsi[2], ct, -(phif64 + math.Pi/2))

	vshPhi = vsh.GetPhi(2, -1, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-15, 8, vshPhi[1], ct, -(phif64 + math.Pi/2))
	verify(15, 8, vshPhi[2], 2*ct*ct-1, -phif64)

	// (l, m) = (3, 1)
	vshY = vsh.GetY(3, 1, theta, phi)
	verify(-21, 64, vshY[0], st*(-1+5*ct*ct), phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(3, 1, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(-21, 64, vshPsi[1], 15*ct*ct*ct-11*ct, phif64)
	verify(-21, 64, vshPsi[2], -1+5*ct*ct, phif64+math.Pi/2)

	vshPhi = vsh.GetPhi(3, 1, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(21, 64, vshPhi[1], -1+5*ct*ct, phif64+math.Pi/2)
	verify(-21, 64, vshPhi[2], 15*ct*ct*ct-11*ct, phif64)

	// (l, m) = (3, -1)
	vshY = vsh.GetY(3, -1, theta, phi)
	verify(21, 64, vshY[0], st*(-1+5*ct*ct), -phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(3, -1, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(21, 64, vshPsi[1], 15*ct*ct*ct-11*ct, -phif64)
	verify(21, 64, vshPsi[2], -1+5*ct*ct, -(phif64 + math.Pi/2))

	vshPhi = vsh.GetPhi(3, -1, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-21, 64, vshPhi[1], -1+5*ct*ct, -(phif64 + math.Pi/2))
	verify(21, 64, vshPhi[2], 15*ct*ct*ct-11*ct, -phif64)

	// (l, m) = (2, 2)
	vshY = vsh.GetY(2, 2, theta, phi)
	verify(15, 32, vshY[0], st*st, 2*phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(2, 2, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(15, 8, vshPsi[1], st*ct, 2*phif64)
	verify(15, 8, vshPsi[2], st, 2*phif64+math.Pi/2)

	vshPhi = vsh.GetPhi(2, 2, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-15, 8, vshPhi[1], st, 2*phif64+math.Pi/2)
	verify(15, 8, vshPhi[2], st*ct, 2*phif64)

	// (l, m) = (2, -2)
	vshY = vsh.GetY(2, -2, theta, phi)
	verify(15, 32, vshY[0], st*st, -2*phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(2, -2, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(15, 8, vshPsi[1], st*ct, -2*phif64)
	verify(15, 8, vshPsi[2], st, -(2*phif64 + math.Pi/2))

	vshPhi = vsh.GetPhi(2, -2, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-15, 8, vshPhi[1], st, -(2*phif64 + math.Pi/2))
	verify(15, 8, vshPhi[2], st*ct, -2*phif64)

	// (l, m) = (3, 2)
	vshY = vsh.GetY(3, 2, theta, phi)
	verify(105, 32, vshY[0], st*st*ct, 2*phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(3, 2, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(105, 32, vshPsi[1], 2*st*ct*ct-st*st*st, 2*phif64)
	verify(105, 32, vshPsi[2], 2*st*ct, 2*phif64+math.Pi/2)

	vshPhi = vsh.GetPhi(3, 2, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-105, 32, vshPhi[1], 2*st*ct, 2*phif64+math.Pi/2)
	verify(105, 32, vshPhi[2], 2*st*ct*ct-st*st*st, 2*phif64)

	// (l, m) = (3, -2)
	vshY = vsh.GetY(3, -2, theta, phi)
	verify(105, 32, vshY[0], st*st*ct, -2*phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(3, -2, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(105, 32, vshPsi[1], 2*st*ct*ct-st*st*st, -2*phif64)
	verify(105, 32, vshPsi[2], 2*st*ct, -(2*phif64 + math.Pi/2))

	vshPhi = vsh.GetPhi(3, -2, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-105, 32, vshPhi[1], 2*st*ct, -(2*phif64 + math.Pi/2))
	verify(105, 32, vshPhi[2], 2*st*ct*ct-st*st*st, -2*phif64)

	// (l, m) = (3, 3)
	vshY = vsh.GetY(3, 3, theta, phi)
	verify(-35, 64, vshY[0], st*st*st, 3*phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(3, 3, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(-35, 64, vshPsi[1], 3*st*st*ct, 3*phif64)
	verify(-35, 64, vshPsi[2], 3*st*st, 3*phif64+math.Pi/2)

	vshPhi = vsh.GetPhi(3, 3, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(35, 64, vshPhi[1], 3*st*st, 3*phif64+math.Pi/2)
	verify(-35, 64, vshPhi[2], 3*st*st*ct, 3*phif64)

	// (l, m) = (3, -3)
	vshY = vsh.GetY(3, -3, theta, phi)
	verify(35, 64, vshY[0], st*st*st, -3*phif64)
	assert.Nil(t, vshY[1])
	assert.Nil(t, vshY[2])

	vshPsi = vsh.GetPsi(3, -3, theta, phi)
	assert.Nil(t, vshPsi[0])
	verify(35, 64, vshPsi[1], 3*st*st*ct, -3*phif64)
	verify(35, 64, vshPsi[2], 3*st*st, -(3*phif64 + math.Pi/2))

	vshPhi = vsh.GetPhi(3, -3, theta, phi)
	assert.Nil(t, vshPhi[0])
	verify(-35, 64, vshPhi[1], 3*st*st, -(3*phif64 + math.Pi/2))
	verify(35, 64, vshPhi[2], 3*st*st*ct, -3*phif64)
}
