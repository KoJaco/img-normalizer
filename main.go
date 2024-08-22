package main

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"

	// "image/draw"
	"image/jpeg" // Use "image/png" for PNGs
	"image/png"
	"math"
	"os"

	"github.com/chai2010/webp"
)

func findBestAspectRatio(width, height int, tolerance int) (int, int) {
	aspectRatios := []struct {
		w int
		h int
	}{
		{1, 1}, {4, 3}, {4, 5}, {3, 2}, {16, 9}, {21, 9}, {9, 16},
	}

	var bestRatio struct{ w, h int }
	minDiff := math.MaxFloat64

	for _, ratio := range aspectRatios {
		targetWidth, targetHeight := calculateTargetDimensions(width, height, ratio.w, ratio.h)

		widthDiff := width - targetWidth
		heightDiff := height - targetHeight

		// Ensure that both width and height adjustments are within tolerance
		if widthDiff <= tolerance && heightDiff <= tolerance {
			// Calculate the total pixel reduction and check if it's the smallest found
			totalDiff := math.Abs(float64(widthDiff)) + math.Abs(float64(heightDiff))
			if totalDiff < minDiff {
				minDiff = totalDiff
				bestRatio = ratio
			}
		}
	}

	return bestRatio.w, bestRatio.h
}

func calculateTargetDimensions(originalWidth, originalHeight, targetAspectW, targetAspectH int) (int, int) {
	aspectRatio := float64(targetAspectW) / float64(targetAspectH)

	maxWidth := int(math.Min(float64(originalWidth), float64(originalHeight)*aspectRatio))
	maxHeight := int(math.Min(float64(originalHeight), float64(originalWidth)/aspectRatio))

	return maxWidth, maxHeight
}

// scaleImage scales the image to the target width and height.
func scaleImage(img image.Image, targetWidth, targetHeight int) image.Image {
	fmt.Println("Scaled Image")
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

func cropImage(img image.Image, targetWidth, targetHeight int) image.Image {
	fmt.Println("Cropped Image")
	bounds := img.Bounds()
	widthDiff := bounds.Dx() - targetWidth
	heightDiff := bounds.Dy() - targetHeight

	// Center cropping
	startX := widthDiff / 2
	startY := heightDiff / 2

	cropRect := image.Rect(startX, startY, startX+targetWidth, startY+targetHeight)
	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(cropRect)
}

func adjustImageToAspectRatio(img image.Image, targetWidth, targetHeight int, tolerance int) image.Image {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	widthDiff := originalWidth - targetWidth
	heightDiff := originalHeight - targetHeight

	// TODO: how to balance between scaling and cropping?
	// Determine if scaling is possible within the tolerance
	if widthDiff <= tolerance && heightDiff <= tolerance {
		// Perform scaling to fit the aspect ratio within the tolerance
		return scaleImage(img, targetWidth, targetHeight)
	}

	// If scaling is not possible, check if cropping is within tolerance
	if widthDiff <= tolerance || heightDiff <= tolerance {
		// Perform cropping to fit the aspect ratio within the tolerance
		return cropImage(img, targetWidth, targetHeight)
	}

	// If neither scaling nor cropping are within the tolerance, return the original image
	return img
}

func processImage(path, destDir string) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Printf("Error opening image: %s: %v\n", path, err)
		return
	}

	defer file.Close()

	ext := strings.ToLower(filepath.Ext(path))
	var img image.Image

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".webp":
		img, err = webp.Decode(file)
	default:
		fmt.Printf("Unsupported file type: %s\n", path)
		return
	}

	if err != nil {
		fmt.Printf("Error decoding image %s: %v\n", path, err)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	tolerance := 100

	aspectW, aspectH := findBestAspectRatio(width, height, tolerance)
	if aspectW == 0 || aspectH == 0 {
		fmt.Printf("No suitable aspect ratio found within tolerance for image %s.\n", path)
		return
	}

	targetW, targetH := calculateTargetDimensions(width, height, aspectW, aspectH)
	resultImg := adjustImageToAspectRatio(img, targetW, targetH, tolerance)

	fileName := filepath.Base(path)
	baseName := strings.TrimSuffix(fileName, ext)
	newFileName := fmt.Sprintf("%s_%dx%d%s", baseName, targetW, targetH, ext)
	newFilePath := filepath.Join(destDir, newFileName)

	outFile, err := os.Create(newFilePath)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", newFilePath, err)
		return
	}
	defer outFile.Close()

	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outFile, resultImg, nil)
	case ".png":
		err = png.Encode(outFile, resultImg)
	case ".webp":
		err = webp.Encode(outFile, resultImg, nil)
	}

	if err != nil {
		fmt.Printf("Error saving image %s: %v\n", newFilePath, err)
		return
	}

	fmt.Printf("Processed image saved as %s\n", newFilePath)

}

func main() {

	srcDir := "assets/testing/src"
	destDir := "assets/testing/out"

	// Ensure the destination directory exists
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		fmt.Println("Error creating destination directory:", err)
		return
	}

	// Walk through all files in the source directory
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process files (skip directories)
		if !info.IsDir() {
			processImage(path, destDir)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through the directory:", err)
	}
}
