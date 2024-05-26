package graphix

import "math"

// axisAngleRotation represents a rotation around an axis by an angle.
type axisAngleRotation struct {
	mat [3]Vec3
}

var _ Transform = (*axisAngleRotation)(nil)

// NewAxisAngleRotation creates a transform that rotates a point around axis n by theta.
// Caller is responsible for passing in a normalized n.
func NewAxisAngleRotation(n *Vec3, theta float64) Transform {
	ct, st := math.Cos(theta), math.Sin(theta)
	aar := &axisAngleRotation{}
	// See https://en.wikipedia.org/wiki/Rotation_matrix
	nxx, nyy, nzz := n[0]*n[0], n[1]*n[1], n[2]*n[2]
	nyz, nzx, nxy := n[1]*n[2], n[2]*n[0], n[0]*n[1]
	m := &aar.mat
	m[0][0] = ct + nxx*(1-ct)
	m[0][1] = nxy*(1-ct) - n[2]*st
	m[0][2] = nzx*(1-ct) + n[1]*st

	m[1][0] = nxy*(1-ct) + n[2]*st
	m[1][1] = ct + nyy*(1-ct)
	m[1][2] = nyz*(1-ct) - n[0]*st

	m[2][0] = nzx*(1-ct) - n[1]*st
	m[2][1] = nyz*(1-ct) + n[0]*st
	m[2][2] = ct + nzz*(1-ct)
	return aar
}

func (aar *axisAngleRotation) Apply(v, u *Vec3) *Vec3 {
	v[0], v[1], v[2] = aar.mat[0].Dot(u), aar.mat[1].Dot(u), aar.mat[2].Dot(u)
	return v
}
