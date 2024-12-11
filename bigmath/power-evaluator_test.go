package bigmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPowerEvaluator(t *testing.T) {
	pEval := NewPowerEvaluator(NewFloat(2), 60)
	assert.Equal(t, len(pEval.powers), 7)
	exp := NewFloat(1)
	for p := 0; p <= 60; p++ {
		assert.Equal(t, exp, pEval.pow(p))
		exp.Mul(exp, NewFloat(2))
	}
}
