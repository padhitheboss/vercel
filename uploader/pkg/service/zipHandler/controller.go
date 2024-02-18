package ziphandler

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateZipFile(zipFileName, folderPath string) (string, error) {
	err := os.MkdirAll(filepath.Dir(zipFileName), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return "", err
	}
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		fmt.Println("Error creating zip file:", err)
		return "", err
	}
	defer zipFile.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the directory and add files to the zip
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}
		// Create a new file in the zip archive
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Open the original file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the contents of the file to the zip archive
		_, err = io.Copy(zipFile, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error zipping directory:", err)
		return "", err
	}
	fmt.Println("Directory zipped successfully")
	return zipFileName, nil
}

func UnzipFile(zipFilePath, outputDir string) string {
	// Open the zipped file
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return ""
	}
	defer zipFile.Close()

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return ""
	}

	// Extract files from the zip archive
	for _, file := range zipFile.File {
		// Open the file from the zip archive
		zipFile, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file from zip:", err)
			return ""
		}
		defer zipFile.Close()

		// Create the output file in the output directory
		outputFilePath := filepath.Join(outputDir, file.Name)
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return ""
		}
		defer outputFile.Close()

		// Copy the contents of the file from the zip archive to the output file
		_, err = io.Copy(outputFile, zipFile)
		if err != nil {
			fmt.Println("Error copying file contents:", err)
			return ""
		}
	}

	fmt.Println("Files unzipped successfully to:", outputDir)
	return outputDir
}
