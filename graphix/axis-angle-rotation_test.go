package graphix

import (
	"math"
	"testing"
)

func TestAxisAngleRotation(t *testing.T) {
	aar := NewAxisAngleRotation(NewVec3(0, 1, 0), math.Pi/4)
	verifyTransform(t, aar, 0, 1, 0, 0, 1, 0, 1e-8)
	verifyTransform(t, aar, 1/math.Sqrt(2), 0, 1/math.Sqrt(2), 0, 0, 1, 1e-8)
	verifyTransform(t, aar, 1/math.Sqrt(2), 0, -1/math.Sqrt(2), 1, 0, 0, 1e-8)

	n := NewVec3(1, 1, 1)
	aar = NewAxisAngleRotation(n.Normalize(n), 2*math.Pi/3)
	verifyTransform(t, aar, 0, 0, 1, 0, 1, 0, 1e-8)
	verifyTransform(t, aar, 1, 0, 0, 0, 0, 1, 1e-8)
	verifyTransform(t, aar, 0, 1, 0, 1, 0, 0, 1e-8)
}
