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

func (c *BlobConfig) UploadFile(localFilePath, relativePath string) (string, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open local file: %v", err)
	}
	contentType, _ := detectContentType(localFilePath)
	defer file.Close()
	// Upload the file to Azure Blob Storage
	uploadPath := path.Join(c.blobName, relativePath)
	_, err = c.client.UploadFile(context.TODO(), c.containerName, uploadPath, file, &azblob.UploadFileOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr(contentType),
		},
	})
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(detectContentType(localFilePath))
	fmt.Printf("Uploaded file %s\n", localFilePath)
	return fmt.Sprint(strings.TrimPrefix(path.Join(c.containerName, uploadPath), os.Getenv("BLOB_CONTAINER_NAME")+"/")), nil
}

func (c *BlobConfig) UploadFolder(localFolderPath string) (string, error) {
	err := filepath.Walk(localFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Get the relative path of the file
			relativePath := strings.TrimPrefix(path, localFolderPath)
			_, err = c.UploadFile(path, relativePath)
			return err
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to walk local folder: %v", err)
	}

	fmt.Printf("Uploaded folder %s\n", localFolderPath)
	return fmt.Sprint(strings.TrimPrefix(path.Join(c.containerName, c.blobName), os.Getenv("BLOB_CONTAINER_NAME")+"/")), nil
}

func (c *BlobConfig) DeleteFolder() {
	_, err := c.client.DeleteBlob(context.TODO(), c.containerName, c.blobName, nil)
	log.Println(err)
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
