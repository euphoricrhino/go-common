package bigmath

import "math/big"

// Polynomial represents a polynomial with big.Float coefficients.
type Polynomial struct {
	coeff []*big.Float
}

// NewPolynomial returns a new Polynomial with the given coefficients, where coeff[k] is the coefficient for x^k, nil coefficients means zero.
func NewPolynomial(coeff []*big.Float) *Polynomial {
	return &Polynomial{coeff: coeff}
}

func (poly *Polynomial) eval(x *big.Float) *big.Float {
	if len(poly.coeff) == 0 {
		return BlankFloat()
	}
	pEval := NewPowerEvaluator(x, len(poly.coeff)-1)
	ans := BlankFloat()
	for k := 0; k < len(poly.coeff); k++ {
		if poly.coeff[k] != nil {
			ans.Add(ans, BlankFloat().Mul(poly.coeff[k], pEval.pow(k)))
		}
	}
	return ans
}
