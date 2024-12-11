package bigmath

import (
	"fmt"
	"math/big"
)

// SphericalBessel represents a family of spherical Bessel functions up to the maximum order.
type SphericalBessel struct {
	// Kind of the spherical Bessel function: 1 for j_l, 2 for n_l, 3 for h_l(1), 4 for h_l(2).
	kind int
	// Maximum order of the spherical Bessel function.
	maxL int

	a [][]*big.Int
	b [][]*big.Int
	c [][]*big.Int
	d [][]*big.Int
}

// NewSphericalBessel creates the family of spherical Bessel functions up to the maximum order (maxL).
// kind is the kind of spherical Bessel function, 1 for j_l, 2 for n_l, 3 for h_l(1), 4 for h_l(2).
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
