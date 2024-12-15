package bigmath

import (
	"math"
	"math/big"
)

// SphericalHarmonics represents a family of spherical harmonics with l up to maxL and m in the given list of m.
type SphericalHarmonics struct {
	al map[int]*AssocLegendre
	// Untruncated normalization constant - (2l+1)(l-m)!/4(l+m)!
	c map[int][]*big.Rat
}

// NewSphericalHarmonics creates a family of spherical harmonics with l up to maxL and m in the given list of m.
func NewSphericalHarmonics(maxL int, m []int) *SphericalHarmonics {
	le := NewLegendre(maxL)
	al := make(map[int]*AssocLegendre)
	c := make(map[int][]*big.Rat)
	for _, mm := range m {
		al[mm] = NewAssocLegendre(mm, le)
		c[mm] = make([]*big.Rat, maxL+1)
		accum := (*big.Rat).Quo
		absm := abs(mm)
		if mm < 0 {
			accum = (*big.Rat).Mul
		}
		for l := absm; l <= le.maxL; l++ {
			r := NewRat(2*l+1, 4)
			for k := absm; k > -absm; k-- {
				accum(r, r, NewRat(l+k, 1))
			}
			c[mm][l] = r
		}
	}
	sh := &SphericalHarmonics{
		al: al,
		c:  c,
	}

	return sh
}

// Get returns the spherical harmonic value Y_{l,m}(theta, phi).
func (sh *SphericalHarmonics) Get(l, m int, theta, phi *big.Float) *ModArg {
	thetaf64, _ := theta.Float64()
	// We may have precision loss for using float64 sin/cos and constant Pi here.
	ctf64 := math.Cos(thetaf64)
	ct := NewFloat(ctf64, theta.Prec())

	r := sh.al[m].Get(l, ct)
	r.Mul(r, sh.norm(l, m, theta.Prec()))

	arg := NewFloat(float64(m), theta.Prec())
	arg.Mul(arg, phi)
	return NewModArg(r, arg)
}

// The normalization factor for Y_{l,m}.
func (sh *SphericalHarmonics) norm(l, m int, prec uint) *big.Float {
	coeff := BlankFloat(prec).SetRat(sh.c[m][l])
	coeff.Quo(coeff, NewFloat(math.Pi, prec))
	return coeff.Sqrt(coeff)
}
