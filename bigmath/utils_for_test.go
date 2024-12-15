package bigmath

import (
	"math"
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

func assertModArgEqualF64(t *testing.T, expMod, expArg float64, actual *ModArg, errBnd float64) {
	argf64, _ := actual.Arg.Float64()
	actualRe := NewFloat(math.Cos(argf64), actual.Arg.Prec())
	actualRe.Mul(actualRe, actual.Mod)
	actualIm := NewFloat(math.Sin(argf64), actual.Arg.Prec())
	actualIm.Mul(actualIm, actual.Mod)
	assertFloatEqualF64(t, expMod*math.Cos(expArg), actualRe, errBnd)
	assertFloatEqualF64(t, expMod*math.Sin(expArg), actualIm, errBnd)
}
