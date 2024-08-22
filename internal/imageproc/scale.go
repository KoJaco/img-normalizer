package imageproc

import "image"

func ScaleImage(img image.Image, targetWidth, targetHeight int) image.Image {
	rect := image.Rect(0, 0, targetWidth, targetHeight)
	scaledImg := image.NewRGBA(rect)

	xRatio := float64(img.Bounds().Dx()) / float64(targetWidth)
	yRatio := float64(img.Bounds().Dy()) / float64(targetHeight)

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)
			color := img.At(srcX, srcY)
			scaledImg.Set(x, y, color)
		}
	}

	return scaledImg
}
