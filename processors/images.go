package processors

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bibauporto/photosOrganizer/exif"
	"github.com/bibauporto/photosOrganizer/helpers"
)

// ProcessImage processes and renames an image file based on its EXIF data or filename.
func ProcessImage(file, folderPath string) error {
	ext := strings.ToLower(filepath.Ext(file))
	filePath := filepath.Join(folderPath, file)

	// Check if the file is already named correctly
	if helpers.CorrectNameRegex.MatchString(file) {
		fmt.Printf("Skipping already named image: %s\n", file)
		return nil
	}

	// Try to get the date taken from EXIF data
	dateTaken, err := exif.GetExifDateTaken(filePath)
	if err == nil {
		// Rename the file based on EXIF date taken
		baseName := formatDateTime(dateTaken)
		return renameFile(filePath, folderPath, baseName, ext)
	}

	// Try to parse the date from the filename
	match := helpers.DateParserRegex.FindStringSubmatch(file)
	if match == nil {
		fmt.Printf("No date in filename: %s\n", file)
		return nil
	}

	year, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	day, _ := strconv.Atoi(match[3])
	hour := helpers.MatchOrDefault(match[4], "14")
	minute := helpers.MatchOrDefault(match[5], "00")
	second := helpers.MatchOrDefault(match[6], "00")

	if !helpers.IsValidDate(year, month, day) {
		fmt.Printf("Invalid date in filename: %s\n", file)
		return nil
	}

	baseName := fmt.Sprintf("%04d-%02d-%02d %02s.%02s.%02s", year, month, day, hour, minute, second)
	return renameFile(filePath, folderPath, baseName, ext)
}

// formatDateTime formats the given time.Time into a string for the new filename.
func formatDateTime(dateTaken time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d.%02d.%02d",
		dateTaken.Year(),
		int(dateTaken.Month()),
		dateTaken.Day(),
		dateTaken.Hour(),
		dateTaken.Minute(),
		dateTaken.Second(),
	)
}

// renameFile renames the file based on the new base name and ensures uniqueness.
func renameFile(filePath, folderPath, baseName, ext string) error {
	newFileName, err := helpers.GenerateUniqueFileName(folderPath, baseName, ext)
	if err != nil {
		return err
	}

	newFilePath := filepath.Join(folderPath, newFileName+ext)
	if err := os.Rename(filePath, newFilePath); err != nil {
		return fmt.Errorf("error renaming file: %v", err)
	}

	fmt.Printf("Renamed image: %s -> %s\n", filepath.Base(filePath), newFileName+ext)
	return nil
}
