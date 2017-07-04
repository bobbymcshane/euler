package lib

import "math"

func NDigits(number int) int {
	return int(math.Log10(float64(number))) + 1
}
