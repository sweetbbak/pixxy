package main

import (
	"fmt"
	"image/color"
)

func clampMax[T float64 | int | uint8](value, max T) T {
	if value > max {
		return max
	}
	return value
}

func clamp(x float64) uint8 {
	v := int64(x + 0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint8(v)
	}
	return 0
}

// maxFloat compares three float64 values and returns the largest one.
func maxNum[T float64 | float32 | int | uint8](a, b, c T) T {
	if (a >= b) && (a >= c) {
		return a
	} else if (b >= a) && (b >= c) {
		return b
	} else {
		return c
	}
}

// maxFloat compares three float64 values and returns the largest one.
func minNum[T float64 | float32 | int | uint8](a, b, c T) T {
	if (a <= b) && (a <= c) {
		return a
	} else if (b <= a) && (b <= c) {
		return b
	} else {
		return c
	}
}

// convert default RGBA colors to 0-255 range
func ConvertTo255(clr color.Color) color.Color {
	r, g, b, a := clr.RGBA()
	r, g, b, a = r>>8, g>>8, b>>8, a>>8
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// convert 0-65353 color range to 0-1 float
func RGBtoFloat(clr color.Color) (float64, float64, float64) {
	r, g, b, _ := clr.RGBA()
	r, g, b = r>>8, g>>8, b>>8
	rf, gf, bf := float64(r)/255.0, float64(g)/255, float64(b)/255
	return rf, gf, bf
}

// returns a 0-360 degrees HUE
func GetHUE(clr color.Color) float64 {
	r, g, b := RGBtoFloat(clr)

	largest := maxNum(r, g, b)
	smallest := minNum(r, g, b)

	fmt.Printf("min %f - max %f\n", smallest, largest)
	luma := (largest + smallest) / 2
	fmt.Printf("luma %f - luma2 %f\n", luma, luminance64(r, g, b))

	var hue float64
	if r >= g && r >= b {
		hue = (g - b) / (largest - smallest)
	} else if g >= r && g >= b {
		hue = 2.0 + (b-r)/(largest-smallest)
	} else if b >= r && b >= g {
		hue = 4.0 + (r-g)/(largest-smallest)
	}

	hue *= 60
	if hue < 0 {
		hue += 360
	}

	return hue
}

func luminance(r, g, b float32) float32 {
	return r*0.299 + g*0.587 + b*0.114
}

func luminance64(r, g, b float64) float64 {
	return r*0.299 + g*0.587 + b*0.114
}
