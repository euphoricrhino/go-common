package bigmath

import (
	"fmt"
	"math/big"
)

// Legendre stores P_l(x) at x for l=0..maxL.
type Legendre struct {
	// The maximum degree of the Legendre polynomial.
	maxL   int
	prec   uint
	x      *big.Float
	values []*big.Float
}

// EvalLegendre returns the evaluation of Legendre polynomial for l=0..maxL at x.
func EvalLegendre(x *big.Float, maxL int) *Legendre {
	if maxL < 0 {
		panic(fmt.Sprintf("maxL must be non-negative, got %v", maxL))
	}
	le := &Legendre{
		maxL:   maxL,
		prec:   x.Prec(),
		x:      CopyFloat(x),
		values: make([]*big.Float, maxL+1),
	}
	le.values[0] = NewFloat(1, le.prec)
	if maxL > 0 {
		le.values[1] = CopyFloat(x)
	}
	// Recurrence relation: P_l(x)=(2l-1)/l*x*P_{l-1}(x)-(l-1)/l*P_{l-2}(x)
	for l := 2; l <= maxL; l++ {
		v1 := NewFloatFromRat(2*l-1, l, le.prec)
		v1.Mul(v1, le.x)
		v1.Mul(v1, le.values[l-1])
		v2 := NewFloatFromRat(l-1, l, le.prec)
		v2.Mul(v2, le.values[l-2])
		le.values[l] = v1.Sub(v1, v2)
	}
	return le
}

// Get returns the Legendre polynomial value for P_l(x).
func (le *Legendre) Get(l int) *big.Float {
	return le.values[l]
}

// AssocLegendre stores P_l^m(x) for l=m..maxL and given m.
type AssocLegendre struct {
	le       *Legendre
	negative bool
	m        int
	values   []*big.Float
}

// EvalAssocLegendre returns the evaluation of associated Legendre polynomial for l=m..maxL at x. Requires 0 < m <= le.maxL
func EvalAssocLegendre(m int, le *Legendre) *AssocLegendre {
	if m < -le.maxL || m > le.maxL {
		panic(fmt.Sprintf("m must be within range [%v, %v], got %v", -le.maxL, le.maxL, m))
	}

	al := &AssocLegendre{
		le:     le,
		m:      m,
		values: make([]*big.Float, le.maxL+1),
	}

	// Copy values from le if m=0.
	if m == 0 {
		for l := 0; l <= le.maxL; l++ {
			al.values[l] = CopyFloat(le.values[l])
		}
		return al
	}

	absm := m
	if m < 0 {
		absm = -m
	}

	// Initial value: P_m^m(x)=(-1)^m(2m-1)!!(1-x^2)^(m/2)
	v1 := BlankFloat(le.prec).SetInt(Fact2(2*absm - 1))
	v2 := BlankFloat(le.prec).Mul(le.x, le.x)
	v2.Sub(NewFloat(1, le.prec), v2)
	v2.Sqrt(v2)
	v1.Mul(v1, PowerN(v2, absm))
	if absm%2 == 1 {
		v1.Neg(v1)
	}
	al.values[absm] = v1

	if le.maxL > absm {
		// Initial value: P_{m+1}^m(x)=(2m+1)xP_m^m(x)
		v1 = NewFloatFromInt(2*absm+1, le.prec)
		v1.Mul(v1, al.values[absm])
		v1.Mul(v1, le.x)
		al.values[absm+1] = v1
	}

	// Recurrence relation: P_l^m(x)=(2l-1)xP_{l-1}^m(x)/(l-m)-(l+m-1)P_{l-2}^m(x)/(l-m)
	for l := absm + 2; l <= le.maxL; l++ {
		v1 = NewFloatFromRat(2*l-1, l-absm, le.prec)
		v1.Mul(v1, le.x)
		v1.Mul(v1, al.values[l-1])
		v2 = NewFloatFromRat(l+absm-1, l-absm, le.prec)
		v2.Mul(v2, al.values[l-2])
		al.values[l] = v1.Sub(v1, v2)
	}

	// Deal with negative m.
	if m < 0 {
		// P_l^{-m}(x)=(-1)^m(l-m)!/(l+m)!P_l^m(x)
		f := BlankRat().SetFrac(NewInt(1), Fact(2*absm))
		if absm%2 == 1 {
			f.Neg(f)
		}
		al.values[absm].Mul(al.values[absm], BlankFloat(le.prec).SetRat(f))
		for l := absm + 1; l <= le.maxL; l++ {
			f.Mul(f, NewRat(l-absm, l+absm))
			al.values[l].Mul(al.values[l], BlankFloat(le.prec).SetRat(f))
		}
	}
	return al
}

// Get returns the associated Legendre polynomial value for P_l^m(x).
func (ale *AssocLegendre) Get(l int) *big.Float {
	return ale.values[l]
}
