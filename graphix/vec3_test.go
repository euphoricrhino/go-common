package graphix

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVec3Clear(t *testing.T) {
	v := NewVec3(1, 2, 3)
	assert.Same(t, v, v.Clear())
	assertVec3Equal(t, 0, 0, 0, v, 1e-8)
}
func TestVec3Copy(t *testing.T) {
	v := NewVec3(1, 2, 3)
	u := NewVec3(4, 5, 6)
	assert.Same(t, v, v.Copy(u))
	assertVec3Equal(t, 4, 5, 6, v, 1e-8)
}

func TestVec3Add(t *testing.T) {
	v := BlankVec3()

	u := NewVec3(4, 5, 6)
	w := NewVec3(1, 2, 3)
	assert.Same(t, v, v.Add(u, w))
	assertVec3Equal(t, 5, 7, 9, v, 1e-8)

	// In place op.
	assert.Same(t, u, u.Add(u, w))
	assertVec3Equal(t, 5, 7, 9, u, 1e-8)
}

func TestVec3Sub(t *testing.T) {
	v := BlankVec3()

	u := NewVec3(4, 5, 6)
	w := NewVec3(1, 2, 3)
	assert.Same(t, v, v.Sub(u, w))
	assertVec3Equal(t, 3, 3, 3, v, 1e-8)

	// In place op.
	assert.Same(t, u, u.Sub(u, w))
	assertVec3Equal(t, 3, 3, 3, u, 1e-8)
}

func TestVec3Scale(t *testing.T) {
	v := BlankVec3()

	u := NewVec3(4, 5, 6)
	assert.Same(t, v, v.Scale(u, 2))
	assertVec3Equal(t, 8, 10, 12, v, 1e-8)

	// In place op.
	assert.Same(t, u, u.Scale(u, 2))
	assertVec3Equal(t, 8, 10, 12, u, 1e-8)
}

func TestVec3Dot(t *testing.T) {
	v := NewVec3(4, 5, 6)
	u := NewVec3(1, 2, 3)
	assert.Equal(t, 32.0, v.Dot(u))
}

func TestVec3Cross(t *testing.T) {
	v := BlankVec3()
	u := NewVec3(4, 5, 6)
	w := NewVec3(1, 2, 3)
	assert.Same(t, v, v.Cross(u, w))
	assertVec3Equal(t, 3, -6, 3, v, 1e-8)

	// In place op.
	assert.Same(t, u, u.Cross(u, w))
	assertVec3Equal(t, 3, -6, 3, u, 1e-8)
}

func TestVec3Norm(t *testing.T) {
	v := NewVec3(4, 5, 6)
	assert.Equal(t, math.Sqrt(4*4+5*5+6*6), v.Norm())
}

func TestVec3Normalize(t *testing.T) {
	u := NewVec3(4, 5, 6)
	r := math.Sqrt(4*4 + 5*5 + 6*6)
	v := BlankVec3()
	assert.Same(t, v, v.Normalize(u))
	assertVec3Equal(t, 4/r, 5/r, 6/r, v, 1e-8)

	// In place op.
	assert.Same(t, u, u.Normalize(u))
	assertVec3Equal(t, 4/r, 5/r, 6/r, u, 1e-8)
}
