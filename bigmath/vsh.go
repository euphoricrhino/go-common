package bigmath

import (
	"math"
	"math/big"
)

// VSH represents a family of vector spherical harmonics with l up to maxL and m in the given list of m.
type VSH struct {
	y *SphericalHarmonics
	// Coefficients of derivatives of the polynomial part of P_l^m(x).
	der map[int][][]*big.Rat
}

// NewVSH creates a family of vector spherical harmonics with l up to maxL and m in the given list of m.
func NewVSH(maxL int, m []int) *VSH {
	vsh := &VSH{
		y:   NewSphericalHarmonics(maxL, m),
		der: make(map[int][][]*big.Rat),
	}

	for _, mm := range m {
		vsh.der[mm] = make([][]*big.Rat, maxL+1)
		for l := abs(mm) + 1; l <= maxL; l++ {
			poly := vsh.y.al[mm].lm[l]
			derlm := make([]*big.Rat, len(poly)-1)
			for k := 1; k < len(poly); k++ {
				if poly[k] != nil {
					d := NewRat(k, 1)
					d.Mul(d, poly[k])
					derlm[k-1] = d
				}
			}
			vsh.der[mm][l] = derlm
		}
	}
	return vsh
}

// GetY returns the Y vector of the VSH Y_{l,m}(theta, phi). The three components of the return value are along the r/theta/phi direction respectively.
func (vsh *VSH) GetY(l, m int, theta, phi *big.Float) [3]*ModArg {
	return [3]*ModArg{vsh.y.Get(l, m, theta, phi), nil, nil}
}

// GetPsi returns the Psi vector of the VSH Psi_{l,m}(theta, phi). The three components of the return value are along the r/theta/phi direction respectively.
func (vsh *VSH) GetPsi(l, m int, theta, phi *big.Float) [3]*ModArg {
	thetaf64, _ := theta.Float64()
	// We may have precision loss for using float64 sin/cos here.
	ctf64 := math.Cos(thetaf64)
	stf64 := math.Sin(thetaf64)
	ct := NewFloat(ctf64, theta.Prec())
	st := NewFloat(stf64, theta.Prec())
	absm := abs(m)
	vt := evalRatSkipPolynomial(ct, vsh.der[m][l])
	vt.Mul(vt, PowerN(st, absm+1))
	vt.Neg(vt)
	var poly, stm1 *big.Float
	if absm > 0 {
		tmp := NewFloatFromInt(absm, theta.Prec())
		tmp.Mul(tmp, ct)
		stm1 = PowerN(st, absm-1)
		tmp.Mul(tmp, stm1)
		poly = evalRatSkipPolynomial(ct, vsh.y.al[m].lm[l])
		tmp.Mul(tmp, poly)
		vt.Add(vt, tmp)
	}

	norm := vsh.y.norm(l, m, theta.Prec())
	vt.Mul(vt, norm)

	var vp *big.Float
	if absm > 0 {
		vp = NewFloatFromInt(m, theta.Prec())
		vp.Mul(vp, poly)
		vp.Mul(vp, stm1)
		vp.Mul(vp, norm)
	} else {
		vp = BlankFloat(theta.Prec())
	}

	argt := NewFloat(float64(m), phi.Prec())
	argt.Mul(argt, phi)

	argp := NewFloat(math.Pi/2, phi.Prec())
	argp.Add(argp, argt)

	return [3]*ModArg{
		nil,
		NewModArg(vt, argt),
		NewModArg(vp, argp),
	}
}

// GetPhi returns the Phi vector of the VSH Phi_{l,m}(theta, phi). The three components of the return value are along the r/theta/phi direction respectively.
func (vsh *VSH) GetPhi(l, m int, theta, phi *big.Float) [3]*ModArg {
	psi := vsh.GetPsi(l, m, theta, phi)
	return [3]*ModArg{nil, NewModArg(psi[2].Mod.Neg(psi[2].Mod), psi[2].Arg), psi[1]}
}
