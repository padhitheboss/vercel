package azblob

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type BlobConfig struct {
	accountName   string
	accountKey    string
	containerName string
	client        *azblob.Client
}

func CreateConfig() BlobConfig {
	var b BlobConfig
	b.accountName = os.Getenv("BLOB_ACCOUNT_NAME")
	b.accountKey = os.Getenv("BLOB_ACCOUNT_KEY")
	// b.folderPath = os.Getenv("BLOB_ROOT_FOLDER_PATH")
	b.containerName = os.Getenv("BLOB_CONTAINER_NAME")
	cred, err := azblob.NewSharedKeyCredential(b.accountName, b.accountKey)
	if err != nil {
		log.Panicf("failed to create shared key credential: %v", err)
	}
	b.client, err = azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net", b.accountName), cred, nil)
	if err != nil {
		log.Panicf("failed to create shared key credential: %v", err)
	}
	return b
}

func (b *BlobConfig) DownloadStream(filePath string) (io.ReadCloser, error) {
	res, err := b.client.DownloadStream(context.TODO(), b.containerName, filePath, nil)
	if err != nil {
		fmt.Printf("error downloading blob: %v", err)
		return nil,err
	}
	// defer res.Body.Close()
	return res.Body, err
}
