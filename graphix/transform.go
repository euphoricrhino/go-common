package graphix

// Transform defines a 3D transformation.
type Transform interface {
	// Apply applies the transform to u and stores the result into v and returns v.
	// Implementation must correctly handle the situation where u and v are the same pointer.
	Apply(v, u *Vec3) *Vec3
}

type TransformFunc func(v, u *Vec3) *Vec3

var _ Transform = (TransformFunc)(nil)

func (tf TransformFunc) Apply(v, u *Vec3) *Vec3 { return tf(v, u) }

// IdentityTransform returns an identity transform of a vector.
func IdentityTransform() Transform {
	return TransformFunc(func(v, u *Vec3) *Vec3 {
		return v.Copy(u)
	})
}
