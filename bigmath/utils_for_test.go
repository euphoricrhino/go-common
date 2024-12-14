package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func evalPolynomial(x float64, c []float64) float64 {
	sum, power := 0.0, 1.0
	for _, v := range c {
		sum += v * power
		power *= x
	}
	return sum
}

func assertFloatEqualF64(t *testing.T, expected float64, actual *big.Float, errBnd float64) {
	actf64, _ := actual.Float64()
	assert.InDelta(t, expected, actf64, errBnd)
}

func assertFloatEqual(t *testing.T, expected, actual *big.Float, errBnd float64) {
	expf64, _ := expected.Float64()
	assertFloatEqualF64(t, expf64, actual, errBnd)
}
