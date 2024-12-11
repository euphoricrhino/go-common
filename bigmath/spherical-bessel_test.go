package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSphericalBessel(t *testing.T) {
	j := NewSphericalBessel(1, 5)
	n := NewSphericalBessel(2, 5)
	// h1 := NewSphericalBessel(3, 5)
	// h2 := NewSphericalBessel(4, 5)

	assert.Nil(t, j.c)
	assert.Nil(t, j.d)
	assert.Nil(t, n.a)
	assert.Nil(t, n.b)
	assert.Equal(t, []*big.Int{NewInt(1)}, j.a[0])
	assert.Equal(t, []*big.Int{NewInt(0)}, j.b[0])
	assert.Equal(t, []*big.Int{nil, NewInt(1)}, j.a[1])
	assert.Equal(t, []*big.Int{nil, NewInt(-1)}, j.b[1])
	assert.Equal(t, []*big.Int{NewInt(-1), nil, NewInt(3)}, j.a[2])
	assert.Equal(t, []*big.Int{NewInt(0), nil, NewInt(-3)}, j.b[2])
	assert.Equal(t, []*big.Int{nil, NewInt(-6), nil, NewInt(15)}, j.a[3])
	assert.Equal(t, []*big.Int{nil, NewInt(1), nil, NewInt(-15)}, j.b[3])
	assert.Equal(t, []*big.Int{NewInt(1), nil, NewInt(-45), nil, NewInt(105)}, j.a[4])
	assert.Equal(t, []*big.Int{NewInt(0), nil, NewInt(10), nil, NewInt(-105)}, j.b[4])
	assert.Equal(t, []*big.Int{nil, NewInt(15), nil, NewInt(-420), nil, NewInt(945)}, j.a[5])
	assert.Equal(t, []*big.Int{nil, NewInt(-1), nil, NewInt(105), nil, NewInt(-945)}, j.b[5])

	assert.Equal(t, []*big.Int{NewInt(0)}, n.c[0])
	assert.Equal(t, []*big.Int{NewInt(-1)}, n.d[0])
	assert.Equal(t, []*big.Int{nil, NewInt(-1)}, n.c[1])
	assert.Equal(t, []*big.Int{nil, NewInt(-1)}, n.d[1])
	assert.Equal(t, []*big.Int{NewInt(0), nil, NewInt(-3)}, n.c[2])
	assert.Equal(t, []*big.Int{NewInt(1), nil, NewInt(-3)}, n.d[2])
	assert.Equal(t, []*big.Int{nil, NewInt(1), nil, NewInt(-15)}, n.c[3])
	assert.Equal(t, []*big.Int{nil, NewInt(6), nil, NewInt(-15)}, n.d[3])
	assert.Equal(t, []*big.Int{NewInt(0), nil, NewInt(10), nil, NewInt(-105)}, n.c[4])
	assert.Equal(t, []*big.Int{NewInt(-1), nil, NewInt(45), nil, NewInt(-105)}, n.d[4])
	assert.Equal(t, []*big.Int{nil, NewInt(-1), nil, NewInt(105), nil, NewInt(-945)}, n.c[5])
	assert.Equal(t, []*big.Int{nil, NewInt(-15), nil, NewInt(420), nil, NewInt(-945)}, n.d[5])
}
