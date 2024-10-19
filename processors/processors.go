package processors

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bibauporto/photosOrganizer/exif"
	"github.com/bibauporto/photosOrganizer/helpers"
)

// ProcessImage processes and renames an image file based on its EXIF data or filename.
func ProcessImage(file, folderPath string) error {
	filePath := filepath.Join(folderPath, file)

	if helpers.CorrectNameRegex.MatchString(file) {
		fmt.Printf("Skipping already named image: %s\n", file)
		return nil
	}

	// Try to get the date taken from EXIF data
	if dateTaken, err := exif.GetExifDateTaken(filePath); err == nil {
		return renameFile(filePath, folderPath, formatDateTime(dateTaken))
	}

	// Try to parse the date from the filename
	return tryRenameByFilename(file, folderPath, filePath)
}

// ProcessVideo processes and renames a video file based on its filename or modified date.
func ProcessVideo(file, folderPath string) error {
	filePath := filepath.Join(folderPath, file)

	if helpers.CorrectNameRegex.MatchString(file) {
		fmt.Printf("Skipping already named video: %s\n", file)
		return nil
	}

	// Try to parse the date from the filename
	if err := tryRenameByFilename(file, folderPath, filePath); err == nil {
		return nil
	}

	// Use the file modification time as the fallback
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	modTime := fileInfo.ModTime()
	return renameFile(filePath, folderPath, formatDateTime(modTime))
}

// Try to rename a file based on the parsed date from its filename.
func tryRenameByFilename(file, folderPath, filePath string) error {
	match := helpers.DateParserRegex.FindStringSubmatch(file)
	if match == nil {
		fmt.Printf("No date in filename: %s\n", file)
		return nil
	}

	year, month, day := match[1], match[2], match[3]
	hour, minute, second := helpers.MatchOrDefault(match[4], "14"), helpers.MatchOrDefault(match[5], "00"), helpers.MatchOrDefault(match[6], "00")
	baseName := fmt.Sprintf("%s-%s-%s %s.%s.%s", year, month, day, hour, minute, second)

	return renameFile(filePath, folderPath, baseName)
}

// Rename a file and ensure the new name is unique.
func renameFile(filePath, folderPath, baseName string) error {
	ext := filepath.Ext(filePath)
	newFileName, err := helpers.GenerateUniqueFileName(folderPath, baseName, ext)
	if err != nil {
		return fmt.Errorf("error generating unique name: %w", err)
	}

	newFilePath := filepath.Join(folderPath, newFileName+ext)
	if err := os.Rename(filePath, newFilePath); err != nil {
		return fmt.Errorf("error renaming file: %w", err)
	}

	fmt.Printf("Renamed file: %s -> %s\n", filepath.Base(filePath), newFileName+ext)
	return nil
}

// Format time.Time to a string for renaming files.
func formatDateTime(date time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d.%02d.%02d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
}
