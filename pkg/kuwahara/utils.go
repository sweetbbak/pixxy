package kuwahara

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

func getAperture(axisValue, axisMax, apertureMin, apertureMax int) (int, int) {
	start, end := 0, axisMax
	if axisValue+apertureMin > 0 {
		start = axisValue + apertureMin
	}
	if axisValue+apertureMax < axisMax {
		end = axisValue + apertureMax
	}
	return start, end
}

func getCoordinates(pixel, width int) (int, int) {
	return pixel % width, int(math.Floor(float64(pixel) / float64(width)))
}

func getGradientPoint(axisValue, shift, axisLength int) int {
	if (axisValue + shift) >= axisLength {
		return axisLength - axisValue - 1
	}
	return shift
}

func getPixel(x, y, width int) int {
	return ((y * width) + x) * 4
}

func getPixPerThread(pixLen, threads int) int {
	pixPerThreadRaw := float64(pixLen) / float64(threads)
	module := math.Mod(pixPerThreadRaw, 4.0)
	if module == 0 {
		return int(pixPerThreadRaw)
	}
	return int(pixPerThreadRaw + (float64(threads) - math.Mod(pixPerThreadRaw, 4.0)))
}

func openfile(path string) *image.RGBA {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open the file: ", err)
	}
	defer file.Close()

	content, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("Could not decode the file: ", err)
	}

	rect := content.Bounds()
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), content, rect.Min, draw.Src)
	return img
}

func open(path string) (*image.RGBA, string, int, int) {
	now := math.Round(float64(time.Now().UnixNano()) / 1000000)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open the file: ", err)
	}
	defer file.Close()

	openMS := int(math.Round(float64(time.Now().UnixNano())/1000000) - now)
	now2 := math.Round(float64(time.Now().UnixNano()) / 1000000)

	content, format, err := image.Decode(file)
	if err != nil {
		log.Fatal("Could not decode the file: ", err)
	}

	rect := content.Bounds()
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), content, rect.Min, draw.Src)
	convertMS := int(math.Round(float64(time.Now().UnixNano())/1000000) - now2)
	return img, format, openMS, convertMS
}

func save(img *image.RGBA, format string) int {
	name := fmt.Sprintf(`file-%d.%s`, time.Now().Unix(), format)
	now := math.Round(float64(time.Now().UnixNano()) / 1000000)

	newFile, err := os.Create(name)
	if err != nil {
		log.Fatal("Could not save the file")
	}
	defer newFile.Close()

	if format == "png" {
		png.Encode(newFile, img.SubImage(img.Rect))
	} else {
		jpeg.Encode(
			newFile,
			img.SubImage(img.Rect),
			&jpeg.Options{
				Quality: 100,
			},
		)
	}

	return int(math.Round(float64(time.Now().UnixNano())/1000000) - now)
}
