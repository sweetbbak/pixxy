package glitch

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"

	// dither2 "github.com/makeworld-the-better-one/dither/v2"
	"pix/pkg/glitch/dither"
	"pix/pkg/glitch/effects"
	"pix/pkg/glitch/utils"
	"pix/pkg/quantize"
)

var debug = func(string, ...interface{}) {}

type glitch_options struct {
	brightness   float64
	glitchFactor float64
	scanlines    bool
	seed         string
	frames       int
	gif          bool
	colors       []color.Color
	frameDelay   int
}

type GlitchOption func(args *glitch_options) error

func GlitchBrightness(c float64) GlitchOption {
	return func(args *glitch_options) error {
		args.brightness = c
		return nil
	}
}

func GlitchFactor(c float64) GlitchOption {
	return func(args *glitch_options) error {
		if c <= 0 || c >= 100 {
			return fmt.Errorf("glitch factor must be between 0 and 100")
		}
		args.glitchFactor = c
		return nil
	}
}

func GlitchUseScanlines(b bool) GlitchOption {
	return func(args *glitch_options) error {
		args.scanlines = b
		return nil
	}
}

func GlitchSeed(s string) GlitchOption {
	return func(args *glitch_options) error {
		args.seed = s
		return nil
	}
}

func GlitchPalette(c []color.Color) GlitchOption {
	return func(args *glitch_options) error {
		args.colors = c
		return nil
	}
}

func GlitchFrames(c int) GlitchOption {
	return func(args *glitch_options) error {
		if c <= 0 {
			return fmt.Errorf("frames must be a positive integer")
		}
		args.frames = c
		return nil
	}
}

func GlitchFrameDelay(c int) GlitchOption {
	return func(args *glitch_options) error {
		if c <= 0 {
			return fmt.Errorf("frames must be a positive integer")
		}
		args.frameDelay = c
		return nil
	}
}

// generate a random seed from a str value
func randseed(seed string) int64 {
	hasher := md5.New()
	hasher.Write([]byte(seed))
	hash := hasher.Sum(nil)

	length := len(hash)
	var seedInt int64

	for i, hashByte := range hash {
		// get byte shift offset
		shift := uint64((length - i - length) * 8)
		// binary OR the shifted byte onto the return value
		seedInt |= int64(hashByte) << shift
	}
	return seedInt
}

// glitch an image with the defualt options
func Glitch(srcImg image.Image) (image.Image, error) {
	return GlitchWithOpts(srcImg)
}

// glitch an image with the defualt options
func GlitchGif(srcImg image.Image, writer io.Writer, opts ...GlitchOption) (*gif.GIF, error) {
	seed, err := os.Hostname()
	if err != nil {
		seed = "unknown_333"
	}

	rand.New(rand.NewSource(randseed(seed)))

	defaultOpts := &glitch_options{
		brightness:   5.0,
		glitchFactor: 5.0,
		scanlines:    true,
		seed:         seed,
		frames:       7,
		frameDelay:   0,
	}

	for _, setter := range opts {
		if setter == nil {
			return nil, fmt.Errorf("option supplied is nil")
		}

		err := setter(defaultOpts)
		if err != nil {
			return nil, err
		}
	}

	output, err := defaultOpts.GlitchImage(srcImg)
	if err != nil {
		return nil, err
	}

	outputgif := &gif.GIF{}
	bounds := srcImg.Bounds()

	pal := color.Palette{}
	var palettedImage *image.Paletted

	if len(defaultOpts.colors) > 0 {
		for _, c := range defaultOpts.colors {
			pal = append(pal, c)
		}

		palettedImage = image.NewPaletted(bounds, pal)
	} else {
		pal = palette.Plan9[:256]
		palettedImage = image.NewPaletted(bounds, palette.Plan9[:256])
	}

	draw.FloydSteinberg.Draw(palettedImage, bounds, output, image.Pt(0, 0))

	outputgif.Image = append(outputgif.Image, palettedImage)
	outputgif.Delay = append(outputgif.Delay, defaultOpts.frameDelay)
	defaultOpts.frames--

	for {
		palettedImage = image.NewPaletted(bounds, pal)
		draw.FloydSteinberg.Draw(palettedImage, bounds, output, image.Pt(0, 0))

		// Add new frame to animated GIF
		outputgif.Image = append(outputgif.Image, palettedImage)
		outputgif.Delay = append(outputgif.Delay, defaultOpts.frameDelay)

		defaultOpts.frames--

		if defaultOpts.frames == 0 {
			break
		}

		output, err = defaultOpts.GlitchImage(srcImg)
	}

	if writer != nil {
		err = gif.EncodeAll(writer, outputgif)
		return outputgif, nil
	}

	return outputgif, nil
}

