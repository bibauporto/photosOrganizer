package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	helpers "github.com/bibauporto/photosOrganizer/helpers"
	processors "github.com/bibauporto/photosOrganizer/processors"

	"github.com/rwcarlsen/goexif/exif"
)

// Set EXIF DateTimeOriginal in a JPEG file
func setExifDateTaken(filePath, dateTime string) error {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening image: %v", err)
	}
	defer imgFile.Close()

	// Decode the image
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("error decoding image: %v", err)
	}

	// Create a new EXIF data structure
	x := exif.New()
	x.Set("Exif.Image.DateTimeOriginal", dateTime)

	// Save the new EXIF data
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Write the image with new EXIF data
	if err := jpeg.Encode(outFile, img, nil); err != nil {
		return fmt.Errorf("error encoding image: %v", err)
	}

	return nil
}
// Get EXIF DateTimeOriginal
func getExifDateTaken(filePath string) (string, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening image: %v", err)
	}
	defer imgFile.Close()

	// Read EXIF data
	exifData, err := exif.Decode(imgFile)
	if err != nil {
		return "", nil
	}

	// Get DateTimeOriginal
	val, err := exifData.Get(exif.DateTimeOriginal)
	if err != nil {
		return "", nil
	}

	return val.String(), nil
}

// Process Image Files
func processImage(file, folderPath string) error {
	ext := strings.ToLower(filepath.Ext(file))
	filePath := filepath.Join(folderPath, file)

	exifDate, err := getExifDateTaken(filePath)
	if err == nil && exifDate != "" {
		fmt.Printf("Skipping already named image: %s\n", file)
		return nil
	}

	match := helpers.DateParserRegex.FindStringSubmatch(file)
	if match == nil {
		fmt.Printf("No date in filename: %s\n", file)
		return nil
	}

	year, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	day, _ := strconv.Atoi(match[3])
	hour := matchOrDefault(match[4], "14")
	minute := matchOrDefault(match[5], "00")
	second := matchOrDefault(match[6], "00")

	if !helpers.IsValidDate(year, month, day) {
		fmt.Printf("Invalid date in filename: %s\n", file)
		return nil
	}

	baseName := fmt.Sprintf("%04d-%02d-%02d %02s.%02s.%02s", year, month, day, hour, minute, second)
	newFileName, err := helpers.GenerateUniqueFileName(folderPath, baseName, ext)
	if err != nil {
		return err
	}

	newFilePath := filepath.Join(folderPath, newFileName+ext)
	dateTime := fmt.Sprintf("%04d:%02d:%02d %02s:%02s:%02s", year, month, day, hour, minute, second)

	if err := setExifDateTaken(filePath, dateTime); err != nil {
		return err
	}
	if err := os.Rename(filePath, newFilePath); err != nil {
		return err
	}

	fmt.Printf("Renamed image: %s -> %s\n", file, newFileName+ext)
	return nil
}

// Match or provide a default value
func matchOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Process all files in a directory
func processFiles(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if contains(helpers.IMAGE_EXTENSIONS, ext) {
			if err := processImage(file.Name(), folderPath); err != nil {
				return err
			}
		} else if contains(helpers.VIDEO_EXTENSIONS, ext) {
			if err := processors.ProcessVideo(file.Name(), folderPath); err != nil {
				return err
			}
		} else {
			fmt.Printf("Skipping unsupported file type: %s\n", file.Name())
		}
	}
	return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
}

func main() {
	folderPath, _ := os.Getwd()
	fmt.Println("Starting processing...")
	if err := processFiles(folderPath); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("Processing complete.")
}