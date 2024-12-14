package bigmath

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLegendre(t *testing.T) {
	prec := uint(100)
	errBnd := 0.0
	xf64 := 0.5
	x := NewFloat(xf64, prec)
	le := EvalLegendre(x, 10)

	eval := func(c ...float64) float64 { return evalPolynomial(xf64, c) }

	// https://en.wikipedia.org/wiki/Legendre_polynomials
	assertFloatEqualF64(t, eval(1), le.Get(0), errBnd)
	assertFloatEqualF64(t, eval(0, 1), le.Get(1), errBnd)
	assertFloatEqualF64(t, eval(-1.0/2, 0, 3.0/2), le.Get(2), errBnd)
	assertFloatEqualF64(t, eval(0, -3.0/2, 0, 5.0/2), le.Get(3), errBnd)
	assertFloatEqualF64(t, eval(3.0/8, 0, -30.0/8, 0, 35.0/8), le.Get(4), errBnd)
	assertFloatEqualF64(t, eval(0, 15.0/8, 0, -70.0/8, 0, 63.0/8), le.Get(5), errBnd)
	assertFloatEqualF64(t, eval(-5.0/16, 0, 105.0/16, 0, -315.0/16, 0, 231.0/16), le.Get(6), errBnd)
	assertFloatEqualF64(
		t,
		eval(0, -35.0/16, 0, 315.0/16, 0, -693.0/16, 0, 429.0/16),
		le.Get(7),
		errBnd,
	)
	assertFloatEqualF64(
		t,
		eval(35.0/128, 0, -1260.0/128, 0, 6930.0/128, 0, -12012.0/128, 0, 6435.0/128),
		le.Get(8),
		errBnd,
	)
	assertFloatEqualF64(
		t,
		eval(0, 315.0/128, 0, -4620.0/128, 0, 18018.0/128, 0, -25740.0/128, 0, 12155.0/128),
		le.Get(9),
		errBnd,
	)
	assertFloatEqualF64(
		t,
		eval(
			-63.0/256,
			0,
			3465.0/256,
			0,
			-30030.0/256,
			0,
			90090.0/256,
			0,
			-109395.0/256,
			0,
			46189.0/256,
		),
		le.Get(10),
		errBnd,
	)
}

func TestAssocLegendre(t *testing.T) {
	prec := uint(100)
	errBnd := 1e-13
	xf64 := 0.5
	x := NewFloat(xf64, prec)
	y := NewFloat(xf64, prec)
	y.Mul(y, y)
	y.Sub(NewFloat(1, prec), y)
	y.Sqrt(y)

	eval := func(m int, c ...float64) float64 {
		poly := evalPolynomial(xf64, c)
		return poly * math.Pow(math.Sqrt(1.0-xf64*xf64), float64(m))
	}

	scale := func(n, d int, v *big.Float) *big.Float {
		return BlankFloat(prec).Mul(NewFloatFromRat(n, d, prec), v)
	}

	le := EvalLegendre(x, 4)

	// m=0.
	al0 := EvalAssocLegendre(0, le)
	assert.Equal(t, le.Get(0), al0.Get(0))
	assert.Equal(t, le.Get(1), al0.Get(1))
	assert.Equal(t, le.Get(2), al0.Get(2))
	assert.Equal(t, le.Get(3), al0.Get(3))
	assert.Equal(t, le.Get(4), al0.Get(4))

	// m=1.
	al1 := EvalAssocLegendre(1, le)
	assert.Nil(t, al1.values[0])
	assertFloatEqualF64(t, eval(1, -1), al1.Get(1), errBnd)
	assertFloatEqualF64(t, eval(1, 0, -3), al1.Get(2), errBnd)
	assertFloatEqualF64(t, eval(1, 3.0/2, 0, -15.0/2), al1.Get(3), errBnd)
	assertFloatEqualF64(t, eval(1, 0, 15.0/2, 0, -35.0/2), al1.Get(4), errBnd)

	// m=-1.
	aln1 := EvalAssocLegendre(-1, le)
	assert.Nil(t, aln1.values[0])
	assertFloatEqual(t, scale(-1, 2, al1.Get(1)), aln1.Get(1), errBnd)
	assertFloatEqual(t, scale(-1, 6, al1.Get(2)), aln1.Get(2), errBnd)
	assertFloatEqual(t, scale(-1, 12, al1.Get(3)), aln1.Get(3), errBnd)
	assertFloatEqual(t, scale(-1, 20, al1.Get(4)), aln1.Get(4), errBnd)

	// m=2.
	al2 := EvalAssocLegendre(2, le)
	assert.Nil(t, al2.values[0])
	assert.Nil(t, al2.values[1])
	assertFloatEqualF64(t, eval(2, 3), al2.Get(2), errBnd)
	assertFloatEqualF64(t, eval(2, 0, 15), al2.Get(3), errBnd)
	assertFloatEqualF64(t, eval(2, -15.0/2, 0, 105.0/2), al2.Get(4), errBnd)

	// m=-2.
	aln2 := EvalAssocLegendre(-2, le)
	assert.Nil(t, aln2.values[0])
	assert.Nil(t, aln2.values[1])
	assertFloatEqual(t, scale(1, 24, al2.Get(2)), aln2.Get(2), errBnd)
	assertFloatEqual(t, scale(1, 120, al2.Get(3)), aln2.Get(3), errBnd)
	assertFloatEqual(t, scale(1, 360, al2.Get(4)), aln2.Get(4), errBnd)

	// m=3.
	al3 := EvalAssocLegendre(3, le)
	assert.Nil(t, al3.values[0])
	assert.Nil(t, al3.values[1])
	assert.Nil(t, al3.values[2])
	assertFloatEqualF64(t, eval(3, -15), al3.Get(3), errBnd)
	assertFloatEqualF64(t, eval(3, 0, -105), al3.Get(4), errBnd)

	// m=-3.
	aln3 := EvalAssocLegendre(-3, le)
	assert.Nil(t, aln3.values[0])
	assert.Nil(t, aln3.values[1])
	assert.Nil(t, aln3.values[1])
	assertFloatEqual(t, scale(-1, 720, al3.Get(3)), aln3.Get(3), errBnd)
	assertFloatEqual(t, scale(-1, 5040, al3.Get(4)), aln3.Get(4), errBnd)

	// m=4.
	al4 := EvalAssocLegendre(4, le)
	assert.Nil(t, al4.values[0])
	assert.Nil(t, al4.values[1])
	assert.Nil(t, al4.values[2])
	assert.Nil(t, al4.values[3])
	assertFloatEqualF64(t, eval(4, 105), al4.Get(4), errBnd)

	// m=-4.
	aln4 := EvalAssocLegendre(-4, le)
	assert.Nil(t, aln4.values[0])
	assert.Nil(t, aln4.values[1])
	assert.Nil(t, aln4.values[1])
	assert.Nil(t, aln4.values[1])
	assertFloatEqual(t, scale(1, 40320, al4.Get(4)), aln4.Get(4), errBnd)
}
