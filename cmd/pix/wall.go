package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"pix/pkg/filters"

	"github.com/disintegration/gift"
)

func RGBtoAnsi(rgb color.RGBA) (foreground, background string) {
	cstr := strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m"
	foreground = "\x1b[38;2;" + cstr
	background = "\x1b[48;2;" + cstr
	return
}

func ColorToAnsi(clr color.Color) (foreground, background string) {
	c := color.RGBA{}
	r, g, b, a := clr.RGBA()
	c.R, c.G, c.B, c.A = uint8(r), uint8(g), uint8(b), uint8(a)

	return RGBtoAnsi(c)
}

// overlay an image as a square over a 1920x1080 color bg picked from the image
func wallpaperOverlay(img image.Image) image.Image {
	xy := img.Bounds().Dx()

	filter := gift.CropToSize(xy, xy, gift.LeftAnchor)
	g := gift.New(filter)

	newimg := image.NewNRGBA(g.Bounds(img.Bounds()))
	g.Draw(newimg, img)

	pal := filters.GetColorPalette(newimg, 3)

	for _, p := range pal {
		_, bg := ColorToAnsi(p)
		fmt.Printf("%s  %s", bg, "\x1b[0m")
	}

	width := 1920
	height := 1080

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	ximg := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < ximg.Bounds().Dx(); x++ {
		for y := 0; y < ximg.Bounds().Dy(); y++ {
			color := pal[3]
			ximg.Set(x, y, color)
		}
	}

	x1 := ximg.Bounds().Min.X
	y1 := ximg.Bounds().Min.Y
	w1 := ximg.Bounds().Dx()
	h1 := ximg.Bounds().Dy()

	w2 := newimg.Bounds().Dx()
	h2 := newimg.Bounds().Dy()

	x := -(x1 + ((w1 - w2) / 2))
	y := -(y1 + ((h1 - h2) / 2))

	// Create a blank canvas of the same size as the large image
	canvas := image.NewRGBA(image.Rect(0, 0, ximg.Bounds().Dx(), ximg.Bounds().Dy()))
	// Draw the large image onto the canvas
	draw.Draw(canvas, canvas.Bounds(), ximg, image.Point{}, draw.Src)
	// Overlay the small image onto the canvas at the calculated position
	draw.Draw(canvas, canvas.Bounds(), newimg, image.Point{x, y}, draw.Over)
	return canvas
}
