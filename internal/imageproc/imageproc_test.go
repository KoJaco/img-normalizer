package imageproc

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"testing"
)

// TODO: Expand tests.. these are very basic.

func TestCalculateTargetDimensions(t *testing.T) {
	tests := []struct {
		width, height        int
		aspectW, aspectH     int
		expectedW, expectedH int
	}{
		// Exact matches
		{1271, 799, 3, 2, 1199, 799},
		{1920, 1080, 16, 9, 1920, 1080},

		// Slightly off aspect ratios
		{1280, 853, 16, 9, 1280, 720},
		{1400, 1050, 4, 3, 1400, 1050},

		// Edge cases
		{1000, 500, 16, 9, 888, 500}, // width needs to be reduced significantly
		{999, 500, 2, 1, 1000, 500},  // closest is 2:1, but image height already fits
		{500, 1000, 1, 2, 500, 1000}, // portrait-oriented image

		// Large image
		{5000, 3333, 3, 2, 5000, 3333},
	}

	tol := 1 // tolerance for rounding in calculation (casting to integer).
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			targetW, targetH := CalculateTargetDimensions(test.width, test.height, test.aspectW, test.aspectH)
			if targetW > test.expectedW+tol || targetW < test.expectedW-tol || targetH > test.expectedH+tol || targetH < test.expectedH-tol {
				t.Errorf("Expected %dx%d, got %dx%d", test.expectedW, test.expectedH, targetW, targetH)
			}
		})
	}
}

func TestScaleSampleImage(t *testing.T) {
	tests := []struct {
		imgPath      string
		targetWidth  int
		targetHeight int
		outputDir    string
	}{
		{"./test-assets/sample.jpg", 1280, 720, "./test-assets/out/"}, // scaled beyond bounds... realistically this should not be possible.
		{"./test-assets/sample.jpg", 640, 480, "./test-assets/out/"},
		{"./test-assets/sample.jpg", 320, 240, "./test-assets/out/"},
	}

	for _, test := range tests {
		img, err := LoadImage(test.imgPath)
		if err != nil {
			t.Fatalf("Failed to load image: %v", err)
		}

		resultImg := ScaleImage(img, test.targetWidth, test.targetHeight)

		if resultImg.Bounds().Dx() != test.targetWidth || resultImg.Bounds().Dy() != test.targetHeight {
			t.Errorf("Expected size %dx%d, got %dx%d", test.targetWidth, test.targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
		}

		outputFileName := fmt.Sprintf("scaled_sample_%dx%d.jpg", test.targetWidth, test.targetHeight)
		outputFilePath := filepath.Join(test.outputDir, outputFileName)

		err = SaveImage(resultImg, outputFilePath)
		if err != nil {
			t.Fatalf("Failed to save scaled image: %v", err)
		}
	}
}

func TestScaleImage(t *testing.T) {
	tests := []struct {
		width, height             int
		targetWidth, targetHeight int
	}{
		{800, 600, 400, 300},    // Standard scaling
		{1600, 1200, 800, 600},  // Large scaling down
		{320, 240, 160, 120},    // Small scaling down
		{1280, 853, 1280, 720},  // Slightly off aspect ratio
		{1920, 1080, 1280, 720}, // Common aspect ratio
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			img := createTestImage(test.width, test.height)
			resultImg := ScaleImage(img, test.targetWidth, test.targetHeight)
			if resultImg.Bounds().Dx() != test.targetWidth || resultImg.Bounds().Dy() != test.targetHeight {
				t.Errorf("Expected size %dx%d, got %dx%d", test.targetWidth, test.targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
			}
		})
	}
}

func TestCropSampleImage(t *testing.T) {
	tests := []struct {
		imgPath      string
		targetWidth  int
		targetHeight int
		outputDir    string
	}{
		{"./test-assets/sample.jpg", 640, 480, "./test-assets/out/"},
		{"./test-assets/sample.jpg", 320, 240, "./test-assets/out/"},
	}

	for _, test := range tests {
		img, err := LoadImage(test.imgPath)
		if err != nil {
			t.Fatalf("Failed to load image: %v", err)
		}

		resultImg := CropImage(img, test.targetWidth, test.targetHeight)

		if resultImg.Bounds().Dx() != test.targetWidth || resultImg.Bounds().Dy() != test.targetHeight {
			t.Errorf("Expected size %dx%d, got %dx%d", test.targetWidth, test.targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
		}

		outputFileName := fmt.Sprintf("cropped_sample_%dx%d.jpg", test.targetWidth, test.targetHeight)
		outputFilePath := filepath.Join(test.outputDir, outputFileName)

		err = SaveImage(resultImg, outputFilePath)
		if err != nil {
			t.Fatalf("Failed to save cropped image: %v", err)
		}
	}
}

func TestCropImage(t *testing.T) {
	tests := []struct {
		width, height             int
		targetWidth, targetHeight int
	}{
		{800, 600, 400, 300},    // Standard crop
		{1600, 1200, 800, 600},  // Large crop
		{1280, 853, 1280, 720},  // Slightly off aspect ratio
		{1920, 1080, 1280, 720}, // Common aspect ratio
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			img := createTestImage(test.width, test.height)
			resultImg := CropImage(img, test.targetWidth, test.targetHeight)
			if resultImg.Bounds().Dx() != test.targetWidth || resultImg.Bounds().Dy() != test.targetHeight {
				t.Errorf("Expected size %dx%d, got %dx%d", test.targetWidth, test.targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
			}
		})
	}
}

func TestFindBestAspectRatio(t *testing.T) {
	tests := []struct {
		width, height        int
		tolerance            int
		expectedW, expectedH int
	}{
		// TODO: extend these test cases...

		// Exact matches
		{1920, 1080, 1, 16, 9}, // one pixel width, exact... should still be good.

		// Slightly off aspect ratios
		{1280, 853, 100, 3, 2},
		{1271, 799, 100, 3, 2},

		// tolerances
		{952, 932, 5, 0, 0}, // square but tolerance is too low, so return 0, 0

		// Edge cases
		{1023, 768, 20, 4, 3}, // Close to 4:3
		{1000, 1001, 5, 1, 1}, // Nearly square
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			aspectW, aspectH := FindBestAspectRatio(test.width, test.height, test.tolerance)
			if aspectW != test.expectedW || aspectH != test.expectedH {
				t.Errorf("Expected aspect ratio %d:%d, got %d:%d", test.expectedW, test.expectedH, aspectW, aspectH)
			}
		})
	}
}

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	blue := color.RGBA{0, 0, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, blue)
		}
	}
	return img
}
