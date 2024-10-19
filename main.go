package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bibauporto/photosOrganizer/helpers"
	"github.com/bibauporto/photosOrganizer/processors"
)

// Process all folders and files recursively
func processFiles(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %w", folderPath, err)
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		if file.IsDir() {
			// Recursively process subfolders
			if err := processFiles(filePath); err != nil {
				return fmt.Errorf("error processing folder %s: %w", filePath, err)
			}
		} else {
			// Process files based on extension
			ext := strings.ToLower(filepath.Ext(file.Name()))
			switch {
			case helpers.Contains(helpers.IMAGE_EXTENSIONS, ext):
				if err := processors.ProcessImage(file.Name(), folderPath); err != nil {
					return fmt.Errorf("error processing image %s: %w", file.Name(), err)
				}
			case helpers.Contains(helpers.VIDEO_EXTENSIONS, ext):
				if err := processors.ProcessVideo(file.Name(), folderPath); err != nil {
					return fmt.Errorf("error processing video %s: %w", file.Name(), err)
				}
			default:
				fmt.Printf("Skipping unsupported file type: %s\n", file.Name())
			}
		}
	}
	return nil
}

func main() {
	folderPath, _ := os.Getwd() 
	fmt.Println("Starting processing...")
	if err := processFiles(folderPath); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("Processing complete.")
}
