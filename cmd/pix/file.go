package main

import (
	"fmt"
	"image"
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
