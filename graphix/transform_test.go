package graphix

import "testing"

func TestIdentityTransform(t *testing.T) {
	id := IdentityTransform()
	u := NewVec3(1, 2, 3)
	v := id.Apply(BlankVec3(), u)
	assertVec3Equal(t, 1, 2, 3, v, 1e-8)
}
