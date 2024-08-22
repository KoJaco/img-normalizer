package imageproc

import (
	"image"
	"math"
)

// AdjustToAspectRatio adjusts the image to match the closest aspect ratio
// It either crops or scales the image depending on how close it is to the target ratio.
func AdjustToAspectRatio(img image.Image, targetWidth, targetHeight int, tolerance float64) image.Image {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// Calculate the target dimensions
	adjustedWidth, adjustedHeight := calculateAdjustedDimensions(originalWidth, originalHeight, targetWidth, targetHeight)

	// Ensure we are only scaling or cropping downward
	if adjustedWidth > originalWidth || adjustedHeight > originalHeight {
		adjustedWidth, adjustedHeight = scaleDown(adjustedWidth, adjustedHeight, originalWidth, originalHeight)
	}

	// Determine if scaling or cropping is required based on tolerance
	if shouldScale(originalWidth, originalHeight, adjustedWidth, adjustedHeight, tolerance) {
		return scaleImage(img, adjustedWidth, adjustedHeight)
	}
	return cropImage(img, adjustedWidth, adjustedHeight)
}

// calculateAdjustedDimensions calculates the dimensions that match the target aspect ratio.
func calculateAdjustedDimensions(originalWidth, originalHeight, targetWidth, targetHeight int) (int, int) {
	// Calculate the aspect ratios
	originalRatio := float64(originalWidth) / float64(originalHeight)
	targetRatio := float64(targetWidth) / float64(targetHeight)

	// Adjust dimensions based on the target aspect ratio
	if originalRatio > targetRatio {
		// The original image is wider than the target ratio, so adjust width
		return int(float64(originalHeight) * targetRatio), originalHeight
	} else {
		// The original image is taller than the target ratio, so adjust height
		return originalWidth, int(float64(originalWidth) / targetRatio)
	}
}

// scaleDown ensures the new dimensions are not larger than the original ones
func scaleDown(adjustedWidth, adjustedHeight, originalWidth, originalHeight int) (int, int) {
	return min(adjustedWidth, originalWidth), min(adjustedHeight, originalHeight)
}

func shouldScale(originalWidth, originalHeight, targetWidth, targetHeight int, tolerance float64) bool {
	originalRatio := float64(originalWidth) / float64(originalHeight)
	targetRatio := float64(targetWidth) / float64(targetHeight)

	return math.Abs(originalRatio-targetRatio) <= tolerance
}

// scaleImage scales the image to the target width and height.
func scaleImage(img image.Image, targetWidth, targetHeight int) image.Image {
	// create new image with target dimensions
	rect := image.Rect(0, 0, targetWidth, targetHeight)
	scaledImg := image.NewRGBA(rect)

	// scaling factors
	xRatio := float64(img.Bounds().Dx()) / float64(targetWidth)
	yRatio := float64(img.Bounds().Dy()) / float64(targetHeight)

	// nearest-neighbor scaling (https://courses.cs.vt.edu/~masc1044/L17-Rotation/ScalingNN.html)
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// calc nearest pixel in original image
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)

			color := img.At(srcX, srcY)
			scaledImg.Set(x, y, color)
		}
	}

	return scaledImg
}

// cropImage crops the image to the target width and height, centered.
func cropImage(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	widthDiff := bounds.Dx() - targetWidth
	heightDiff := bounds.Dy() - targetHeight

	// starting point for crop (centered)
	startX := widthDiff / 2
	startY := heightDiff / 2

	cropRect := image.Rect(startX, startY, startX+targetWidth, startY+targetHeight)
	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(cropRect)
}
