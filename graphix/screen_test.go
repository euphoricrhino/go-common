package graphix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreen(t *testing.T) {
	sc := NewScreen(200, 400, 0.1, 0.7, 0.5, 0.9)
	assert.Equal(t, 200, sc.Width())
	assert.Equal(t, 400, sc.Height())

	p := BlankProjection()
	q := NewProjection(.3, .8, 0)
	assert.Same(t, p, sc.Map(p, q))
	assertProjectionEqual(t, 100, 200, 0, p, 1e-8)

	// In place op.
	assert.Same(t, q, sc.Map(q, q))
	assertProjectionEqual(t, 100, 200, 0, q, 1e-8)

	q = NewProjection(.2, .85, 0)
	assert.Same(t, p, sc.Map(p, q))
	assertProjectionEqual(t, 50, 100, 0, p, 1e-8)
	// In place op.
	assert.Same(t, q, sc.Map(q, q))
	assertProjectionEqual(t, 50, 100, 0, q, 1e-8)
}
