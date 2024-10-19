package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	helpers "github.com/bibauporto/photosOrganizer/helpers"
	processors "github.com/bibauporto/photosOrganizer/processors"
)

// Process all files in a directory
func processFiles(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if helpers.Contains(helpers.IMAGE_EXTENSIONS, ext) {
			if err := processors.ProcessImage(file.Name(), folderPath); err != nil {
				return err
			}
		} else if helpers.Contains(helpers.VIDEO_EXTENSIONS, ext) {
			if err := processors.ProcessVideo(file.Name(), folderPath); err != nil {
				return err
			}
		} else {
			fmt.Printf("Skipping unsupported file type: %s\n", file.Name())
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