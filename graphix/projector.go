package graphix

// Projection defines a projected Vec3. Its [0] and [1] are the two coordinates (e.g., x and y)
// in the projected plane, while its [2] keeps the z-distance of the original Vec3 with respect to the camera
// (orthographic or perspective).
type Projection [3]float64

func BlankProjection() *Projection                 { return NewProjection(0, 0, 0) }
func NewProjection(v0, v1, v2 float64) *Projection { return &Projection{v0, v1, v2} }

// Projector defines an interface projecting a Vec3 into a Projection.
type Projector interface {
	Project(p *Projection, v *Vec3) *Projection
}

// Orthographic defines an orthographic projector with respect to the canonical camera position, i.e.,
// the camera is positioned at origin, forward is -z, up is +y.
type Orthographic struct{}

var _ Projector = (*Orthographic)(nil)

func (*Orthographic) Project(p *Projection, v *Vec3) *Projection {
	// Use -z as distance since camera is looking at the -z direction.
	p[0], p[1], p[2] = v[0], v[1], -v[2]
	return p
}

// Perspective defines a perspective projector.
type Perspective struct {
	d float64
}

var _ Projector = (*Perspective)(nil)

// NewPerspective returns a perspective projector with respect to the canonical camera position, i.e.,
// the camera is positioned at origin, forward is -z, up is +y.
// E.g., dist=2 means the projection plane is 2 units in front of the camera (located at z=-2).
func NewPerspective(dist float64) *Perspective {
	return &Perspective{d: dist}
}

func (per *Perspective) Project(p *Projection, v *Vec3) *Projection {
	ratio := -per.d / v[2]
	p[0], p[1], p[2] = v[0]*ratio, v[1]*ratio, -v[2]
	return p
}
