package bigmath

import (
	"fmt"
	"math/big"
)

// Evaluator for x raised to some power, after construction, all power evaluations can be run in logarithmic time.
type PowerEvaluator struct {
	x *big.Float

	maxPower int
	// Precomputed all powers of form x^(2^n),
	powers []*big.Float
}

// NewPowerEvaluator returns a new PowerEvaluator for x raised to up to the maxPower (inclusive).
func NewPowerEvaluator(x *big.Float, maxPower int) *PowerEvaluator {
	bits := 0
	n := maxPower
	for n != 0 {
		n = n >> 1
		bits++
	}
	pEval := &PowerEvaluator{
		x:        x,
		maxPower: maxPower,
		powers:   make([]*big.Float, bits+1),
	}
	// powers[0] was never explicitly used by pow() below, so let it remain nil.
	if bits == 0 {
		return pEval
	}
	pEval.powers[1] = BlankFloat().Set(x)
	for k := 2; k <= bits; k++ {
		pEval.powers[k] = BlankFloat().Mul(pEval.powers[k-1], pEval.powers[k-1])
	}

	return pEval
}

func (pEval *PowerEvaluator) pow(n int) *big.Float {
	if n > pEval.maxPower {
		panic(fmt.Sprintf("power n is greater than maxPower %v", pEval.maxPower))
	}
	ans := NewFloat(1.0)
	shift := 1
	for n != 0 {
		if n&1 == 1 {
			ans.Mul(ans, pEval.powers[shift])
		}
		shift++
		n = n >> 1
	}
	return ans
}
