package quantize

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func copyPalette(p []color.Color) []color.Color {
	ret := make([]color.Color, len(p))
	for i, c := range p {
		r, g, b, a := c.RGBA()
		ret[i] = color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
	}
	return ret
}

func copyOfImage(img image.Image) *image.RGBA {
	dst := image.NewRGBA(img.Bounds())
	copyImage(dst, img)
	return dst
}

func copyImage(dst draw.Image, src image.Image) {
	draw.Draw(dst, src.Bounds(), src, src.Bounds().Min, draw.Src)
}

// samePalette returns true if both palettes contain the same colors,
// regardless of order.
func samePalette(p1 []color.Color, p2 []color.Color) bool {
	if len(p1) != len(p2) {
		return false
	}

	// Modified from: https://stackoverflow.com/a/36000696/7361270

	diff := make(map[color.Color]int, len(p1))
	for _, x := range p1 {
		// 0 value for int is 0, so just increment a counter for the string
		diff[x]++
	}
	for _, y := range p2 {
		// If _y is not in diff bail out early
		if _, ok := diff[y]; !ok {
			return false
		}
		diff[y] -= 1
		if diff[y] == 0 {
			delete(diff, y)
		}
	}
	return len(diff) == 0
}

// unpremultAndLinearize unpremultiplies the provided color, and returns the
// linearized RGB values, as well as the unchanged alpha value.
func unpremultAndLinearize(c color.Color) (uint16, uint16, uint16, uint16) {
	// alpha
	var a uint16

	// Optimize for different color types
	// Opaque colors are fast-tracked
	// Non-premultiplied colors aren't unpremulted, and all others are
	switch v := c.(type) {
	case color.Gray:
		a = 0xffff
	case color.Gray16:
		a = 0xffff
	case color.NRGBA:
		// (1/255)*65535 = 257
		// This converts 8-bit color into 16-bit
		a = uint16(v.A) * 257
	case color.NRGBA64:
		a = v.A
	default:
		c = color.NRGBA64Model.Convert(c)
		_, _, _, x := c.RGBA()
		a = uint16(x)
	}

	r, g, b := toLinearRGB(c)
	return r, g, b, a
}

