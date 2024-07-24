package main

import (
	"fmt"
)

// gcd calculates the greatest common divisor of two integers.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// simplifyFraction simplifies a fraction represented as a numerator and denominator.
func simplifyFraction(numerator, denominator int64) (int64, int64) {
	g := gcd(numerator, denominator)
	return numerator / g, denominator / g
}

// floatToFraction converts a float32 to a simplified fraction string.
func floatToFraction(f float32) (string, error) {
	if f == 0 {
		return "0/1", nil
	}

	// Convert float32 to int64, rounding towards zero.
	n := int64(float64(f) * 10000) // Multiply by 10000 to get more precision for conversion.
	d := int64(10000)              // Denominator is always 10^4 for simplicity.

	// Simplify the fraction.
	simplifiedN, simplifiedD := simplifyFraction(n, d)

	// Convert the simplified fraction back to a string.
	fractionStr := fmt.Sprintf("%d/%d", simplifiedN, simplifiedD)

	return fractionStr, nil
}
