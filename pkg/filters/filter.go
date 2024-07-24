package filters

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	_ "image/png" // For frame decoding
	"os"

	"pix/pkg/imaging"
	"pix/pkg/quantize"

	"github.com/anthonynsimon/bild/effect"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
)

// Get the color palette of an image. The higher the level the more colors are returned
func GetColorPalette(img image.Image, level int) []color.Color {
	pal := []color.Color{}
	colors := quantize.Palette(img, level)
	for _, c := range colors {
		pal = append(pal, c)
	}
	return pal
}

// Dither an image and adjust colors to match the given color palette
func DitherColor(img image.Image, colors []color.Color) image.Image {
	d := dither.NewDitherer(colors)
	d.Matrix = dither.FloydSteinberg
	return d.Dither(img)
}

// bayer dither
func DitherBayer(img image.Image, colors []color.Color) image.Image {
	d := dither.NewDitherer(colors)
	d.Mapper = dither.Bayer(8, 8, 1.0) // 8x8 Bayer matrix at 100% strength
	return d.Dither(img)
}

// Dither a gif / array of images
func DitherGIF(imgs []image.Image, palette []color.Color) (gif.GIF, error) {
	d := dither.NewDitherer(palette)
	d.Matrix = dither.FloydSteinberg // Why not?
	numFrames := len(imgs)

	if numFrames < 2 {
		return gif.GIF{}, fmt.Errorf("must have more than one frame to create a gif")
	}

	// Decode first frame and get image.Config for use in gif.GIF.
	// gif.GIF requires *image.Paletted is used, so DitherPaletted
	// is called instead of Dither.
	firstFrame, config := d.DitherPalettedConfig(imgs[0])
	frames := make([]*image.Paletted, numFrames)
	frames[0] = firstFrame

	// Decode other frames
	for i := 1; i < numFrames; i++ {
		frames[i] = d.DitherPaletted(imgs[i])
	}

	// Frame delay - same for each frame
	delays := make([]int, numFrames)
	for i := range delays {
		delays[i] = 7
	}

	// Setup GIF and encode
	g := gif.GIF{
		Image: frames,
		Delay: delays,
		// By specifying a Config, we can set a global color table for the GIF.
		// This is more efficient then each frame having its own color table, which
		// is the default when there's no config.
		Config: config,
	}
	return g, nil
}

// dither an array of images into a gif on the filesystem
func DitherGIFtoFile(imgs []image.Image, palette []color.Color, filename string) error {
	g, err := DitherGIF(imgs, palette)
	if err != nil {
		return err
	}

	f2, err := os.Create(filename)
	if err != nil {
		return err
	}

	err = gif.EncodeAll(f2, &g)
	if err != nil {
		return err
	}
	return nil
}

func SmartCrop(img image.Image, w, h int, resize bool) image.Image {
	width, height := getCropDimensions(img, w, h)
	resizer := nfnt.NewDefaultResizer()
	analyzer := smartcrop.NewAnalyzer(resizer)
	topCrop, _ := analyzer.FindBestCrop(img, 250, 250)

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	img = img.(SubImager).SubImage(topCrop)
	if resize && (img.Bounds().Dx() != width || img.Bounds().Dy() != height) {
		img = resizer.Resize(img, uint(width), uint(height))
	}

	return img
}

func getCropDimensions(img image.Image, width, height int) (int, int) {
	// if we don't have width or height set use the smaller image dimension as both width and height
	if width == 0 && height == 0 {
		bounds := img.Bounds()
		x := bounds.Dx()
		y := bounds.Dy()
		if x < y {
			width = x
			height = x
		} else {
			width = y
			height = y
		}
	}
	return width, height
}

// Bloom applies a bloom effect on the given image.
// Because of the nature of the effect, a larger image is returned.
// 10px padding is added to each side of the image, growing it by
// 20px on X and 20px on Y.
func Bloom(img image.Image) image.Image {
	// dilate the image to have a bigger source of light
	dilated := effect.Dilate(img, 5)

	// blur the image
	bloomed := imaging.Blur(dilated, 5.0)

	bounds := bloomed.Bounds()

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			//get one of r, g, b on the mask image ...
			r, g, b, a := bloomed.At(x, y).RGBA()
			//... and set it as the alpha value on the mask.
			//Assuming that white is your transparency, subtract it from 255
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			// a = a / 2
			if a < 0 {
				a = 0
			}
			if a > 255 {
				a = 255
			}
			bloomed.Set(x, y, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(float32(a) / 2)})
		}
	}

	canvas := image.NewRGBA(img.Bounds())
	// Draw the large image onto the canvas
	draw.Draw(canvas, canvas.Bounds(), img, image.Point{}, draw.Src)
	// Overlay the small image onto the canvas at the calculated position
	draw.Draw(canvas, canvas.Bounds(), bloomed, image.Point{0, 0}, draw.Over)

	return canvas
}

func Emboss(img image.Image) image.Image {
	return imaging.Convolve3x3(
		img,
		[9]float64{
			-1, -1, 0,
			-1, 1, 1,
			0, 1, 1,
		},
		nil,
	)
}

// translateImage copies the src image applying the given offset on a new Image
// bounds is the size of the resulting image.
func translateImage(src image.Image, bounds image.Rectangle, xOffset, yOffset int) image.Image {
	rv := image.NewRGBA(bounds)
	size := src.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			rv.Set(xOffset+x, yOffset+y, src.At(x, y))
		}
	}
	return rv
}
