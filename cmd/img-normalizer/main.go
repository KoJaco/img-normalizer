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

	srcPath := flag.String("src", "", "(Required) Source directory or image file path")
	destDir := flag.String("dest", "", "(Required) Destination directory for processed images")
	tolerance := flag.Int("tolerance", 100, "(Optional) Maximum allowed pixel difference for adjustment")
	flag.Parse()

	if *srcPath == "" || *destDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(*destDir, os.ModePerm); err != nil {
		fmt.Println("Error creating destination directory:", err)
		return
	}

	// check if srcPath is a file or directory

	info, err := os.Stat(*srcPath)
	if err != nil {
		fmt.Printf("Error accessing source path: %v\n", err)
		return
	}

	var logEntries []imageproc.LogEntry

	if info.IsDir() {
		// Process directory (including nested)
		err := filepath.Walk(*srcPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Only process files (skip directories)
			if !info.IsDir() {
				// Calculate the relative path from the source directory
				relativePath, err := filepath.Rel(*srcPath, path)
				if err != nil {
					return err
				}

				// Create the corresponding output directory structure
				outputPath := filepath.Join(*destDir, filepath.Dir(relativePath))
				if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
					return err
				}

				// Process the image and save it to the corresponding location
				imageproc.ProcessImage(path, outputPath, *tolerance, &logEntries)
			}
			return nil
		})

		if err != nil {
			fmt.Println("Error walking through the directory:", err)
		}
	} else {
		// Process a single image file
		imageproc.ProcessImage(*srcPath, *destDir, *tolerance, &logEntries)
	}

	// Save the log file
	if err := imageproc.SaveLog(logEntries, *destDir); err != nil {
		fmt.Printf("Error saving log file: %v\n", err)
	}

}
