package main

import (
	"context"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
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

// Decode reads and analyzes the given reader as a GIF image
func SplitAnimatedGIF(gif *gif.GIF) (err error) {
	imgWidth, imgHeight := getGifDimensions(gif)

	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), gif.Image[0], image.Point{0, 0}, draw.Src)

	for i, srcImg := range gif.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.Point{0, 0}, draw.Over)

		// save current frame "stack". This will overwrite an existing file with that name
		file, err := os.Create(fmt.Sprintf("%s%d%s", "frame", i, ".png"))
		if err != nil {
			return err
		}

		err = png.Encode(file, overpaintImage)
		if err != nil {
			return err
		}

		file.Close()
	}

	return nil
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

	imgWidth, imgHeight := getGifDimensions(g)
	newg := &gif.GIF{}

	var pimg *image.Paletted
	pall := palette.Plan9

	bounds := image.Rect(0, 0, imgWidth, imgHeight)
	pimg = image.NewPaletted(bounds, pall)

	// draw base gif image to regular image
	overpaintImage := image.NewRGBA(bounds)
	draw.Draw(overpaintImage, overpaintImage.Bounds(), g.Image[0], image.Point{0, 0}, draw.Src)

	// ascii conversion
	asciiimg, err := ascii.ConvertWithOpts(overpaintImage, optset...)
	if err != nil {
		return err
	}

	// create the initial frame
	draw.Draw(pimg, overpaintImage.Bounds(), asciiimg, image.Point{0, 0}, draw.Over)

	newg.Image = append(newg.Image, pimg)
	newg.Delay = append(newg.Delay, g.Delay[0])

	// maxGoroutines := 10
	// guard := make(chan struct{}, maxGoroutines)

	for i, srcImg := range g.Image {
		// guard <- struct{}{}
		// go func() {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.Point{0, 0}, draw.Over)
		pimg := image.NewPaletted(bounds, pall)

		asciiimg, err := ascii.ConvertWithOpts(overpaintImage, optset...)
		if err != nil {
			// return err
			log.Println(err)
		}

		draw.Draw(pimg, overpaintImage.Bounds(), asciiimg, image.Point{0, 0}, draw.Over)
		// draw.Draw(pimg, overpaintImage.Bounds(), overpaintImage, image.ZP, draw.Over)

		newg.Image = append(newg.Image, pimg)
		newg.Delay = append(newg.Delay, g.Delay[i])
		// <-guard
		// }()
	}
	SaveAsGif(newg, "test.gif")
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
