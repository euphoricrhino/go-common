package bigmath

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSphericalBessel(t *testing.T) {
	prec := uint(100)
	cpr := newBigFloatComparator(t, prec)
	xf64 := 0.5
	x := NewFloat(xf64, prec)

	invx := NewFloat(1/xf64, prec)
	eval := func(s, c float64, u, v []int) *big.Float {
		polyu, polyv := evalPolynomial(
			fromInts(u, prec),
			invx,
			prec,
		), evalPolynomial(
			fromInts(v, prec),
			invx,
			prec,
		)
		polyu.Mul(polyu, NewFloat(s, prec))
		polyv.Mul(polyv, NewFloat(c, prec))
		return polyu.Add(polyu, polyv)
	}

	s, c := math.Sin(xf64), math.Cos(xf64)
	sx, cx := s/xf64, c/xf64

	j := EvalSphericalBessel(x, 5, 1)
	re, im := j.Get(0)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(sx, c, []int{1}, []int{0}), re)
	re, im = j.Get(1)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(sx, c, []int{0, 1}, []int{0, -1}), re)
	re, im = j.Get(2)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(sx, c, []int{-1, 0, 3}, []int{0, 0, -3}), re)
	re, im = j.Get(3)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(sx, c, []int{0, -6, 0, 15}, []int{0, 1, 0, -15}), re)
	re, im = j.Get(4)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(sx, c, []int{1, 0, -45, 0, 105}, []int{0, 0, 10, 0, -105}), re)
	re, im = j.Get(5)
	assert.Nil(t, im)
	cpr.assertFloatEqual(
		eval(sx, c, []int{0, 15, 0, -420, 0, 945}, []int{0, -1, 0, 105, 0, -945}),
		re,
	)

	n := EvalSphericalBessel(x, 5, 2)
	re, im = n.Get(0)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(s, cx, []int{0}, []int{-1}), re)
	re, im = n.Get(1)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(s, cx, []int{0, -1}, []int{0, -1}), re)
	re, im = n.Get(2)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(s, cx, []int{0, 0, -3}, []int{1, 0, -3}), re)
	re, im = n.Get(3)
	assert.Nil(t, im)
	cpr.assertFloatEqual(eval(s, cx, []int{0, 1, 0, -15}, []int{0, 6, 0, -15}), re)
	re, im = n.Get(4)
	assert.Nil(t, im)
	cpr.assertFloatEqual(
		eval(s, cx, []int{0, 0, 10, 0, -105}, []int{-1, 0, 45, 0, -105}),
		re,
	)
	re, im = n.Get(5)
	assert.Nil(t, im)
	cpr.assertFloatEqual(
		eval(s, cx, []int{0, -1, 0, 105, 0, -945}, []int{0, -15, 0, 420, 0, -945}),
		re,
	)
	h1 := EvalSphericalBessel(x, 5, 3)
	h2 := EvalSphericalBessel(x, 5, 4)
	for l := 0; l < 5; l++ {
		cpr.assertFloatEqual(j.re[l], h1.re[l])
		cpr.assertFloatEqual(j.re[l], h2.re[l])
		cpr.assertFloatEqual(n.re[l], h1.im[l])
		cpr.assertFloatEqual(BlankFloat(prec).Neg(n.re[l]), h2.im[l])
	}
}
