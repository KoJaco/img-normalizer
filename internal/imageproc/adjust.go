package imageproc

import (
	"image"
	"math"
)

func CalculateTargetDimensions(originalWidth, originalHeight, targetAspectW, targetAspectH int) (int, int) {
	aspectRatio := float64(targetAspectW) / float64(targetAspectH)
	maxWidth := int(math.Min(float64(originalWidth), float64(originalHeight)*aspectRatio))
	maxHeight := int(math.Min(float64(originalHeight), float64(originalWidth)/aspectRatio))
	return maxWidth, maxHeight
}

func AdjustImageToAspectRatio(img image.Image, targetWidth, targetHeight int, tolerance int) image.Image {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	widthDiff := originalWidth - targetWidth
	heightDiff := originalHeight - targetHeight

	if widthDiff <= tolerance && heightDiff <= tolerance {
		return ScaleImage(img, targetWidth, targetHeight)
	}

	if widthDiff <= tolerance || heightDiff <= tolerance {
		return CropImage(img, targetWidth, targetHeight)
	}

	return img
}
