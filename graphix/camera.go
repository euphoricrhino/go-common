package graphix

import "math"

// Camera is a simple wrapper of a view transform, a projector and a screen.
// This series of transforms takes a point in world coordinates and converts it to a 2D pixel coordinates
// on the screen, while keeping the z-depth info.
type Camera struct {
	vt Transform
	pr Projector
	sc *Screen
}

// NewCamera creates a new Camera given view transform, projector and screen.
func NewCamera(vt Transform, pr Projector, sc *Screen) *Camera {
	return &Camera{
		vt: vt,
		pr: pr,
		sc: sc,
	}
}

func (cam *Camera) ViewTransform() Transform { return cam.vt }
func (cam *Camera) Projector() Projector     { return cam.pr }
func (cam *Camera) Screen() *Screen          { return cam.sc }

// CameraOrbit represents an orbiting camera with a finite number of positions defined by the position index.
type CameraOrbit interface {
	// NumPositions returns the total number of camera positions along the orbit.
	NumPositions() int
	// GetCamera gets the camera at the ith position along the orbit.
	// It is up to the implementation to define the mapping of i and the camera configuration.
	GetCamera(i int) *Camera
}

// A simple implementation of CameraOrbit. The cameras are positioned along a circular
// orbit centered at the origin with a configurable normal vector.
type circularCameraOrbit struct {
	// The unit normal of the circular orbit.
	n *Vec3
	// The camera's initial configuration.
	pos     *Vec3
	forward *Vec3
	up      *Vec3

	numPositions int
	angleOffset  float64
	angleInc     float64

	pr Projector
	sc *Screen
}

var _ CameraOrbit = (*circularCameraOrbit)(nil)

// NewCircularCameraOrbit creates a circular camera orbit with normal n, and numPositions positions equally distributed
// along the circular orbit. At position i, a rotation transform Ri will be applied to the view transform defined
// by pos/forward/up, where Ri is a rotation around n by angle θ which is equal to angleOffset+i*2π/numPositions.
// Caller is responsible for passing in arguments satisfying the following requirements:
// - n, forward, up must be normalized;
// - n and pos must be orthogonal to each other;
// - up and forward must be orthogonal to each other.
func NewCircularCameraOrbit(
	n *Vec3,
	pos *Vec3,
	forward *Vec3,
	up *Vec3,
	numPositions int,
	angleOffset float64,
	pr Projector,
	sc *Screen,
) CameraOrbit {
	return &circularCameraOrbit{
		n:            NewCopyVec3(n),
		pos:          NewCopyVec3(pos),
		forward:      NewCopyVec3(forward),
		up:           NewCopyVec3(up),
		numPositions: numPositions,
		angleOffset:  angleOffset,
		angleInc:     2 * math.Pi / float64(numPositions),
		pr:           pr,
		sc:           sc,
	}
}

func (cir *circularCameraOrbit) NumPositions() int { return cir.numPositions }

func (cir *circularCameraOrbit) GetCamera(i int) *Camera {
	i = i % cir.numPositions
	theta := cir.angleOffset + cir.angleInc*float64(i)
	rot := NewAxisAngleRotation(cir.n, theta)
	vt := NewViewTransform(
		rot.Apply(BlankVec3(), cir.pos),
		rot.Apply(BlankVec3(), cir.forward),
		rot.Apply(BlankVec3(), cir.up),
	)
	return NewCamera(vt, cir.pr, cir.sc)
}
