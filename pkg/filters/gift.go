package filters

import (
	"image"
	"image/color"

	"github.com/disintegration/gift"
	"pix/pkg/imaging"
)

type Anchor int

// Anchor point positions.
const (
	CenterAnchor Anchor = iota
	TopLeftAnchor
	TopAnchor
	TopRightAnchor
	LeftAnchor
	RightAnchor
	BottomLeftAnchor
	BottomAnchor
	BottomRightAnchor
)

// Apply a GIFT filter to src image and return a new image
func Filter(filter gift.Filter, src image.Image) image.Image {
	g := gift.New(filter)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}

func Crop(img image.Image, rect image.Rectangle, anchor gift.Anchor) image.Image {
	filter := gift.Crop(rect)
	return Filter(filter, img)
}

func CropToSize(img image.Image, x, y int, anchor gift.Anchor) image.Image {
	filter := gift.CropToSize(x, y, anchor)
	return Filter(filter, img)
}

func ResizeToFill(img image.Image, x, y int) image.Image {
	rf := gift.ResizeToFill(x, y, gift.LanczosResampling, gift.Anchor(0))
	return Filter(rf, img)
}

func Blur(img image.Image, size float32) image.Image {
	rf := gift.GaussianBlur(size)
	return Filter(rf, img)
}

// adjust hue -100 to 100
func Hue(img image.Image, shift float64) image.Image {
	return imaging.AdjustHue(img, shift)
}

func Gamma(img image.Image, shift float64) image.Image {
	return imaging.AdjustGamma(img, shift)
}

func Contrast(img image.Image, pct float64) image.Image {
	return imaging.AdjustContrast(img, pct)
}

func Sharpen(img image.Image, pct float64) image.Image {
	return imaging.Sharpen(img, pct)
}

func OverlayCenter(bg, fg image.Image, opacity float64) image.Image {
	return imaging.OverlayCenter(bg, fg, opacity)
}

func CalculateLuminance(r, g, b uint32) float64 {
	// taken from here
	// https://stackoverflow.com/questions/596216/formula-to-determine-perceived-brightness-of-rgb-color
	r, g, b = r>>8, g>>8, b>>8

	// Get a brightness value of the colour from
	// here: https://www.w3.org/TR/AERT/#color-contrast
	//
	// This method isn't super-accurate since the
	// standards used are dated but this is negligible
	// in terms of the final image produced
	bright := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return bright
}

// turns grey
func Greyscale(img image.Image) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	// total := float32(width * height)
	// averageLuminance := float32(0)
	out := image.NewRGBA(img.Bounds())

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			l := CalculateLuminance(r, g, b)
			// averageLuminance += float32(l / float64(total))
			newclr := color.RGBA{uint8(l), uint8(l), uint8(l), uint8(a)}
			out.Set(x, y, newclr)
		}
	}
	return out
}
