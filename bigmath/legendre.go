package bigmath

import (
	"fmt"
	"math/big"
)

// Legendre represents a family of Legendre polynomials up to the maximum degree (maxL).
type Legendre struct {
	// Maximum degree of the Legendre polynomial.
	maxL int

	// Polynomial coefficients for P_l.
	coeff [][]*big.Rat
}

// NewLegendre creates the family of Legendre polynomials up to degree maxL.
func NewLegendre(maxL int) *Legendre {
	if maxL < 0 {
		panic(fmt.Sprintf("maxL must be non-negative, got %v", maxL))
	}
	le := &Legendre{
		maxL:  maxL,
		coeff: make([][]*big.Rat, maxL+1),
	}
	// Base case for recursion.
	le.coeff[0] = []*big.Rat{NewRat(1, 1)}
	if maxL > 0 {
		le.coeff[1] = []*big.Rat{nil, NewRat(1, 1)}
	}
	pprev, prev := le.coeff[0], le.coeff[1]
	for l := 2; l <= maxL; l++ {
		cur := make([]*big.Rat, l+1)
		le.coeff[l] = cur
		f1, f2 := NewRat(2*l-1, l), NewRat(l-1, l)
		for k := l; k >= 0; k -= 2 {
			le.coeff[l][k] = BlankRat()
			if k > 0 {
				cur[k].Mul(prev[k-1], f1)
			}
			if k <= l-2 {
				cur[k].Sub(cur[k], BlankRat().Mul(pprev[k], f2))
			}
		}
		pprev, prev = prev, cur
	}
	return le
}

// AssocLegendre represents a family of associated Legendre functions of a given order m, up to maximum degree (maxL).
type AssocLegendre struct {
	m        int
	negative bool
	// Coefficients for the m-th derivative of P_l.
	lm [][]*big.Rat
}

func NewAssocLegendre(m int, le *Legendre) *AssocLegendre {
	absM := m
	if m < 0 {
		absM = -m
	}
	if absM > le.maxL {
		panic(fmt.Sprintf("|m| must be no larger than maxL (%v), got %v", absM, le.maxL))
	}
	ale := &AssocLegendre{
		m:        absM,
		negative: m < 0,
		lm:       make([][]*big.Rat, le.maxL+1),
	}
	for l := absM; l <= le.maxL; l++ {
		lm := make([]*big.Rat, l-absM+1)
		ale.lm[l] = lm
		for k := l; k >= absM; k -= 2 {
			if absM%2 == 1 {
				lm[k-absM] = BlankRat().Neg(le.coeff[l][k])
			} else {
				lm[k-absM] = BlankRat().Set(le.coeff[l][k])
			}
			for j := 0; j < absM; j++ {
				lm[k-absM].Mul(lm[k-absM], NewRat(k-j, 1))
			}
		}
	}
	return ale
}
