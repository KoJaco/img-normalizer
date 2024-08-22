package imageproc

import (
	"math"
)

// Predefined aspect ratios
var commonAspectRatios = []struct {
	Width, Height int
}{
	{1, 1},  // Square
	{4, 3},  // Standard photo
	{3, 2},  // 35mm film
	{16, 9}, // Widescreen
	{21, 9}, // Ultrawide
	{9, 16}, // Portrait
}

// FindBestAspectRatio determines the closest common aspect ratio to the image's aspect ratio.
func FindBestAspectRatio(width, height int, tolerance float64) (int, int) {
	currentRatio := float64(width) / float64(height)
	var bestRatio struct{ Width, Height int }
	smallestDiff := math.MaxFloat64

	for _, ratio := range commonAspectRatios {
		ratioValue := float64(ratio.Width) / float64(ratio.Height)
		diff := math.Abs(currentRatio - ratioValue)
		if diff < smallestDiff && diff <= tolerance {
			smallestDiff = diff
			bestRatio = ratio
		}
	}

	return bestRatio.Width, bestRatio.Height
}
