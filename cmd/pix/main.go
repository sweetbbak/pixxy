package main

import (
	"image/color"
	_ "image/jpeg"
	"log"
	"os"

	// "pix/pkg/filters"
	// "pix/pkg/glitch"
	"pix/pkg/colors"
	// "pix/pkg/filters"
	"pix/pkg/quantize"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	ParseFile string `short:"p" long:"parse" description:"file with hex colors to parse to use as a color palette"`
	Verbose   bool   `short:"v" long:"verbose" description:"print debugging information and verbose output"`
}

var Glitch struct {
	Gif    bool    `short:"g" long:"gif" description:"create a gif"`
	Amount float64 `short:"t" long:"threshold" description:"glitch threshold"`
}

var debug = func(string, ...interface{}) {}

func Pixxy(args []string) error {
	img, err := openImage(args[0])
	if err != nil {
		return err
	}

	var pal []color.Color
	if opts.ParseFile != "" {
		file, err := os.Open(opts.ParseFile)
		if err != nil {
			return err
		}
		defer file.Close()

		parser := colors.NewParser()
		err = parser.ParseFile(file)
		if err != nil {
			return err
		}
		pal = parser.Colors
		println(len(pal))
	} else {
		pal = quantize.GetColorPalette(img, 5)
	}

	output := quantize.ApplyQuantization(img, pal)
	// filters.DitherGIF()

	// f, err := os.Create("out.gif")
	// defer f.Close()

	// output1 := filters.Bloom(img)
	// output := quantize.ApplyBayerDither(img, pal, 0.5)
	SaveImageToPNG(output, "output.png")

	gifin := wallpaperOverlay(output)
	SaveImageToPNG(gifin, "wall.png")

	// _, err = glitch.GlitchGif(
	// 	gifin,
	// 	f,
	// 	glitch.GlitchPalette(pal),
	// 	glitch.GlitchSeed("sweet33"),
	// 	glitch.GlitchFactor(1.0),
	// 	glitch.GlitchFrameDelay(100),
	// 	glitch.GlitchFrames(10),
	// )
	// if err != nil {
	// 	return err
	// }

	// out, err := glitch.GlitchWithOpts(
	// 	img,
	// 	glitch.GlitchPalette(pal),
	// 	glitch.GlitchSeed("cumpy"),
	// 	glitch.GlitchFactor(1.0),
	// )
	// if err != nil {
	// 	return err
	// }

	// out3 := filters.Bloom(img)
	// SaveImageToPNG(out, "output.png")
	// SaveImageToPNG(out2, "pal.png")
	// SaveImageToPNG(out3, "bloom.png")

	return nil
}

// * lets just lower the scope to turning an arbitrary image into a wallpaper *//
func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	if opts.Verbose {
		debug = log.Printf
	}

	if err := Pixxy(args); err != nil {
		log.Fatal(err)
	}
}