// premult takes the current position in the image and the dithered
// color for that position, and returns a color that's corrected to
// take into account the alpha value of the original image at that
// position -- premultipling it.
func premult(c color.RGBA64, x, y int, img image.Image) color.RGBA64 {
	// Algorithm described in #8
	// https://github.com/makeworld-the-better-one/dither/issues/8

	_, _, _, a := img.At(x, y).RGBA()
	if a == 0 {
		// Transparent, no color values are held
		return color.RGBA64{0, 0, 0, 0}
	}
	if a == 0xffff {
		// Pixel is opaque, no alpha math needed
		return c
	}
	// Multiply RGB by alpha value - return premultiplied color
	// Adapted from https://github.com/golang/go/blob/go1.16.4/src/image/color/color.go#L84
	r := uint32(c.R)
	r *= a
	r /= 0xffff
	g := uint32(c.G)
	g *= a
	g /= 0xffff
	b := uint32(c.B)
	b *= a
	b /= 0xffff

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

func sqDiff(v1 uint16, v2 uint16) uint32 {
	// This optimization is copied from Go stdlib, see
	// https://github.com/golang/go/blob/go1.15.7/src/image/color/color.go#L314

	d := uint32(v1) - uint32(v2)
	return (d * d) >> 2
}

// closestColor returns the index of the color in the palette that's closest to
// the provided one, using Euclidean distance in linear RGB space. The provided
// RGB values must be linear RGB.
func closestColor(r, g, b uint16, pal [][3]uint16) int {
	// Go through each color and find the closest one
	color, best := 0, uint32(math.MaxUint32)
	for i, c := range pal {

		// Euclidean distance, but the square root part is removed
		// Weight by luminance value to approximate radiant power / luminance
		// as humans perceive it.
		//
		// These values were taken from Wikipedia:
		// https://en.wikipedia.org/wiki/Grayscale#Colorimetric_(perceptual_luminance-preserving)_conversion_to_grayscale
		// 0.2126, 0.7152, 0.0722
		// The are changed to fractions here to keep everything in integer math:
		//     1063/5000, 447/625, 361/5000
		// Unfortunately this requires promoting them to uint64 to prevent overflow

		dist := uint32(
			1063*uint64(sqDiff(r, c[0]))/5000 +
				447*uint64(sqDiff(g, c[1]))/625 +
				361*uint64(sqDiff(b, c[2]))/5000,
		)

		if dist < best {
			if dist == 0 {
				return i
			}
			color, best = i, dist
		}
	}
	return color
}

type PixelMapper func(x, y int, r, g, b uint16) (uint16, uint16, uint16)

func Bayer(x, y uint, strength float32) PixelMapper {
	var matrix [][]uint

	if x == 0 || y == 0 {
		panic("dither: Bayer: neither x or y can be zero")
	}
	if x == 3 && y == 3 {
		matrix = [][]uint{
			{0, 5, 2},
			{3, 8, 7},
			{6, 1, 4},
		}
	} else if x == 5 && y == 3 {
		matrix = [][]uint{
			{0, 12, 7, 3, 9},
			{14, 8, 1, 5, 11},
			{6, 4, 10, 13, 2},
		}
	} else if x == 3 && y == 5 {
		matrix = [][]uint{
			{0, 14, 16},
			{12, 8, 4},
			{7, 1, 10},
			{3, 5, 13},
			{9, 11, 2},
		}
	} else if (x&(x-1)) == 0 && (y&(y-1)) == 0 {
		// Both are powers of two
		matrix = bayerMatrix(x, y)
	} else {
		// Neither are powers of two
		panic("dither: Bayer: dimensions aren't both a power of two")
	}

	// Create precalculated matrix
	scale := 65535.0 * strength
	max := x * y

	precalc := make([][]float32, y)
	for i := uint(0); i < y; i++ {
		precalc[i] = make([]float32, x)
		for j := uint(0); j < x; j++ {
			precalc[i][j] = convThresholdToAddition(scale, matrix[i][j], max)
		}
	}

	return PixelMapper(func(xx, yy int, r, g, b uint16) (uint16, uint16, uint16) {
		return RoundClamp(float32(r) + precalc[yy%int(y)][xx%int(x)]),
			RoundClamp(float32(g) + precalc[yy%int(y)][xx%int(x)]),
			RoundClamp(float32(b) + precalc[yy%int(y)][xx%int(x)])
	})
}

func RoundClamp(i float32) uint16 {
	if i < 0 {
		return 0
	}
	if i > 65535 {
		return 65535
	}
	return uint16(math.RoundToEven(float64(i)))
}

func log2(v uint) uint {
	// Sources:
	// https://graphics.stanford.edu/~seander/bithacks.html#IntegerLogObvious
	// https://stackoverflow.com/a/18139978/7361270

	var r uint
	v = v >> 1
	for v != 0 {
		r++
		v = v >> 1
	}
	return r
}

func bayerMatrix(xdim, ydim uint) [][]uint {
	// Bit math algorithm is used to calculate each cell of matrix individually.
	// This allows for easy generation of non-square matrices, as long as side
	// lengths are powers of two.
	//
	// Source for this bit math algorithm:
	// https://bisqwit.iki.fi/story/howto/dither/jy/#Appendix%202ThresholdMatrix
	//
	// The second code example on that part of the page is what this was based off
	// of, the one that works for rectangular matrices.
	//
	// The code was re-implemented exactly and tested to make sure the results
	// are the same. No algorithmic changes were made. The only code change was
	// to create a 2D slice to store and return results.

	M := log2(xdim)
	L := log2(ydim)

	matrix := make([][]uint, ydim)

	for y := uint(0); y < ydim; y++ {
		matrix[y] = make([]uint, xdim)
		for x := uint(0); x < xdim; x++ {

			var v, offset uint
			xmask := M
			ymask := L

			if M == 0 || (M > L && L != 0) {
				xc := x ^ ((y << M) >> L)
				yc := y
				for bit := uint(0); bit < M+L; {
					ymask--
					v |= ((yc >> ymask) & 1) << bit
					bit++
					for offset += M; offset >= L; offset -= L {
						xmask--
						v |= ((xc >> xmask) & 1) << bit
						bit++
					}
				}
			} else {
				xc := x
				yc := y ^ ((x << L) >> M)
				for bit := uint(0); bit < M+L; {
					xmask--
					v |= ((xc >> xmask) & 1) << bit
					bit++
					for offset += L; offset >= M; offset -= M {
						ymask--
						v |= ((yc >> ymask) & 1) << bit
						bit++
					}
				}
			}

			matrix[y][x] = v
		}
	}
	return matrix
}

// convThresholdToAddition takes a value from a matrix usually used for thresholding,
// and returns a value that can be added to a color instead of thresholded.
//
// scale is the number that's multiplied at the end, usually you want this to be
// 65535 to scale to match the color value range. value is the cell of the matrix.
// max is the divisor of the cell value, usually this is the product of the matrix
// dimensions.
func convThresholdToAddition(scale float32, value uint, max uint) float32 {
	// See:
	// https://en.wikipedia.org/wiki/Ordered_dithering
	// https://en.wikipedia.org/wiki/Talk:Ordered_dithering#Sources

	// 0.50000006 is next possible float32 value after 0.5. This is to correct
	// a rounding error that occurs when the number is exactly 0.5, which results
	// in pure black being dithered when it should be left alone.
	return scale * (float32(value+1.0)/float32(max) - 0.50000006)
}

func ApplyQuantization(src image.Image, pal []color.Color) image.Image {
	palette := copyPalette(pal)
	linearPalette := make([][3]uint16, len(palette))
	for i := range linearPalette {
		r, g, b := toLinearRGB(palette[i])
		linearPalette[i] = [3]uint16{r, g, b}
	}

	var img draw.Image

	if pi, ok := src.(*image.Paletted); ok {
		if !samePalette(palette, pi.Palette) {
			// Can't use this because it will change image colors
			// Instead make a copy, and return that later
			img = copyOfImage(src)
		}
	} else if img, ok = src.(draw.Image); !ok {
		// Can't be changed
		// Instead make a copy and dither and return that
		img = copyOfImage(src)
	}

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := unpremultAndLinearize(c)
			if a == 0 {
				img.Set(x, y, c)
				continue
			}

			// Use PixelMapper -> find closest palette color -> get that color
			// -> cast to color.RGBA64
			// Comes from d.palette so this cast will always work
			// strength := float32(1.0)
			// pmfunc := Bayer(uint(3), uint(3), strength)
			// a1, a2, a3 := pmfunc(x, y, r, g, b)
			// idx := closestColor(a1, a2, a3, linearPalette)
			idx := closestColor(r, g, b, linearPalette)
			outColor := premult(palette[idx].(color.RGBA64), x, y, img)

			img.Set(x, y, outColor)
		}
	}
	return img
}

