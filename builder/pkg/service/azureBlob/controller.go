package azureBlob

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/gabriel-vasile/mimetype"
)

func (c *BlobConfig) UploadFile(localFilePath, relativePath string) error {
	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	contentType, _ := detectContentType(localFilePath)
	defer file.Close()
	// Upload the file to Azure Blob Storage
	uploadPath := c.blobName + relativePath
	_, err = c.client.UploadFile(context.TODO(), c.destContainerName, uploadPath, file, &azblob.UploadFileOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr(contentType),
		},
	})
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(detectContentType(localFilePath))
	fmt.Printf("Uploaded file %s\n", localFilePath)
	return nil
}

func (c *BlobConfig) DownloadFile(blobPath, destinationPath string) (string, error) {
	// Download the file from Azure Blob Storage
	if err := os.MkdirAll(destinationPath, os.ModePerm); err != nil {
		fmt.Println("Error creating parent directories:", err)
		return "", nil
	}
	target := path.Join(destinationPath, path.Base(blobPath))
	file, err := os.Create(target)
	if err != nil {
		return "", err
	}
	fmt.Println(target, c.srcContainerName, blobPath)
	defer file.Close()
	_, err = c.client.DownloadFile(context.TODO(), c.srcContainerName, blobPath, file, nil)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(detectContentType(localFilePath))
	fmt.Printf("Downloaded file %s\n", target)
	return target, nil
}

func (c *BlobConfig) UploadFolder(localFolderPath string) error {
	err := filepath.Walk(localFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Get the relative path of the file
			relativePath := strings.TrimPrefix(path, localFolderPath+"/")
			return c.UploadFile(path, relativePath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk local folder: %v", err)
	}

	fmt.Printf("Uploaded folder %s\n", localFolderPath)
	return nil
}

func detectContentType(filePath string) (string, error) {
	mime, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", err
	}
	cType := mime.String()
	s := strings.Split(filePath, ".")
	extention := s[len(s)-1]
	if extention == "js" {
		cType = "text/javascript"
	} else if extention == "css" {
		cType = "text/css"
	} else if extention == "html" {
		cType = "text/html"
	}
	return cType, nil
}
