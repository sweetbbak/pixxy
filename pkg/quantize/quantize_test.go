package quantize

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"testing"
)

func __RGBtoAnsi(rgb color.RGBA) (foreground, background string) {
	cstr := strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m"
	foreground = "\x1b[38;2;" + cstr
	background = "\x1b[48;2;" + cstr
	return
}

func TestPalette(t *testing.T) {
	tests := []struct {
		path   string
		levels int
	}{
		{"testdata/img.jpg", 3},
		{"testdata/img2.jpg", 2},
		{"testdata/img3.jpg", 3},
	}

	for _, tc := range tests {
		t.Run("HUE", func(t *testing.T) {

			f, err := os.Open(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Fatal(err)
			}

			rect := img.Bounds()
			pixels := make([]color.RGBA, 0, rect.Max.X*rect.Max.Y)

			for x := rect.Min.X; x < rect.Max.X; x++ {
				for y := rect.Min.Y; y < rect.Max.Y; y++ {
					r, g, b, _ := img.At(x, y).RGBA()

					// bit shifted to convert integer range to 0-255
					pixel := color.RGBA{
						uint8(r >> 8),
						uint8(g >> 8),
						uint8(b >> 8),
						0xFF,
					}

					pixels = append(pixels, pixel)
				}
			}

			partitions := [][]color.RGBA{
				pixels,
			}

			for iteration := 0; iteration < tc.levels; iteration++ {

				next := [][]color.RGBA{}

				for _, partition := range partitions {
					left, right := Partition(partition)
					t.Logf("pixels left [%v] right [%v]\n", len(left), len(right))
					next = append(next, left, right)
				}

				partitions = next
			}

			averages := make([]color.RGBA, len(partitions))
			t.Logf("pixels partitions [%v]\n", len(partitions))

			for index, partition := range partitions {
				averages[index] = Average(partition)
				t.Logf("pixels partitions [%v]\n", averages[index])
				_, bg := __RGBtoAnsi(averages[index])
				fmt.Printf("%s  %s", bg, "\x1b[0m")
			}
		})
	}
}
