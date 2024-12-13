package bigmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFact(t *testing.T) {
	assert.Equal(t, NewInt(1), Fact(0))
	assert.Equal(t, NewInt(1), Fact(1))
	assert.Equal(t, NewInt(1307674368000), Fact(15))

	assert.Equal(t, NewInt(1), Fact2(0))
	assert.Equal(t, NewInt(1), Fact2(1))
	assert.Equal(t, NewInt(2), Fact2(2))
	assert.Equal(t, NewInt(1961990553600), Fact2(24))
	assert.Equal(t, NewInt(7905853580625), Fact2(25))
}

func TestPowerN(t *testing.T) {
	prec := uint(100)
	cpr := newBigFloatComparator(t, prec)
	cpr.assertFloatEqual(NewFloat(1, prec), PowerN(NewFloat(2, prec), 0))
	cpr.assertFloatEqual(NewFloat(1152921504606846976, prec), PowerN(NewFloat(2, prec), 60))
	cpr.assertFloatEqual(NewFloat(1024, prec), PowerN(NewFloat(2, prec), 10))
	cpr.assertFloatEqual(NewFloat(1.0/1024, prec), PowerN(NewFloat(2, prec), -10))
}
