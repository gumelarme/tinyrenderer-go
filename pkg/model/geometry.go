package model

import "math"

type Vec2 struct {
	X, Y int
}

type Vec2f struct {
	X, Y float32
}

func (v *Vec2) Set(index, value int) {
	switch index {
	case 0:
		v.X = value
	case 1:
		v.Y = value
	default:
		panic("Index is out of bound")

	}
}

func (v Vec2) Get(index int) int {
	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	default:
		panic("Index is out of bound")

	}

}

func (v Vec2) Add(x Vec2) Vec2 {
	return Vec2{X: v.X + x.X, Y: v.Y + x.Y}
}

func (v Vec2) Subtract(x Vec2) Vec2 {
	return Vec2{X: v.X - x.X, Y: v.Y - x.Y}
}

func (v Vec2) MultiplyFloat32(f float32) Vec2 {
	return Vec2{X: int(float32(v.X) * f), Y: int(float32(v.Y) * f)}
}

func (v Vec2) ToVec2f() Vec2f {
	return Vec2f{float32(v.X), float32(v.Y)}
}

type Vec3f struct {
	X, Y, Z float64
}

func (v Vec3f) Subtract(x Vec3f) Vec3f {
	return Vec3f{X: v.X - x.X, Y: v.Y - x.Y, Z: v.Z - x.Z}
}

func (v Vec3f) MultiplyVec(rhs Vec3f) float64 {
	return v.X*rhs.X + v.Y*rhs.Y + v.Z*rhs.Z
}

func (v Vec3f) norm() float64 {
	norm2 := v.MultiplyVec(v)
	return math.Sqrt(float64(norm2))
}

func (v *Vec3f) Normalize() {
	norm := v.norm()
	v.X /= norm
	v.Y /= norm
	v.Z /= norm
}

func (v *Vec3f) Power(rhs Vec3f) Vec3f {
	return Vec3f{
		v.Y*rhs.Z - v.Z*rhs.Y,
		v.Z*rhs.X - v.X*rhs.Z,
		v.X*rhs.Y - v.Y*rhs.X,
	}
}

func Cross(v1, v2 Vec3f) Vec3f {
	return Vec3f{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
}
