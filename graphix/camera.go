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

// NewOrtho2DCamera returns a camera in canonical position (sitting at +z, looking at -z with y as up), given the screen mapping.
func NewOrtho2DCamera(sc *Screen) *Camera {
	return NewCamera(
		NewViewTransform(NewVec3(0, 0, 1), NewVec3(0, 0, -1), NewVec3(0, 1, 0)),
		NewOrthographic(),
		sc,
	)
}

// CameraOrbit represents a series of camera configurations.
type CameraOrbit interface {
	// Frames returns the total number of frames of camera orbit.
	Frames() int
	// GetCamera gets the camera at the fth frame along the orbit.
	// It is up to the implementation to define the mapping of f and the camera configuration.
	GetCamera(f int) *Camera
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

	frames      int
	angleOffset float64
	angleInc    float64

	pr Projector
	sc *Screen
}

var _ CameraOrbit = (*circularCameraOrbit)(nil)

// NewCircularCameraOrbit creates a circular camera orbit with normal n, and frames positions equally distributed
// along the circular orbit. At position i, a rotation transform Ri will be applied to the view transform defined
// by pos/forward/up, where Ri is a rotation around n by angle θ which is equal to angleOffset+i*2π/frames.
// Caller is responsible for passing in arguments satisfying the following requirements:
// - n, forward, up must be normalized;
// - n and pos must be orthogonal to each other;
// - up and forward must be orthogonal to each other.
func NewCircularCameraOrbit(
	n *Vec3,
	pos *Vec3,
	forward *Vec3,
	up *Vec3,
	frames int,
	angleOffset float64,
	pr Projector,
	sc *Screen,
) CameraOrbit {
	return &circularCameraOrbit{
		n:           NewCopyVec3(n),
		pos:         NewCopyVec3(pos),
		forward:     NewCopyVec3(forward),
		up:          NewCopyVec3(up),
		frames:      frames,
		angleOffset: angleOffset,
		angleInc:    2 * math.Pi / float64(frames),
		pr:          pr,
		sc:          sc,
	}
}

func (cir *circularCameraOrbit) Frames() int { return cir.frames }

func (cir *circularCameraOrbit) GetCamera(i int) *Camera {
	i = i % cir.frames
	theta := cir.angleOffset + cir.angleInc*float64(i)
	rot := NewAxisAngleRotation(cir.n, theta)
	vt := NewViewTransform(
		rot.Apply(BlankVec3(), cir.pos),
		rot.Apply(BlankVec3(), cir.forward),
		rot.Apply(BlankVec3(), cir.up),
	)
	return NewCamera(vt, cir.pr, cir.sc)
}

type stationaryCameraOrbit struct {
	cam    *Camera
	frames int
}

var _ CameraOrbit = (*stationaryCameraOrbit)(nil)

// NewStationaryCamera returns a CameraOrbit with only one position.
func NewStationaryCamera(cam *Camera, frames int) CameraOrbit {
	return &stationaryCameraOrbit{cam: cam, frames: frames}
}

func (st *stationaryCameraOrbit) Frames() int             { return st.frames }
func (st *stationaryCameraOrbit) GetCamera(i int) *Camera { return st.cam }
