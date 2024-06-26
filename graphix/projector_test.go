package graphix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrthographicProjector(t *testing.T) {
	o := NewOrthographic()
	assert.Equal(t, 0.0, o.NearZClip())
	v := NewVec3(3, 5, -2)
	p := BlankProjection()
	assert.Same(t, p, o.Project(p, v))
	assertProjectionEqual(t, 3, 5, 2, p, 1e-8)
}

func TestPerspectiveProjector(t *testing.T) {
	per := NewPerspective(2)
	assert.Equal(t, 0.4, per.NearZClip())

	v := NewVec3(16, 30, -8)
	p := BlankProjection()
	assert.Same(t, p, per.Project(p, v))
	assertProjectionEqual(t, 4, 7.5, 8, p, 1e-8)
}
