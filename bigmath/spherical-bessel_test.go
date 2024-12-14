package bigmath

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSphericalBessel(t *testing.T) {
	prec := uint(100)
	errBnd := 1e-11
	xf64 := 0.5
	x := NewFloat(xf64, prec)

	eval := func(s, c float64, u, v []float64) float64 {
		return s*evalPolynomial(1/xf64, u) + c*evalPolynomial(1/xf64, v)
	}

	s, c := math.Sin(xf64), math.Cos(xf64)
	sx, cx := s/xf64, c/xf64

	j := EvalSphericalBessel(x, 5, 1)
	re, im := j.Get(0)
	assert.Nil(t, im)
	assertFloatEqualF64(t, eval(sx, c, []float64{1}, []float64{0}), re, errBnd)
	re, im = j.Get(1)
	assert.Nil(t, im)
	assertFloatEqualF64(t, eval(sx, c, []float64{0, 1}, []float64{0, -1}), re, errBnd)
	re, im = j.Get(2)
	assert.Nil(t, im)
	assertFloatEqualF64(t, eval(sx, c, []float64{-1, 0, 3}, []float64{0, 0, -3}), re, errBnd)
	re, im = j.Get(3)
	assert.Nil(t, im)
	assertFloatEqualF64(
		t,
		eval(sx, c, []float64{0, -6, 0, 15}, []float64{0, 1, 0, -15}),
		re,
		errBnd,
	)
	re, im = j.Get(4)
	assert.Nil(t, im)
	assertFloatEqualF64(
		t,
		eval(sx, c, []float64{1, 0, -45, 0, 105}, []float64{0, 0, 10, 0, -105}),
		re,
		errBnd,
	)
	re, im = j.Get(5)
	assert.Nil(t, im)
	assertFloatEqualF64(
		t,
		eval(sx, c, []float64{0, 15, 0, -420, 0, 945}, []float64{0, -1, 0, 105, 0, -945}),
		re,
		errBnd,
	)

	n := EvalSphericalBessel(x, 5, 2)
	re, im = n.Get(0)
	assert.Nil(t, im)
	assertFloatEqualF64(t, eval(s, cx, []float64{0}, []float64{-1}), re, errBnd)
	re, im = n.Get(1)
	assert.Nil(t, im)
	assertFloatEqualF64(t, eval(s, cx, []float64{0, -1}, []float64{0, -1}), re, errBnd)
	re, im = n.Get(2)
	assert.Nil(t, im)
	assertFloatEqualF64(t, eval(s, cx, []float64{0, 0, -3}, []float64{1, 0, -3}), re, errBnd)
	re, im = n.Get(3)
	assert.Nil(t, im)
	assertFloatEqualF64(
		t,
		eval(s, cx, []float64{0, 1, 0, -15}, []float64{0, 6, 0, -15}),
		re,
		errBnd,
	)
	re, im = n.Get(4)
	assert.Nil(t, im)
	assertFloatEqualF64(
		t,
		eval(s, cx, []float64{0, 0, 10, 0, -105}, []float64{-1, 0, 45, 0, -105}),
		re,
		errBnd,
	)
	re, im = n.Get(5)
	assert.Nil(t, im)
	assertFloatEqualF64(
		t,
		eval(s, cx, []float64{0, -1, 0, 105, 0, -945}, []float64{0, -15, 0, 420, 0, -945}),
		re,
		errBnd,
	)
	h1 := EvalSphericalBessel(x, 5, 3)
	h2 := EvalSphericalBessel(x, 5, 4)
	for l := 0; l < 5; l++ {
		assertFloatEqual(t, j.re[l], h1.re[l], errBnd)
		assertFloatEqual(t, j.re[l], h2.re[l], errBnd)
		assertFloatEqual(t, n.re[l], h1.im[l], errBnd)
		assertFloatEqual(t, BlankFloat(prec).Neg(n.re[l]), h2.im[l], errBnd)
	}
}
