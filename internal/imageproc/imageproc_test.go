package imageproc

import (
	"image"
	"image/color"
	"testing"
)

// TODO: Expand tests.. these are very basic.

func TestCalculateTargetDimensions(t *testing.T) {
	width, height := 1271, 799
	aspectW, aspectH := 3, 2

	var tol = 1 // tolerance for rounding in calculation (casting to integer).

	targetW, targetH := CalculateTargetDimensions(width, height, aspectW, aspectH)

	expectedW, expectedH := 1199, 799
	if targetW > expectedW+tol || targetW < expectedW-tol || targetH > expectedH+tol || targetH < expectedH-tol {
		t.Errorf("Expected %dx%d, got %dx%d", expectedW, expectedH, targetW, targetH)
	}
}

func TestScaleImage(t *testing.T) {
	img := createTestImage(800, 600)
	targetWidth, targetHeight := 400, 300

	resultImg := ScaleImage(img, targetWidth, targetHeight)

	if resultImg.Bounds().Dx() != targetWidth || resultImg.Bounds().Dy() != targetHeight {
		t.Errorf("Expected size %dx%d, got %dx%d", targetWidth, targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
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

// import (
// 	"fmt"
// 	"image"

// 	"image/color"
// 	"image/draw"
// 	"image/jpeg"
// 	"os"
// 	"path/filepath"
// 	"testing"
// )

// // TODO: add test for calculate dimensions (ensure always scaling down)

// // testing with source images in `assets/tests/src``, save results in `assets/tests/results`

// func TestLoadImage(t *testing.T) {
// 	absPath, err := filepath.Abs("assets/test/expected_3x2.jpg")

// 	if err != nil {
// 		t.Fatalf("Failed to get absolute path: %v", err)
// 	}

// 	fmt.Println("Absolute path: ", absPath)
// }

// func TestAdjustToAspectRatio_Scale_ActualImage(t *testing.T) {
// 	img, err := loadImage("assets/test/expected_3x2.jpg")
// 	if err != nil {
// 		t.Fatalf("Failed to load image: %v", err)
// 	}

// 	// expect_3x2.jpg
// 	targetWidth, targetHeight := 1271, 847
// 	tolerance := 0.01

// 	resultImg := AdjustToAspectRatio(img, targetWidth, targetHeight, tolerance)

// 	if resultImg.Bounds().Dx() != targetWidth || resultImg.Bounds().Dy() != targetHeight {
// 		t.Errorf("Expected size %dx%d, got %dx%d", targetWidth, targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
// 	}

// 	err = saveImage(resultImg, "assets/res/adjusted_scale_expected_3x2.jpg")
// 	if err != nil {
// 		t.Fatalf("Failed to save scaled image: %v", err)
// 	}
// }

// func TestAdjustToAspectRatio_Crop_ActualImage(t *testing.T) {
// 	img, err := loadImage("assets/test/expected_3x2.jpg")
// 	if err != nil {
// 		t.Fatalf("Failed to load image: %v", err)
// 	}

// 	// expected_3x2.jpg
// 	targetWidth, targetHeight := 1199, 799
// 	tolerance := 0.01

// 	resultImg := AdjustToAspectRatio(img, targetWidth, targetHeight, tolerance)

// 	if resultImg.Bounds().Dx() != targetWidth || resultImg.Bounds().Dy() != targetHeight {
// 		t.Errorf("Expected size %dx%d, got %dx%d", targetWidth, targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
// 	}

// 	err = saveImage(resultImg, "assets/res/adjusted_crop_expected_3x2.jpg")
// 	if err != nil {
// 		t.Fatalf("Failed to save cropped image: %v", err)
// 	}
// }

// // Helper function to load an image from a file
// func loadImage(filename string) (image.Image, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	img, err := jpeg.Decode(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return img, nil
// }

// // Testing with created images (does not save images)

// func TestAdjustToAspectRatio_Scale_CreatedImage(t *testing.T) {
// 	img := createTestImage(800, 600)
// 	targetWidth, targetHeight := 400, 300
// 	tolerance := 0.01

// 	resultImg := AdjustToAspectRatio(img, targetWidth, targetHeight, tolerance)

// 	if resultImg.Bounds().Dx() != targetWidth || resultImg.Bounds().Dy() != targetHeight {
// 		t.Errorf("Expected size %dx%d, got %dx%d", targetWidth, targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
// 	}
// }

// func TestAdjustToAspectRatio_Crop_CreatedImage(t *testing.T) {
// 	img := createTestImage(800, 800)
// 	targetWidth, targetHeight := 400, 300
// 	tolerance := 0.01

// 	resultImg := AdjustToAspectRatio(img, targetWidth, targetHeight, tolerance)

// 	if resultImg.Bounds().Dx() != targetWidth || resultImg.Bounds().Dy() != targetHeight {
// 		t.Errorf("Expected size %dx%d, got %dx%d", targetWidth, targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
// 	}
// }

// func TestScaleImage(t *testing.T) {
// 	img := createTestImage(800, 600)
// 	targetWidth, targetHeight := 400, 300

// 	resultImg := scaleImage(img, targetWidth, targetHeight)

// 	if resultImg.Bounds().Dx() != targetWidth || resultImg.Bounds().Dy() != targetHeight {
// 		t.Errorf("Expected size %dx%d, got %dx%d", targetWidth, targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
// 	}
// }

// func TestCropImage(t *testing.T) {
// 	img := createTestImage(800, 600)
// 	targetWidth, targetHeight := 400, 300

// 	resultImg := cropImage(img, targetWidth, targetHeight)

// 	if resultImg.Bounds().Dx() != targetWidth || resultImg.Bounds().Dy() != targetHeight {
// 		t.Errorf("Expected size %dx%d, got %dx%d", targetWidth, targetHeight, resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
// 	}
// }

// // Helper function to create a test image with a solid color
// func createTestImage(width, height int) image.Image {
// 	img := image.NewRGBA(image.Rect(0, 0, width, height))
// 	blue := color.RGBA{0, 0, 255, 255}
// 	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)
// 	return img
// }

// // Helper function to save an image to a file
// func saveImage(img image.Image, filename string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	return jpeg.Encode(file, img, nil)
// }
