package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"pix/pkg/filters"
	"pix/pkg/glitch"
	dither2 "pix/pkg/glitch/dither"
	fx "pix/pkg/glitch/dither"
	"pix/pkg/imaging"

	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/sahilm/fuzzy"
)

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

var (
	DitherList []string
	MatrixList []string
)

func RandomDither() (string, dither.ErrorDiffusionMatrix, bool) {
	r := rand.New(rand.NewSource(1))
	min := 0
	max := len(DitherList)
	randInt := r.Intn(max - min)
	name := DitherList[randInt]
	d, ok := ditherers[name]
	return name, d, ok
}

func RandomMatrix() (string, dither.OrderedDitherMatrix, bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 0
	max := len(MatrixList) - 1
	randInt := r.Intn(max-min) + min
	name := MatrixList[randInt]
	m, ok := odmName[name]
	return name, m, ok
}

func getDlist() []string {
	var list []string
	for k := range ditherers {
		list = append(list, k)
	}
	return list
}

func getMlist() []string {
	var list []string
	for k := range odmName {
		list = append(list, k)
	}
	return list
}

func (d *Dither) ListMaps() {
	for item := range odmName {
		fmt.Fprintln(os.Stdout, item)
	}
}

func (d *Dither) ListDitherers() {
	for item := range ditherers {
		fmt.Fprintln(os.Stdout, item)
	}
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

	if d.Verbose {
		debug = log.Printf
	}

	// open image file
	var inputfile string
	if d.Input != "" {
		inputfile = string(d.Input)
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
		debug("using color palette from provided colors: %v", pal)
	}

	if d.PaletteFile != "" {
		c2, err := ParsePalette(d.PaletteFile)
		if err != nil {
			return err
		}
		pal = append(pal, c2...)
		debug("using color palette from file colors: %v", pal)
	}

	// userProvidedPallete := (len(d.Palette) < 0 && d.PaletteFile == "" && d.ColorDepth < 1)
	userProvidedPallete := (len(d.Palette) > 1 || d.PaletteFile != "")

	// if no pallette, use image
	if d.ColorDepth > 0 && !userProvidedPallete {
		pal = glitch.GetColorPalette(img, d.ColorDepth)
		debug("using color palette from quantization: %v", pal)
	}

	if len(pal) < 1 {
		return fmt.Errorf("pallette empty")
	}

	bounds := img.Bounds()
	if d.Scale {
		var sfact int
		if d.ScaleFactor > 0 {
			sfact = d.ScaleFactor
		} else {
			sfact = 2
		}
		img = imaging.Resize(img, bounds.Dx()/sfact, bounds.Dy()/sfact, imaging.NearestNeighbor)
	}

	dx := dither.NewDitherer(pal)

	for _, input := range d.DitherType {
		userInput := strings.ReplaceAll(strings.ToLower(input), "-", "_")
		var dt dither.ErrorDiffusionMatrix
		var name string
		DitherList = getDlist()

		if userInput == "rand" || userInput == "random" {
			var ok bool
			name, dt, ok = RandomDither()
			if !ok {
				return fmt.Errorf("idk what the fuck is happening sis: %v %v %v", name, d, ok)
			}
		} else {
			matches := fuzzy.Find(userInput, DitherList)
			debug("score: %v", matches[0].Score)
			name = matches[0].Str

			var ok bool
			dt, ok = ditherers[matches[0].Str]
			if !ok {
				return fmt.Errorf("ditherer not recognized: %v\naccepted values: %v", input, DitherList)
			}
		}

		fmt.Printf("running dither: %v\n", name)
		dx.Matrix = dt
		dx.Serpentine = true
		img = dx.Dither(img)
		dx.Matrix = nil
	}

	for _, input := range d.ODM {
		userInput := strings.ReplaceAll(strings.ToLower(input), "-", "_")

		var matrix dither.OrderedDitherMatrix
		var name string
		MatrixList = getMlist()

		if userInput == "rand" || userInput == "random" {
			var ok bool
			name, matrix, ok = RandomMatrix()
			if !ok {
				return fmt.Errorf("idk what the fuck is happening sis: %v %v %v", name, matrix, ok)
			}
		} else {
			var ok bool
			matches := fuzzy.Find(userInput, MatrixList)
			name = matches[0].Str

			matrix, ok = odmName[name]
			if !ok {
				return fmt.Errorf("matrix type not found: %v", input)
			}
		}

		fmt.Printf("running dither matrix: %v\n", name)
		dx.Mapper = dither.PixelMapperFromMatrix(matrix, float32(d.Threshold))
		img = dx.Dither(img)
	}

	if d.Bayer {
		img = filters.DitherBayer(img, float32(d.Threshold), pal)
	}

	if d.Halftone {
		ximg := imageToRGBA(img)
		dither2.Halftone(ximg, uint16(d.Threshold*255))
		img = ximg
	}

	if d.EightBit {
		ximg := imageToRGBA(img)
		fx.EightBit(ximg, int(d.Threshold*255))
		img = ximg
	}

	if d.Scale {
		img = imaging.Resize(img, bounds.Dx(), bounds.Dy(), imaging.NearestNeighbor)
	}

	if d.Output == "-" {
		WriteImageToStdout(img)
		return nil
	}

	if d.Output != "" {
		f, err := os.Create(string(d.Output))
		if err != nil {
			return err
		}
		defer f.Close()

		debug("saving image: %s", d.Output)
		SaveImageToPNG(img, string(d.Output))
	}

	return nil
}

func (d *Dither) DitherImage() error {
	if d.ListMatrices {
		d.ListMaps()
		return nil
	}

	if d.ListDithers {
		d.ListDitherers()
		return nil
	}

	return d.DitherF()
}
