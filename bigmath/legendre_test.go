package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLegendre(t *testing.T) {
	prec := uint(100)
	cpr := newBigFloatComparator(t, prec)
	xf64 := 0.5
	x := NewFloat(xf64, prec)
	le := EvalLegendre(x, 10)

	eval := func(c ...int) *big.Float { return evalPolynomial(fromRationals(c, prec), x, prec) }

	// https://en.wikipedia.org/wiki/Legendre_polynomials
	cpr.assertFloatEqual(eval(1, 1), le.Get(0))
	cpr.assertFloatEqual(eval(0, 1, 1, 1), le.Get(1))
	cpr.assertFloatEqual(eval(-1, 2, 0, 1, 3, 2), le.Get(2))
	cpr.assertFloatEqual(eval(0, 1, -3.0, 2, 0, 1, 5, 2), le.Get(3))
	cpr.assertFloatEqual(eval(3, 8, 0, 1, -30, 8, 0, 1, 35, 8), le.Get(4))
	cpr.assertFloatEqual(eval(0, 1, 15, 8, 0, 1, -70, 8, 0, 1, 63, 8), le.Get(5))
	cpr.assertFloatEqual(
		eval(-5, 16, 0, 1, 105, 16, 0, 1, -315, 16, 0, 1, 231, 16),
		le.Get(6),
	)
	cpr.assertFloatEqual(
		eval(0, 1, -35, 16, 0, 1, 315, 16, 0, 1, -693, 16, 0, 1, 429, 16),
		le.Get(7),
	)
	cpr.assertFloatEqual(
		eval(
			35, 128,
			0, 1,
			-1260, 128,
			0, 1,
			6930, 128,
			0, 1,
			-12012, 128,
			0, 1,
			6435, 128,
		), le.Get(8),
	)
	cpr.assertFloatEqual(
		eval(
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
		), le.Get(9),
	)
	cpr.assertFloatEqual(
		eval(
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
		), le.Get(10),
	)
}

func TestAssocLegendre(t *testing.T) {
	prec := uint(100)
	cpr := newBigFloatComparator(t, prec)
	xf64 := 0.5
	x := NewFloat(xf64, prec)
	y := NewFloat(xf64, prec)
	y.Mul(y, y)
	y.Sub(NewFloat(1, prec), y)
	y.Sqrt(y)

	eval := func(m int, prec uint, c ...int) *big.Float {
		poly := evalPolynomial(fromRationals(c, prec), x, prec)
		return BlankFloat(prec).Mul(poly, PowerN(y, m))
	}

	scale := func(n, d int, v *big.Float) *big.Float {
		return BlankFloat(prec).Mul(NewFloatFromRat(n, d, prec), v)
	}

	le := EvalLegendre(x, 4)

	// m=0.
	al0 := EvalAssocLegendre(0, le)
	cpr.assertFloatEqual(le.Get(0), al0.Get(0))
	cpr.assertFloatEqual(le.Get(1), al0.Get(1))
	cpr.assertFloatEqual(le.Get(2), al0.Get(2))
	cpr.assertFloatEqual(le.Get(3), al0.Get(3))
	cpr.assertFloatEqual(le.Get(4), al0.Get(4))

	// m=1.
	al1 := EvalAssocLegendre(1, le)
	assert.Nil(t, al1.values[0])
	cpr.assertFloatEqual(eval(1, prec, -1, 1), al1.Get(1))
	cpr.assertFloatEqual(eval(1, prec, 0, 1, -3, 1), al1.Get(2))
	cpr.assertFloatEqual(eval(1, prec, 3, 2, 0, 1, -15, 2), al1.Get(3))
	cpr.assertFloatEqual(eval(1, prec, 0, 1, 15, 2, 0, 1, -35, 2), al1.Get(4))

	// m=-1.
	aln1 := EvalAssocLegendre(-1, le)
	assert.Nil(t, aln1.values[0])
	cpr.assertFloatEqual(scale(-1, 2, al1.Get(1)), aln1.Get(1))
	cpr.assertFloatEqual(scale(-1, 6, al1.Get(2)), aln1.Get(2))
	cpr.assertFloatEqual(scale(-1, 12, al1.Get(3)), aln1.Get(3))
	cpr.assertFloatEqual(scale(-1, 20, al1.Get(4)), aln1.Get(4))

	// m=2.
	al2 := EvalAssocLegendre(2, le)
	assert.Nil(t, al2.values[0])
	assert.Nil(t, al2.values[1])
	cpr.assertFloatEqual(eval(2, prec, 3, 1), al2.Get(2))
	cpr.assertFloatEqual(eval(2, prec, 0, 1, 15, 1), al2.Get(3))
	cpr.assertFloatEqual(eval(2, prec, -15, 2, 0, 1, 105, 2), al2.Get(4))

	// m=-2.
	aln2 := EvalAssocLegendre(-2, le)
	assert.Nil(t, aln2.values[0])
	assert.Nil(t, aln2.values[1])
	cpr.assertFloatEqual(scale(1, 24, al2.Get(2)), aln2.Get(2))
	cpr.assertFloatEqual(scale(1, 120, al2.Get(3)), aln2.Get(3))
	cpr.assertFloatEqual(scale(1, 360, al2.Get(4)), aln2.Get(4))

	// m=3.
	al3 := EvalAssocLegendre(3, le)
	assert.Nil(t, al3.values[0])
	assert.Nil(t, al3.values[1])
	assert.Nil(t, al3.values[2])
	cpr.assertFloatEqual(eval(3, prec, -15, 1), al3.Get(3))
	cpr.assertFloatEqual(eval(3, prec, 0, 1, -105, 1), al3.Get(4))

	// m=-3.
	aln3 := EvalAssocLegendre(-3, le)
	assert.Nil(t, aln3.values[0])
	assert.Nil(t, aln3.values[1])
	assert.Nil(t, aln3.values[1])
	cpr.assertFloatEqual(scale(-1, 720, al3.Get(3)), aln3.Get(3))
	cpr.assertFloatEqual(scale(-1, 5040, al3.Get(4)), aln3.Get(4))

	// m=4.
	al4 := EvalAssocLegendre(4, le)
	assert.Nil(t, al4.values[0])
	assert.Nil(t, al4.values[1])
	assert.Nil(t, al4.values[2])
	assert.Nil(t, al4.values[3])
	cpr.assertFloatEqual(eval(4, prec, 105, 1), al4.Get(4))

	// m=-4.
	aln4 := EvalAssocLegendre(-4, le)
	assert.Nil(t, aln4.values[0])
	assert.Nil(t, aln4.values[1])
	assert.Nil(t, aln4.values[1])
	assert.Nil(t, aln4.values[1])
	cpr.assertFloatEqual(scale(1, 40320, al4.Get(4)), aln4.Get(4))
}
