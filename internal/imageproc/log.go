package imageproc

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

type LogEntry struct {
	ImagePath         string
	OriginalDim       string
	ChosenAspectRatio string
	NewDim            string
	Status            string
}

// SaveLog saves the log entries to a CSV file in the dest directory.

func SaveLog(logEntries []LogEntry, destDir string) error {
	logFilePath := filepath.Join(destDir, "image_process_log.csv")

	file, err := os.Create(logFilePath)

	if err != nil {
		return fmt.Errorf("could not create log file: %v", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Image Path", "Original Dimensions", "Chosen Aspect Ratio", "New Dimensions", "Status"})

	for _, entry := range logEntries {
		writer.Write([]string{
			entry.ImagePath,
			entry.OriginalDim,
			entry.ChosenAspectRatio,
			entry.NewDim,
			entry.Status,
		})
	}

	return nil
}
