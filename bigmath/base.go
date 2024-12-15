package bigmath

import (
	"math/big"
)

// BlankFloat returns a new big.Float at the given precision.
func BlankFloat(prec uint) *big.Float { return big.NewFloat(0).SetPrec(prec) }

// CopyFloat returns a new copy of x.
func CopyFloat(x *big.Float) *big.Float {
	return big.NewFloat(0).Copy(x)
}

// NewFloat returns a new big.Float with value x and at the given precision.
func NewFloat(x float64, prec uint) *big.Float {
	return big.NewFloat(x).SetPrec(prec)
}

// NewFloatFromInt returns a new big.Float with int value x and at the given precision.
func NewFloatFromInt(x int, prec uint) *big.Float {
	return BlankFloat(prec).SetInt64(int64(x))
}

// NewFloatFromRat returns a new big.Float with value n/d and at the given precision.
func NewFloatFromRat(n, d int, prec uint) *big.Float {
	return BlankFloat(prec).SetRat(big.NewRat(int64(n), int64(d)))
}

// BlankRat returns a new big.Rat with value 0.
func BlankRat() *big.Rat { return big.NewRat(0, 1) }

// NewRat returns a new big.Rat with value n/d.
func NewRat(n, d int) *big.Rat { return big.NewRat(int64(n), int64(d)) }

// BlankInt returns a new big.Int with value 0.
func BlankInt() *big.Int { return big.NewInt(0) }

// NewInt returns a new big.Int with value x.
func NewInt(x int) *big.Int { return big.NewInt(int64(x)) }

// Fact returns n!.
func Fact(n int) *big.Int {
	if n < 0 {
		panic("n must be non-negative")
	}
	ans := NewInt(1)
	for i := 2; i <= n; i++ {
		ans.Mul(ans, NewInt(i))
	}
	return ans
}

// Fact2 returns n!!.
func Fact2(n int) *big.Int {
	if n < 0 {
		panic("n must be non-negative")
	}
	ans := NewInt(1)
	for i := n; i > 1; i -= 2 {
		ans.Mul(ans, NewInt(i))
	}
	return ans
}

// PowerN returns base^n.
func PowerN(base *big.Float, n int) *big.Float {
	result := NewFloat(1, base.Prec())
	absn := abs(n)

	r := CopyFloat(base)
	// Perform repeated multiplication
	for absn > 0 {
		if absn%2 == 1 {
			result.Mul(result, r)
		}
		r.Mul(r, r)
		absn /= 2
	}

	if n < 0 {
		result.Quo(NewFloat(1, base.Prec()), result)
	}

	return result
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
