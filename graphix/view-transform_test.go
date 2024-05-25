package graphix

import (
	"math"
	"testing"
)

func TestViewTransform(t *testing.T) {
	// View from (0,0,-1) into +z, keeping +y.
	vt := NewViewTransform(NewVec3(0, 0, -1), NewVec3(0, 0, 1), NewVec3(0, 1, 0))
	verifyTransform(t, vt, 0, 0, -2, 0, 0, 1, 1e-8)
	verifyTransform(t, vt, -1, 0, -1, 1, 0, 0, 1e-8)
	verifyTransform(t, vt, 0, 1, -1, 0, 1, 0, 1e-8)

	// View from (0,1,0) into -y with +x as up.
	vt = NewViewTransform(NewVec3(0, 1, 0), NewVec3(0, -1, 0), NewVec3(1, 0, 0))
	verifyTransform(t, vt, 0, 0, -1, 0, 0, 0, 1e-8)
	verifyTransform(t, vt, 0, 1, -1, 1, 0, 0, 1e-8)
	verifyTransform(t, vt, 1, 0, -1, 0, 0, 1, 1e-8)

	// View from (1,1,0) into +z with (1/√2,1/√2,0) as up.
	vt = NewViewTransform(NewVec3(1, 1, 0), NewVec3(0, 0, 1), NewVec3(1/math.Sqrt(2), 1/math.Sqrt(2), 0))
	verifyTransform(t, vt, 0, -math.Sqrt(2), -1, 0, 0, 1, 1e-8)
	verifyTransform(t, vt, -1/math.Sqrt(2), -1/math.Sqrt(2), 0, 1, 0, 0, 1e-8)
	verifyTransform(t, vt, 1/math.Sqrt(2), -1/math.Sqrt(2), 0, 0, 1, 0, 1e-8)
}
