package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bibauporto/photosOrganizer/helpers"
	"github.com/bibauporto/photosOrganizer/processors"
	"github.com/schollz/progressbar/v3"
)

// Process all folders and files recursively with progress
func processFiles(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %w", folderPath, err)
	}

	// Calculate total files to process for the progress bar
	var totalFiles int
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalFiles++
		}
		return nil
	})

	// Initialize progress bar
	bar := progressbar.New(totalFiles)

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
				// fmt.Printf("Skipping unsupported file type: %s\n", file.Name())
			}
			// Update progress bar
			bar.Add(1)
		}
	}
	return nil
}

// Delete duplicate files based on MD5 hash with progress
func deleteDuplicates(folderPath string) error {
	hashes := make(map[string]string)
	var totalFiles int

	// Count total files to process for the progress bar
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalFiles++
		}
		return nil
	})

	// Initialize progress bar
	bar := progressbar.New(totalFiles)

err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}

	// Open the file within a limited scope
	func() error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			return err
		}
		hashStr := fmt.Sprintf("%x", hash.Sum(nil))

		if _, found := hashes[hashStr]; found {
			// fmt.Printf("Duplicate found: %s and %s\n", existingPath, path)
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to delete file %s: %w", path, err)
			}
		} else {
			hashes[hashStr] = path
		}
		return nil
	}()

	// Update progress bar
	bar.Add(1)
	return nil
})

	return err
}

func main() {
	for  {
	var choice int
	fmt.Println("Select an option:")
	fmt.Println("1. Parse and rename photos and videos")
	fmt.Println("2. Delete duplicates")
	fmt.Println("3. Exit")
	fmt.Scan(&choice)

	folderPath, err := os.Getwd()

	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}


	switch choice {
	case 1:
		fmt.Println("Starting processing...")
		if err := processFiles(folderPath); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("Processing complete.")
	case 2:
		fmt.Println("Deleting duplicates...")
		if err := deleteDuplicates(folderPath); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("Duplicate deletion complete.")
	case 3:
		fmt.Println("Exiting program.")
		os.Exit(0)
	default:
		fmt.Println("Invalid option. Exiting program.")
		os.Exit(1)
	}

	fmt.Println("---------------------------------")
}
}
