package bigmath

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fromInts(v []int, prec uint) []*big.Float {
	ret := make([]*big.Float, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = NewFloatFromInt(v[i], prec)
	}
	return ret
}

func fromRationals(v []int, prec uint) []*big.Float {
	if len(v)%2 != 0 {
		panic("fromRationals requires an even number of arguments")
	}
	ret := make([]*big.Float, len(v)/2)
	for i := 0; i < len(v); i += 2 {
		ret[i/2] = NewFloatFromRat(v[i], v[i+1], prec)
	}
	return ret
}

func evalPolynomial(c []*big.Float, x *big.Float, prec uint) *big.Float {
	sum, power := BlankFloat(prec), NewFloat(1, prec)
	for _, v := range c {
		sum.Add(sum, BlankFloat(prec).Mul(v, power))
		power.Mul(power, x)
	}
	return sum
}

type bigFloatComparator struct {
	t      *testing.T
	prec   uint
	loPrec uint
}

func newBigFloatComparator(t *testing.T, prec uint) bigFloatComparator {
	// We test the equality of two big.Floats up to a 10% degradation of precision.
	return bigFloatComparator{t: t, prec: prec, loPrec: prec * 3 / 10}
}

func (cpr *bigFloatComparator) assertFloatEqual(expected, actual *big.Float) {
	pr := cpr.prec
	// We test the equality of two floats up to a 10% degradation of precision.
	for pr > cpr.loPrec {
		exp := BlankFloat(pr).Set(expected)
		act := BlankFloat(pr).Set(actual)
		if exp.Cmp(act) == 0 {
			return
		}
		pr--
		fmt.Printf("degrading precision to %v\n", pr)
	}
	assert.Fail(cpr.t, fmt.Sprintf("expected %v, got %v", expected, actual))
}
