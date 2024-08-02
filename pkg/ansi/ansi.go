package ansi

import (
	"image"
	"image/color"
	"strconv"

	"pix/pkg/quantize"
)

const CLEAR = "\x1b[0m"

type Ansi struct {
	Fg string
	Bg string
}

// get image color palette. depth for dithering is best at [1-3] and anything higher than [5-10] will return
// a lot of colors and have no real effect on the images color
func GetColorPalette(img image.Image, level int) []color.Color {
	pal := []color.Color{}
	colors := quantize.Palette(img, level)
	for _, c := range colors {
		pal = append(pal, c)
	}
	return pal
}

func ColorsToAnsi(colors []color.Color) []Ansi {
	ansicolors := make([]Ansi, len(colors))

	for i, clr := range colors {
		fg, bg := ColorToAnsi(clr)
		ansicolors[i].Fg = fg
		ansicolors[i].Bg = bg
	}
	return ansicolors
}

func HextoRGBA(hex string) color.RGBA {
	if hex[0:1] == "#" {
		hex = hex[1:]
	}

	r := hex[0:2]
	g := hex[2:4]
	b := hex[4:6]

	var A int64

	if len(hex) == 8 {
		a := hex[6:8]
		A, _ = strconv.ParseInt(a, 16, 0)
	} else {
		A = 255
	}

	R, _ := strconv.ParseInt(r, 16, 0)
	G, _ := strconv.ParseInt(g, 16, 0)
	B, _ := strconv.ParseInt(b, 16, 0)

	return color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)}

}

func HextoAnsi(hex string) (foreground, background string) {
	rgb := HextoRGBA(hex)
	cstr := strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m"
	foreground = "\x1b[38;2;" + cstr
	background = "\x1b[48;2;" + cstr
	return // implicit return
}

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
