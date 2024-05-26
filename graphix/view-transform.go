package graphix

// ViewTransform is defined with 3 parameters all of which are in world frame:
// - pos defines the coordinates of the camera's position;
// - forward defines the unit forward direction along which the camera is facing;
// - up defines the unit up direction of the camera.
//
// The transformation thus defined will transform pos into the origin, forward into the unit -z direction and up into the unit +y direction.
type viewTransform struct {
	pos *Vec3
	// Unit right/up/forward vector in world coordinates.
	ux *Vec3
	uy *Vec3
	uz *Vec3
}

var _ Transform = (*viewTransform)(nil)

// NewViewTransform creates a ViewTransform with camera position at pos, looking into the forward direction and pointing up to the up direction.
// It is the caller's responsibility to ensure that forward and up are mutually orthogonal and normalized.
func NewViewTransform(pos, forward, up *Vec3) Transform {
	return &viewTransform{
		pos: NewCopyVec3(pos),
		ux:  BlankVec3().Cross(forward, up),
		uy:  NewCopyVec3(up),
		uz:  BlankVec3().Scale(forward, -1),
	}
}

func (vt *viewTransform) Apply(v, u *Vec3) *Vec3 {
	v.Sub(u, vt.pos)
	v[0], v[1], v[2] = v.Dot(vt.ux), v.Dot(vt.uy), v.Dot(vt.uz)
	return v
}
