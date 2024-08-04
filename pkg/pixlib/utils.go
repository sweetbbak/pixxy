package pixlib

import ()

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

// turns a uint32 into its apprxoimate 0-1 value
func c(a uint32) uint8 {
	return uint8((float64(a) / MAXC) * 255)
}