// glitch an image with the specified options
func GlitchWithOpts(srcImg image.Image, opts ...GlitchOption) (image.Image, error) {
	seed, err := os.Hostname()
	if err != nil {
		seed = "unknown_333"
	}

	rand.New(rand.NewSource(randseed(seed)))

	defaultOpts := &glitch_options{
		brightness:   5.0,
		glitchFactor: 5.0,
		scanlines:    true,
		seed:         seed,
		frames:       0,
	}

	for _, setter := range opts {
		if setter == nil {
			return nil, fmt.Errorf("option supplied is nil")
		}

		err := setter(defaultOpts)
		if err != nil {
			return nil, err
		}
	}

	return defaultOpts.GlitchImage(srcImg)
}

func (g *glitch_options) GlitchImage(img image.Image) (image.Image, error) {
	bounds := img.Bounds()

	input := image.NewRGBA(bounds)
	draw.Draw(input, bounds, img, bounds.Min, draw.Src)

	output := image.NewRGBA(bounds)
	draw.Draw(output, bounds, img, bounds.Min, draw.Src)

	if len(g.colors) > 0 {
		imgq := quantize.ApplyQuantization(input, g.colors)
		draw.Draw(input, bounds, imgq, bounds.Min, draw.Src)
	}

	glitchify(input, output, bounds, g.glitchFactor)
	effects.ApplyBrightness(output, g.brightness)

	if g.scanlines {
		effects.ApplyScanlines(output)
	}

	return output, nil
}

func GetColorPalette(img image.Image, level int) []color.Color {
	pal := []color.Color{}
	colors := quantize.Palette(img, level)
	for _, c := range colors {
		pal = append(pal, c)
	}
	return pal
}

