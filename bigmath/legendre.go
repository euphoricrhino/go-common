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

// Get returns P_l(x).
func (le *Legendre) Get(l int, x *big.Float) *big.Float {
	return evalRatSkipPolynomial(x, le.coeff[l])
}

// AssocLegendre represents a family of associated Legendre functions of a given order m, up to maximum degree (maxL).
type AssocLegendre struct {
	m  int
	le *Legendre
	// Coefficients for the m-th derivative of P_l.
	lm [][]*big.Rat
}

// NewAssocLegendre creates the family of associated Legendre functions of a given order m up to degree maxL.
func NewAssocLegendre(m int, le *Legendre) *AssocLegendre {
	absm := m
	if m < 0 {
		absm = -m
	}
	if absm > le.maxL {
		panic(fmt.Sprintf("|m| must be no larger than maxL (%v), got %v", le.maxL, m))
	}
	al := &AssocLegendre{
		m:  m,
		le: le,
		lm: make([][]*big.Rat, le.maxL+1),
	}
	for l := absm; l <= le.maxL; l++ {
		lm := make([]*big.Rat, l-absm+1)
		al.lm[l] = lm
		for k := l; k >= absm; k -= 2 {
			if absm%2 == 1 {
				lm[k-absm] = BlankRat().Neg(le.coeff[l][k])
			} else {
				lm[k-absm] = BlankRat().Set(le.coeff[l][k])
			}
			for j := 0; j < absm; j++ {
				lm[k-absm].Mul(lm[k-absm], NewRat(k-j, 1))
			}
		}
	}
	if m < 0 {
		// P_l^{-m}(x)=(-1)^m(l-m)!/(l+m)!P_l^m(x)
		f := BlankRat().SetFrac(NewInt(1), Fact(2*absm))
		if absm%2 == 1 {
			f.Neg(f)
		}
		for _, c := range al.lm[absm] {
			if c != nil {
				c.Mul(c, f)
			}
		}
		for l := absm + 1; l <= le.maxL; l++ {
			f.Mul(f, NewRat(l-absm, l+absm))
			for _, c := range al.lm[l] {
				if c != nil {
					c.Mul(c, f)
				}
			}
		}
	}
	return al
}

// Get returns P_l^m(x).
func (al *AssocLegendre) Get(l int, x *big.Float) *big.Float {
	p := evalRatSkipPolynomial(x, al.lm[l])
	s := BlankFloat(x.Prec()).Mul(x, x)
	s.Sub(NewFloat(1, x.Prec()), s)
	s.Sqrt(s)
	absm := al.m
	if al.m < 0 {
		absm = -al.m
	}
	return p.Mul(p, PowerN(s, absm))
}

func evalRatSkipPolynomial(x *big.Float, c []*big.Rat) *big.Float {
	sum := BlankFloat(x.Prec())
	power := NewFloat(1, x.Prec())
	k := 0
	if len(c)%2 == 0 {
		k = 1
		power = CopyFloat(x)
	}
	x2 := BlankFloat(x.Prec()).Mul(x, x)
	for ; k < len(c); k += 2 {
		sum.Add(sum, BlankFloat(x.Prec()).Mul(power, BlankFloat(x.Prec()).SetRat(c[k])))
		power.Mul(power, x2)
	}
	return sum
}
