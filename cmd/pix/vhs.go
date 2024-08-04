package main

import (
	"fmt"
	"image"
	"image/color"

	"pix/pkg/imaging"
	"pix/pkg/pixlib"
)

var fuck = fmt.Printf

// ApplyScanlines applies scanlines
func ApplyScanlines(destImage *image.RGBA) {
	bounds := destImage.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			clr := destImage.At(x, y)

			v := pixlib.ColorToVec3(clr)
			v2 := pixlib.ColorToVec3(color.Black)
			vv := v.Lerp(v2, 0.33)

			newClr := pixlib.Vec3ToColor(vv)
			destImage.Set(x, y, newClr)
			// destImage.Set(x, y, color.Black)
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
	bounds2 := img2.Bounds()
	outimg := imageToRGBA(img)

	if v.Scale {
		var sfact float64
		if v.ScaleFactor > 0 {
			sfact = v.ScaleFactor
		} else {
			sfact = 2
		}

		img2 = imaging.Resize(img2, int(float64(bounds.Dx())*sfact), int(float64(bounds.Dy())*sfact), imaging.NearestNeighbor)
		SaveImageToPNG(img2, "testing.png")
	}

	bounds2 = img2.Bounds()

	xF := bounds2.Min.X
	yF := bounds2.Min.Y

	fuck("%v %v\n", xF, yF)

	red := color.RGBA{0xff, 0, 0, 0xff}
	green := color.RGBA{0, 0xff, 0, 0xff}
	blue := color.RGBA{0, 0, 0xff, 0xff}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			src := img.At(x, y)
			src2 := img2.At(xF, yF)

			v := pixlib.ColorToVec3(src)
			v2 := pixlib.ColorToVec3(src2)
			// vv := v.Sub(v2)
			// vv = v.Add(v2)
			vv := v.Lerp(v2, 0.333)

			var c color.Color
			if x%5 == 0 || x%6 == 0 {
				c = blue
			} else if x%3 == 0 || x%4 == 0 {
				c = green
			} else if x%1 == 0 || x%2 == 0 {
				c = red
			} else {
				c = red
			}

			cv := pixlib.ColorToVec3(c)
			cv2 := cv.Lerp(v, 0.9)

			newClr := pixlib.Vec3ToColor(vv)
			outimg.Set(x, y, newClr)

			rgb := pixlib.Vec3ToColor(cv2)
			outimg.Set(x, y, rgb)

			// if xF < bounds2.Max.X {
			// 	xF++
			// } else if xF > bounds2.Max.X-1 {
			// 	xF = bounds2.Min.X + 4
			// }
			//
			// if yF < bounds2.Max.Y {
			// 	yF++
			// } else if yF > bounds2.Max.Y-1 {
			// 	yF = bounds2.Min.Y + 4
			// }
		}
	}

	ApplyScanlines(outimg)

	return SaveImageToPNG(outimg, "output.png")
}
