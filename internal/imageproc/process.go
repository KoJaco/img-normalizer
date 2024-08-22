package imageproc

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
)

func ProcessImage(path, destDir string, tolerance int) {
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

	aspectW, aspectH := FindBestAspectRatio(width, height, tolerance)
	if aspectW == 0 || aspectH == 0 {
		fmt.Printf("No suitable aspect ratio found within tolerance for image %s.\n", path)
		return
	}

	targetW, targetH := CalculateTargetDimensions(width, height, aspectW, aspectH)
	resultImg := AdjustImageToAspectRatio(img, targetW, targetH, tolerance)

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
