package imageproc

import "math"

func FindBestAspectRatio(width, height int, tolerance int) (int, int) {
	aspectRatios := []struct {
		w int
		h int
	}{
		{1, 1}, {4, 3}, {4, 5}, {3, 2}, {16, 9}, {21, 9}, {9, 16},
	}

	var bestRatio struct{ w, h int }
	minDiff := math.MaxFloat64

	for _, ratio := range aspectRatios {
		targetWidth, targetHeight := CalculateTargetDimensions(width, height, ratio.w, ratio.h)

		widthDiff := width - targetWidth
		heightDiff := height - targetHeight

		if widthDiff <= tolerance && heightDiff <= tolerance {
			totalDiff := math.Abs(float64(widthDiff)) + math.Abs(float64(heightDiff))
			if totalDiff < minDiff {
				minDiff = totalDiff
				bestRatio = ratio
			}
		}
	}

	return bestRatio.w, bestRatio.h
}
