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

func ProcessImage(path, destDir string, tolerance int, logEntries *[]LogEntry) {
	file, err := os.Open(path)
	if err != nil {
		// append log
		*logEntries = append(*logEntries, LogEntry{
			ImagePath: path,
			Status:    fmt.Sprintf("Error opening image: %v", err),
		})
		// fmt.Printf("Error opening image: %s: %v\n", path, err)
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
		// append log
		*logEntries = append(*logEntries, LogEntry{
			ImagePath: path,
			Status:    "Unsupported file type",
		})
		// fmt.Printf("Unsupported file type: %s\n", path)
		return
	}

	if err != nil {
		// append log
		*logEntries = append(*logEntries, LogEntry{
			ImagePath: path,
			Status:    fmt.Sprintf("Error decoding image: %v", err),
		})
		// fmt.Printf("Error decoding image %s: %v\n", path, err)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	originalDim := fmt.Sprintf("%dx%d", width, height)

	aspectW, aspectH := FindBestAspectRatio(width, height, tolerance)
	if aspectW == 0 || aspectH == 0 {
		// append log
		*logEntries = append(*logEntries, LogEntry{
			ImagePath:   path,
			OriginalDim: originalDim,
			Status:      fmt.Sprintf("No suitable aspect ratio found within the tolerance of %d.", tolerance),
		})
		// fmt.Printf("No suitable aspect ratio found within tolerance for image %s.\n", path)
		return
	}

	targetW, targetH := CalculateTargetDimensions(width, height, aspectW, aspectH)

	// Only process image if the new dimensions differ from the original
	if targetW == width && targetH == height {
		*logEntries = append(*logEntries, LogEntry{
			ImagePath:         path,
			OriginalDim:       originalDim,
			ChosenAspectRatio: fmt.Sprintf("%d:%d", aspectW, aspectH),
			Status:            "Image dimensions already match target aspect ratio.",
		})
		return
	}

	newDim := fmt.Sprintf("%dx%d", targetW, targetH)

	resultImg := AdjustImageToAspectRatio(img, targetW, targetH, tolerance)

	if resultImg == img {
		*logEntries = append(*logEntries, LogEntry{
			ImagePath:         path,
			OriginalDim:       originalDim,
			ChosenAspectRatio: fmt.Sprintf("%d:%d", aspectW, aspectH),
			Status:            "Image not adjusted due to tolerance",
		})
		return
	}

	fileName := filepath.Base(path)
	baseName := strings.TrimSuffix(fileName, ext)
	newFileName := fmt.Sprintf("%s_%dx%d%s", baseName, targetW, targetH, ext)
	newFilePath := filepath.Join(destDir, newFileName)

	outFile, err := os.Create(newFilePath)
	if err != nil {

		*logEntries = append(*logEntries, LogEntry{
			ImagePath:         path,
			OriginalDim:       originalDim,
			ChosenAspectRatio: fmt.Sprintf("%d:%d", aspectW, aspectH),
			Status:            fmt.Sprintf("Error creating output file: %v", err),
		})
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
		*logEntries = append(*logEntries, LogEntry{
			ImagePath:         path,
			OriginalDim:       originalDim,
			ChosenAspectRatio: fmt.Sprintf("%d:%d", aspectW, aspectH),
			Status:            fmt.Sprintf("Error saving image: %v", err),
		})
		return
	}

	*logEntries = append(*logEntries, LogEntry{
		ImagePath:         path,
		OriginalDim:       originalDim,
		ChosenAspectRatio: fmt.Sprintf("%d:%d", aspectW, aspectH),
		NewDim:            newDim,
		Status:            "Image processed successfully",
	})
}
