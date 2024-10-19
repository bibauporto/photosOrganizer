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