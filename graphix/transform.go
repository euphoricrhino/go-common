package graphix

// Transform defines a 3D transformation.
type Transform interface {
	// Apply applies the transform to u and stores the result into v and returns v.
	// Implementation must correctly handle the situation where u and v are the same pointer.
	Apply(v, u *Vec3) *Vec3
}
