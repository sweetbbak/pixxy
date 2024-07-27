package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"

	"pix/pkg/filters"
	"pix/pkg/glitch"
	dither2 "pix/pkg/glitch/dither"
	"pix/pkg/imaging"

	"github.com/makeworld-the-better-one/dither/v2"
)

func noOp() {
	_ = draw.Src
}

var ditherers = map[string]dither.ErrorDiffusionMatrix{
	"simple2d":            dither.Simple2D,
	"floyd":               dither.FloydSteinberg,
	"falsefloydsteinberg": dither.FalseFloydSteinberg,
	"jarvisjudiceninke":   dither.JarvisJudiceNinke,
	"atkinson":            dither.Atkinson,
	"stucki":              dither.Stucki,
	"burkes":              dither.Burkes,
	"sierra":              dither.Sierra,
	"sierra3":             dither.Sierra3,
	"tworowsierra":        dither.TwoRowSierra,
	"sierralite":          dither.SierraLite,
	"sierra2_4a":          dither.Sierra2_4A,
	"stevenpigeon":        dither.StevenPigeon,
}

var odmName = map[string]dither.OrderedDitherMatrix{
	"clustereddot4x4":            dither.ClusteredDot4x4,
	"clustereddotdiagonal8x8":    dither.ClusteredDotDiagonal8x8,
	"vertical5x3":                dither.Vertical5x3,
	"horizontal3x5":              dither.Horizontal3x5,
	"clustereddotdiagonal6x6":    dither.ClusteredDotDiagonal6x6,
	"clustereddotdiagonal8x8_2":  dither.ClusteredDotDiagonal8x8_2,
	"clustereddotdiagonal16x16":  dither.ClusteredDotDiagonal16x16,
	"clustereddot6x6":            dither.ClusteredDot6x6,
	"clustereddotspiral5x5":      dither.ClusteredDotSpiral5x5,
	"clustereddothorizontalline": dither.ClusteredDotHorizontalLine,
	"clustereddotverticalline":   dither.ClusteredDotVerticalLine,
	"clustereddot8x8":            dither.ClusteredDot8x8,
	"clustereddot6x6_2":          dither.ClusteredDot6x6_2,
	"clustereddot6x6_3":          dither.ClusteredDot6x6_3,
	"clustereddotdiagonal8x8_3":  dither.ClusteredDotDiagonal8x8_3,
}

func (d *Dither) DitherODM(img image.Image, op string, pal []color.Color) (image.Image, error) {
	matrix, ok := odmName[strings.ReplaceAll(strings.ToLower(op), "-", "_")]
	if !ok {
		return nil, fmt.Errorf("matrix type not found")
	}

	dx := dither.NewDitherer(pal)
	dx.Mapper = dither.PixelMapperFromMatrix(matrix, float32(d.Threshold))
	return dx.Dither(img), nil
}

func (d *Dither) DitherF() error {
	var pal color.Palette
	var img image.Image

	// open image file
	var inputfile string
	if d.Input != "" {
		inputfile = d.Input
	} else if d.Args.Image != "" {
		inputfile = d.Args.Image
	} else {
		return fmt.Errorf("no image supplied")
	}

	var err error
	img, err = openImage(inputfile)
	if err != nil {
		return err
	}

	if len(d.Palette) > 0 {
		c1, err := ParsePaletteString(strings.Join(d.Palette, " "))
		if err != nil {
			return err
		}
		pal = c1
	}

	if d.PaletteFile != "" {
		c2, err := ParsePalette(d.PaletteFile)
		if err != nil {
			return err
		}
		pal = append(pal, c2...)
	}

	// if no pallette, use image
	if d.ColorDepth > 0 {
		pal = glitch.GetColorPalette(img, d.ColorDepth)
	}

	bounds := img.Bounds()
	if d.Scale {
		img = imaging.Resize(img, bounds.Dx()/2, bounds.Dy()/2, imaging.Lanczos)
	}

	dx := dither.NewDitherer(pal)

	// rimg := imageToRGBA(img)
	// newImg := image.NewRGBA(img.Bounds())
	// copy(newImg.Pix, rimg.Pix)

	for _, input := range d.DitherType {
		dt, ok := ditherers[strings.ReplaceAll(strings.ToLower(input), "-", "_")]
		if !ok {
			return fmt.Errorf("ditherer not recognized: %v", input)
		}

		fmt.Printf("running %v\n", input)
		dx.Matrix = dt
		// dx.Serpentine = true
		img = dx.Dither(img)
		dx.Matrix = nil
	}

	// dx2 := dither.NewDitherer(pal)
	for _, input := range d.ODM {
		matrix, ok := odmName[strings.ReplaceAll(strings.ToLower(input), "-", "_")]
		if !ok {
			return fmt.Errorf("matrix type not found: %v", input)
		}

		dx.Mapper = dither.PixelMapperFromMatrix(matrix, float32(d.Threshold))
		img = dx.Dither(img)
	}

	if d.Scale {
		img = imaging.Resize(img, bounds.Dx(), bounds.Dy(), imaging.Lanczos)
	}

	if d.Bayer {
		img = filters.DitherBayer(img, float32(d.Threshold), pal)
	}

	if d.Halftone {
		ximg := imageToRGBA(img)
		dither2.Halftone(ximg, uint16(d.Threshold*100))
		img = ximg
	}

	if d.Output != "" {
		f, err := os.Create(d.Output)
		if err != nil {
			return err
		}
		defer f.Close()

		SaveImageToPNG(img, d.Output)
	}
	return nil
}

