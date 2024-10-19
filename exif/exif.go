package exif

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// openImageFile opens an image file and returns the file handle.
func openImageFile(filePath string) (*os.File, error) {
    imgFile, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("error opening image: %v", err)
    }
    return imgFile, nil
}

// decodeExifData decodes EXIF data from an image file.
func decodeExifData(imgFile *os.File) (*exif.Exif, error) {
    exifData, err := exif.Decode(imgFile)
    if err != nil {
        return nil, fmt.Errorf("error decoding EXIF data: %v", err)
    }

    return exifData, nil
}


// GetExifDateTaken gets the EXIF DateTimeOriginal from a JPEG file.
func GetExifDateTaken(filePath string) (time.Time, error) {
    imgFile, err := openImageFile(filePath)
    if err != nil {
        return time.Time{}, err
    }
    defer imgFile.Close()

    // Read EXIF data
    exifData, err := decodeExifData(imgFile)
    if err != nil {
        return time.Time{}, err
    }

    // Get DateTimeOriginal
    val, err := exifData.DateTime()
    if err != nil {
        return time.Time{}, fmt.Errorf("error getting DateTimeOriginal: %v", err)
    }

    return val, nil
}