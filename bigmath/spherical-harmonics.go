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
		absm := mm
		if mm < 0 {
			absm = -mm
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

// Get returns the real part and imaginary part of the spherical harmonic value Y_{l,m}(theta, phi).
func (sh *SphericalHarmonics) Get(l, m int, theta, phi *big.Float) (*big.Float, *big.Float) {
	thetaf64, _ := theta.Float64()
	// We may have precision loss for using float64 sin/cos and constant Pi here.
	ctf64 := math.Cos(thetaf64)
	ct := NewFloat(ctf64, theta.Prec())
	coeff := BlankFloat(theta.Prec()).SetRat(sh.c[m][l])
	coeff.Quo(coeff, NewFloat(math.Pi, theta.Prec()))
	coeff.Sqrt(coeff)

	r := sh.al[m].Get(l, ct)
	r.Mul(r, coeff)

	phif64, _ := phi.Float64()
	mphi := float64(m) * phif64
	cp := NewFloat(math.Cos(mphi), phi.Prec())
	sp := NewFloat(math.Sin(mphi), phi.Prec())
	return cp.Mul(cp, r), sp.Mul(sp, r)
}
