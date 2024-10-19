package exif

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

// Set EXIF DateTimeOriginal in a JPEG file
func SetExifDateTaken(filePath, dateTime string) error {
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
func GetExifDateTaken(filePath string) (string, error) {
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
