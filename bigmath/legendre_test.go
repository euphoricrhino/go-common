package bigmath

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLegendre(t *testing.T) {
	prec := uint(100)
	errBnd := 0.0
	xf64 := 0.5
	le := NewLegendre(10)

	x := NewFloat(xf64, prec)
	verify := func(l int, c ...int) {
		assert.Equal(t, len(c)/2, len(le.coeff[l]))
		f := make([]float64, len(c)/2)
		for i, v := range le.coeff[l] {
			if v != nil {
				assert.Equal(t, NewRat(c[2*i], c[2*i+1]), v)
			}
			f[i] = float64(c[2*i]) / float64(c[2*i+1])
		}
		assertFloatEqualF64(t, evalPolynomial(xf64, f), le.Get(l, x), errBnd)
	}

	// https://en.wikipedia.org/wiki/Legendre_polynomials
	verify(0, 1, 1)
	verify(1, 0, 1, 1, 1)
	verify(2, -1, 2, 0, 1, 3, 2)
	verify(3, 0, 1, -3, 2, 0, 1, 5, 2)
	verify(4, 3, 8, 0, 1, -30, 8, 0, 1, 35, 8)
	verify(5, 0, 1, 15, 8, 0, 1, -70, 8, 0, 1, 63, 8)
	verify(6, -5, 16, 0, 1, 105, 16, 0, 1, -315, 16, 0, 1, 231, 16)
	verify(7, 0, 1, -35, 16, 0, 1, 315, 16, 0, 1, -693, 16, 0, 1, 429, 16)
	verify(8, 35, 128, 0, 1, -1260, 128, 0, 1, 6930, 128, 0, 1, -12012, 128, 0, 1, 6435, 128)
	verify(
		9,
		0, 1,
		315, 128,
		0, 1,
		-4620, 128,
		0, 1,
		18018, 128,
		0, 1,
		-25740, 128,
		0, 1,
		12155, 128,
	)
	verify(
		10,
		-63, 256,
		0, 1,
		3465, 256,
		0, 1,
		-30030, 256,
		0, 1,
		90090, 256,
		0, 1,
		-109395, 256,
		0, 1,
		46189, 256,
	)
}

func TestAssocLegendre(t *testing.T) {
	prec := uint(100)
	errBnd := 1e-13
	xf64 := 0.5
	x := NewFloat(xf64, prec)

	verify := func(al *AssocLegendre, l, m int, c ...int) {
		assert.Equal(t, len(c)/2, len(al.lm[l]))
		f := make([]float64, len(c)/2)
		for i, v := range al.lm[l] {
			if v != nil {
				assert.Equal(t, NewRat(c[2*i], c[2*i+1]), v)
			}
			f[i] = float64(c[2*i]) / float64(c[2*i+1])
		}
		assertFloatEqualF64(
			t,
			evalPolynomial(xf64, f)*math.Pow(math.Sqrt(1-xf64*xf64), math.Abs(float64(m))),
			al.Get(l, x),
			errBnd,
		)
	}

	verifyScaled := func(al *AssocLegendre, l, m int, n, d int, c ...int) {
		for i := 0; i < len(c); i += 2 {
			c[i] *= n
			c[i+1] *= d
		}
		verify(al, l, m, c...)
	}

	le := NewLegendre(4)

	// m=0.
	al0 := NewAssocLegendre(0, le)
	verify(al0, 0, 0, 1, 1)
	verify(al0, 1, 0, 0, 1, 1, 1)
	verify(al0, 2, 0, -1, 2, 0, 1, 3, 2)
	verify(al0, 3, 0, 0, 1, -3, 2, 0, 1, 5, 2)
	verify(al0, 4, 0, 3, 8, 0, 1, -30, 8, 0, 1, 35, 8)

	// m=1.
	al1 := NewAssocLegendre(1, le)
	assert.Nil(t, al1.lm[0])
	verify(al1, 1, 1, -1, 1)
	verify(al1, 2, 1, 0, 1, -3, 1)
	verify(al1, 3, 1, 3, 2, 0, 1, -15, 2)
	verify(al1, 4, 1, 0, 1, 15, 2, 0, 1, -35, 2)

	// m=-1.
	aln1 := NewAssocLegendre(-1, le)
	assert.Nil(t, aln1.lm[0])
	verifyScaled(aln1, 1, 1, -1, 2, -1, 1)
	verifyScaled(aln1, 2, 1, -1, 6, 0, 1, -3, 1)
	verifyScaled(aln1, 3, 1, -1, 12, 3, 2, 0, 1, -15, 2)
	verifyScaled(aln1, 4, 1, -1, 20, 0, 1, 15, 2, 0, 1, -35, 2)

	// m=2.
	al2 := NewAssocLegendre(2, le)
	assert.Nil(t, al2.lm[0])
	assert.Nil(t, al2.lm[1])
	verify(al2, 2, 2, 3, 1)
	verify(al2, 3, 2, 0, 1, 15, 1)
	verify(al2, 4, 2, -15, 2, 0, 1, 105, 2)

	// m=-2.
	aln2 := NewAssocLegendre(-2, le)
	assert.Nil(t, aln2.lm[0])
	assert.Nil(t, aln2.lm[1])
	verifyScaled(aln2, 2, 2, 1, 24, 3, 1)
	verifyScaled(aln2, 3, 2, 1, 120, 0, 1, 15, 1)
	verifyScaled(aln2, 4, 2, 1, 360, -15, 2, 0, 1, 105, 2)

	// m=3.
	al3 := NewAssocLegendre(3, le)
	assert.Nil(t, al3.lm[0])
	assert.Nil(t, al3.lm[1])
	assert.Nil(t, al3.lm[2])
	verify(al3, 3, 3, -15, 1)
	verify(al3, 4, 3, 0, 1, -105, 1)

	// m=-3.
	aln3 := NewAssocLegendre(-3, le)
	assert.Nil(t, aln3.lm[0])
	assert.Nil(t, aln3.lm[1])
	assert.Nil(t, aln3.lm[1])
	verifyScaled(aln3, 3, 3, -1, 720, -15, 1)
	verifyScaled(aln3, 4, 3, -1, 5040, 0, 1, -105, 1)

	// m=4.
	al4 := NewAssocLegendre(4, le)
	assert.Nil(t, al4.lm[0])
	assert.Nil(t, al4.lm[1])
	assert.Nil(t, al4.lm[2])
	assert.Nil(t, al4.lm[3])
	verify(al4, 4, 4, 105, 1)

	// m=-4.
	aln4 := NewAssocLegendre(-4, le)
	assert.Nil(t, aln4.lm[0])
	assert.Nil(t, aln4.lm[1])
	assert.Nil(t, aln4.lm[1])
	assert.Nil(t, aln4.lm[1])
	verifyScaled(aln4, 4, 4, 1, 40320, 105, 1)
}
