package graphix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertVec3Equal(t *testing.T, exp0, exp1, exp2 float64, v *Vec3, delta float64) {
	assert.InDelta(t, exp0, v[0], delta)
	assert.InDelta(t, exp1, v[1], delta)
	assert.InDelta(t, exp2, v[2], delta)
}

func verifyTransform(t *testing.T, tr Transform, expTo0, expTo1, expTo2, from0, from1, from2, delta float64) {
	u := NewVec3(from0, from1, from2)
	assertVec3Equal(t, expTo0, expTo1, expTo2, tr.Apply(BlankVec3(), u), delta)
	// In place op.
	assert.Equal(t, u, tr.Apply(u, u))
	assertVec3Equal(t, expTo0, expTo1, expTo2, u, delta)
}
