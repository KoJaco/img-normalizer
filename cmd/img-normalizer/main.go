package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"img-normalizer/internal/imageproc"
)

func main() {

	// Define flags
	// TODO: I want the user to define a directory with a flat structure or a nested structure and the it will walk that structure... or just a single image file.
	srcDir := flag.String("src", "assets/testing/src", "Source directory of images")
	destDir := flag.String("dest", "assets/testing/out", "Destination directory for processed images")
	tolerance := flag.Int("tolerance", 100, "Maximum allowed pixel difference for adjustment")
	flag.Parse()

	// Ensure destination directory exists
	if err := os.MkdirAll(*destDir, os.ModePerm); err != nil {
		fmt.Println("Error creating destination directory:", err)
		return
	}

	// Walk through all files in the source directory
	err := filepath.Walk(*srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process files (skip directories)
		if !info.IsDir() {
			imageproc.ProcessImage(path, *destDir, *tolerance)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through the directory:", err)
	}
}
