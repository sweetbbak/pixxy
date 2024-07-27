package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

const MAXC = (1 << 16) - 1

type Channel int

const (
	Red Channel = iota
	Green
	Blue
	Alpha
)

func imageToRGBA(src image.Image) *image.RGBA {
	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func ConvertToRGBA(src image.Image) *image.RGBA {
	// Create a new RGBA image with the same dimensions as the source image
	dst := &image.RGBA{
		Pix:    make([]byte, src.Bounds().Dx()*4*src.Bounds().Dy()),
		Stride: src.Bounds().Dx() * 4,
		Rect:   src.Bounds(),
	}

	// Copy pixels from src to dst
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			oldColor := src.At(x, y)
			r, g, b, _ := oldColor.RGBA()

			// Set the RGBA value for the destination pixel
			dst.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
		}
	}

	return dst
}

func RandomChannel() Channel {
	r := rand.Float32()
	if r < 0.33 {
		return Green
	} else if r < 0.66 {
		return Red
	}
	return Blue
}

// turns a uint32 into its apprxoimate 0-1 value
func c(a uint32) uint8 {
	return uint8((float64(a) / MAXC) * 255)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid llength, must be 7 or 4")
	}
	return
}

func shiftColor(in color.Color, left int) (out color.RGBA) {
	var shiftedColor color.RGBA
	r, g, b, a := in.RGBA()

	shiftedColor = color.RGBA{
		R: uint8(b),
		G: uint8(r),
		B: uint8(g),
		A: uint8(a),
	}

	if left == 1 {
		shiftedColor = color.RGBA{
			R: uint8(g),
			G: uint8(b),
			B: uint8(r),
			A: uint8(a),
		}
	}

	return shiftedColor
}
