package main

import (
	// "image/color"
	_ "image/jpeg"
	"log"
	"os"

	// "pix/pkg/colors"
	// "pix/pkg/quantize"

	"github.com/jessevdk/go-flags"
)

var (
	opts       Options
	pixopts    Pixels
	ditheropts Dither
	glitchopts Glitch
	filteropts Filters
	asciiopts  Ascii
	coloropts  Pally
	vhsopts    VHS
)

var parser = flags.NewParser(&opts, flags.Default)
var osargs []string // stand-in for os.Args to allow for an empty sub-command

var debug = func(string, ...interface{}) {}

func Pixxy(args []string) error {
	switch parser.Active.Name {
	case "glitch":
		return glitchopts.GlitchImage()
	case "dither":
		ditheropts.Threshold = 0.333 // set default
		return ditheropts.DitherImage()
	case "ascii":
		return asciiopts.RunAscii()
	case "color":
		return coloropts.GetColors()
	case "vhs":
		return vhsopts.Run()
	default:
		return nil
	}
}

func init() {
	p, err := parser.AddCommand("pixel", "modify image pixels", "", &pixopts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = parser.AddCommand("dither", "dither an image", "", &ditheropts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = parser.AddCommand("glitch", "glitch an image", "", &glitchopts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = parser.AddCommand("filter", "apply a set of filters to an image", "", &filteropts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = parser.AddCommand("ascii", "turn an image into ascii characters", "", &asciiopts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = parser.AddCommand("color", "turn an image into ascii characters", "", &coloropts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = parser.AddCommand("vhs", "turn an image into ascii characters", "", &vhsopts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = parser.AddCommand("version", "print version and debugging info", "print version and debugging info", &opts)
	if err != nil {
		log.Fatal(err)
	}

	p.Aliases = []string{"pix", "p"}

	// if len(os.Args) == 1 {
	// 	osargs = append(osargs, "run")
	// } else {
	// 	osargs = os.Args[1:]
	// }
}

// * lets just lower the scope to turning an arbitrary image into a wallpaper *//
func main() {
	args, err := parser.Parse()
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		} else {
			log.Fatal(err)
		}
	}

	if opts.Verbose {
		debug = log.Printf
	}

	if err := Pixxy(args); err != nil {
		log.Fatal(err)
	}
}
