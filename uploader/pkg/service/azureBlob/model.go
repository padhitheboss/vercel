package azureBlob

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type BlobConfig struct {
	accountName   string
	accountKey    string
	containerName string
	blobName      string
	client        *azblob.Client
	folderPath    string
}

func CreateConfig(remotePath string) BlobConfig {
	var b BlobConfig
	b.accountName = os.Getenv("BLOB_ACCOUNT_NAME")
	b.accountKey = os.Getenv("BLOB_ACCOUNT_KEY")
	b.folderPath = path.Join(os.Getenv("BLOB_ROOT_FOLDER_PATH"), remotePath)
	b.containerName = path.Join(os.Getenv("BLOB_CONTAINER_NAME"), b.folderPath)
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
