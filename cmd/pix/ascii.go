package main

import (
	"context"
	"fmt"
	"image"
	"os"
	"strings"

	"pix/pkg/ascii"
	"pix/pkg/ascii/video"

	"golang.org/x/image/font/opentype"
)

func OpenFont(path string) (*opentype.Font, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (a *Ascii) ParseCharset() ascii.Charset {
	// a.ParseCharset()
	// limited|extended|block
	switch a.CharsetPreset {
	case "extended", "ext", "e", "x":
		return ascii.CharsetExtended
	case "block", "b":
		return ascii.CharsetBlock
	default:
		return ascii.CharsetLimited
	}
}

func (a *Ascii) OpenImage() (image.Image, error) {
	var inputfile string
	if a.Input != "" {
		inputfile = a.Input
	} else if a.Args.Image != "" {
		inputfile = a.Args.Image
	} else {
		return nil, fmt.Errorf("no image supplied")
	}

	file, err := openImage(inputfile)
	if err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func (a *Ascii) CreateVideo(opts []ascii.Option, args []string) error {
	err := video.Convert(context.Background(), a.Input, a.Output, opts, args...)
	if err != nil {
		return err
	}
	return nil
}

func (a *Ascii) RunAscii() error {
	var optSet []ascii.Option
	mem := &ascii.Memory{}

	if a.FontPT > 0.0 {
		optSet = append(optSet, ascii.FontPts(a.FontPT))
	}

	if a.CharsetPreset != "" {
		charset := a.ParseCharset()
		optSet = append(optSet, ascii.CSet(charset))
	}

	if a.Charset != "" {
		optSet = append(optSet, ascii.CSet(ascii.Charset(a.Charset)))
	}

	if a.Font != "" {
		font, err := OpenFont(a.Font)
		if err != nil {
			return fmt.Errorf("error opening font: %v", err)
		}
		optSet = append(optSet, ascii.Font(font))
	}

	if a.Interpolate {
		optSet = append(optSet, ascii.Interpolate(mem))
	}

	if a.Noise > 0 {
		optSet = append(optSet, ascii.Noise(a.Noise))
	}

	if a.Video {
		args := strings.Split(a.FFMpegArgs, " ")
		return a.CreateVideo(optSet, args)
	}

	img, err := a.OpenImage()
	if err != nil {
		return err
	}

	var outname string
	if a.Output != "" {
		outname = a.Output
	} else {
		outname = "ascii.png"
	}

	asciiimg, err := ascii.ConvertWithOpts(img, optSet...)
	if err != nil {
		return err
	}

	return SaveImageToPNG(asciiimg, outname)
}
