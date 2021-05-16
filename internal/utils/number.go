package utils

import (
	"math"
)

func RoundTo(n float64, decimals int) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}
