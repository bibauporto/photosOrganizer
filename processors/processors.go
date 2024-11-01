package processors

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bibauporto/photosOrganizer/exif"
	"github.com/bibauporto/photosOrganizer/helpers"
)

// ProcessImage processes and renames an image file based on its EXIF data or filename.
func ProcessImage(file, folderPath string) error {
	filePath := filepath.Join(folderPath, file)

	// Skip renaming if the file is already named correctly
	if helpers.CorrectNameRegex.MatchString(file) {
		return nil
	}

	// Try to get the date taken from EXIF data
	if dateTaken, err := exif.GetExifDateTaken(filePath); err == nil {
		// Rename the file using EXIF date
		return renameFile(filePath, folderPath, formatDateTime(dateTaken))
	}

	// Try to parse the date from the filename if EXIF data is unavailable
	return tryRenameByFilename(file, folderPath, filePath)
}

// ProcessVideo processes and renames a video file based on its filename or modified date.
func ProcessVideo(file, folderPath string) error {
	filePath := filepath.Join(folderPath, file)

	// Check if the file is already named correctly
	if helpers.CorrectNameRegex.MatchString(file) {
		// Extract date from the filename and update the modified date if necessary
		if parsedDate, err := parseDateFromFilename(file); err == nil {
			return updateModifiedDateIfNeeded(filePath, parsedDate)
		} else {
			fmt.Printf("Error parsing date from correctly named video: %s\n", file)
			return err
		}
	}

	// Try to parse the date from the filename, if it's not already in the correct format
	if err := tryRenameByFilename(file, folderPath, filePath); err == nil {
		return nil
	}

	// Use the file's modified date to rename the file as a fallback
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	modTime := fileInfo.ModTime()
	return renameFile(filePath, folderPath, formatDateTime(modTime))
}

// Try to rename a file based on the parsed date from its filename.
func tryRenameByFilename(file, folderPath, filePath string) error {
	// Parse date from the filename
	parsedDate, err := parseDateFromFilename(file)
	if err != nil {
		return nil
	}

	// Update the modified date if needed, then rename the file
	if err := updateModifiedDateIfNeeded(filePath, parsedDate); err != nil {
		return fmt.Errorf("error updating modified date: %w", err)
	}

	// Format and rename the file using the parsed date
	baseName := formatDateTime(parsedDate)
	return renameFile(filePath, folderPath, baseName)
}

// Update the modified date of the file if it differs from the parsed date.
func updateModifiedDateIfNeeded(filePath string, parsedDate time.Time) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	modTime := fileInfo.ModTime()
	if !modTime.Equal(parsedDate) {
		// Update the modified date to match the parsed date
		if err := os.Chtimes(filePath, modTime, parsedDate); err != nil {
			return fmt.Errorf("error setting modified date: %w", err)
		}
	}
	return nil
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

	return nil
}

// Parse the date from the filename.
func parseDateFromFilename(file string) (time.Time, error) {
	match := helpers.DateParserRegex.FindStringSubmatch(file)
	if match == nil {
		return time.Time{}, fmt.Errorf("no date found in filename")
	}

	year, month, day := match[1], match[2], match[3]
	hour, minute, second := helpers.MatchOrDefault(match[4], "14"), helpers.MatchOrDefault(match[5], "00"), helpers.MatchOrDefault(match[6], "00")

	// Convert the extracted date and time strings to integers
	parsedDate, err := constructDateFromParts(year, month, day, hour, minute, second)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}

// Construct a time.Time object from string parts.
func constructDateFromParts(year, month, day, hour, minute, second string) (time.Time, error) {
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return time.Time{}, err
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return time.Time{}, err
	}
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return time.Time{}, err
	}
	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		return time.Time{}, err
	}
	minuteInt, err := strconv.Atoi(minute)
	if err != nil {
		return time.Time{}, err
	}
	secondInt, err := strconv.Atoi(second)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minuteInt, secondInt, 0, time.Local), nil
}

// Format time.Time to a string for renaming files.
func formatDateTime(date time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d.%02d.%02d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
}
