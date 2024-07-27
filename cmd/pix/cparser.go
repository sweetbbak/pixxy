package main

import (
	"image/color"
	"os"

	"pix/pkg/colors"
)

func ParsePaletteString(p string) (color.Palette, error) {
	cparser := colors.NewParser()
	var pal color.Palette

	err := cparser.ParseString(p)
	if err != nil {
		return nil, err
	}

	pal = append(pal, cparser.Colors...)
	return pal, nil
}

func ParsePalette(paletteFile string) (color.Palette, error) {
	cparser := colors.NewParser()
	var pal color.Palette

	f, err := os.Open(paletteFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = cparser.ParseFile(f)
	if err != nil {
		return nil, err
	}

	pal = append(pal, cparser.Colors...)
	return pal, nil
}
