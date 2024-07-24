package main

import (
	"image"
	"image/color"
)

func toNRGBA(img image.Image) *image.NRGBA {
	if img, ok := img.(*image.NRGBA); ok {
		return &image.NRGBA{
			Pix:    img.Pix,
			Stride: img.Stride,
			Rect:   img.Rect.Sub(img.Rect.Min),
		}
	}
	return Clone(img)
}

// Clone returns a copy of the given image.
func Clone(img image.Image) *image.NRGBA {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	dst := image.NewNRGBA(img.Bounds())

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			clr := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			dst.Set(x, y, clr)
		}
	}
	return dst
}
