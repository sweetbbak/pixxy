package pixlib

import (
	"fmt"
	"image"
	"math"
)

// ZV is a zero vector.
var ZV = Vec2{0, 0}

// a line / vector2 from x to y
type Vec2 struct {
	X, Y float64
}

// V returns a new 2D vector with the given coordinates.
func V(x, y float64) Vec2 {
	return Vec2{x, y}
}

// IV returns a new 2d vector based on the given int x, y pair.
func IV(x, y int) Vec2 {
	return Vec2{
		float64(x),
		float64(y),
	}
}

// PV returns a new 2D vector based on the given image.Point.
func PointToVec(p image.Point) Vec2 {
	return IV(p.X, p.Y)
}

// Unit returns a vector of length 1 facing the given angle.
func Unit(angle float64) Vec2 {
	return Vec2{1, 0}.Rotated(angle)
}

// String returns the string representation of the vector u.
//
//	u := gfx.V(4.5, -1.3)
//	u.String()     // returns "gfx.V(4.5, -1.3)"
//	fmt.Println(u) // gfx.V(4.5, -1.3)
func (u Vec2) String() string {
	return fmt.Sprintf("gfx.V(%.8f, %.8f)", u.X, u.Y)
}

// XY returns the components of the vector in two return values.
func (u Vec2) XY() (x, y float64) {
	return u.X, u.Y
}

// Eq checks the equality of two vectors.
func (u Vec2) Eq(v Vec2) bool {
	return u.X == v.X && u.Y == v.Y
}

// Add returns the sum of vectors u and v.
func (u Vec2) Add(v Vec2) Vec2 {
	return Vec2{
		u.X + v.X,
		u.Y + v.Y,
	}
}

// AddXY returns the sum of x and y added to v.
func (u Vec2) AddXY(x, y float64) Vec2 {
	return Vec2{
		u.X + x,
		u.Y + y,
	}
}

// Sub returns the difference betweeen vectors u and v.
func (u Vec2) Sub(v Vec2) Vec2 {
	return Vec2{
		u.X - v.X,
		u.Y - v.Y,
	}
}

// To returns the vector from u to v. Equivalent to v.Sub(u).
func (u Vec2) To(v Vec2) Vec2 {
	return Vec2{
		v.X - u.X,
		v.Y - u.Y,
	}
}

// Mod returns the floating-point remainder vector of x and y.
func (u Vec2) Mod(v Vec2) Vec2 {
	return Vec2{
		math.Mod(u.X, v.X),
		math.Mod(u.Y, v.Y),
	}
}

// Scaled returns the vector u multiplied by c.
func (u Vec2) Scaled(c float64) Vec2 {
	return Vec2{
		u.X * c,
		u.Y * c,
	}
}

// ScaledXY returns the vector u multiplied by the vector v component-wise.
func (u Vec2) ScaledXY(v Vec2) Vec2 {
	return Vec2{
		u.X * v.X,
		u.Y * v.Y,
	}
}

// Len returns the length of the vector u.
func (u Vec2) Len() float64 {
	return math.Hypot(u.X, u.Y)
}

// Angle returns the angle between the vector u and the x-axis. The result is in range [-Pi, Pi].
func (u Vec2) Angle() float64 {
	return math.Atan2(u.Y, u.X)
}

// Unit returns a vector of length 1 facing the direction of u (has the same angle).
func (u Vec2) Unit() Vec2 {
	if u.X == 0 && u.Y == 0 {
		return Vec2{1, 0}
	}

	return u.Scaled(1 / u.Len())
}

// Abs returns the absolute vector of the vector u.
func (u Vec2) Abs() Vec2 {
	return Vec2{
		math.Abs(u.X),
		math.Abs(u.Y),
	}
}

// Max returns the maximum vector of u and v.
func (u Vec2) Max(v Vec2) Vec2 {
	return Vec2{
		math.Max(u.X, v.X),
		math.Max(u.Y, v.Y),
	}
}

// Min returns the minimum vector of u and v.
func (u Vec2) Min(v Vec2) Vec2 {
	return Vec2{
		math.Min(u.X, v.X),
		math.Min(u.Y, v.Y),
	}
}

// Rotated returns the vector u rotated by the given angle in radians.
func (u Vec2) Rotated(angle float64) Vec2 {
	sin, cos := math.Sincos(angle)

	return Vec2{
		u.X*cos - u.Y*sin,
		u.X*sin + u.Y*cos,
	}
}

// Normal returns a vector normal to u. Equivalent to u.Rotated(math.Pi / 2), but faster.
func (u Vec2) Normal() Vec2 {
	return Vec2{
		-u.Y,
		u.X,
	}
}

// Dot returns the dot product of vectors u and v.
func (u Vec2) Dot(v Vec2) float64 {
	return u.X*v.X + u.Y*v.Y
}

// Cross return the cross product of vectors u and v.
func (u Vec2) Cross(v Vec2) float64 {
	return u.X*v.Y - v.X*u.Y
}

// Project returns a projection (or component) of vector u in the direction of vector v.
//
// Behaviour is undefined if v is a zero vector.
func (u Vec2) Project(v Vec2) Vec2 {
	len := u.Dot(v) / v.Len()

	return v.Unit().Scaled(len)
}

// Map applies the function f to both x and y components of the vector u and returns the modified
// vector.
//
//	u := gfx.V(10.5, -1.5)
//	v := u.Map(math.Floor)   // v is gfx.V(10, -2), both components of u floored
func (u Vec2) Map(f func(float64) float64) Vec2 {
	return Vec2{
		f(u.X),
		f(u.Y),
	}
}

// Vec3 converts the vector into a Vec3.
func (u Vec2) Vec3(z float64) Vec3 {
	return Vec3{u.X, u.Y, z}
}

// Pt returns the image.Point for the vector.
func (u Vec2) Pt() image.Point {
	return image.Pt(int(u.X), int(u.Y))
}

// R creates a new Rect for the vectors u and v.
//
// Note that the returned rectangle is not automatically normalized.
func (u Vec2) R(v Vec2) Rect {
	return NewRect(u, v)
}

// B creates a new image.Rectangle for the vectors2 u and v.
func (u Vec2) B(v Vec2) image.Rectangle {
	return u.R(v).Bounds()
}

// Rect constructs a Rect around the vector2 based on the provided Left, Top, Right, Bottom values.
func (u Vec2) Rect(l, t, r, b float64) Rect {
	return R(u.X+l, u.Y+t, u.X+r, u.Y+b)
}

// Bounds returns the bounds around the vector2 based on the provided Left, Top, Right, Bottom values.
func (u Vec2) Bounds(l, t, r, b float64) image.Rectangle {
	return u.Rect(l, t, r, b).Bounds()
}

// Lerp returns a linear interpolation between vectors2 u and v.
//
// This function basically returns a point along the line between a and b and t chooses which one.
// If t is 0, then a will be returned, if t is 1, b will be returned. Anything between 0 and 1 will
// return the appropriate point between a and b and so on.
func (u Vec2) Lerp(v Vec2, t float64) Vec2 {
	return u.Scaled(1 - t).Add(v.Scaled(t))
}

// Centroid returns the centroid O of three vectors2.
func Centroid(a, b, c Vec2) Vec2 {
	return V(
		(a.X+b.X+c.X)/3,
		(a.Y+b.Y+c.Y)/3,
	)
}
