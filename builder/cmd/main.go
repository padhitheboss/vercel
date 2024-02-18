package main

import (
	"fmt"
	"log"
	"os"

	"github.com/padhitheboss/code-builder/pkg/service/azureBlob"
	pkgbuilder "github.com/padhitheboss/code-builder/pkg/service/buildStatic"
	ziphandler "github.com/padhitheboss/code-builder/pkg/service/zipHandler"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println(err)
	// }
	// queuemodel.StartQueue()
	// reqFromQueue, _ := queuemodel.CollectRequest()
	// fmt.Println(reqFromQueue)
	// var g git.GitService = git.CreateGitService(reqFromQueue)
	// err = g.Clone()
	// if err != nil {
	// 	log.Panicf("unable to clone git repo %v", err)
	// 	return
	// }
	// b := azureBlob.CreateConfig(reqFromQueue.RequestId)
	b := azureBlob.CreateConfig("", os.Getenv("REQUEST_ID"))
	downloadPath, err := b.DownloadFile(os.Getenv("BLOB_INPUT_FOLDER_PATH"), os.Getenv("DOWNLOAD_FOLDER_PATH"))
	if err != nil {
		log.Panicf("unable to clone git repo %v", err)
		return
	}
	fmt.Println(downloadPath)
	unzipFilePath := ziphandler.UnzipFile(downloadPath, os.Getenv("BUILD_FOLDER_PATH"))
	pkg := pkgbuilder.CreateConfig(os.Getenv("PROJECT_TYPE"), unzipFilePath)
	pkg.BuildStatic()
	fmt.Println(pkg.GetOutputPath())
	b.UploadFolder(pkg.GetOutputPath())
	// // b.DeleteFolder()
	if err != nil {
		log.Println(err)
	}
}