func (d *Dither) DitherImage() error {
	var inputfile string
	if d.Input != "" {
		inputfile = d.Input
	} else if d.Args.Image != "" {
		inputfile = d.Args.Image
	} else {
		return fmt.Errorf("no image supplied")
	}

	var (
		err error
		img image.Image
	)

	img, err = openImage(inputfile)
	if err != nil {
		return err
	}

	var rimg *image.RGBA
	rimg = imageToRGBA(img)

	bounds := rimg.Bounds()

	halftone := image.NewRGBA(bounds)
	copy(halftone.Pix, rimg.Pix)
	if d.Halftone {
		dither2.Halftone(halftone, uint16(d.Threshold))
	}

	if d.Halftone {
		filters.DitherBayer(rimg, float32(d.Threshold), filters.GetColorPalette(rimg, 3))
	}

	eightBitted := image.NewRGBA(bounds)
	copy(eightBitted.Pix, rimg.Pix)
	dither2.EightBit(eightBitted, int(192))

	atkinsons := image.NewRGBA(bounds)
	copy(atkinsons.Pix, rimg.Pix)
	output, _ := dither2.AtkinsonsGrey(atkinsons)

	// alphaMask := image.NewAlpha(bounds)
	// copy alpha from src
	// for i := range alphaMask.Pix {
	// alphaMask.Pix[i] = rimg.Pix[i*4]
	// }

	// output := image.NewRGBA(bounds)
	// copy(output.Pix, rimg.Pix)
	// dither.Atkinsons(output, 122)
	// for i := range alphaMask.Pix {
	// alphaMask.Pix[i] = output.Pix[i*4]
	// }

	// draw.DrawMask(output, bounds, atkinsons, bounds.Min, alphaMask, bounds.Min, draw.Over)
	// draw.DrawMask(atkinsons, bounds, rimg, bounds.Min, alphaMask, bounds.Min, draw.Over)

	// newIn := image.NewRGBA(bounds)
	// copy(newIn.Pix, atkinsons.Pix)
	// dither.Atkinsons(newIn, 192)

	// for i := range alphaMask.Pix {
	// 	alphaMask.Pix[i] = newIn.Pix[i*4]
	// }

	bayer := image.NewRGBA(bounds)
	copy(bayer.Pix, rimg.Pix)
	dither2.Bayer(bayer)

	floydsteinberg := image.NewRGBA(bounds)
	copy(floydsteinberg.Pix, rimg.Pix)
	dither2.FloydSteinberg(floydsteinberg, uint8(d.Threshold))

	if d.Output != "" {
		f, err := os.Create(d.Output)
		if err != nil {
			return err
		}
		defer f.Close()

		SaveImageToPNG(halftone, d.Output)
		SaveImageToPNG(eightBitted, "8bit.png")
		SaveImageToPNG(atkinsons, "atk.png")
		SaveImageToPNG(output, "atk2.png")
		SaveImageToPNG(bayer, "bayer.png")
		SaveImageToPNG(floydsteinberg, "floydsteinberg.png")
	}
	return nil
}
