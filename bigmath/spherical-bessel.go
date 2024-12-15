package bigmath

import (
	"fmt"
	"math"
	"math/big"
)

// SphericalBessel represents a family of spherical Bessel functions up to the maximum order.
type SphericalBessel struct {
	// Kind of the spherical Bessel function: 1 for j_l, 2 for n_l, 3 for h_l(1), 4 for h_l(2).
	kind int
	// Maximum order of the spherical Bessel function.
	maxL int

	// Polynomial coefficients in 1/x for the spherical Bessel function.
	// j_l(x) = a_l(1/x) * sin(x)/x + b_l(1/x) * cos(x)
	// n_l(x) = c_l(1/x) * sin(x) + d_l(1/x) * cos(x)/x
	a [][]*big.Int
	b [][]*big.Int
	c [][]*big.Int
	d [][]*big.Int
}

// NewSphericalBessel creates a family of spherical Bessel functions up to the maximum order (maxL).
// kind is the kind of the spherical Bessel function: 1 for j_l, 2 for n_l, 3 for h_l(1), 4 for h_l(2).
func NewSphericalBessel(kind, maxL int) *SphericalBessel {
	if maxL < 0 {
		panic(fmt.Sprintf("maxL must be non-negative, got %v", maxL))
	}
	if kind < 1 || kind > 4 {
		panic(fmt.Sprintf("kind must be in [1, 4], got %v", kind))
	}
	generate := func(do bool, v0, v1 *big.Int) [][]*big.Int {
		if !do {
			return nil
		}
		v := make([][]*big.Int, maxL+1)
		v[0] = []*big.Int{v0}
		if maxL > 0 {
			v[1] = []*big.Int{nil, v1}
		}
		for l := 2; l <= maxL; l++ {
			vl := make([]*big.Int, l+1)
			v[l] = vl
			f := NewInt(2*l - 1)
			for k := l; k >= 0; k -= 2 {
				vl[k] = BlankInt()
				if k > 0 {
					vl[k].Mul(v[l-1][k-1], f)
				}
				if k <= l-2 {
					vl[k].Sub(vl[k], v[l-2][k])
				}
			}
		}
		return v
	}

	return &SphericalBessel{
		kind: kind,
		maxL: maxL,
		a:    generate(kind != 2, NewInt(1), NewInt(1)),
		b:    generate(kind != 2, NewInt(0), NewInt(-1)),
		c:    generate(kind != 1, NewInt(0), NewInt(-1)),
		d:    generate(kind != 1, NewInt(-1), NewInt(-1)),
	}
}

// Get returns z_l(x) where z is the spherical Bessel function of the corresponding kind.
func (sb *SphericalBessel) Get(l int, x *big.Float) *ReIm {
	xf64, _ := x.Float64()
	// We may have precision loss for using float64 sin/cos here.
	sf64, cf64 := math.Sin(xf64), math.Cos(xf64)
	s, c := NewFloat(sf64, x.Prec()), NewFloat(cf64, x.Prec())
	sx := BlankFloat(x.Prec()).Quo(s, x)
	cx := BlankFloat(x.Prec()).Quo(c, x)
	invx := NewFloat(1, x.Prec())
	invx.Quo(invx, x)
	var va, vb, vc, vd *big.Float
	if sb.kind != 2 {
		va = evalIntSkipPolynomial(invx, sb.a[l])
		va.Mul(va, sx)
		vb = evalIntSkipPolynomial(invx, sb.b[l])
		vb.Mul(vb, c)
	}
	if sb.kind != 1 {
		vc = evalIntSkipPolynomial(invx, sb.c[l])
		vc.Mul(vc, s)
		vd = evalIntSkipPolynomial(invx, sb.d[l])
		vd.Mul(vd, cx)
	}

	switch sb.kind {
	case 1:
		return NewReIm(va.Add(va, vb), nil)
	case 2:
		return NewReIm(vc.Add(vc, vd), nil)
	case 3:
		return NewReIm(va.Add(va, vb), vc.Add(vc, vd))
	case 4:
		vc.Add(vc, vd)
		vc.Neg(vc)
		return NewReIm(va.Add(va, vb), vc)
	}
	panic("unreachable")
}

func evalIntSkipPolynomial(x *big.Float, c []*big.Int) *big.Float {
	sum := BlankFloat(x.Prec())
	power := NewFloat(1, x.Prec())
	k := 0
	if len(c)%2 == 0 {
		k = 1
		power = CopyFloat(x)
	}
	x2 := BlankFloat(x.Prec()).Mul(x, x)
	term := BlankFloat(x.Prec())
	for ; k < len(c); k += 2 {
		term.SetInt(c[k])
		term.Mul(term, power)
		sum.Add(sum, term)
		power.Mul(power, x2)
	}
	return sum
}
