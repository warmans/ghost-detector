package util

// Normalize take an int with a scale of e.g. 0-1000 and converts it to an int
// in the scale of 0-20.
func Normalize(val, oldMax, newMax int) int {
	return int(float64(newMax) / float64(oldMax) * float64(val))
}
