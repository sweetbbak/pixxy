package main

import (
	"fmt"
	"image/color"
	"testing"
)

func TestGetHUE(t *testing.T) {
	tests := []struct {
		name string
		clr  color.RGBA
	}{
		{"color1", color.RGBA{100, 17, 199, 255}},
		{"color2", color.RGBA{10, 255, 87, 255}},
		{"color3", color.RGBA{133, 100, 201, 255}},
	}

	for _, ts := range tests {
		t.Run("HUE", func(t *testing.T) {
			hue := GetHUE(ts.clr)
			fmt.Printf("hue value: %f of color [%v]\n", hue, ts.clr)
			t.Logf("hue value: %f of color [%v]\n", hue, ts.clr)
		})
	}
}