func glitchify(input, output *image.RGBA, bounds image.Rectangle, glitchFactor float64) {
	copyInput := image.NewRGBA(bounds)
	copy(copyInput.Pix, input.Pix)

	eightBitted := image.NewRGBA(bounds)
	copy(eightBitted.Pix, input.Pix)
	dither.EightBit(eightBitted, utils.Random(0, 255))

	atkinsons := image.NewRGBA(bounds)
	copy(atkinsons.Pix, input.Pix)
	dither.Atkinsons(atkinsons, uint8(utils.Random(0, 255)))

	bayer := image.NewRGBA(bounds)
	copy(bayer.Pix, input.Pix)
	dither.Bayer(bayer)

	halftone := image.NewRGBA(bounds)
	copy(halftone.Pix, input.Pix)
	dither.Halftone(halftone, uint16(utils.Random(0, 255)))

	floydsteinberg := image.NewRGBA(bounds)
	copy(floydsteinberg.Pix, input.Pix)
	dither.FloydSteinberg(floydsteinberg, uint8(utils.Random(0, 255)))

	redOnly := image.NewRGBA(bounds)
	effects.CopyChannel(redOnly, input, utils.Red)

	greenOnly := image.NewRGBA(bounds)
	effects.CopyChannel(greenOnly, input, utils.Green)

	blueOnly := image.NewRGBA(bounds)
	effects.CopyChannel(blueOnly, input, utils.Blue)

	alphaMask := image.NewAlpha(bounds)
	for i := range alphaMask.Pix {
		alphaMask.Pix[i] = input.Pix[i*4]
	}

	srcs := []*image.RGBA{
		eightBitted,
		halftone,
		redOnly,
		greenOnly,
		blueOnly,
		copyInput,
	}
	srcNames := []string{
		"8bit",
		"halftone",
		"red",
		"green",
		"blue",
		"original",
	}

	wrapSlice := func(in, out *image.RGBA, op draw.Op) {
		width, height := bounds.Max.X, bounds.Max.Y
		maxOffset := int(glitchFactor / 100.0 * float64(width))

		// Random image slice offsetting
		for i := 0.0; i < glitchFactor; i++ {
			startY := utils.Random(0, height)
			chunkHeight := int(math.Min(float64(height-startY), float64(utils.Random(1, int(float64(height/2)*glitchFactor/100.0)))))
			offset := utils.Random(-maxOffset, maxOffset)
			effects.WrapSlice(out, in, offset, startY, chunkHeight, alphaMask, op)
		}
	}

	transforms := []func(in, out *image.RGBA){
		func(in, out *image.RGBA) {
			newIn := image.NewRGBA(bounds)
			copy(newIn.Pix, in.Pix)
			dither.Atkinsons(newIn, uint8(utils.Random(64, 192)))
			for i := range alphaMask.Pix {
				alphaMask.Pix[i] = newIn.Pix[i*4]
			}
			wrapSlice(newIn, out, draw.Over)
		},
		func(in, out *image.RGBA) {
			newIn := image.NewRGBA(bounds)
			copy(newIn.Pix, in.Pix)
			dither.EightBit(newIn, utils.Random(64, 192))
			for i := range alphaMask.Pix {
				alphaMask.Pix[i] = newIn.Pix[i*4]
			}
			wrapSlice(newIn, out, draw.Over)
		},
		func(in, out *image.RGBA) {
			newIn := image.NewRGBA(bounds)
			copy(newIn.Pix, in.Pix)
			dither.Bayer(newIn)
			for i := range alphaMask.Pix {
				alphaMask.Pix[i] = newIn.Pix[i*4]
			}
			wrapSlice(newIn, out, draw.Over)
		},
		func(in, out *image.RGBA) {
			newIn := image.NewRGBA(bounds)
			copy(newIn.Pix, in.Pix)
			dither.Halftone(newIn, uint16(utils.Random(64, 192)))
			for i := range alphaMask.Pix {
				alphaMask.Pix[i] = newIn.Pix[i*4]
			}
			wrapSlice(newIn, out, draw.Over)
		},
		func(in, out *image.RGBA) {
			newIn := image.NewRGBA(bounds)
			copy(newIn.Pix, in.Pix)
			dither.FloydSteinberg(newIn, uint8(utils.Random(64, 192)))
			for i := range alphaMask.Pix {
				alphaMask.Pix[i] = newIn.Pix[i*4]
			}
			wrapSlice(newIn, out, draw.Over)
		},
		func(in, out *image.RGBA) { wrapSlice(in, out, draw.Over) },
		func(in, out *image.RGBA) { wrapSlice(in, out, draw.Src) },
		func(in, out *image.RGBA) { effects.CopyChannel(out, in, utils.Red) },
		func(in, out *image.RGBA) { effects.CopyChannel(out, in, utils.Green) },
		func(in, out *image.RGBA) { effects.CopyChannel(out, in, utils.Blue) },
		func(in, out *image.RGBA) {
			for i := range alphaMask.Pix {
				alphaMask.Pix[i] = in.Pix[i*4]
			}
		},
	}
	transformNames := []string{
		"atkinsons",
		"8bit",
		"bayer",
		"halftone",
		"floydsteinberg",
		"wrapOver",
		"wrapSrc",
		"copyRed",
		"copyGreen",
		"copyBlue",
		"copyAlpha",
	}

	i := len(transforms)
	for i > 0 {
		destIdx := utils.Random(0, len(srcs))
		srcIdx := utils.Random(0, len(srcs))
		fIdx := utils.Random(0, len(transforms))
		transforms[fIdx](srcs[srcIdx], srcs[destIdx])
		debug("transform[%v] %v -> %v\n", transformNames[fIdx], srcNames[srcIdx], srcNames[destIdx])
		destIdx = utils.Random(0, len(srcs))
		fIdx = utils.Random(0, len(transforms))
		transforms[fIdx](input, srcs[destIdx])

		i--
	}

	for i, src := range srcs {
		debug("transform[wrapOver] %v -> output\n", srcNames[i])
		wrapSlice(src, output, draw.Over)
	}

	debug("reset alpha mask")
	for i := range alphaMask.Pix {
		alphaMask.Pix[i] = 255
	}

	finalOutput := image.NewRGBA(bounds)
	copy(finalOutput.Pix, output.Pix)
	debug("imageglitcher for final output")
	imageglitcher(finalOutput, output, bounds, glitchFactor)
}

// The imageglitcher algorithm from airtight interactive
func imageglitcher(inputData, outputData *image.RGBA, bounds image.Rectangle, glitchFactor float64) {
	width, height := bounds.Max.X, bounds.Max.Y
	maxOffset := int(glitchFactor / 100.0 * float64(width))
	mask := image.NewUniform(color.Alpha{A: 255})

	// Random image slice offsetting
	for i := 0.0; i < glitchFactor*2; i++ {
		startY := utils.Random(0, height)
		chunkHeight := int(math.Min(float64(height-startY), float64(utils.Random(1, height/4))))
		offset := utils.Random(-maxOffset, maxOffset)

		effects.WrapSlice(outputData, inputData, offset, startY, chunkHeight, mask, draw.Src)
	}

	// Copy a random channel from the pristene original input data onto the slice-offsetted output data
	effects.CopyChannel(outputData, inputData, utils.RandomChannel())
}
