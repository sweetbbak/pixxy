package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
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

func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}

func (a *Ascii) CreateGif(optset []ascii.Option) error {
	var inputfile string
	if a.Input != "" {
		inputfile = a.Input
	} else if a.Args.Image != "" {
		inputfile = a.Args.Image
	} else {
		return fmt.Errorf("no image supplied")
	}

	g, err := openGif(inputfile)
	if err != nil {
		return err
	}

	asciiimg, err := ascii.ConvertWithOpts(g.Image[0], optset...)
	if err != nil {
		return err
	}

	// imgWidth, imgHeight := getGifDimensions(g)
	newGif := &gif.GIF{}
	bounds := asciiimg.Bounds()

	var palettedImage *image.Paletted
	pal := color.Palette{}
	pal = palette.Plan9[:256]

	// bounds := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	println(bounds.Max.X, bounds.Max.Y)

	palettedImage = image.NewPaletted(bounds, pal)
	draw.FloydSteinberg.Draw(palettedImage, bounds, asciiimg, image.Pt(0, 0))

	frameDelay := g.Delay[0]
	frames := len(g.Image)

	newGif.Image = append(newGif.Image, palettedImage)
	newGif.Delay = append(newGif.Delay, frameDelay)
	frames--

	// maxGoroutines := 10
	// guard := make(chan struct{}, maxGoroutines)

	for _, img := range g.Image {
		// guard <- struct{}{}
		// go func() {
		asciiimg, err = ascii.ConvertWithOpts(img, optset...)
		if err != nil {
			// return err
		}

		println("processing frame ", frames)
		frameDelay := g.Delay[frames]

		palettedImage = image.NewPaletted(bounds, pal)
		draw.FloydSteinberg.Draw(palettedImage, bounds, asciiimg, image.Pt(0, 0))

		newGif.Image = append(newGif.Image, palettedImage)
		newGif.Delay = append(newGif.Delay, frameDelay)
		frames--
		// <-guard
		// }()
	}

	SaveAsGif(newGif, a.Output)
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

	if a.Gif {
		return a.CreateGif(optSet)
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
