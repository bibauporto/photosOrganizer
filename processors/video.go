package processors

import (
	"fmt"
	"os"
	"path/filepath"

	helpers "github.com/bibauporto/photosOrganizer/helpers"
)

func ProcessVideo(file, folderPath string) error {
	// try to parse the name of the file
	match := helpers.DateParserRegex.FindStringSubmatch(file)
	if match == nil {
		fmt.Printf("No date in filename: %s\n", file)
		return nil
	} else {
		// use modified date of the file to set the name of the file
		filePath := filepath.Join(folderPath, file)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		year := fileInfo.ModTime().Year()
		month := int(fileInfo.ModTime().Month())
		day := fileInfo.ModTime().Day()
		hour := fileInfo.ModTime().Hour()
		minute := fileInfo.ModTime().Minute()
		second := fileInfo.ModTime().Second()

		baseName := fmt.Sprintf("%04d-%02d-%02d %02d.%02d.%02d", year, month, day, hour, minute, second)
		newFileName, err := helpers.GenerateUniqueFileName(folderPath, baseName, filepath.Ext(file))
		if err != nil {
			return err
		}

		newFilePath := filepath.Join(folderPath, newFileName+filepath.Ext(file))
		if err := os.Rename(filePath, newFilePath); err != nil {
			return err
		}

	}
	return nil
}
