package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLegendre(t *testing.T) {
	le := NewLegendre(10)
	assert.Equal(t, 11, len(le.coeff))

	assert.Equal(t, []*big.Rat{NewRat(1, 1)}, le.coeff[0])
	assert.Equal(t, []*big.Rat{nil, NewRat(1, 1)}, le.coeff[1])
	assert.Equal(t, []*big.Rat{NewRat(-1, 2), nil, NewRat(3, 2)}, le.coeff[2])
	assert.Equal(t, []*big.Rat{nil, NewRat(-3, 2), nil, NewRat(5, 2)}, le.coeff[3])
	assert.Equal(t, []*big.Rat{NewRat(3, 8), nil, NewRat(-30, 8), nil, NewRat(35, 8)}, le.coeff[4])
	assert.Equal(
		t,
		[]*big.Rat{nil, NewRat(15, 8), nil, NewRat(-70, 8), nil, NewRat(63, 8)},
		le.coeff[5],
	)
	assert.Equal(
		t,
		[]*big.Rat{
			NewRat(-5, 16),
			nil,
			NewRat(105, 16),
			nil,
			NewRat(-315, 16),
			nil,
			NewRat(231, 16),
		},
		le.coeff[6],
	)
	assert.Equal(
		t,
		[]*big.Rat{
			nil,
			NewRat(-35, 16),
			nil,
			NewRat(315, 16),
			nil,
			NewRat(-693, 16),
			nil,
			NewRat(429, 16),
		},
		le.coeff[7],
	)
	assert.Equal(
		t,
		[]*big.Rat{
			NewRat(35, 128),
			nil,
			NewRat(-1260, 128),
			nil,
			NewRat(6930, 128),
			nil,
			NewRat(-12012, 128),
			nil,
			NewRat(6435, 128),
		},
		le.coeff[8],
	)
	assert.Equal(
		t,
		[]*big.Rat{
			nil,
			NewRat(315, 128),
			nil,
			NewRat(-4620, 128),
			nil,
			NewRat(18018, 128),
			nil,
			NewRat(-25740, 128),
			nil,
			NewRat(12155, 128),
		},
		le.coeff[9],
	)
	assert.Equal(
		t,
		[]*big.Rat{
			NewRat(-63, 256),
			nil,
			NewRat(3465, 256),
			nil,
			NewRat(-30030, 256),
			nil,
			NewRat(90090, 256),
			nil,
			NewRat(-109395, 256),
			nil,
			NewRat(46189, 256),
		},
		le.coeff[10],
	)
}

func TestAssocLegendre(t *testing.T) {
	le := NewLegendre(4)

	// m=0
	ale0 := NewAssocLegendre(0, le)
	assert.Equal(t, 0, ale0.m)
	assert.False(t, ale0.negative)
	assert.Equal(t, 5, len(ale0.lm))
	assert.Equal(t, []*big.Rat{NewRat(1, 1)}, ale0.lm[0])
	assert.Equal(t, []*big.Rat{nil, NewRat(1, 1)}, ale0.lm[1])
	assert.Equal(t, []*big.Rat{NewRat(-1, 2), nil, NewRat(3, 2)}, ale0.lm[2])
	assert.Equal(t, []*big.Rat{nil, NewRat(-3, 2), nil, NewRat(5, 2)}, ale0.lm[3])
	assert.Equal(t, []*big.Rat{NewRat(3, 8), nil, NewRat(-30, 8), nil, NewRat(35, 8)}, ale0.lm[4])

	// m=1
	ale1 := NewAssocLegendre(1, le)
	assert.Equal(t, 1, ale1.m)
	assert.False(t, ale1.negative)
	assert.Equal(t, 5, len(ale1.lm))
	assert.Nil(t, ale1.lm[0])
	assert.Equal(t, []*big.Rat{NewRat(-1, 1)}, ale1.lm[1])
	assert.Equal(t, []*big.Rat{nil, NewRat(-3, 1)}, ale1.lm[2])
	assert.Equal(t, []*big.Rat{NewRat(3, 2), nil, NewRat(-15, 2)}, ale1.lm[3])
	assert.Equal(t, []*big.Rat{nil, NewRat(15, 2), nil, NewRat(-35, 2)}, ale1.lm[4])
	// m=-1
	alen1 := NewAssocLegendre(-1, le)
	assert.Equal(t, 1, alen1.m)
	assert.True(t, alen1.negative)
	assert.Equal(t, ale1.lm, alen1.lm)

	// m=2
	ale2 := NewAssocLegendre(2, le)
	assert.Equal(t, 2, ale2.m)
	assert.False(t, ale2.negative)
	assert.Equal(t, 5, len(ale2.lm))
	assert.Nil(t, ale2.lm[0])
	assert.Nil(t, ale2.lm[1])
	assert.Equal(t, []*big.Rat{NewRat(3, 1)}, ale2.lm[2])
	assert.Equal(t, []*big.Rat{nil, NewRat(15, 1)}, ale2.lm[3])
	assert.Equal(t, []*big.Rat{NewRat(-15, 2), nil, NewRat(105, 2)}, ale2.lm[4])
	// m=-2
	alen2 := NewAssocLegendre(-2, le)
	assert.Equal(t, 2, alen2.m)
	assert.True(t, alen2.negative)

	// m=3
	ale3 := NewAssocLegendre(3, le)
	assert.Equal(t, 3, ale3.m)
	assert.False(t, ale3.negative)
	assert.Equal(t, 5, len(ale3.lm))
	assert.Nil(t, ale3.lm[0])
	assert.Nil(t, ale3.lm[1])
	assert.Nil(t, ale3.lm[2])
	assert.Equal(t, []*big.Rat{NewRat(-15, 1)}, ale3.lm[3])
	assert.Equal(t, []*big.Rat{nil, NewRat(-105, 1)}, ale3.lm[4])
	// m=-3
	alen3 := NewAssocLegendre(-3, le)
	assert.Equal(t, 3, alen3.m)
	assert.True(t, alen3.negative)

	// m=4
	ale4 := NewAssocLegendre(4, le)
	assert.Equal(t, 4, ale4.m)
	assert.False(t, ale4.negative)
	assert.Equal(t, 5, len(ale4.lm))
	assert.Nil(t, ale4.lm[0])
	assert.Nil(t, ale4.lm[1])
	assert.Nil(t, ale4.lm[2])
	assert.Nil(t, ale4.lm[3])
	assert.Equal(t, []*big.Rat{NewRat(105, 1)}, ale4.lm[4])
	// m=-4
	alen4 := NewAssocLegendre(-4, le)
	assert.Equal(t, 4, alen4.m)
	assert.True(t, alen4.negative)
}
