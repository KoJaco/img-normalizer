package imageproc

import "image"

func CropImage(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	widthDiff := bounds.Dx() - targetWidth
	heightDiff := bounds.Dy() - targetHeight

	startX := widthDiff / 2
	startY := heightDiff / 2

	cropRect := image.Rect(startX, startY, startX+targetWidth, startY+targetHeight)
	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(cropRect)
}
