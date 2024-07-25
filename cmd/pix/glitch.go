package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
	"path"
	"strings"

	"pix/pkg/colors"
	"pix/pkg/glitch"
)

func (g *Glitch) GlitchImage() error {
	var pal color.Palette
	var img image.Image

	cparser := colors.NewParser()

	// open image file
	var inputfile string
	if g.Input != "" {
		inputfile = g.Input
	} else if g.Args.Image != "" {
		inputfile = g.Args.Image
	} else {
		return fmt.Errorf("no image supplied")
	}

	var err error
	img, err = openImage(inputfile)
	if err != nil {
		return err
	}

	if g.Verbose {
		glitch.GlitchSetDebug(true)
	}

	if len(g.Palette) > 0 {
		err := cparser.ParseString(strings.Join(g.Palette, " "))
		if err != nil {
			return err
		}
		pal = cparser.Colors
	}

	if g.PaletteFile != "" {
		f, err := os.Open(g.PaletteFile)
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

	// if no pallette, use image
	if g.ColorDepth > 0 {
		pal = glitch.GetColorPalette(img, g.ColorDepth)
	}

	var oppys []glitch.GlitchOption

	if len(pal) > 0 {
		oppys = append(oppys, glitch.GlitchPalette(pal))
	}

	if g.Seed != "" {
		oppys = append(oppys, glitch.GlitchSeed(g.Seed))
	}

	if g.Factor > 0 {
		oppys = append(oppys, glitch.GlitchFactor(g.Factor))
	}

	if g.FrameCount > 0 {
		oppys = append(oppys, glitch.GlitchFrames(g.FrameCount))
	}

	if g.FrameDelay > 0 {
		oppys = append(oppys, glitch.GlitchFrameDelay(g.FrameDelay))
	}

	var outname string
	if g.Output != "" {
		outname = g.Output
	}

	if g.Gif {
		if outname == "" {
			outname = "output.gif"
		}

		if path.Ext(outname) != ".gif" {
			outname = strings.TrimSuffix(outname, path.Ext(outname)) + ".gif"
		}

		outfile, err := os.Create(outname)
		if err != nil {
			return err
		}
		defer outfile.Close()

		_, err = glitch.GlitchGif(
			img,
			outfile,
			oppys...,
		)
		if err != nil {
			return err
		}

	} else {
		if outname == "" {
			outname = "output.png"
		}

		out, err := glitch.GlitchWithOpts(
			img,
			oppys...,
		)
		if err != nil {
			return err
		}

		SaveImageToPNG(out, outname)
	}

	return nil
}
