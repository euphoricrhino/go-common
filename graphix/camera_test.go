package graphix

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrtho2DCamera(t *testing.T) {
	sc := &Screen{}
	cam := NewOrtho2DCamera(sc)
	assert.Same(t, sc, cam.Screen())
	v := NewVec3(0, 1, 2)
	assertVec3Equal(t, 0, 1, 1, cam.ViewTransform().Apply(BlankVec3(), v), 1e-8)
}

func TestCircularCameraOrbit(t *testing.T) {
	pr := NewOrthographic()
	sc := &Screen{}
	cir := NewCircularCameraOrbit(
		NewVec3(-1, 0, 0),
		NewVec3(0, 0, -1),
		NewVec3(0, 0, 1),
		NewVec3(1, 0, 0),
		4,
		math.Pi/2,
		pr,
		sc,
	)

	assert.Equal(t, 4, cir.Frames())

	// Test the camera position by verifying the view transform of v.
	v := NewVec3(1, 0, 1)
	cam := cir.GetCamera(0)
	assert.Same(t, pr, cam.Projector())
	assert.Same(t, sc, cam.Screen())
	assertVec3Equal(t, -1, 1, -1, cam.ViewTransform().Apply(BlankVec3(), v), 1e-8)

	cam = cir.GetCamera(1)
	assert.Same(t, pr, cam.Projector())
	assert.Same(t, sc, cam.Screen())
	assertVec3Equal(t, 0, 1, 0, cam.ViewTransform().Apply(BlankVec3(), v), 1e-8)

	cam = cir.GetCamera(2)
	assert.Same(t, pr, cam.Projector())
	assert.Same(t, sc, cam.Screen())
	assertVec3Equal(t, 1, 1, -1, cam.ViewTransform().Apply(BlankVec3(), v), 1e-8)

	cam = cir.GetCamera(3)
	assert.Same(t, pr, cam.Projector())
	assert.Same(t, sc, cam.Screen())
	assertVec3Equal(t, 0, 1, -2, cam.ViewTransform().Apply(BlankVec3(), v), 1e-8)

	// Should come back to position 0.
	cam = cir.GetCamera(4)
	assert.Same(t, pr, cam.Projector())
	assert.Same(t, sc, cam.Screen())
	assertVec3Equal(t, -1, 1, -1, cam.ViewTransform().Apply(BlankVec3(), v), 1e-8)
}

func TestStationaryCamera(t *testing.T) {
	cam := &Camera{}
	st := NewStationaryCamera(cam, 5)
	assert.Equal(t, 5, st.Frames())
	assert.Same(t, cam, st.GetCamera(15))
}
