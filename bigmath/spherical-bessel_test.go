package bigmath

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSphericalBessel(t *testing.T) {
	prec := uint(100)
	errBnd := 1e-12
	xf64 := 0.5
	x := NewFloat(xf64, prec)

	verify := func(s, c float64, expu, expv []int, u, v []*big.Int, res *big.Float) {
		assert.Equal(t, len(expu), len(u))
		assert.Equal(t, len(expv), len(v))
		uf := make([]float64, len(u))
		vf := make([]float64, len(v))
		for i := 0; i < len(expu); i++ {
			if u[i] != nil {
				assert.Equal(t, NewInt(expu[i]), u[i])
			}
			uf[i] = float64(expu[i])
		}
		for i := 0; i < len(expv); i++ {
			if v[i] != nil {
				assert.Equal(t, NewInt(expv[i]), v[i])
			}
			vf[i] = float64(expv[i])
		}
		assertFloatEqualF64(
			t,
			s*evalPolynomial(1/xf64, uf)+c*evalPolynomial(1/xf64, vf),
			res,
			errBnd,
		)
	}

	s, c := math.Sin(xf64), math.Cos(xf64)
	sx, cx := s/xf64, c/xf64

	j := NewSphericalBessel(1, 5)
	re, im := j.Get(0, x)
	assert.Nil(t, im)
	verify(sx, c, []int{1}, []int{0}, j.a[0], j.b[0], re)
	re, im = j.Get(1, x)
	assert.Nil(t, im)
	verify(sx, c, []int{0, 1}, []int{0, -1}, j.a[1], j.b[1], re)
	re, im = j.Get(2, x)
	assert.Nil(t, im)
	verify(sx, c, []int{-1, 0, 3}, []int{0, 0, -3}, j.a[2], j.b[2], re)
	re, im = j.Get(3, x)
	assert.Nil(t, im)
	verify(sx, c, []int{0, -6, 0, 15}, []int{0, 1, 0, -15}, j.a[3], j.b[3], re)
	re, im = j.Get(4, x)
	assert.Nil(t, im)
	verify(sx, c, []int{1, 0, -45, 0, 105}, []int{0, 0, 10, 0, -105}, j.a[4], j.b[4], re)
	re, im = j.Get(5, x)
	assert.Nil(t, im)
	verify(sx, c, []int{0, 15, 0, -420, 0, 945}, []int{0, -1, 0, 105, 0, -945}, j.a[5], j.b[5], re)

	n := NewSphericalBessel(2, 5)
	re, im = n.Get(0, x)
	assert.Nil(t, im)
	verify(s, cx, []int{0}, []int{-1}, n.c[0], n.d[0], re)
	re, im = n.Get(1, x)
	assert.Nil(t, im)
	verify(s, cx, []int{0, -1}, []int{0, -1}, n.c[1], n.d[1], re)
	re, im = n.Get(2, x)
	assert.Nil(t, im)
	verify(s, cx, []int{0, 0, -3}, []int{1, 0, -3}, n.c[2], n.d[2], re)
	re, im = n.Get(3, x)
	assert.Nil(t, im)
	verify(s, cx, []int{0, 1, 0, -15}, []int{0, 6, 0, -15}, n.c[3], n.d[3], re)
	re, im = n.Get(4, x)
	assert.Nil(t, im)
	verify(s, cx, []int{0, 0, 10, 0, -105}, []int{-1, 0, 45, 0, -105}, n.c[4], n.d[4], re)
	re, im = n.Get(5, x)
	assert.Nil(t, im)
	verify(s, cx, []int{0, -1, 0, 105, 0, -945}, []int{0, -15, 0, 420, 0, -945}, n.c[5], n.d[5], re)

	h1 := NewSphericalBessel(3, 5)
	h2 := NewSphericalBessel(4, 5)
	for l := 0; l <= 5; l++ {
		assert.Equal(t, j.a[l], h1.a[l])
		assert.Equal(t, j.b[l], h1.b[l])
		assert.Equal(t, j.a[l], h2.a[l])
		assert.Equal(t, j.b[l], h2.b[l])
		assert.Equal(t, n.c[l], h1.c[l])
		assert.Equal(t, n.d[l], h1.d[l])
		assert.Equal(t, n.c[l], h2.c[l])
		assert.Equal(t, n.d[l], h2.d[l])

		jre, _ := j.Get(l, x)
		nre, _ := n.Get(l, x)
		h1re, h1im := h1.Get(l, x)
		h2re, h2im := h2.Get(l, x)
		assertFloatEqual(t, jre, h1re, errBnd)
		assertFloatEqual(t, nre, h1im, errBnd)
		assertFloatEqual(t, jre, h2re, errBnd)
		assertFloatEqual(t, BlankFloat(x.Prec()).Neg(nre), h2im, errBnd)
	}
}
