package pixlib

import (
	"image/color"
	"math"
)

// max 16 bit integer used by color.Color.RGBA()
const MAXC = (1 << 16) - 1
const MAX8 = (1 << 8) - 1

type Pixel struct {
	R, G, B, A float64
}

// convert 0-65353 color range to 0-1 float
func (p *Pixel) Floatize(clr color.Color) Pixel {
	r, g, b, a := clr.RGBA()
	r, g, b, a = r>>8, g>>8, b>>8, a>>8
	rf, gf, bf, af := float64(r)/255.0, float64(g)/255, float64(b)/255, float64(a)/255
	return Pixel{rf, gf, bf, af}
}

// convert 0-1 RGBA color to 0-255 color
func (p *Pixel) To8Bit() (r uint32, g uint32, b uint32, a uint32) {
	linearize := func(f float64) uint32 { return uint32(math.Round((f * MAX8))) }
	r = linearize(p.R)
	g = linearize(p.G)
	b = linearize(p.B)
	a = linearize(p.A)
	return
}

// convert 0-1 RGBA color to 0-65353 color
func (p *Pixel) To16Bit() (r uint32, g uint32, b uint32, a uint32) {
	linearize := func(f float64) uint32 { return uint32(math.Round((f * MAXC))) }
	r = linearize(p.R)
	g = linearize(p.G)
	b = linearize(p.B)
	a = linearize(p.A)
	return
}

// gets pixels grey value aka luminance or light intensity
// fast but less accurate WE3 spec function
func (p *Pixel) Luma() float64 {
	return p.R*0.299 + p.G*0.587 + p.B*0.114
}

// gets pixels grey value aka luminance or light intensity
// fast but less accurate WE3 spec function
func Luma(r, g, b float64) float64 {
	return r*0.299 + g*0.587 + b*0.114
}

// satisfy color.Color interface
func (p *Pixel) RGBA() (r uint32, g uint32, b uint32, a uint32) {
	return p.To16Bit()
}

// returns a 0-360 degrees HUE
func (p *Pixel) GetHUE() float64 {
	r, g, b := p.R, p.G, p.B
	largest := maxNum(r, g, b)
	smallest := minNum(r, g, b)

	var hue float64
	if r >= g && r >= b {
		hue = (g - b) / (largest - smallest)
	} else if g >= r && g >= b {
		hue = 2.0 + (b-r)/(largest-smallest)
	} else if b >= r && b >= g {
		hue = 4.0 + (r-g)/(largest-smallest)
	}

	hue *= 60
	if hue < 0 {
		hue += 360
	}

	return hue
}
