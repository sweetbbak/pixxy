package pixlib

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Draw draws src on dst, at the zero point using draw.Src.
func Draw(dst draw.Image, r image.Rectangle, src image.Image) {
	draw.Draw(dst, r, src, ZP, draw.Src)
}

// DrawColor draws an image.Rectangle of uniform color on dst.
func DrawColor(dst draw.Image, r image.Rectangle, c color.Color) {
	draw.Draw(dst, r, NewUniform(c), ZP, draw.Src)
}

// DrawColorOver draws an image.Rectangle of uniform color over dst.
func DrawColorOver(dst draw.Image, r image.Rectangle, c color.Color) {
	draw.Draw(dst, r, NewUniform(c), ZP, draw.Over)
}

// DrawSrc draws src on dst.
func DrawSrc(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	draw.Draw(dst, r, src, sp, draw.Src)
}

// DrawOver draws src over dst.
func DrawOver(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	draw.Draw(dst, r, src, sp, draw.Over)
}

// DrawPalettedImage draws a PalettedImage over a PalettedDrawImage.
// func DrawPalettedImage(dst PalettedDrawImage, r image.Rectangle, src PalettedImage) {
// 	w, h, m := r.Dx(), r.Dy(), r.Min
//
// 	for x := m.X; x != w; x++ {
// 		for y := m.Y; y != h; y++ {
// 			if src.AlphaAt(x, y) > 0 {
// 				dst.SetColorIndex(x, y, src.ColorIndexAt(x, y))
// 			}
// 		}
// 	}
// }
//
// // DrawPalettedLayer draws a *Layer over a *Paletted.
// // (slightly faster than using the generic DrawPalettedImage)
// func DrawPalettedLayer(dst *Paletted, r image.Rectangle, src *Layer) {
// 	w, h, m := r.Dx(), r.Dy(), r.Min
//
// 	for x := m.X; x != w; x++ {
// 		for y := m.Y; y != h; y++ {
// 			if src.AlphaAt(x, y) > 0 {
// 				dst.SetColorIndex(x, y, src.ColorIndexAt(x, y))
// 			}
// 		}
// 	}
// }

// DrawLine draws a line of the given color.
// A thickness of <= 1 is drawn using DrawBresenhamLine.
func DrawLine(dst draw.Image, from, to Vec2, thickness float64, c color.Color) {
	if thickness <= 1 {
		DrawLineBresenham(dst, from, to, c)
		return
	}

	polylineFromTo(from, to, thickness).Fill(dst, c)
}

// DrawCircle draws a circle with radius and thickness. (filled if thickness == 0)
func DrawCircle(dst draw.Image, u Vec2, radius, thickness float64, c color.Color) {
	if thickness == 0 {
		DrawCircleFilled(dst, u, radius, c)
		return
	}

	bounds := IR(int(u.X-radius), int(u.Y-radius), int(u.X+radius), int(u.Y+radius))

	EachPixel(dst.Bounds().Intersect(bounds), func(x, y int) {
		v := V(float64(x), float64(y))

		l := u.To(v).Len() + 0.5

		if l < radius && l > radius-thickness {
			Mix(dst, x, y, c)
		}
	})
}

// DrawCircleFilled draws a filled circle.
func DrawCircleFilled(dst draw.Image, u Vec2, radius float64, c color.Color) {
	bounds := IR(int(u.X-radius+1), int(u.Y-radius+1), int(u.X+radius+1), int(u.Y+radius+1))

	EachPixel(dst.Bounds().Intersect(bounds), func(x, y int) {
		v := V(float64(x), float64(y))

		if u.To(v).Len() < radius {
			Mix(dst, x, y, c)
		}
	})
}

// DrawCicleFast draws a (crude) filled circle.
func DrawCicleFast(dst draw.Image, u Vec2, radius float64, c color.Color) {
	ir := int(radius)
	r2 := ir * ir
	pt := u.Pt()

	for y := -ir; y <= ir; y++ {
		for x := -ir; x <= ir; x++ {
			if x*x+y*y <= r2 {
				SetPoint(dst, pt.Add(Pt(x, y)), c)
			}
		}
	}
}

// DrawLineBresenham draws a line using Bresenham's line algorithm.
//
// http://en.wikipedia.org/wiki/Bresenham's_line_algorithm
func DrawLineBresenham(dst draw.Image, from, to Vec2, c color.Color) {
	x0, y0 := from.XY()
	x1, y1 := to.XY()

	steep := math.Abs(y0-y1) > math.Abs(x0-x1)

	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}

	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	dx := x1 - x0
	dy := math.Abs(y1 - y0)
	e := dx / 2
	y := y0

	var ystep float64 = -1

	if y0 < y1 {
		ystep = 1
	}

	for x := x0; x <= x1; x++ {
		if steep {
			Mix(dst, int(y), int(x), c)
		} else {
			Mix(dst, int(x), int(y), c)
		}

		e -= dy

		if e < 0 {
			y += ystep
			e += dx
		}
	}
}
