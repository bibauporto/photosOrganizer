package main


import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	
	helpers "github.com/bibauporto/photosOrganizer/helpers"
	"github.com/rwcarlsen/goexif/exif"
)

// Supported file extensions
var IMAGE_EXTENSIONS = []string{".jpg", ".jpeg", ".heic"}
var VIDEO_EXTENSIONS = []string{".mp4", ".mov"}

var dateRegex = regexp.MustCompile(`(\d{4})[._-]?(\d{2})[._-]?(\d{2})(?:[._-]?(\d{2}))?(?:[._-]?(\d{2}))?(?:[._-]?(\d{2}))?`)

// Helper function to check if the date is valid
func isValidDate(year, month, day int) bool {
	return year >= 1970 && year <= 2050 && month >= 1 && month <= 12 && day >= 1 && day <= 31
}

// Generate a unique file name if a file with the same name already exists
func generateUniqueFileName(folderPath, baseName, ext string) (string, error) {
	uniqueName := baseName
	counter := 1
	for {
		if _, err := os.Stat(filepath.Join(folderPath, uniqueName+ext)); os.IsNotExist(err) {
			break
		}
		uniqueName = fmt.Sprintf("%s_%d", baseName, counter)
		counter++
	}
	return uniqueName, nil
}

// Set EXIF DateTimeOriginal in a JPEG file
func setExifDateTaken(filePath, dateTime string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	exifData, err := exif.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding EXIF: %v", err)
	}

	// Set the EXIF DateTimeOriginal
	err = exifData.Set(exif.DateTimeOriginal, dateTime)
	if err != nil {
		return fmt.Errorf("error setting EXIF date: %v", err)
	}

	// Save the updated EXIF data
	// Note: Saving EXIF data back to a JPEG file in Go may require using an additional library
	// or rewriting the JPEG with updated EXIF headers.
	// This is a simplified placeholder.
	return nil
}

// Process Image Files
func processImage(file, folderPath string) error {
	ext := strings.ToLower(filepath.Ext(file))
	filePath := filepath.Join(folderPath, file)

	// Check for EXIF data
	exifDate, err := getExifDateTaken(filePath)
	if err == nil && exifDate != "" {
		fmt.Printf("Skipping already named image: %s\n", file)
		return nil
	}

	match := dateRegex.FindStringSubmatch(file)
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

	if !isValidDate(year, month, day) {
		fmt.Printf("Invalid date in filename: %s\n", file)
		return nil
	}

	baseName := fmt.Sprintf("%04d-%02d-%02d %02s.%02s.%02s", year, month, day, hour, minute, second)
	newFileName, err := generateUniqueFileName(folderPath, baseName, ext)
	if err != nil {
		return err
	}

	newFilePath := filepath.Join(folderPath, newFileName+ext)
	dateTime := fmt.Sprintf("%04d:%02d:%02d %02s:%02s:%02s", year, month, day, hour, minute, second)

	// Set EXIF DateTaken and rename the file
	if err := setExifDateTaken(filePath, dateTime); err != nil {
		return err
	}
	if err := os.Rename(filePath, newFilePath); err != nil {
		return err
	}

	fmt.Printf("Renamed image: %s -> %s\n", file, newFileName+ext)
	return nil
}

// Process Video Files
func processVideo(file, folderPath string) error {
	ext := strings.ToLower(filepath.Ext(file))
	filePath := filepath.Join(folderPath, file)

	match := dateRegex.FindStringSubmatch(file)

	var parsedDate time.Time
	if match != nil {
		year, _ := strconv.Atoi(match[1])
		month, _ := strconv.Atoi(match[2])
		day, _ := strconv.Atoi(match[3])
		hour := matchOrDefault(match[4], "14")
		minute := matchOrDefault(match[5], "00")
		second := matchOrDefault(match[6], "00")

		if isValidDate(year, month, day) {
			parsedDate = time.Date(year, time.Month(month), day, atoi(hour), atoi(minute), atoi(second), 0, time.UTC)
		} else {
			fmt.Printf("Invalid date in filename: %s. Using file's modified date.\n", file)
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return err
			}
			parsedDate = fileInfo.ModTime()
		}
	} else {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		parsedDate = fileInfo.ModTime()
		fmt.Printf("Using file's modified date for video: %s\n", file)
	}

	year := parsedDate.Year()
	month := fmt.Sprintf("%02d", parsedDate.Month())
	day := fmt.Sprintf("%02d", parsedDate.Day())
	hour := fmt.Sprintf("%02d", parsedDate.Hour())
	minute := fmt.Sprintf("%02d", parsedDate.Minute())
	second := fmt.Sprintf("%02d", parsedDate.Second())

	baseName := fmt.Sprintf("%d-%s-%s %s.%s.%s", year, month, day, hour, minute, second)
	newFileName, err := generateUniqueFileName(folderPath, baseName, ext)
	if err != nil {
		return err
	}

	newFilePath := filepath.Join(folderPath, newFileName+ext)
	if err := os.Rename(filePath, newFilePath); err != nil {
		return err
	}

	fmt.Printf("Renamed video: %s -> %s\n", file, newFileName+ext)
	return nil
}

// Get EXIF DateTimeOriginal (for illustration, actual implementation depends on the library you use)
func getExifDateTaken(filePath string) (string, error) {
	// Implement EXIF extraction logic based on Go's EXIF library of your choice.
	// Returning an empty string for the sake of the example.
	return "", nil
}

// Match or provide a default value
func matchOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// atoi is a simple helper to convert a string to an int
func atoi(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

// Process files in directory
func processFiles(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if contains(IMAGE_EXTENSIONS, ext) {
			if err := processImage(file.Name(), folderPath); err != nil {
				return err
			}
		} else if contains(VIDEO_EXTENSIONS, ext) {
			if err := processVideo(file.Name(), folderPath); err != nil {
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


