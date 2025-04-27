package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

// Helper function to convert byte buffer to binary string
func BufferToBinaryString(buffer []byte) string {
	return string(buffer)
}

// Generate a unique file name
func GenerateUniqueFileName(folderPath, baseName, ext string) (string, error) {
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


// Check if the date is valid
func IsValidDate(year, month, day int) bool {
	return year >= 1970 && year <= 2050 && month >= 1 && month <= 12 && day >= 1 && day <= 31
}


// Helper function to check if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
}

// Match or provide a default value
func MatchOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
