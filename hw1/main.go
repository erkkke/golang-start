package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z, prev, diff := 1.0, 0.0, 1e-9
	for math.Abs(prev - z) > diff {
		prev = z
		z -= (z * z - x) / (2 * z)
	}
	return z
}

func main() {
	result := Sqrt(2)
	answer := math.Sqrt(2)
	fmt.Printf("Result: %v \nCorrect answer: %v \nDifference: %v\n", result, answer, math.Abs(answer - result))
}