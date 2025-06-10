package main

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
)

func SaveImageToPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

func WriteImageToStdout(img image.Image) error {
	return png.Encode(os.Stdout, img)
}

func SaveAsGif(g *gif.GIF, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gif.EncodeAll(file, g)
	if err != nil {
		return err
	}

	return nil
}

func openGif(imgpath string) (*gif.GIF, error) {
	file, err := os.Open(imgpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// gifs are messed up af
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while decoding: %s", r)
		}
	}()

	img, err := gif.DecodeAll(file)
	if err != nil {
		return nil, fmt.Errorf("gif decode: %v", err)
	}

	return img, nil
}

func openImage(imgpath string) (image.Image, error) {
	file, err := os.Open(imgpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("%w - known image formats are jpeg and png", err)
	}

	return img, nil
}
