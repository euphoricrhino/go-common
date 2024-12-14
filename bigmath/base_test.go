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
	assertFloatEqualF64(t, 1, PowerN(NewFloat(2, prec), 0), 0)
	assertFloatEqualF64(t, 1152921504606846976, PowerN(NewFloat(2, prec), 60), 0)
	assertFloatEqualF64(t, 1024, PowerN(NewFloat(2, prec), 10), 0)
	assertFloatEqualF64(t, 1.0/1024, PowerN(NewFloat(2, prec), -10), 0)
}
