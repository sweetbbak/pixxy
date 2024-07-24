package quantize

import (
	"image"
	"image/color"
	"sort"
)

func min(first uint8, second uint8) uint8 {
	if first < second {
		return first
	}
	return second
}

func max(first uint8, second uint8) uint8 {
	if first > second {
		return first
	}
	return second
}

// Spread takes in a slice of RGB pixels, and returns the delta across the red,
// green, & blue components of all pixels.
func Spread(pixels []color.RGBA) (uint8, uint8, uint8) {

	// If there are no pixels, then the spread must be zero
	if len(pixels) == 0 {
		return 0, 0, 0
	}

	var (
		minRed, maxRed     = pixels[0].R, pixels[0].R
		minGreen, maxGreen = pixels[0].G, pixels[0].G
		minBlue, maxBlue   = pixels[0].B, pixels[0].B
	)

	for _, pixel := range pixels {
		r, g, b := pixel.R, pixel.G, pixel.B

		// Minmax the red component
		minRed = min(minRed, r)
		maxRed = max(maxRed, r)

		// Minmax the green component
		minGreen = min(minGreen, g)
		maxGreen = max(maxGreen, g)

		// Minmax the blue component
		minBlue = min(minBlue, b)
		maxBlue = max(maxBlue, b)
	}

	return maxRed - minRed, maxGreen - minGreen, maxBlue - minBlue
}

// Partition takes in a slice of RGB pixels, and divides the slice into two equal parts with respect
// to the color component with the largest spread.
func Partition(pixels []color.RGBA) ([]color.RGBA, []color.RGBA) {

	if len(pixels) == 0 {
		return []color.RGBA{}, []color.RGBA{}
	}

	deltaR, deltaG, deltaB := Spread(pixels)

	var less func(int, int) bool

	switch {
	// Does the red component have the largest spread?
	case deltaR >= deltaG && deltaR >= deltaB:
		less = func(i int, j int) bool {
			return pixels[i].R < pixels[j].R
		}

	// Does the green component have the largest spread?
	case deltaG >= deltaR && deltaG >= deltaB:
		less = func(i int, j int) bool {
			return pixels[i].G < pixels[j].G
		}

	// Does the blue component have the largest spread?
	case deltaB >= deltaR && deltaB >= deltaG:
		less = func(i int, j int) bool {
			return pixels[i].B < pixels[j].B
		}
	}

	// Sort pixels by the component with the largest spread
	sort.SliceStable(pixels, less)

	return pixels[:len(pixels)/2], pixels[len(pixels)/2:]
}

// Average takes in a slice of RGB pixels, and returns the average across the
// red, green, & blue components of all pixels.
func Average(pixels []color.RGBA) color.RGBA {
	var totalR int
	var totalG int
	var totalB int

	if len(pixels) == 0 {
		return color.RGBA{0, 0, 0, 0xFF}
	}

	for _, pixel := range pixels {
		totalR += int(pixel.R)
		totalG += int(pixel.G)
		totalB += int(pixel.B)
	}

	return color.RGBA{
		uint8(totalR / len(pixels)),
		uint8(totalG / len(pixels)),
		uint8(totalB / len(pixels)),
		0xFF,
	}
}

func Pixels(pixels []color.RGBA, levels int) []color.RGBA {
	partitions := [][]color.RGBA{
		pixels,
	}

	for iteration := 0; iteration < levels; iteration++ {

		next := [][]color.RGBA{}

		for _, partition := range partitions {
			left, right := Partition(partition)
			next = append(next, left, right)
		}

		partitions = next
	}

	averages := make([]color.RGBA, len(partitions))

	for index, partition := range partitions {
		averages[index] = Average(partition)
	}

	return averages
}

func GetColorPalette(img image.Image, level int) []color.Color {
	pal := []color.Color{}
	colors := Palette(img, level)
	for _, c := range colors {
		pal = append(pal, c)
	}
	return pal
}

// Returns a color Palette of N colors, the higher the levels the more colors and the more
// accurate the color palette is. uses the "median cut algorithm" we recursively cut our data (pixels)
// at the median average of the pixels color value
func Palette(img image.Image, levels int) []color.RGBA {
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
	return Pixels(pixels, levels)
}
