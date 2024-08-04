package main

import (
	"fmt"
	"image"
	"image/color"

	"pix/pkg/pixlib"
)

var fuck = fmt.Print

// ApplyScanlines applies scanlines
func ApplyScanlines(destImage *image.RGBA) {
	bounds := destImage.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			destImage.Set(x, y, color.Black)
		}
	}
}

func (v *VHS) Run() error {
	// if v.Input == "" {
	// 	return fmt.Errorf("provide an input image")
	// }

	img, err := openImage(v.Input[0])
	if err != nil {
		return err
	}

	img2, err := openImage(v.Input[1])
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	outimg := imageToRGBA(img)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			src := img2.At(x, y)
			src2 := img.At(x, y)

			v := pixlib.ColorToVec3(src)
			v2 := pixlib.ColorToVec3(src2)
			// vv := v.Sub(v2)
			// vv := v.Add(v2)
			vv := v.Lerp(v2, 0.333)

			newClr := pixlib.Vec3ToColor(vv)
			outimg.Set(x, y, newClr)
		}
	}

	return SaveImageToPNG(outimg, "output.png")
}
