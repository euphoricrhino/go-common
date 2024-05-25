package graphix

import "math"

// Vec3 represents a 3D vector.
type Vec3 [3]float64

func BlankVec3() *Vec3                 { return NewVec3(0, 0, 0) }
func NewVec3(v0, v1, v2 float64) *Vec3 { return &Vec3{v0, v1, v2} }
func NewCopyVec3(v *Vec3) *Vec3        { return &Vec3{v[0], v[1], v[2]} }

// Add adds u and w and stores the sum into v then returns v.
func (v *Vec3) Add(u, w *Vec3) *Vec3 {
	v[0] = u[0] + w[0]
	v[1] = u[1] + w[1]
	v[2] = u[2] + w[2]
	return v
}

// Sub subtracts w from u and stores the difference into v then returns v.
func (v *Vec3) Sub(u, w *Vec3) *Vec3 {
	v[0] = u[0] - w[0]
	v[1] = u[1] - w[1]
	v[2] = u[2] - w[2]
	return v
}

// Scale scales u by s and stores the result into v then returns v.
func (v *Vec3) Scale(u *Vec3, s float64) *Vec3 {
	v[0] = u[0] * s
	v[1] = u[1] * s
	v[2] = u[2] * s
	return v
}

// Dot returns the dot product of u and v.
func (v *Vec3) Dot(u *Vec3) float64 {
	return v[0]*u[0] + v[1]*u[1] + v[2]*u[2]
}

// Cross stores the cross product u×w into v then returns v.
func (v *Vec3) Cross(u, w *Vec3) *Vec3 {
	v[0], v[1], v[2] = u[1]*w[2]-w[1]*u[2], u[2]*w[0]-w[2]*u[0], u[0]*w[1]-w[0]*u[1]
	return v
}

// Norm returns the L2-norm of v.
func (v *Vec3) Norm() float64 { return math.Sqrt(v.Dot(v)) }

// Normalize normalizes v and returns v.
func (v *Vec3) Normalize() *Vec3 {
	return v.Scale(v, 1.0/v.Norm())
}
