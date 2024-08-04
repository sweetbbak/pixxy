package pixlib

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Polyline is a slice of polygons forming a line.
type Polyline []Polygon

// DrawPolygon filled or as line polygons if the thickness is >= 1.
func DrawPolygon(dst draw.Image, p Polygon, thickness float64, c color.Color) {
	n := len(p)

	if n < 3 {
		return
	}

	switch {
	case thickness < 1:
		p.Fill(dst, c)
	default:
		for i := 0; i < n; i++ {
			if i+1 == n {
				polylineFromTo(p[n-1], p[0], thickness).Fill(dst, c)
			} else {
				polylineFromTo(p[i], p[i+1], thickness).Fill(dst, c)
			}
		}
	}
}

// DrawPolyline draws a polyline with the given color and thickness.
func DrawPolyline(dst draw.Image, pl Polyline, thickness float64, c color.Color) {
	for _, p := range pl {
		DrawPolygon(dst, p, thickness, c)
	}
}

// Polygon is represented by a list of vectors.
type Polygon []Vec2

// Bounds return the bounds of the polygon rectangle.
func (p Polygon) Bounds() image.Rectangle {
	return p.Rect().Bounds()
}

// Rect is the polygon rectangle.
func (p Polygon) Rect() Rect {
	r := R(math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)

	for _, u := range p {
		x, y := u.XY()

		if x > r.Max.X {
			r.Max.X = x
		}

		if y > r.Max.Y {
			r.Max.Y = y
		}

		if x < r.Min.X {
			r.Min.X = x
		}

		if y < r.Min.Y {
			r.Min.Y = y
		}
	}

	return r
}

// Project creates a new Polygon with all vertexes projected through the given Matrix.
func (p Polygon) Project(m Matrix) Polygon {
	pp := make(Polygon, len(p))

	for i, u := range p {
		pp[i] = m.Project(u)
	}

	return pp
}

// EachPixel calls the provided function for each pixel
// in the polygon rectangle bounds.
func (p Polygon) EachPixel(m image.Image, fn func(x, y int)) {
	if len(p) < 3 {
		return
	}

	b := p.Bounds()

	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if IV(x, y).In(p) {
				fn(x, y)
			}
		}
	}
}

// Fill polygon on the image with the given color.
func (p Polygon) Fill(dst draw.Image, c color.Color) (drawCount int) {
	if len(p) < 3 {
		return
	}

	b := p.Bounds()

	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if IV(x, y).In(p) {
				Mix(dst, x, y, c)
				drawCount++
			}
		}
	}

	return drawCount
}

// Outline draws an outline of the polygon on dst.
func (p Polygon) Outline(dst draw.Image, thickness float64, c color.Color) {
	for i := 1; i < len(p); i++ {
		DrawLine(dst, p[i-1], p[i], thickness, c)
	}
}

// In returns true if the vector is inside the given polygon.
func (u Vec2) In(p Polygon) bool {
	if len(p) < 3 {
		return false
	}

	a := p[0]

	in := rayIntersectsSegment(u, p[len(p)-1], a)

	for _, b := range p[1:] {
		if rayIntersectsSegment(u, a, b) {
			in = !in
		}

		a = b
	}

	return in
}

// Points are a list of points.
type Points []image.Point

// Polygon based on the points.
func (pts Points) Polygon() Polygon {
	var p Polygon

	for i := range pts {
		p = append(p, PointToVec(pts[i]))
	}

	return p
}

// Segment intersect expression from
// https://www.ecse.rpi.edu/Homepages/wrf/Research/Short_Notes/pnpoly.html
//
// Currently the compiler inlines the function by default.
func rayIntersectsSegment(u, a, b Vec2) bool {
	return (a.Y > u.Y) != (b.Y > u.Y) &&
		u.X < (b.X-a.X)*(u.Y-a.Y)/(b.Y-a.Y)+a.X
}

// NewPolyline constructs a slice of line polygons.
func NewPolyline(p Polygon, t float64) Polyline {
	l := len(p)

	if l < 2 {
		return []Polygon{}
	}

	var pl Polyline

	for i := range p[:l-1] {
		pl = append(pl, newLinePolygon(p[i], p[i+1], t))
	}

	return pl
}

func polylineFromTo(from, to Vec2, t float64) Polygon {
	return NewPolyline(Polygon{from, to}, t)[0]
}

func newLinePolygon(from, to Vec2, t float64) Polygon {
	a := from.To(to).Angle()

	return Polygon{
		V(from.X+t*math.Cos(a+math.Pi/2), from.Y+t*math.Sin(a+math.Pi/2)),
		V(from.X+t*math.Cos(a-math.Pi/2), from.Y+t*math.Sin(a-math.Pi/2)),

		V(to.X+t*math.Cos(a-math.Pi/2), to.Y+t*math.Sin(a-math.Pi/2)),
		V(to.X+t*math.Cos(a+math.Pi/2), to.Y+t*math.Sin(a+math.Pi/2)),
	}
}
