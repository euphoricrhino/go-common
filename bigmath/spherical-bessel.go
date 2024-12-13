package bigmath

import (
	"fmt"
	"math"
	"math/big"
)

// SphericalBessel stores values of spherical Bessel functions at x for l=0..maxL.
type SphericalBessel struct {
	kind int
	maxL int
	re   []*big.Float
	im   []*big.Float
}

// EvalSphericalBessel returns the evaluation of spherical Bessel function for l=0..maxL at x.
// kind represents the kind of spherical Bessel function: 1 for j_l(x), 2 for n_l(x), 3 for h_l^(1)(x), 4 for h_l^(2)(x).
func EvalSphericalBessel(x *big.Float, maxL int, kind int) *SphericalBessel {
	if maxL < 0 {
		panic(fmt.Sprintf("maxL must be non-negative, got %v", maxL))
	}
	sb := &SphericalBessel{
		kind: kind,
		maxL: maxL,
	}

	coeff := func(v0, v1 *big.Float) []*big.Float {
		v := make([]*big.Float, maxL+1)
		v[0] = v0
		if maxL > 0 {
			v[1] = v1
		}
		for l := 2; l <= maxL; l++ {
			f := NewFloatFromInt(2*l-1, x.Prec())
			v[l] = BlankFloat(x.Prec()).Mul(f, v[l-1])
			v[l].Quo(v[l], x)
			v[l].Sub(v[l], v[l-2])
		}
		return v
	}

	val := func(u, v []*big.Float, s, c *big.Float) []*big.Float {
		ret := make([]*big.Float, maxL+1)
		for l := 0; l <= maxL; l++ {
			ret[l] = BlankFloat(x.Prec()).Mul(u[l], s)
			ret[l].Add(ret[l], BlankFloat(x.Prec()).Mul(v[l], c))
		}
		return ret
	}

	invx := func(sign float64) *big.Float {
		invx := NewFloat(sign, x.Prec())
		invx.Quo(invx, x)
		return invx
	}

	sc := func() (*big.Float, *big.Float) {
		// For the lack of Sin/Cos on big.Float, we use the float64 version, which might have precision loss.
		x64f, _ := x.Float64()
		return NewFloat(math.Sin(x64f), x.Prec()), NewFloat(math.Cos(x64f), x.Prec())
	}

	switch kind {
	case 1:
		// 1, 1/x
		u := coeff(NewFloat(1, x.Prec()), invx(1))
		// 0, -1/x
		v := coeff(NewFloat(0, x.Prec()), invx(-1))
		s, c := sc()
		s.Quo(s, x)
		// sin(x)/x, cos(x)
		sb.re = val(u, v, s, c)
	case 2:
		// 0, -1/x
		u := coeff(NewFloat(0, x.Prec()), invx(-1))
		// -1, -1/x
		v := coeff(NewFloat(-1, x.Prec()), invx(-1))
		s, c := sc()
		c.Quo(c, x)
		// sin(x), cos(x)/x
		sb.re = val(u, v, s, c)
	case 3:
		fallthrough
	case 4:
		// Re: 1, 1/x; 0, -1/x
		reu := coeff(NewFloat(1, x.Prec()), invx(1))
		rev := coeff(NewFloat(0, x.Prec()), invx(-1))
		// sin(x), cos(x)/x
		res, rec := sc()
		res.Quo(res, x)
		sb.re = val(reu, rev, res, rec)

		// Im: 0, -1/x; -1, -1/x
		imu := coeff(NewFloat(0, x.Prec()), invx(-1))
		imv := coeff(NewFloat(-1, x.Prec()), invx(-1))
		ims, imc := sc()
		imc.Quo(imc, x)
		sb.im = val(imu, imv, ims, imc)
		// For h_l^(2)(x), negate the imaginary part.
		if kind == 4 {
			for l := 0; l <= maxL; l++ {
				sb.im[l].Neg(sb.im[l])
			}
		}
	default:
		panic(fmt.Sprintf("kind must be within range [1, 4], got %v", kind))
	}
	return sb
}

// Get returns the real and imaginary parts (or nil when not applicable) of the spherical Bessel function for l.
func (sb *SphericalBessel) Get(l int) (*big.Float, *big.Float) {
	if sb.kind <= 2 {
		return sb.re[l], nil
	}
	return sb.re[l], sb.im[l]
}
