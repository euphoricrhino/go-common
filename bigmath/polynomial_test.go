package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolynomial(t *testing.T) {
	// Trivial case.
	assert.Equal(t, NewFloat(0), NewPolynomial(nil).eval(NewFloat(1)))
	assert.Equal(t, NewFloat(0), NewPolynomial([]*big.Float{}).eval(NewFloat(1)))
	assert.Equal(t, NewFloat(0), NewPolynomial([]*big.Float{nil}).eval(NewFloat(1)))

	// Polynomial: 13+2x+3x^7+5x^12
	coeff := make([]*big.Float, 13)
	coeff[0] = NewFloat(13)
	coeff[1] = NewFloat(2)
	coeff[7] = NewFloat(3)
	coeff[12] = NewFloat(5)
	poly := NewPolynomial(coeff)

	assert.Equal(t, NewFloat(13+2*2+3*128+5*4096), poly.eval(NewFloat(2)))
	assert.Equal(t, NewFloat(13+2*3+3*2187+5*531441), poly.eval(NewFloat(3)))
}
