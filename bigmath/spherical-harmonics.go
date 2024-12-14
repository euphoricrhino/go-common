package bigmath

import (
	"math"
	"math/big"
)

// SphericalHarmonics stores the spherical harmonics Y_{lm}(theta, phi) for l value up to maxL, and m values in the designated list.
type SphericalHarmonics struct {
	al map[int]*AssocLegendre
	// Normalization constants, map key is m, then index to the slice is l.
	c   map[int][]*big.Float
	phi *big.Float
}

func NewSphericalHarmonics(theta, phi *big.Float, maxL int, m []int) *SphericalHarmonics {
	thetaf64, _ := theta.Float64()
	ct := NewFloat(math.Cos(thetaf64), theta.Prec())
	le := EvalLegendre(ct, maxL)
	al := make(map[int]*AssocLegendre)
	c := make(map[int][]*big.Float)
	pi := NewFloat(math.Pi, le.prec)
	for _, mm := range m {
		al[mm] = EvalAssocLegendre(mm, le)
		c[mm] = make([]*big.Float, maxL+1)
		absm := mm
		accum := (*big.Float).Quo
		if mm < 0 {
			absm = -mm
			accum = (*big.Float).Mul
		}
		for l := absm; l <= le.maxL; l++ {
			// 2l+1/4pi*sqrt((l-m)!/(l+m)!)
			f := NewFloatFromRat(2*l+1, 4, le.prec)
			f.Quo(f, pi)
			for k := absm; k > -absm; k-- {
				accum(f, f, NewFloatFromInt(l+k, le.prec))
			}
			c[mm][l] = f.Sqrt(f)
		}
	}
	sh := &SphericalHarmonics{
		al:  al,
		c:   c,
		phi: CopyFloat(phi),
	}

	return sh
}

func (sh *SphericalHarmonics) Get(l, m int) (*big.Float, *big.Float) {
	phif64, _ := sh.phi.Float64()
	// c_{lm}P_l^m(cos(theta))
	al := sh.al[m]
	r := BlankFloat(al.le.prec).Mul(sh.c[m][l], al.values[l])
	mphi := float64(m) * phif64
	cp := NewFloat(math.Cos(mphi), sh.phi.Prec())
	sp := NewFloat(math.Sin(mphi), sh.phi.Prec())
	return cp.Mul(cp, r), sp.Mul(sp, r)
}
