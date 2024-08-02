package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"strings"

	"pix/pkg/ansi"
	"pix/pkg/colors"
	"pix/pkg/quantize"
)

func (p *Pally) GetColors() error {
	var pal color.Palette
	var img image.Image

	cparser := colors.NewParser()

	// open image file
	var inputfile string
	if p.Input != "" {
		inputfile = p.Input
	} else if p.Args.Image != "" {
		inputfile = p.Args.Image
	}

	if inputfile != "" {
		var err error
		img, err = openImage(inputfile)
		if err != nil {
			return err
		}

		// if no pallette, use image
		if p.ColorDepth > 0 {
			pal = ansi.GetColorPalette(img, p.ColorDepth)
		}
	}

	if len(p.Palette) > 0 {
		err := cparser.ParseString(strings.Join(p.Palette, " "))
		if err != nil {
			return err
		}

		pal = append(pal, cparser.Colors...)
	}

	if p.PaletteFile != "" {
		f, err := os.Open(p.PaletteFile)
		if err != nil {
			return err
		}
		defer f.Close()

		err = cparser.ParseFile(f)
		if err != nil {
			return err
		}

		pal = append(pal, cparser.Colors...)
	}

	if stdinOpen() {
		b, err := io.ReadAll(os.Stdin)
		line := string(b)

		err = cparser.ParseString(line)
		if err != nil {
			return err
		}

		pal = append(pal, cparser.Colors...)
	}

	pal = removeDuplicate(pal)

	if len(pal) == 0 {
		return fmt.Errorf("no colors were found")
	}

	for _, c := range pal {
		fg, bg := ansi.ColorToAnsi(c)
		r, g, b, _ := c.RGBA()
		fmt.Fprintf(os.Stderr, "%s  %s ", bg, ansi.CLEAR)

		fmt.Fprintf(os.Stdout, "#%x%x%x", r/255, g/255, b/255)

		if p.PrintAnsi {
			fmt.Fprintf(os.Stderr, "fg: \\e%s bg: \\e%s", fg[1:], bg[1:])
		}
		println()
	}

	if p.ApplyColor || (p.Input != "" && p.Output != "") {
		if inputfile == "" {
			return fmt.Errorf("no image supplied to apply a color palette to")
		}

		var err error
		img, err = openImage(inputfile)
		if err != nil {
			return err
		}

		var outname string
		if p.Output != "" {
			outname = p.Output
		} else {
			outname = "output.png"
		}

		output := quantize.ApplyQuantization(img, pal)
		SaveImageToPNG(output, outname)
	}

	return nil
}