func ApplyBayerDither(src image.Image, pal []color.Color, strength float32) image.Image {
	palette := copyPalette(pal)
	linearPalette := make([][3]uint16, len(palette))
	for i := range linearPalette {
		r, g, b := toLinearRGB(palette[i])
		linearPalette[i] = [3]uint16{r, g, b}
	}

	var img draw.Image

	if pi, ok := src.(*image.Paletted); ok {
		if !samePalette(palette, pi.Palette) {
			// Can't use this because it will change image colors
			// Instead make a copy, and return that later
			img = copyOfImage(src)
		}
	} else if img, ok = src.(draw.Image); !ok {
		// Can't be changed
		// Instead make a copy and dither and return that
		img = copyOfImage(src)
	}

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := unpremultAndLinearize(c)
			if a == 0 {
				img.Set(x, y, c)
				continue
			}

			// Use PixelMapper -> find closest palette color -> get that color
			// -> cast to color.RGBA64
			// Comes from d.palette so this cast will always work
			pmfunc := Bayer(uint(3), uint(3), strength)
			a1, a2, a3 := pmfunc(x, y, r, g, b)
			idx := closestColor(a1, a2, a3, linearPalette)
			outColor := premult(palette[idx].(color.RGBA64), x, y, img)

			img.Set(x, y, outColor)
		}
	}
	return img
}
