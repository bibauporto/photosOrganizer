package processors

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bibauporto/photosOrganizer/exif"
	helpers "github.com/bibauporto/photosOrganizer/helpers"
)

// Process Image Files
func ProcessImage(file, folderPath string) error {
	ext := strings.ToLower(filepath.Ext(file))
	filePath := filepath.Join(folderPath, file)

	exifDate, err := exif.GetExifDateTaken(filePath)
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
	hour := helpers.MatchOrDefault(match[4], "14")
	minute := helpers.MatchOrDefault(match[5], "00")
	second := helpers.MatchOrDefault(match[6], "00")

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

	if err := exif.SetExifDateTaken(filePath, dateTime); err != nil {
		return err
	}
	if err := os.Rename(filePath, newFilePath); err != nil {
		return err
	}

	fmt.Printf("Renamed image: %s -> %s\n", file, newFileName+ext)
	return nil
}